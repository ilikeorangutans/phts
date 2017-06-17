package model

import (
	"fmt"
	"regexp"
	"strings"
)

var slugFilterRegexp = regexp.MustCompile("[^a-z0-9_]+")

func SlugFromString(input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", fmt.Errorf("empty slugs not allowed")
	}

	input = strings.ToLower(input)

	input = slugFilterRegexp.ReplaceAllString(input, "-")
	if len(input) > 128 {
		return "", fmt.Errorf("slug longer than 128 characters not allowed")
	}

	return input, nil
}
