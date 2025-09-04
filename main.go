package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	customDNSFile = "custom_dns.json"
)

var dnsServers = map[string][]string{
	"403":       {"10.202.10.202", "10.202.10.102"},
	"shecan":    {"178.22.122.100", "185.51.200.2"},
	"begzar":    {"185.55.226.26", "185.55.225.25"},
	"electrotm": {"78.157.42.101", "78.157.42.100"},
	"dynx":      {"10.70.95.150", "10.70.95.162"},
	"radar":     {"10.202.10.10", "10.202.10.11"},
	"shatel":    {"85.15.1.14", "85.15.1.15"},
	"level3":    {"209.244.0.3", "209.244.0.4"},
	"shelter":   {"94.103.125.157", "94.103.125.158"},
	"beshkan":   {"181.41.194.177", "181.41.194.186"},
}

type NetworkInterface struct {
	Name        string
	DisplayName string
	IsActive    bool
	Type        string
}

func main() {
	printWelcomeMessage()

	loadCustomDNS()

	for {
		selectedInterface, err := selectNetworkInterface()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if selectedInterface == "" {
			fmt.Println("Thank you for using DNSwitch!")
			return
		}

		fmt.Printf("Selected interface: %s\n\n", selectedInterface)

		manageDNSForInterface(selectedInterface)

		fmt.Println("\nReturning to network selection...")
		fmt.Println(strings.Repeat("-", 50))
	}
}

func printWelcomeMessage() {
	fmt.Println("Welcome to DNSwitch!")
	fmt.Println("Warning: This program requires administrator privileges to change DNS settings.")
	fmt.Println("Please make sure you're running this program as an administrator.")
	fmt.Println()
}

func selectNetworkInterface() (string, error) {
	interfaces, err := getAllNetworkInterfaces()
	if err != nil {
		return "", fmt.Errorf("failed to get network interfaces: %v", err)
	}

	if len(interfaces) == 0 {
		return "", fmt.Errorf("no network interfaces found")
	}

	fmt.Println("Available Network Interfaces:")
	fmt.Println(strings.Repeat("-", 40))

	for i, iface := range interfaces {
		status := "Inactive"
		if iface.IsActive {
			status = "Active"
		}
		fmt.Printf("%d. %s (%s) - %s [%s]\n",
			i+1, iface.DisplayName, iface.Name, iface.Type, status)
	}

	fmt.Printf("%d. Exit\n", len(interfaces)+1)
	fmt.Println()

	choice := getUserChoice(len(interfaces) + 1)

	if choice == len(interfaces)+1 {
		return "", nil
	}

	return interfaces[choice-1].Name, nil
}

func getAllNetworkInterfaces() ([]NetworkInterface, error) {
	switch runtime.GOOS {
	case "windows":
		return getWindowsNetworkInterfaces()
	case "darwin":
		return getMacOSNetworkInterfaces()
	case "linux":
		return getLinuxNetworkInterfaces()
	default:
		return nil, fmt.Errorf("OS %s is not supported", runtime.GOOS)
	}
}

func getWindowsNetworkInterfaces() ([]NetworkInterface, error) {
	cmd := exec.Command("powershell", "-Command",
		"Get-NetAdapter | Select-Object Name,InterfaceDescription,Status,MediaType | ConvertTo-Json")

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var adapters []map[string]interface{}
	if err := json.Unmarshal(output, &adapters); err != nil {
		var adapter map[string]interface{}
		if err := json.Unmarshal(output, &adapter); err != nil {
			return nil, err
		}
		adapters = []map[string]interface{}{adapter}
	}

	var interfaces []NetworkInterface
	for _, adapter := range adapters {
		name := getString(adapter["Name"])
		desc := getString(adapter["InterfaceDescription"])
		status := getString(adapter["Status"])
		mediaType := getString(adapter["MediaType"])

		if name == "" {
			continue
		}

		netType := determineNetworkType(desc, mediaType)
		isActive := strings.ToLower(status) == "up"

		interfaces = append(interfaces, NetworkInterface{
			Name:        name,
			DisplayName: desc,
			IsActive:    isActive,
			Type:        netType,
		})
	}

	return interfaces, nil
}

func getMacOSNetworkInterfaces() ([]NetworkInterface, error) {
	cmd := exec.Command("networksetup", "-listallnetworkservices")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	services := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(services) > 1 {
		services = services[1:]
	}

	var interfaces []NetworkInterface
	for _, service := range services {
		service = strings.TrimSpace(service)
		if service == "" || strings.HasPrefix(service, "*") {
			continue
		}

		isActive := checkMacOSInterfaceActive(service)
		netType := determineMacOSNetworkType(service)

		interfaces = append(interfaces, NetworkInterface{
			Name:        service,
			DisplayName: service,
			IsActive:    isActive,
			Type:        netType,
		})
	}

	return interfaces, nil
}

