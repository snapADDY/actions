// Package msgflag parses flags from git commit messages.
package msgflag

import (
	"strings"
)

// Flags is a collection of all defined git commit message flags.
type Flags map[string][]string

// Parse parses a git commit message and saves all flags with options to a Flag struct.
func Parse(msg string) Flags {
	flagStr := extractFlags(msg)

	return parseFlags(flagStr)
}

func extractFlags(msg string) []string {
	open := -1
	result := make([]string, 0)

	for i, r := range msg {
		if r == '[' {
			open = i + 1
		} else if r == ']' && open != -1 {
			result = append(result, msg[open:i])
			open = -1
		}
	}

	return result
}

func parseFlags(flagStr []string) Flags {
	flags := Flags{}

	for _, str := range flagStr {
		str = strings.TrimSpace(str)
		str = strings.ToLower(str)

		f := strings.Split(str, ":")

		flags[f[0]] = f
	}

	return flags
}
