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
	fmt.Println("Ver. 0.1.x Waiting for 3 minutes before starting the algorithm...")
	time.Sleep(3 * time.Minute)

	for {
		if !pingRouter(routerIP) {
			fmt.Printf("Router %s not reachable. Shutting down the system...\n", routerIP)
			shutdownSystem()
		} else {
			fmt.Printf("Router %s reachable\n", routerIP)
		}
		time.Sleep(1 * time.Minute)
	}
}

func pingRouter(ip string) bool {
	cmd := exec.Command("ping", "-c", "1", ip)
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing ping: %v (IP: %s)\n", err, ip)
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
