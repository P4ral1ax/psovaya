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

// Write Content
func WriteToMemfd(fd int, content string) {
	fmt.Printf("Writing to MemFd Pointer")
	filepath := ""
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

func ExecMemfd() {
	pid, _, _ := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
	if pid == 0 {
		print("Executing Memfd Pointer")
		filepath := ""
		unix.Exec(filepath, nil, os.Environ())
	}
}

func main() {
	content := "funny mode"
	fd := MemfdCreate("psovaya")
	WriteToMemfd(fd, content)

}
