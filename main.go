package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Constants and global variables
const (
	customDNSFile = "custom_dns.json"
)

var dnsServers = map[string][]string{
	"403":       {"10.202.10.202", "10.202.10.102"},
	"Shecan":    {"178.22.122.100", "185.51.200.2"},
	"Begzar":    {"185.55.226.26", "185.55.225.25"},
	"electrotm": {"78.157.42.101", "78.157.42.100"},
}

var dnsOrder = []string{"403", "Shecan", "Begzar", "electrotm"}

// Main function and core logic
func main() {
	printWelcomeMessage()

	interfaceName, err := getActiveInterfaceName()
	if err != nil {
		fmt.Printf("Error detecting active interface: %v\n", err)
		return
	}
	fmt.Printf("Detected active interface: %s\n", interfaceName)

	loadCustomDNS()

	for {
		displayCurrentDNS(interfaceName)
		displayMenu()

		choice := getUserChoice(len(dnsServers) + 4)

		if handleUserChoice(choice, interfaceName) {
			break
		}

		fmt.Println()
		time.Sleep(2 * time.Second)
	}

	fmt.Println("Thank you for using DNSwitch!")
}

func printWelcomeMessage() {
	fmt.Println("Welcome to DNSwitch!")
	fmt.Println("Warning: This program requires administrator privileges to change DNS settings.")
	fmt.Println("Please make sure you're running this program as an administrator.")
	fmt.Println()
}

func displayCurrentDNS(interfaceName string) {
	currentDNSName, currentDNSAddresses := getCurrentDNS(interfaceName)
	fmt.Printf("Current DNS: %s\n", currentDNSName)
	if len(currentDNSAddresses) > 0 {
		fmt.Printf("Addresses: %s\n", strings.Join(currentDNSAddresses, ", "))
	}
	fmt.Println()
}

func displayMenu() {
	fmt.Println("Available options:")
	options := getOrderedDNSOptions()
	for i, name := range options {
		ips := dnsServers[name]
		fmt.Printf("%d. %s: %s\n", i+1, name, strings.Join(ips, ", "))
	}
	fmt.Printf("%d. Add custom DNS\n", len(options)+1)
	fmt.Printf("%d. Remove custom DNS\n", len(options)+2)
	fmt.Printf("%d. Clear all DNS settings\n", len(options)+3)
	fmt.Printf("%d. Exit\n", len(options)+4)
	fmt.Println()
}

func handleUserChoice(choice int, interfaceName string) bool {
	options := getOrderedDNSOptions()

	switch {
	case choice <= len(options):
		selectedDNS := options[choice-1]
		ips := dnsServers[selectedDNS]
		err := setDNS(interfaceName, ips...)
		if err != nil {
			fmt.Printf("Error setting DNS: %v\n", err)
		} else {
			fmt.Printf("DNS changed to %s (%s).\n", selectedDNS, strings.Join(ips, ", "))
		}
	case choice == len(options)+1:
		addCustomDNS()
	case choice == len(options)+2:
		removeCustomDNS()
	case choice == len(options)+3:
		err := clearAllDNS(interfaceName)
		if err != nil {
			fmt.Printf("Error clearing DNS settings: %v\n", err)
		} else {
			fmt.Println("All DNS settings have been cleared.")
		}
	case choice == len(options)+4:
		return true
	}
	return false
}

// DNS-related functions
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

func setDNS(interfaceName string, ips ...string) error {
	switch runtime.GOOS {
	case "windows":
		return setDNSWindows(interfaceName, ips...)
	case "darwin":
		return setDNSMacOS(interfaceName, ips...)
	case "linux":
		return setDNSLinux(ips...)
	default:
		return fmt.Errorf("OS %s is not supported", runtime.GOOS)
	}
}

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

func addCustomDNS() {
	var name, ip1, ip2 string

	fmt.Print("Enter DNS name: ")
	fmt.Scanln(&name)

	fmt.Print("Enter first IP address: ")
	fmt.Scanln(&ip1)

	fmt.Print("Enter second IP address: ")
	fmt.Scanln(&ip2)

	dnsServers[name] = []string{ip1, ip2}
	saveCustomDNS()
	fmt.Printf("Custom DNS '%s' added successfully.\n", name)
}

func removeCustomDNS() {
	customDNS := getCustomDNSList()

	if len(customDNS) == 0 {
		fmt.Println("No custom DNS servers found.")
		return
	}

	fmt.Println("Custom DNS servers:")
	for i, name := range customDNS {
		ips := dnsServers[name]
		fmt.Printf("%d. %s: %s\n", i+1, name, strings.Join(ips, ", "))
	}
	fmt.Printf("%d. Exit\n", len(customDNS)+1)

	choice := getUserChoice(len(customDNS) + 1)

	if choice == len(customDNS)+1 {
		fmt.Println("Exiting remove custom DNS menu.")
		return
	}

	name := customDNS[choice-1]
	delete(dnsServers, name)
	saveCustomDNS()
	fmt.Printf("Custom DNS '%s' removed successfully.\n", name)
}

