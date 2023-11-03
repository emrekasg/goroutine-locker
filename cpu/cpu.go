package cpu

import (
	"runtime"
	"syscall"
)

/*

#if defined(__linux__)
	#define _GNU_SOURCE
	#include "cpuaffinity.h"
	#include <sched.h>
	#include <pthread.h>
#elif defined(__APPLE__)
	#define _POSIX_C_SOURCE
	# error "No implementation for MacOS yet"
#elif defined(_WIN32)
	// #include <windows.h>
	# error "No implementation for Windows yet"
#endif


*/
import (
	"C"
)

type Task func()

type CoreManager struct {
	cores []*Core
}

type Core struct {
	id              int
	taskQueue       chan Task
	totalGoroutines int
	kernelThreads   []KernelThread
}

type KernelThread struct {
	id             int
	goroutineCount int
}

func NewCoreManager() *CoreManager {
	cpuCount := runtime.NumCPU()
	manager := &CoreManager{
		cores: make([]*Core, cpuCount),
	}

	for i := 0; i < cpuCount; i++ {
		manager.cores[i] = newCore(i)
	}

	return manager
}

func newCore(id int) *Core {
	core := &Core{
		id:            id,
		taskQueue:     make(chan Task),
		kernelThreads: make([]KernelThread, 0),
	}
	go func() {
		// locks the calling goroutine to its current operating system thread.
		// with this way, we can't run any other goroutine on this thread.
		runtime.LockOSThread()

		// lock the thread to the given cpu core
		C.lock_thread(C.int(id))

		for task := range core.taskQueue {
			go func() {
				// We're locking the thread again because we're creating a new goroutine
				// and scheduler can run this goroutine on the different core.
				// since we don't want to run it on the different core,
				// lock_thread will move the thread to the given core.
				C.lock_thread(C.int(id))

				// gettid() returns the caller's thread ID (TID).
				// https://man7.org/linux/man-pages/man2/gettid.2.html
				kernelThread := syscall.Gettid()
				core.registerKernelThread(kernelThread)
				task()
			}()
		}
	}()

	return core
}

func (c *Core) registerKernelThread(threadId int) {
	for i, kernelThread := range c.kernelThreads {
		if kernelThread.id == threadId {
			kernelThread.goroutineCount++
			c.kernelThreads[i] = kernelThread
			return
		}
	}

	c.kernelThreads = append(c.kernelThreads, KernelThread{
		id:             threadId,
		goroutineCount: 1,
	})
}

// sched_getcpu() returns the number of the CPU on which the calling thread is currently executing.
// https://man7.org/linux/man-pages/man3/sched_getcpu.3.html
func GetCpuId() int {
	cpu := C.sched_getcpu()
	return int(cpu)
}

func (cm *CoreManager) GetGoRoutineCount() int {
	count := 0
	for _, core := range cm.cores {
		count += core.totalGoroutines
	}
	return count
}

func (cm *CoreManager) GetGoRoutineCountByCpu(cpu int) int {
	return cm.cores[cpu].totalGoroutines
}

func (cm *CoreManager) RunTask(cpu int, task Task) {
	cm.cores[cpu].totalGoroutines++
	cm.cores[cpu].taskQueue <- task
}

func (cm *CoreManager) RunTaskOnAllCores(task Task) {
	for _, core := range cm.cores {
		core.totalGoroutines++
		core.taskQueue <- task
	}
}

type CoreInfo struct {
	CpuId         int
	KernelThreads []KernelThread
}

func (cm *CoreManager) GetCoreInfo() []CoreInfo {
	coreInfo := make([]CoreInfo, 0)

	for _, core := range cm.cores {
		coreInfo = append(coreInfo, CoreInfo{
			CpuId:         core.id,
			KernelThreads: core.kernelThreads,
		})
	}

	return coreInfo
}
