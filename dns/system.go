package dns

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// get dns command
func getDNSCommand(interfaceName string) *exec.Cmd {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("powershell", "-Command", fmt.Sprintf("(Get-DnsClientServerAddress -InterfaceAlias '%s' -AddressFamily IPv4).ServerAddresses", interfaceName))
	case "darwin":
		return exec.Command("networksetup", "-getpkg.DnsServers", interfaceName)
	case "linux":
		return exec.Command("sh", "-c", "grep nameserver /etc/resolv.conf | awk '{print $2}'")
	default:
		return nil
	}
}

// get clear dns command
func getClearDNSCommand(interfaceName string) *exec.Cmd {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("powershell", "-Command", fmt.Sprintf("Set-DnsClientServerAddress -InterfaceAlias '%s' -ResetServerAddresses", interfaceName))
	case "darwin":
		return exec.Command("networksetup", "-setpkg.DnsServers", interfaceName, "Empty")
	case "linux":
		return exec.Command("sudo", "sh", "-c", "echo '' > /etc/resolv.conf")
	default:
		return nil
	}
}

// clear all dns
func clearAllDNS(interfaceName string) error {
	cmd := getClearDNSCommand(interfaceName)
	if cmd == nil {
		return fmt.Errorf("OS %s is not supported", runtime.GOOS)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(output))
	}
	return nil
}

// get current dns
func getCurrentDNS(interfaceName string) (string, []string) {
	cmd := getDNSCommand(interfaceName)
	if cmd == nil {
		return "Unknown", nil
	}

	output, err := cmd.Output()
	if err != nil {
		return "Unknown", nil
	}

	dnsAddresses := removeDuplicates(strings.Fields(string(output)))

	dnsName := "Unknown"
	for name, ips := range dnsServers {
		if containsAny(dnsAddresses, ips) {
			dnsName = name
			break
		}
	}

	return dnsName, dnsAddresses
}
