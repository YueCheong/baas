package utils

import (
	"strings"
)

func SplitContainerName(c string) (string, string) {
	segments := strings.Split(c, ".")
	var name, domain string

	name = segments[0]

	for i, segment := range segments[1:] {
		domain += segment
		if i != len(segments)-2 {
			domain += "."
		}
	}

	return name, domain
}
