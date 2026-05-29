package dns

import (
	"slices"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// contains any
func containsAny(slice []string, items []string) bool {
	for _, item := range items {
		if containsString(slice, item) {
			return true
		}
	}
	return false
}

// contains string
func containsString(slice []string, item string) bool {
	return slices.Contains(slice, item)
}

// remove duplicates
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

// format dns name
func formatDNSName(name string) string {
	if name == "403" {
		return "403"
	}
	return cases.Title(language.English).String(name)
}
