package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/getlantern/systray"
	"github.com/gen2brain/beeep"
)

const (
	routerIP = "192.168.1.4" //"127.0.0.1"
)

var (
	greenIcon  []byte
	redIcon []byte
	yellowIcon []byte
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
	stopItem := systray.AddMenuItem("Stop Ping", "Stop pinging the router")
	quitItem := systray.AddMenuItem("Quit", "Quit the application")

	stopItem.Disable()

	go func() {
		for {
			select {
			case <-startItem.ClickedCh:
				startPinging()
				startItem.Disable()
				stopItem.Enable()
				fmt.Println("startItem")
			case <-stopItem.ClickedCh:
				stopPinging()
				startItem.Enable()
				stopItem.Disable()
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
	notifyUser("Start checking.")
	systray.SetIcon(greenIcon)
	stopChan = make(chan bool)
	pingTicker = time.NewTicker(1 * time.Second)

	go func() {
		firstMinute := true
		for {
			select {
			case <-pingTicker.C:
				if !pingRouter(routerIP) {
					if firstMinute {
						systray.SetIcon(redIcon)
						notifyUser("The router is not reachable. Please check your network connection.")
						firstMinute = false
					} else {
						fmt.Println("Router not reachable. Shutting down the system...")
						//shutdownSystem()
						notifyUser("The router is not reachable. Please check your network connection.")
					}
				} else {
					fmt.Println("Router reachable")
					systray.SetIcon(yellowIcon)
				}
			case <-stopChan:
				pingTicker.Stop()
				return
			}
		}
	}()
}

func stopPinging() {
	systray.SetIcon(yellowIcon)
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
