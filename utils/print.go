package utils

import (
	"time"

	"github.com/pterm/pterm"
)

func PrintfWithTime(msg string) {
	pterm.Printf("[%s] %s\n", time.Now().Format("15:04:05"), msg)
}

func SprintfWithTime(msg string) string {
	return pterm.Sprintf("[%s] %s\n", time.Now().Format("15:04:05"), msg)
}
