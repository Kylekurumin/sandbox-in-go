package main

import (
	"fmt"
	"syscall"
	"time"
)

func main() {
	rlimit := &syscall.Rlimit{
		Cur: 256 << 20,
		Max: 256 << 20,
	}
	err := syscall.Setrlimit(syscall.RLIMIT_DATA, rlimit)
	if err != nil {
		panic(err)
	}
	pid, err := syscall.ForkExec("/usr/bin/python3", []string{"/usr/bin/python3", "hello.py"}, nil)
	if err != nil {
		panic(err)
	}
	var st *syscall.WaitStatus
	var usage = &syscall.Rusage{}

	syscall.Wait4(pid, st, 0, usage)
	fmt.Println(usage.Utime.Nano() / time.Millisecond.Nanoseconds())
	fmt.Println(usage.Maxrss)

}
