package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/getlantern/systray"
)

const (
	routerIP = "127.0.0.1"
)

var (
	pingTicker *time.Ticker
	stopChan   chan bool
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(getIcon("icon.png")) // replace with your icon file
	systray.SetTitle("Ping Router")
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
			case <-stopItem.ClickedCh:
				stopPinging()
				startItem.Enable()
				stopItem.Disable()
			case <-quitItem.ClickedCh:
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
	stopChan = make(chan bool)
	pingTicker = time.NewTicker(1 * time.Minute)

	go func() {
		for {
			select {
			case <-pingTicker.C:
				if !pingRouter(routerIP) {
					fmt.Println("Router not reachable. Shutting down the system...")
					shutdownSystem()
				} else {
					fmt.Println("Router reachable")
				}
			case <-stopChan:
				pingTicker.Stop()
				return
			}
		}
	}()
}

func stopPinging() {
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
