package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// Prompt asks the user for input with a label
func Prompt(label string) string {
	fmt.Print(TitleStyle.Render(label + ": "))
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// PromptPassword asks the user for a password with a label, masking input
func PromptPassword(label string) string {
	fmt.Print(TitleStyle.Render(label + ": "))
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // Move to new line after enter
	if err != nil {
		return ""
	}
	return string(bytePassword)
}
