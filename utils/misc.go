package utils

import (
	"fmt"
)

func Bold(text string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", text)
}

