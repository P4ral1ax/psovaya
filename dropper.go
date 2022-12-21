package main

import (
	"fmt"

	"golang.org/x/sys/unix"
)

// Create a MemFd pointer
func MemfdCreate(name string) int {
	fmt.Printf("Creating memfd %v", name)
	fd, err := unix.MemfdCreate(name, 0)
	if err != nil {
		fmt.Errorf("MemfdCreate: %v", err)
		unix.Exit(0)
	}
	return fd
}

func WriteToMemfd(fd int, c string) {

	return
}

func ExecMemfd() {

	return
}

func main() {
	content := "funny mode"
	fd := MemfdCreate("psovaya")
	WriteToMemfd(fd, content)

}
