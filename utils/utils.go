package utils

import (
	"regexp"
	"strings"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

func GenerateResourceName(str string) string {
	title := strings.ReplaceAll(str, " ", "_")
	title = strings.ReplaceAll(title, "'", "")
	title = strings.ToLower(title)
	return nonAlphanumericRegex.ReplaceAllString(title, "_")
}

func Contains(slice []string, inputValue string) bool {
	for _, sliceValue := range slice {
		if sliceValue == inputValue {
			return true
		}
	}
	return false
}
