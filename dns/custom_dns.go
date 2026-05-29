package dns

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	_ "embed"
)

//go:embed custom_dns.json
var embeddedCustomDNS []byte

// save custom dns
func saveCustomDNS() {
	defaultDNSNames := map[string]bool{
		"403":        true,
		"shecan":     true,
		"shecan-pro": true,
		"begzar":     true,
		"electrotm":  true,
		"dynx":       true,
		"radar":      true,
		"shatel":     true,
		"level3":     true,
		"shelter":    true,
		"beshkan":    true,
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

// load custom dns
func LoadCustomDNS() {
	data, err := os.ReadFile(customDNSFile)
	if err != nil {
		if len(embeddedCustomDNS) == 0 {
			return
		}
		data = embeddedCustomDNS
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

// get custom dns list
func getCustomDNSList() []string {
	defaultDNSNames := map[string]bool{
		"403":        true,
		"shecan":     true,
		"shecan-pro": true,
		"begzar":     true,
		"electrotm":  true,
		"dynx":       true,
		"radar":      true,
		"shatel":     true,
		"level3":     true,
		"shelter":    true,
		"beshkan":    true,
	}

	var customDNS []string
	for name := range dnsServers {
		if !defaultDNSNames[name] {
			customDNS = append(customDNS, name)
		}
	}
	return customDNS
}
