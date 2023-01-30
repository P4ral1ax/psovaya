package watchdog

import "fmt"

var DogHelp string = "watchdog - The watchdog command allows you to check the status of implants on boxes. Can choose to drop the binary if gone\n   status: Display watchdog status of target\n    load: load configuration file for watchdog\n    add: Add single config\n"

type WatchObj struct {
	path   string
	status bool
	name   string
}

func CheckStatus() {

}

func PrintHelp() {
	fmt.Printf(DogHelp)
}
