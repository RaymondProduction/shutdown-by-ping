package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

const (
	routerIP = "192.168.0.1"
)

func main() {
	for {
		if !pingRouter(routerIP) {
			fmt.Println("Router not reachable. Shutting down the system...")
			shutdownSystem()
		} else {
                       fmt.Println("Router reachable")
                }
		time.Sleep(1 * time.Minute)
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
