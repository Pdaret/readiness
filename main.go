package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

func main() {
	// Step 2: Extract ens5 inet IPv4 address
	var ipv4Addr string
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error getting interfaces:", err)
		return
	}

	for _, iface := range interfaces {
		// Skip down or loopback interfaces
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Printf("Error getting addresses for interface %s: %v\n", iface.Name, err)
			continue
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
			}
		}
	}
	if ipv4Addr == "" {
		fmt.Println("ens5 inet IPv4 address not found")
		return
	}

	fmt.Println("Extracted IPv4 Address:", ipv4Addr)

	time.Sleep(5 * time.Minute)

	endpoint := fmt.Sprintf("https://phoenixstatus.com/api/v1/ovh/readiness/%s", ipv4Addr)
	// Step 3: Create HTTP request and send data to server
	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	req = req.WithContext(ctx)

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
