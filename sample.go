package main

import (
	"fmt"
	libseccomp "github.com/seccomp/libseccomp-golang"
	"os"
	"os/exec"
	"syscall"
)

func whiteList(syscalls []string) {

	filter, err := libseccomp.NewFilter(libseccomp.ActErrno.SetReturnCode(int16(syscall.EPERM)))
	if err != nil {
		fmt.Printf("Error creating filter: %s\n", err)
	}
	for _, element := range syscalls {
		fmt.Printf("[+] Whitelisting: %s\n", element)
		syscallID, err := libseccomp.GetSyscallFromName(element)
		if err != nil {
			panic(err)
		}
		filter.AddRule(syscallID, libseccomp.ActAllow)
	}
	filter.Load()
}

func main() {
	cmd := exec.Command("/bin/ls", "-l")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	//var (
	//	wstatus unix.WaitStatus // wait4 wait status
	//	rusage  unix.Rusage     // wait4 rusage
	//)

	//pid, _ = unix.Wait4(pid, &wstatus, unix.WALL, &rusage)

	//fmt.Println("First process exited with status: ", wstatus.ExitStatus())

	var syscalls = []string{
		"rt_sigaction", "mkdirat", "clone", "mmap", "readlinkat", "futex", "rt_sigprocmask",
		"mprotect", "write", "sigaltstack", "gettid", "read", "open", "close", "fstat", "fork",
		"munmap", "brk", "access", "execve", "getrlimit", "arch_prctl", "sched_getaffinity",
		"set_tid_address", "set_robust_list", "exit_group"}

	whiteList(syscalls)

	cmd = exec.Command("/bin/ls", "-l")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

}
