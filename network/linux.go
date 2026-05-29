package network

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
)

// get linux network
func getLinuxNetwork() ([]Network, error) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var interfaces []Network
	for _, iface := range netInterfaces {
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		isActive := iface.Flags&net.FlagUp != 0
		netType := determineLinuxNetworkType(iface.Name)

		interfaces = append(interfaces, Network{
			Name:        iface.Name,
			DisplayName: iface.Name,
			IsActive:    isActive,
			Type:        netType,
		})
	}

	return interfaces, nil
}

// determine linux network type
func determineLinuxNetworkType(interfaceName string) string {
	name := strings.ToLower(interfaceName)
	if strings.HasPrefix(name, "wl") || strings.HasPrefix(name, "wlan") {
		return "Wi-Fi"
	}
	if strings.HasPrefix(name, "eth") || strings.HasPrefix(name, "en") {
		return "Ethernet"
	}
	if strings.HasPrefix(name, "docker") || strings.HasPrefix(name, "br-") ||
		strings.HasPrefix(name, "veth") {
		return "Virtual"
	}
	return "Unknown"
}

// set network dns linux
func setNetworkDNSLinux(ips ...string) error {
	var content strings.Builder
	for _, ip := range ips {
		fmt.Fprintf(&content, "nameserver %s\n", ip)
	}
	cmd := exec.Command("sudo", "sh", "-c", fmt.Sprintf("echo '%s' > /etc/resolv.conf", content.String()))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(output))
	}
	return nil
}
