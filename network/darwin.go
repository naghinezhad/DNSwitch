package network

import (
	"fmt"
	"os/exec"
	"strings"
)

// get darwin network
func getDarwinNetwork() ([]Network, error) {
	cmd := exec.Command("networksetup", "-listallnetworkservices")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	services := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(services) > 1 {
		services = services[1:]
	}

	var interfaces []Network
	for _, service := range services {
		service = strings.TrimSpace(service)
		if service == "" || strings.HasPrefix(service, "*") {
			continue
		}

		isActive := checkDarwinInterfaceActive(service)
		netType := determineDarwinNetworkType(service)

		interfaces = append(interfaces, Network{
			Name:        service,
			DisplayName: service,
			IsActive:    isActive,
			Type:        netType,
		})
	}

	return interfaces, nil
}

// check darwin interface active
func checkDarwinInterfaceActive(serviceName string) bool {
	cmd := exec.Command("networksetup", "-getinfo", serviceName)
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	result := string(output)
	return strings.Contains(result, "IP address:") &&
		!strings.Contains(result, "IP address: (null)") &&
		!strings.Contains(result, "IP address: none")
}

// determine darwin network type
func determineDarwinNetworkType(serviceName string) string {
	name := strings.ToLower(serviceName)
	if strings.Contains(name, "wi-fi") || strings.Contains(name, "wifi") {
		return "Wi-Fi"
	}
	if strings.Contains(name, "ethernet") {
		return "Ethernet"
	}
	if strings.Contains(name, "bluetooth") {
		return "Bluetooth"
	}
	return "Unknown"
}

// set network dns darwin
func setNetworkDNSDarwin(interfaceName string, ips ...string) error {
	args := append([]string{"-setpkg.DnsServers", interfaceName}, ips...)
	cmd := exec.Command("networksetup", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(output))
	}
	return nil
}
