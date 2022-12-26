package main

import (
	"fmt"
	"io/ioutil"

	"github.com/oneNutW0nder/CatTails/cattails"
	"golang.org/x/sys/unix"
)

func banner() {
	b, err := ioutil.ReadFile("ascii.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func main() {
	// Zoi Time
	banner()

	/* Cattails init */
	// Create a BPF vm for filtering
	cattails.CreateBPFVM(cattails.FilterRaw)

	// Create a socket for reading
	readfd := cattails.NewSocket()
	defer unix.Close(readfd)

	fmt.Println("[+] Created sockets")
}
