package tokens

// Palette defines colors for TUI and CLI output.
type Palette struct {
	FG        string // default foreground
	BG        string // background
	Accent    string
	Success   string
	Warn      string
	Danger    string
	Muted     string
	Border    string
	Highlight string
}

// Original: neutral modern.
var Original = Palette{
	FG:        "#e4e4e7",
	BG:        "#18181b",
	Accent:    "#3b82f6",
	Success:   "#22c55e",
	Warn:      "#eab308",
	Danger:    "#ef4444",
	Muted:     "#71717a",
	Border:    "#3f3f46",
	Highlight: "#a1a1aa",
}

// Corporate: clean, subdued, enterprise.
var Corporate = Palette{
	FG:        "#374151",
	BG:        "#f9fafb",
	Accent:    "#1d4ed8",
	Success:   "#047857",
	Warn:      "#b45309",
	Danger:    "#b91c1c",
	Muted:     "#6b7280",
	Border:    "#d1d5db",
	Highlight: "#111827",
}

// Cyber: high-contrast, neon cyber vibe.
var Cyber = Palette{
	FG:        "#e0e0e0",
	BG:        "#0d0d0d",
	Accent:    "#00ff9f",
	Success:   "#00ff41",
	Warn:      "#ffcc00",
	Danger:    "#ff0040",
	Muted:     "#808080",
	Border:    "#00ff9f",
	Highlight: "#00ffff",
}

// ByName returns a palette by config name.
func ByName(name string) Palette {
	switch name {
	case "corporate":
		return Corporate
	case "cyber":
		return Cyber
	default:
		return Original
	}
}
