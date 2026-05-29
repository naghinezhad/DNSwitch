package network

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// get windows network
func getWindowsNetwork() ([]Network, error) {
	cmd := exec.Command("powershell", "-Command",
		"Get-NetAdapter | Select-Object Name,InterfaceDescription,Status,MediaType | ConvertTo-Json")

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var adapters []map[string]any
	if err := json.Unmarshal(output, &adapters); err != nil {
		var adapter map[string]any
		if err := json.Unmarshal(output, &adapter); err != nil {
			return nil, err
		}
		adapters = []map[string]any{adapter}
	}

	var interfaces []Network
	for _, adapter := range adapters {
		name := getString(adapter["Name"])
		desc := getString(adapter["InterfaceDescription"])
		status := getString(adapter["Status"])
		mediaType := getString(adapter["MediaType"])

		if name == "" {
			continue
		}

		netType := determineWindowsNetworkType(desc, mediaType)
		isActive := strings.ToLower(status) == "up"

		interfaces = append(interfaces, Network{
			Name:        name,
			DisplayName: desc,
			IsActive:    isActive,
			Type:        netType,
		})
	}

	return interfaces, nil
}

// get string
func getString(value any) string {
	if value == nil {
		return ""
	}
	if str, ok := value.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", value)
}

// determine windows network type
func determineWindowsNetworkType(description, mediaType string) string {
	desc := strings.ToLower(description)
	media := strings.ToLower(mediaType)

	if strings.Contains(desc, "wireless") || strings.Contains(desc, "wi-fi") ||
		strings.Contains(desc, "wifi") {
		return "Wi-Fi"
	}
	if strings.Contains(desc, "ethernet") || strings.Contains(media, "802.3") {
		return "Ethernet"
	}
	if strings.Contains(desc, "bluetooth") {
		return "Bluetooth"
	}
	if strings.Contains(desc, "virtual") || strings.Contains(desc, "loopback") {
		return "Virtual"
	}
	return "Unknown"
}

// set network dns windows
func setNetworkDNSWindows(interfaceName string, ips ...string) error {
	dnsServers := strings.Join(ips, ",")
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf("Set-DnsClientServerAddress -InterfaceAlias '%s' -ServerAddresses %s", interfaceName, dnsServers))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(output))
	}
	return nil
}
