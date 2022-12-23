package main

import (
	"os"
	"reflect"
	"time"
	"unsafe"
)

func HideName(pname string) {
	argv0str := (*reflect.StringHeader)(unsafe.Pointer(&os.Args[0]))
	argv0 := (*[1 << 30]byte)(unsafe.Pointer(argv0str.Data))[:argv0str.Len]

	n := copy(argv0, pname)
	if n < len(argv0) {
		for i := n; i < len(argv0); i++ {
			argv0[i] = 0
		}
	}
}

func main() {
	HideName("psovaya")
	time.Sleep(10 * time.Second)
}
