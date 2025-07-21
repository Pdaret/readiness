package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

func main() {
	// Step 1: Check marzban status
	for {
		cmd := exec.Command("sudo", "marzban", "status")
		output, err := cmd.Output()
		if err != nil {
			fmt.Println("Error executing marzban command:", err)
			time.Sleep(5 * time.Second) // Wait before retrying
			continue
		}

		if strings.Contains(string(output), "running") {
			break
		}

		fmt.Println("Marzban status not up, retrying...")
		time.Sleep(5 * time.Second) // Wait before retrying
	}

	// Step 2: Extract ens5 inet IPv4 address
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error fetching network interfaces:", err)
		return
	}

	var ipv4Addr string
	for _, iface := range interfaces {
		if iface.Name == "ens5" {
			addrs, err := iface.Addrs()
			if err != nil {
				fmt.Println("Error getting addresses for ens5:", err)
				return
			}

			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}

				if ip != nil && ip.To4() != nil {
					ipv4Addr = ip.String()
					break
				}
			}
			break
		}
	}

	if ipv4Addr == "" {
		fmt.Println("ens5 inet IPv4 address not found")
		return
	}

	fmt.Println("Extracted IPv4 Address:", ipv4Addr)

	endpoint := fmt.Sprintf("https://phoenixstatus.com/api/v1/readiness/%s", ipv4Addr)
	// Step 3: Create HTTP request and send data to server
	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)
	body, _ := bufio.NewReader(resp.Body).ReadString('\n')
	fmt.Println("Response Body:", body)
}
