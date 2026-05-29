package dns

import (
	"fmt"
	"net"
	"sort"
	"strings"
	"time"

	"github.com/naghinezhad/DNSwitch/network"
	"github.com/naghinezhad/DNSwitch/pkg"
)

// manage dns for network
func ManageDNSForNetwork(networkName string) {
	for {
		displayCurrentDNS(networkName)
		displayDNSMenu()

		totalOptions := len(dnsServers) + 4
		choice := pkg.GetUserChoice(totalOptions)

		if handleDNSChoice(choice, networkName) {
			break
		}

		fmt.Println()
		time.Sleep(2 * time.Second)
	}
}

// handle dns Choice
func handleDNSChoice(choice int, interfaceName string) bool {
	options := getOrderedDNSOptions()
	baseIndex := len(options)

	switch {
	case choice <= len(options) && choice >= 1:
		selectedDNS := options[choice-1]
		ips := dnsServers[selectedDNS]
		err := network.SetNetworkDNS(interfaceName, ips...)
		if err != nil {
			fmt.Printf("Error setting DNS: %v\n", err)
		} else {
			displayName := formatDNSName(selectedDNS)
			fmt.Printf("DNS changed to %s (%s).\n", displayName, strings.Join(ips, ", "))
		}
	case choice == baseIndex+1:
		addCustomDNS()
	case choice == baseIndex+2:
		removeCustomDNS()
	case choice == baseIndex+3:
		err := clearAllDNS(interfaceName)
		if err != nil {
			fmt.Printf("Error clearing DNS settings: %v\n", err)
		} else {
			fmt.Println("All DNS settings have been cleared.")
		}
	case choice == baseIndex+4:
		return true
	default:
		fmt.Println("Invalid choice. Please try again.")
	}
	return false
}

// get ordered dns options
func getOrderedDNSOptions() []string {
	var options []string

	for name := range dnsServers {
		options = append(options, name)
	}

	sort.Strings(options)

	return options
}

// display current dns
func displayCurrentDNS(networkName string) {
	currentDNSName, currentDNSAddresses := getCurrentDNS(networkName)
	fmt.Printf("Current DNS for %s: %s\n", networkName, currentDNSName)
	if len(currentDNSAddresses) > 0 {
		fmt.Printf("Addresses: %s\n", strings.Join(currentDNSAddresses, ", "))
	}
	fmt.Println()
}

// display dns menu
func displayDNSMenu() {
	fmt.Println("DNS Options:")
	options := getOrderedDNSOptions()

	for i, name := range options {
		ips := dnsServers[name]
		displayName := formatDNSName(name)
		fmt.Printf("%d. %s: %s\n", i+1, displayName, strings.Join(ips, ", "))
	}

	baseIndex := len(options)
	fmt.Printf("%d. Add custom DNS\n", baseIndex+1)
	fmt.Printf("%d. Remove custom DNS\n", baseIndex+2)
	fmt.Printf("%d. Clear all DNS settings\n", baseIndex+3)
	fmt.Printf("%d. Back to network selection\n", baseIndex+4)
	fmt.Println()
}

// add custom dns
func addCustomDNS() {
	var name, ip1, ip2 string

	fmt.Print("Enter DNS name: ")
	if _, err := fmt.Scanln(&name); err != nil {
		fmt.Printf("Error reading DNS name: %v\n", err)
		return
	}

	fmt.Print("Enter first IP address: ")
	if _, err := fmt.Scanln(&ip1); err != nil {
		fmt.Printf("Error reading first IP address: %v\n", err)
		return
	}

	fmt.Print("Enter second IP address: ")
	if _, err := fmt.Scanln(&ip2); err != nil {
		fmt.Printf("Error reading second IP address: %v\n", err)
		return
	}

	if strings.TrimSpace(name) == "" {
		fmt.Println("DNS name cannot be empty.")
		return
	}

	if strings.TrimSpace(ip1) == "" || strings.TrimSpace(ip2) == "" {
		fmt.Println("IP addresses cannot be empty.")
		return
	}

	if net.ParseIP(ip1) == nil {
		fmt.Printf("Invalid first IP address: %s\n", ip1)
		return
	}

	if net.ParseIP(ip2) == nil {
		fmt.Printf("Invalid second IP address: %s\n", ip2)
		return
	}

	name = strings.ToLower(strings.TrimSpace(name))
	dnsServers[name] = []string{ip1, ip2}
	saveCustomDNS()
	fmt.Printf("Custom DNS '%s' added successfully.\n", formatDNSName(name))
}

// remove custom dns
func removeCustomDNS() {
	customDNS := getCustomDNSList()

	if len(customDNS) == 0 {
		fmt.Println("No custom DNS servers found.")
		return
	}

	fmt.Println("Custom DNS servers:")
	for i, name := range customDNS {
		ips := dnsServers[name]
		displayName := formatDNSName(name)
		fmt.Printf("%d. %s: %s\n", i+1, displayName, strings.Join(ips, ", "))
	}
	fmt.Printf("%d. Exit\n", len(customDNS)+1)

	choice := pkg.GetUserChoice(len(customDNS) + 1)

	if choice == len(customDNS)+1 {
		fmt.Println("Exiting remove custom DNS menu.")
		return
	}

	name := customDNS[choice-1]
	delete(dnsServers, name)
	saveCustomDNS()
	fmt.Printf("Custom DNS '%s' removed successfully.\n", formatDNSName(name))
}
