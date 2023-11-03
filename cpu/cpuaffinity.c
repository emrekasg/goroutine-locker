#define _GNU_SOURCE

/*
    Locks the calling thread to the given CPU.
    cpuid: CPU to lock the thread to.
*/
#if defined(__linux__)
#include <sched.h>
#include <pthread.h>
void lock_thread(int cpuid)
{
    pthread_t tid;
    cpu_set_t cpuset;

    tid = pthread_self();
    CPU_ZERO(&cpuset);
    CPU_SET(cpuid, &cpuset);

    // https://linux.die.net/man/3/pthread_setaffinity_np
    pthread_setaffinity_np(tid, sizeof(cpu_set_t), &cpuset);
}
#elif defined(_WIN32)
#include <windows.h>
int lock_thread(int core_id)
{
    DWORD_PTR mask = 1ULL << core_id;
    return !SetThreadAffinityMask(GetCurrentThread(), mask);
}
#elif defined(__APPLE__)
// MacOS does not support setting thread affinity in its kernel API
// See more: https://developer.apple.com/library/archive/releasenotes/Performance/RN-AffinityAPI/index.html
#error "MacOS is not supported"
#else
#error "Unsupported platform"
#endif
