package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderDiff renders a unified diff with +/- coloring using the given styles.
func RenderDiff(diff string, addStyle, delStyle lipgloss.Style) string {
	var out []string
	for _, line := range strings.Split(diff, "\n") {
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			out = append(out, addStyle.Render(line))
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			out = append(out, delStyle.Render(line))
		} else {
			out = append(out, line)
		}
	}
	return strings.Join(out, "\n")
}
