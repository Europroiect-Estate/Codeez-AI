package ui

import (
	"github.com/Europroiect-Estate/Codeez-AI/internal/ui/tokens"
	"github.com/charmbracelet/lipgloss"
)

// Styles holds Lip Gloss styles driven by the current palette.
type Styles struct {
	Title    lipgloss.Style
	Subtitle lipgloss.Style
	Body     lipgloss.Style
	Muted    lipgloss.Style
	Accent   lipgloss.Style
	Success  lipgloss.Style
	Warn     lipgloss.Style
	Danger   lipgloss.Style
	Border   lipgloss.Style
	Panel    lipgloss.Style
}

// NewStyles returns styles for the given palette name.
func NewStyles(paletteName string) Styles {
	p := tokens.Select(paletteName)
	return Styles{
		Title:    lipgloss.NewStyle().Foreground(lipgloss.Color(p.Accent)).Bold(true),
		Subtitle: lipgloss.NewStyle().Foreground(lipgloss.Color(p.Muted)),
		Body:     lipgloss.NewStyle().Foreground(lipgloss.Color(p.FG)),
		Muted:    lipgloss.NewStyle().Foreground(lipgloss.Color(p.Muted)),
		Accent:   lipgloss.NewStyle().Foreground(lipgloss.Color(p.Accent)),
		Success:  lipgloss.NewStyle().Foreground(lipgloss.Color(p.Success)),
		Warn:     lipgloss.NewStyle().Foreground(lipgloss.Color(p.Warn)),
		Danger:   lipgloss.NewStyle().Foreground(lipgloss.Color(p.Danger)),
		Border:   lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(p.Border)),
		Panel:    lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(p.Border)).Padding(0, 1),
	}
}
