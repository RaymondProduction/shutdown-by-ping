package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/slytomcat/systray"
)

const (
	routerIP = "192.168.1.4" //"127.0.0.1"
)

var (
	// icons
	greenIcon  []byte
	redIcon    []byte
	yellowIcon []byte

	errorPingCounter uint
	timeSec uint

	pingTicker *time.Ticker
	stopChan   chan bool
)

func main() {
	greenIcon = getIcon("greenIcon.png")
	redIcon = getIcon("redIcon.png")
	yellowIcon = getIcon("yellowIcon.png")
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(yellowIcon) // replace with your icon file
	// systray.SetTitle("Ping Router")
	systray.SetTooltip("Ping Router Service")

	startItem := systray.AddMenuItem("Start Ping", "Start pinging the router")
	startItem.SetIcon(greenIcon)
	stopItem := systray.AddMenuItem("Stop Ping", "Stop pinging the router")
	stopItem.SetIcon(redIcon)
	quitItem := systray.AddMenuItem("Quit", "Quit the application")

	stopItem.Disable()

	go func() {
		for {
			select {
			case <-startItem.ClickedCh:
				startPinging()
				fmt.Println("startItem")
			case <-stopItem.ClickedCh:
				stopPinging()
				fmt.Println("stopItem")
			case <-quitItem.ClickedCh:
				fmt.Println("quitItem")
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	// Cleanup here
	stopPinging()
}

func startPinging() {
	startItem.Disable()
	stopItem.Enable()
	notifyUser("Start checking.")
	systray.SetIcon(greenIcon)
	stopChan = make(chan bool)
	pingTicker = time.NewTicker(1 * time.Second)

	go func() {
		for {
			select {
			case <-pingTicker.C:
				timeSec++
				fmt.Printf("Time: %d\n", timeSec)
				if !pingRouter(routerIP) {
					if errorPingCounter < 5 {
						errorPingCounter++
						fmt.Println("The router is unavailable. Pinging will be suspended after", 5 - errorPingCounter, " ping(s)\n")
					} else if errorPingCounter == 5 {
						notifyUser("Pinging was stopped.")
						fmt.Println("Pinging was stopped.")
						stopPinging()
						//fmt.Println("Router not reachable. Shutting down the system...")
					} else {
						//fmt.Println("Router not reachable. Shutting down the system...")
						//shutdownSystem()
						notifyUser("The router is not reachable. Shutting down the system...")
					}
				} else {
					fmt.Println("Router reachable")
					systray.SetIcon(yellowIcon)
				}
			case <-stopChan:
				fmt.Println("STOP!!!")
				pingTicker.Stop()
				return
			}
		}
	}()
}

func stopPinging() {
	errorPingCounter = 0
	systray.SetIcon(yellowIcon)
	startItem.Enable()
	stopItem.Disable()
	if stopChan != nil {
		stopChan <- true
	}
}

func pingRouter(ip string) bool {
	cmd := exec.Command("ping", "-c", "1", ip)
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing ping: %v\n", err)
		return false
	}
	if strings.Contains(string(output), "1 received") {
		return true
	}
	return false
}

func notifyUser(message string) {
	err := beeep.Notify("System shutdown", message, "")
	if err != nil {
		fmt.Printf("Error sending notification: %v\n", err)
	}
}

func shutdownSystem() {
	cmd := exec.Command("shutdown", "-h", "now")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error shutting down the system: %v\n", err)
	}
}

// getIcon reads an icon file from the given path.
func getIcon(filePath string) []byte {
	icon, err := os.ReadFile(filePath)
	if err != nil {
		//log.Fatalf("Error during downloading icon: %v", err)
		fmt.Printf("Error during downloading icon: %v\n", err)
	}
	return icon
}
