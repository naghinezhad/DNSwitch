package main

import (
	"fmt"
	"strings"

	"github.com/naghinezhad/DNSwitch/dns"
	"github.com/naghinezhad/DNSwitch/network"
)

func main() {
	// msg
	fmt.Println("Welcome to DNSwitch!")
	fmt.Println("Warning: This program requires administrator privileges to change DNS settings.")
	fmt.Println("Please make sure you're running this program as an administrator.")
	fmt.Println()

	// load custom dns
	dns.LoadCustomDNS()

	for {
		// select network
		selectedNetwork, err := network.SelectNetwork()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if selectedNetwork == "" {
			fmt.Println("Thank you for using DNSwitch!")
			return
		}

		fmt.Printf("Selected network: %s\n\n", selectedNetwork)

		// manage dns for network
		dns.ManageDNSForNetwork(selectedNetwork)

		fmt.Println("\nReturning to network selection...")
		fmt.Println(strings.Repeat("-", 50))
	}
}