// Helper functions
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

func getMacOSNetworkServiceName() (string, error) {
	cmd := exec.Command("networksetup", "-listallnetworkservices")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to list network services: %v", err)
	}

	services := strings.Split(strings.TrimSpace(string(output)), "\n")

	if len(services) > 1 {
		services = services[1:]
	}

	for _, service := range services {
		service = strings.TrimSpace(service)
		if service == "" || strings.HasPrefix(service, "*") {
			continue
		}

		dnsCmd := exec.Command("networksetup", "-getdnsservers", service)
		dnsOutput, err := dnsCmd.Output()
		if err == nil {
			dnsResult := strings.TrimSpace(string(dnsOutput))
			if !strings.Contains(dnsResult, "There aren't any DNS Servers set") {
				return service, nil
			}
		}

		ipCmd := exec.Command("networksetup", "-getinfo", service)
		ipOutput, err := ipCmd.Output()
		if err == nil {
			ipResult := string(ipOutput)
			if strings.Contains(ipResult, "IP address:") &&
				!strings.Contains(ipResult, "IP address: (null)") &&
				!strings.Contains(ipResult, "IP address: none") {
				return service, nil
			}
		}
	}

	for _, service := range services {
		service = strings.TrimSpace(service)
		if service != "" && !strings.HasPrefix(service, "*") {
			return service, nil
		}
	}

	return "", fmt.Errorf("no active network service found")
}

func getActiveInterfaceName() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		return getMacOSNetworkServiceName()
	default:
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
}

func getOrderedDNSOptions() []string {
	options := make([]string, 0, len(dnsServers))

	// First, add the default DNS options in the specified order
	for _, name := range dnsOrder {
		if _, exists := dnsServers[name]; exists {
			options = append(options, name)
		}
	}

	// Then, add any custom DNS options
	for name := range dnsServers {
		if !isDefaultDNS(name) {
			options = append(options, name)
		}
	}

	return options
}

func getCustomDNSList() []string {
	customDNS := make([]string, 0)
	for name := range dnsServers {
		if !isDefaultDNS(name) {
			customDNS = append(customDNS, name)
		}
	}
	return customDNS
}

func isDefaultDNS(name string) bool {
	return containsString(dnsOrder, name)
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

func containsString(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func containsAny(slice []string, items []string) bool {
	for _, item := range items {
		if containsString(slice, item) {
			return true
		}
	}
	return false
}

// File operations
func loadCustomDNS() {
	data, err := os.ReadFile(customDNSFile)
	if err != nil {
		return
	}

	var customDNS map[string][]string
	err = json.Unmarshal(data, &customDNS)
	if err != nil {
		fmt.Printf("Error loading custom DNS: %v\n", err)
		return
	}

	for name, ips := range customDNS {
		dnsServers[name] = ips
	}
}

func saveCustomDNS() {
	customDNS := make(map[string][]string)
	for name, ips := range dnsServers {
		if !isDefaultDNS(name) {
			customDNS[name] = ips
		}
	}

	data, err := json.Marshal(customDNS)
	if err != nil {
		fmt.Printf("Error saving custom DNS: %v\n", err)
		return
	}

	err = os.WriteFile(customDNSFile, data, 0644)
	if err != nil {
		fmt.Printf("Error writing custom DNS file: %v\n", err)
	}
}

// OS-specific functions
func getDNSCommand(interfaceName string) *exec.Cmd {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("powershell", "-Command", fmt.Sprintf("(Get-DnsClientServerAddress -InterfaceAlias '%s' -AddressFamily IPv4).ServerAddresses", interfaceName))
	case "darwin":
		return exec.Command("networksetup", "-getdnsservers", interfaceName)
	case "linux":
		return exec.Command("sh", "-c", "grep nameserver /etc/resolv.conf | awk '{print $2}'")
	default:
		return nil
	}
}

func getClearDNSCommand(interfaceName string) *exec.Cmd {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("powershell", "-Command", fmt.Sprintf("Set-DnsClientServerAddress -InterfaceAlias '%s' -ResetServerAddresses", interfaceName))
	case "darwin":
		return exec.Command("networksetup", "-setdnsservers", interfaceName, "Empty")
	case "linux":
		return exec.Command("sudo", "sh", "-c", "echo '' > /etc/resolv.conf")
	default:
		return nil
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

func setDNSMacOS(interfaceName string, ips ...string) error {
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
