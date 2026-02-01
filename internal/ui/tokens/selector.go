package tokens

import "strconv"

// Select returns the palette for the given config palette value.
func Select(paletteName string) Palette {
	return ByName(paletteName)
}

// AccentANSI returns an ANSI 24-bit SGR escape for the palette's accent color.
func (p Palette) AccentANSI() string {
	return hexToANSI(p.Accent)
}

// hexToANSI returns \033[38;2;r;g;bm for the given hex color (e.g. #3b82f6).
func hexToANSI(hex string) string {
	if len(hex) >= 7 && hex[0] == '#' {
		hex = hex[1:]
	}
	if len(hex) != 6 {
		return ""
	}
	r, _ := strconv.ParseInt(hex[0:2], 16, 0)
	g, _ := strconv.ParseInt(hex[2:4], 16, 0)
	b, _ := strconv.ParseInt(hex[4:6], 16, 0)
	return "\033[38;2;" + strconv.Itoa(int(r)) + ";" + strconv.Itoa(int(g)) + ";" + strconv.Itoa(int(b)) + "m"
}
