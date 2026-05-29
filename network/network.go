package network

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/naghinezhad/DNSwitch/pkg"
)

type Network struct {
	Name        string
	DisplayName string
	IsActive    bool
	Type        string
}

// select network
func SelectNetwork() (string, error) {
	// get networks
	networks, err := getNetworks()
	if err != nil {
		return "", fmt.Errorf("failed to get network interfaces: %v", err)
	}

	if len(networks) == 0 {
		return "", fmt.Errorf("no network interfaces found")
	}

	fmt.Println("Available Network Interfaces:")
	fmt.Println(strings.Repeat("-", 40))

	for i, iface := range networks {
		status := "Inactive"
		if iface.IsActive {
			status = "Active"
		}
		fmt.Printf("%d. %s (%s) - %s [%s]\n",
			i+1, iface.DisplayName, iface.Name, iface.Type, status)
	}

	fmt.Printf("%d. Exit\n", len(networks)+1)
	fmt.Println()

	// get user choice
	choice := pkg.GetUserChoice(len(networks) + 1)

	if choice == len(networks)+1 {
		return "", nil
	}

	return networks[choice-1].Name, nil
}

// get networks
func getNetworks() ([]Network, error) {
	switch runtime.GOOS {
	case "windows":
		return getWindowsNetwork()
	case "darwin":
		return getDarwinNetwork()
	case "linux":
		return getLinuxNetwork()
	default:
		return nil, fmt.Errorf("OS %s is not supported", runtime.GOOS)
	}
}

// set network dns
func SetNetworkDNS(networkName string, ips ...string) error {
	switch runtime.GOOS {
	case "windows":
		return setNetworkDNSWindows(networkName, ips...)
	case "darwin":
		return setNetworkDNSDarwin(networkName, ips...)
	case "linux":
		return setNetworkDNSLinux(ips...)
	default:
		return fmt.Errorf("OS %s is not supported", runtime.GOOS)
	}
}
