package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/olekukonko/tablewriter"
)

var (
	PrimaryColor   = lipgloss.Color("63")  // Purple-ish
	SecondaryColor = lipgloss.Color("39")  // Blue-ish
	SuccessColor   = lipgloss.Color("42")  // Green
	ErrorColor     = lipgloss.Color("196") // Red
	BorderColor    = lipgloss.Color("240") // Grey

	LogoStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			MarginBottom(1)

	TitleStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			MarginBottom(1)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Bold(true).
			MarginBottom(1)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true)
)

func PrintLogo() {
	logo := `
  ___ __  ___   _____ ___ 
 / _ \  \/  | | / /_ _|__ \
|  __/_/    | |/ / | |/_/ /
 \___/_/|_|___/_/  |___/___/
`
	fmt.Println(LogoStyle.Render(logo))
}

func RenderTable(headers []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetBorder(false)
	table.SetRowLine(false)
	table.SetAutoWrapText(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

	// Style headers
	// Style headers
	// Basic dynamic coloring for headers up to 10 columns
	colors := make([]tablewriter.Colors, len(headers))
	for i := range colors {
		colors[i] = tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor}
	}
	table.SetHeaderColor(colors...)

	table.AppendBulk(data)
	table.Render()
}

func RenderKV(title string, data map[string]string) {
	fmt.Println(TitleStyle.Render(title))

	// Calculate max key length for alignment
	maxKeyLen := 0
	for k := range data {
		if len(k) > maxKeyLen {
			maxKeyLen = len(k)
		}
	}

	for k, v := range data {
		padding := maxKeyLen - len(k)
		padStr := ""
		for i := 0; i < padding; i++ {
			padStr += " "
		}
		key := lipgloss.NewStyle().Foreground(SecondaryColor).Bold(true).Render(k + ": " + padStr)
		val := lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render(v) // White-ish
		fmt.Printf("%s %s\n", key, val)
	}
}

func PrintSuccess(msg string) {
	fmt.Println(SuccessStyle.Render("✓ " + msg))
}

func PrintError(err error) {
	fmt.Println(ErrorStyle.Render("✗ Error: " + err.Error()))
}
