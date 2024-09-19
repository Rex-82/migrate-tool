package utils

import (
	"strings"
)

func PathFormat(s string) string {
	if !strings.HasSuffix(s, "/") {
		s += "/"
	}
	return s
}
