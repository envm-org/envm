package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

func getIconForLabel(label string) string {
	lower := strings.ToLower(label)
	if strings.Contains(lower, "email") || strings.Contains(lower, "mail") {
		return "✉ "
	}
	if strings.Contains(lower, "password") || strings.Contains(lower, "key") {
		return "🔑 "
	}
	if strings.Contains(lower, "name") || strings.Contains(lower, "user") {
		return "👤 "
	}
	if strings.Contains(lower, "project") {
		return "📁 "
	}
	return "📝 "
}

func Prompt(label string) string {
	icon := getIconForLabel(label)
	fmt.Print(HelpCommandStyle.Render(icon) + lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render(label+": "))
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func PromptPassword(label string) string {
	icon := getIconForLabel(label)
	fmt.Print(HelpCommandStyle.Render(icon) + lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render(label+" (typing is hidden): "))
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return ""
	}
	return string(bytePassword)
}
