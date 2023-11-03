
```
This is a very experimental code and should not be used in any production environment. 
```

## Overview
This repository contains experimental code that demonstrates the ability to bind goroutines to specific CPU cores in Go. The primary component of this project, `CoreManager`, is designed to allow to schedule goroutines with affinity to particular cores.

## How It Works

CoreManager works by assigning goroutines to specific cores. When a user wishes to schedule a goroutine with core affinity, the following process occurs:

* The function intended for the goroutine is sent to a dedicated channel.
* A specialized goroutine, which has been pre-bound to the targeted core, retrieves the function from the channel.
* This bound goroutine then executes the function as a goroutine.
* Since we don't know which core a goroutine will be executed on, we use pthread_setaffinity_np to lock the thread to a specific core.

## Core locking mechanism
Every core in CoreManager has a dedicated channel and a goroutine that is locked to that core which uses also runtime.LockOSThread() to lock itself to a kernel thread. But this goroutine is not used to execute any user function. It is only used to pick up the user function from the channel and execute it as a goroutine.

Due to the non-deterministic nature of the Go scheduler, there's no guarantee which core a new goroutine will be executed on. To enforce core affinity, we bypass Go's runtime scheduler using the pthread_setaffinity_np function to lock the thread to a specific core. This ensures that the goroutine will be executed on the desired core. But we don't use runtime.LockOsThread() for new goroutines as it will lock the current goroutine to its kernel thread but no other goroutine will be locked to that thread. If we would have used runtime.LockOsThread() then we would have to create a new kernel thread for every goroutine that we want to lock to a core and using a kernel thread for every goroutine is very expensive and causes C100k, lot of context switching and other problems.
