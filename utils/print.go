package utils

import (
	"time"

	"github.com/pterm/pterm"
)

func PrintInfo(msg string) {
	pterm.Printf("[%s] %s\n", time.Now().Format("15:04:05"), pterm.LightCyan(msg))
}

func PrintSuccess(msg string) {
	pterm.Printf("[%s] %s\n", time.Now().Format("15:04:05"), pterm.Green(msg))
}

func PrintWarn(msg string) {
	pterm.Printf("[%s] %s\n", time.Now().Format("15:04:05"), pterm.Yellow(msg))
}

func PrintError(msg string) {
	pterm.Printf("[%s] %s\n", time.Now().Format("15:04:05"), pterm.Red(msg))
}
