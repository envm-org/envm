package ui

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// ... skipping to PrintError below ...

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

	InfoStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")). // Orange-ish
			Bold(true)

	// CLI Help specific styles
	HelpHeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("220")) // Yellow

	HelpCommandStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("42")) // Green
)

// SetupCobraHelp overrides the default Cobra templates
func SetupCobraHelp(cmd *cobra.Command) {
	cobra.AddTemplateFunc("StyleHeading", func(s string) string {
		return HelpHeaderStyle.Render(s)
	})
	cobra.AddTemplateFunc("StyleCommand", func(s string) string {
		return HelpCommandStyle.Render(s)
	})

	customUsageTemplate := `{{StyleHeading "Usage:"}}
  {{StyleCommand .UseLine}}{{if .HasAvailableSubCommands}}

{{StyleHeading "Available commands:"}}{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{StyleCommand (rpad .Name .NamePadding)}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

{{StyleHeading "Options:"}}
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}
`

	helpTemplate := `{{if .Short}}{{.Short}}
{{end}}{{if .Long}}
{{.Long}}
{{end}}
` + customUsageTemplate

	cmd.SetUsageTemplate(customUsageTemplate)
	cmd.SetHelpTemplate(helpTemplate)

	// Apply to all initialized subcommands
	for _, sub := range cmd.Commands() {
		sub.SetUsageTemplate(customUsageTemplate)
		sub.SetHelpTemplate(helpTemplate)
	}
}

func PrintInfo(msg string) {
	fmt.Println(InfoStyle.Render("ℹ " + msg))
}

func PrintWarning(msg string) {
	fmt.Println(WarningStyle.Render("⚠ " + msg))
}

func PrintLogo() {
	logo := `
  ___ __  ___   _____ ___ 
 / _ \  \/  | | / /_ _|__ \
|  __/_/    | |/ / | |/_/ /
 \___/_/|_|___/_/  |___/___/
`
	fmt.Println(HelpCommandStyle.Render(logo))
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

var validationErrRegex = regexp.MustCompile(`Key: '[^.]+\.([^']+)' Error:Field validation for '[^']+' failed on the '([^']+)' tag`)

func PrintError(err error) {
	errMsg := err.Error()

	// Check for raw validation errors from validator
	if validationErrRegex.MatchString(errMsg) {
		lines := strings.Split(errMsg, "\n")
		for _, line := range lines {
			if matches := validationErrRegex.FindStringSubmatch(line); len(matches) > 2 {
				field := matches[1]
				tag := matches[2]

				formattedErr := fmt.Sprintf("%s failed validation (rule: %s)", field, tag)
				// Enhance common tags with friendlier messages
				switch tag {
				case "required":
					formattedErr = fmt.Sprintf("%s is required.", field)
				case "email":
					formattedErr = fmt.Sprintf("%s must be a valid email address.", field)
				case "min":
					formattedErr = fmt.Sprintf("%s is too short.", field)
				case "max":
					formattedErr = fmt.Sprintf("%s is too long.", field)
				}

				fmt.Println(ErrorStyle.Render("✗ " + formattedErr))
			} else if strings.TrimSpace(line) != "" {
				// Fallback for lines that don't match the regex but exist in the error
				fmt.Println(ErrorStyle.Render("✗ " + strings.TrimSpace(line)))
			}
		}
		return
	}

	if strings.Contains(errMsg, "connection refused") || strings.Contains(errMsg, "dial tcp") {
		errMsg = "unable to connect to the server. Please check your internet connection or try again later."
	} else if strings.Contains(errMsg, "api error: ") {
		errMsg = strings.Replace(errMsg, "api error: ", "", 1)
	}

	fmt.Println(ErrorStyle.Render("✗ " + errMsg))
}
