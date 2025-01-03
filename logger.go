package main

import (
	"log"
	"strings"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
)

// Log functions
func LogError(messages ...string) {
	log.Printf("%s[-]%s %s", ColorRed, strings.Join(messages, " "), ColorReset)
}

func LogWarning(messages ...string) {
	log.Printf("%s[-]%s %s", ColorYellow, strings.Join(messages, " "), ColorReset)
}

func LogOK(messages ...string) {
	log.Printf("%s[+]%s %s", ColorGreen, strings.Join(messages, " "), ColorReset)
}