func getLinuxNetworkInterfaces() ([]NetworkInterface, error) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var interfaces []NetworkInterface
	for _, iface := range netInterfaces {
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		isActive := iface.Flags&net.FlagUp != 0
		netType := determineLinuxNetworkType(iface.Name)

		interfaces = append(interfaces, NetworkInterface{
			Name:        iface.Name,
			DisplayName: iface.Name,
			IsActive:    isActive,
			Type:        netType,
		})
	}

	return interfaces, nil
}

func manageDNSForInterface(interfaceName string) {
	for {
		displayCurrentDNS(interfaceName)
		displayDNSMenu()

		totalOptions := len(dnsServers) + 4
		choice := getUserChoice(totalOptions)

		if handleDNSChoice(choice, interfaceName) {
			break
		}

		fmt.Println()
		time.Sleep(2 * time.Second)
	}
}

func displayCurrentDNS(interfaceName string) {
	currentDNSName, currentDNSAddresses := getCurrentDNS(interfaceName)
	fmt.Printf("Current DNS for %s: %s\n", interfaceName, currentDNSName)
	if len(currentDNSAddresses) > 0 {
		fmt.Printf("Addresses: %s\n", strings.Join(currentDNSAddresses, ", "))
	}
	fmt.Println()
}

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

func handleDNSChoice(choice int, interfaceName string) bool {
	options := getOrderedDNSOptions()
	baseIndex := len(options)

	switch {
	case choice <= len(options) && choice >= 1:
		selectedDNS := options[choice-1]
		ips := dnsServers[selectedDNS]
		err := setDNS(interfaceName, ips...)
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

func formatDNSName(name string) string {
	if name == "403" {
		return "403"
	}
	return strings.Title(name)
}

func getString(value interface{}) string {
	if value == nil {
		return ""
	}
	if str, ok := value.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", value)
}

func determineNetworkType(description, mediaType string) string {
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

func determineMacOSNetworkType(serviceName string) string {
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

func checkMacOSInterfaceActive(serviceName string) bool {
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

	choice := getUserChoice(len(customDNS) + 1)

	if choice == len(customDNS)+1 {
		fmt.Println("Exiting remove custom DNS menu.")
		return
	}

	name := customDNS[choice-1]
	delete(dnsServers, name)
	saveCustomDNS()
	fmt.Printf("Custom DNS '%s' removed successfully.\n", formatDNSName(name))
}

func getUserChoice(maxChoice int) int {
	for {
		var choice string
		fmt.Printf("Please enter your choice (1-%d): ", maxChoice)
		if _, err := fmt.Scanln(&choice); err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}

		choice = strings.TrimSpace(choice)
		if choice == "" {
			fmt.Println("Input cannot be empty. Please enter a valid number.")
			continue
		}

		index, err := strconv.Atoi(choice)
		if err == nil && index >= 1 && index <= maxChoice {
			return index
		}
		fmt.Printf("Invalid input '%s'. Please enter a number between 1 and %d.\n", choice, maxChoice)
	}
}

func getOrderedDNSOptions() []string {
	var options []string

	for name := range dnsServers {
		options = append(options, name)
	}

	sort.Strings(options)

	return options
}

func getCustomDNSList() []string {
	defaultDNSNames := map[string]bool{
		"403":       true,
		"shecan":    true,
		"begzar":    true,
		"electrotm": true,
		"dynx":      true,
		"radar":     true,
		"shatel":    true,
		"level3":    true,
		"shelter":   true,
		"beshkan":   true,
	}

	var customDNS []string
	for name := range dnsServers {
		if !defaultDNSNames[name] {
			customDNS = append(customDNS, name)
		}
	}
	return customDNS
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
		dnsServers[strings.ToLower(name)] = ips
	}
}

func saveCustomDNS() {
	defaultDNSNames := map[string]bool{
		"403":       true,
		"shecan":    true,
		"begzar":    true,
		"electrotm": true,
		"dynx":      true,
		"radar":     true,
		"shatel":    true,
		"level3":    true,
		"shelter":   true,
		"beshkan":   true,
	}

	customDNS := make(map[string][]string)
	for name, ips := range dnsServers {
		if !defaultDNSNames[name] {
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
