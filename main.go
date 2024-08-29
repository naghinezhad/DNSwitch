package main

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var dnsServers = map[string][]string{
	"403": {
		"10.202.10.202",
		"10.202.10.102",
	},
	"Shecan": {
		"178.22.122.100",
		"185.51.200.2",
	},
	"Begzar": {
		"185.55.226.26",
		"185.55.225.25",
	},
}

func main() {
	fmt.Println("Welcome to DNSwitch!")
	fmt.Println("Warning: This program requires administrator privileges to change DNS settings.")
	fmt.Println("Please make sure you're running this program as an administrator.")
	fmt.Println()

	interfaceName, err := getActiveInterfaceName()
	if err != nil {
		fmt.Printf("Error detecting active interface: %v\n", err)
		return
	}
	fmt.Printf("Detected active interface: %s\n", interfaceName)

	for {
		currentDNSName, currentDNSAddresses := getCurrentDNS(interfaceName)
		fmt.Printf("Current DNS: %s\n", currentDNSName)
		if len(currentDNSAddresses) > 0 {
			fmt.Printf("Addresses: %s\n", strings.Join(currentDNSAddresses, ", "))
		}
		fmt.Println()

		fmt.Println("Available options:")
		options := make([]string, 0, len(dnsServers))
		for name := range dnsServers {
			options = append(options, name)
		}
		for i, name := range options {
			ips := dnsServers[name]
			fmt.Printf("%d. %s: %s\n", i+1, name, strings.Join(ips, ", "))
		}
		fmt.Printf("%d. Clear all DNS settings\n", len(options)+1)
		fmt.Printf("%d. Exit\n", len(options)+2)
		fmt.Println()

		choice := getUserChoice(len(options) + 2)

		if choice == len(options)+2 {
			break
		}

		if choice == len(options)+1 {
			err := clearAllDNS(interfaceName)
			if err != nil {
				fmt.Printf("Error clearing DNS settings: %v\n", err)
			} else {
				fmt.Println("All DNS settings have been cleared.")
			}
		} else {
			selectedDNS := options[choice-1]
			ips := dnsServers[selectedDNS]
			err := setDNS(interfaceName, ips...)
			if err != nil {
				fmt.Printf("Error setting DNS: %v\n", err)
			} else {
				fmt.Printf("DNS changed to %s (%s).\n", selectedDNS, strings.Join(ips, ", "))
			}
		}
		fmt.Println()
		time.Sleep(2 * time.Second)
	}

	fmt.Println("Thank you for using DNSwitch!")
}

func getUserChoice(maxChoice int) int {
	for {
		var choice string
		fmt.Print("Please enter the number of your choice: ")
		fmt.Scanln(&choice)

		index, err := strconv.Atoi(choice)
		if err == nil && index >= 1 && index <= maxChoice {
			return index
		}
		fmt.Println("Invalid input. Please enter a valid number.")
	}
}

func getActiveInterfaceName() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			addrs, err := iface.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return iface.Name, nil
					}
				}
			}
		}
	}

	return "", fmt.Errorf("no active network interface found")
}

func getCurrentDNS(interfaceName string) (string, []string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("powershell", "-Command", fmt.Sprintf("(Get-DnsClientServerAddress -InterfaceAlias '%s' -AddressFamily IPv4).ServerAddresses", interfaceName))
	case "darwin":
		cmd = exec.Command("networksetup", "-getdnsservers", interfaceName)
	case "linux":
		cmd = exec.Command("cat", "/etc/resolv.conf")
	default:
		return "Unknown", nil
	}

	output, err := cmd.Output()
	if err != nil {
		return "Unknown", nil
	}

	dnsAddresses := strings.Fields(string(output))
	uniqueDNSAddresses := removeDuplicates(dnsAddresses)

	dnsName := "Unknown"
	for name, ips := range dnsServers {
		for _, ip := range ips {
			if contains(uniqueDNSAddresses, ip) {
				dnsName = name
				break
			}
		}
		if dnsName != "Unknown" {
			break
		}
	}

	return dnsName, uniqueDNSAddresses
}

func setDNS(interfaceName string, ips ...string) error {
	switch runtime.GOOS {
	case "windows":
		return setDNSWindows(interfaceName, ips...)
	case "darwin":
		return setDNSDarwin(interfaceName, ips...)
	case "linux":
		return setDNSLinux(ips...)
	default:
		return fmt.Errorf("OS %s is not supported", runtime.GOOS)
	}
}

func setDNSWindows(interfaceName string, ips ...string) error {
	dnsServers := strings.Join(ips, ",")
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf("Set-DnsClientServerAddress -InterfaceAlias '%s' -ServerAddresses %s", interfaceName, dnsServers))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(output))
	}
	return nil
}

func setDNSDarwin(interfaceName string, ips ...string) error {
	args := append([]string{"-setdnsservers", interfaceName}, ips...)
	cmd := exec.Command("networksetup", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(output))
	}
	return nil
}

func setDNSLinux(ips ...string) error {
	content := ""
	for _, ip := range ips {
		content += fmt.Sprintf("nameserver %s\n", ip)
	}
	cmd := exec.Command("sudo", "sh", "-c", fmt.Sprintf("echo '%s' > /etc/resolv.conf", content))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(output))
	}
	return nil
}

func clearAllDNS(interfaceName string) error {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("powershell", "-Command", fmt.Sprintf("Set-DnsClientServerAddress -InterfaceAlias '%s' -ResetServerAddresses", interfaceName))
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%v: %s", err, string(output))
		}
	case "darwin":
		cmd := exec.Command("networksetup", "-setdnsservers", interfaceName, "Empty")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%v: %s", err, string(output))
		}
	case "linux":
		cmd := exec.Command("sudo", "sh", "-c", "echo '' > /etc/resolv.conf")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%v: %s", err, string(output))
		}
	default:
		return fmt.Errorf("OS %s is not supported", runtime.GOOS)
	}
	return nil
}

func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
