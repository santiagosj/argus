package styles

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func GetBanner() string {
	banner := `
  
  ------------------------------------------------------
         █████╗ ██████╗  ██████╗ ██╗   ██╗███████╗
        ██╔══██╗██╔══██╗██╔════╝ ██║   ██║██╔════╝
        ███████║██████╔╝██║  ███╗██║   ██║███████╗
        ██╔══██║██╔══██╗██║   ██║██║   ██║╚════██║
        ██║  ██║██║  ██║╚██████╔╝╚██████╔╝███████║
        ╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝  ╚═════╝ ╚══════╝
  ------------------------------------------------------
  `

	lines := strings.Split(banner, "\n")
	var coloredBanner strings.Builder

	// Rose Pine palette colors
	colors := []string{
		"#ebbcba", // Mauve
		"#c4a7e7", // Lavender
		"#9ccfd8", // Teal
		"#f6c177", // Peach
		"#eb6f92", // Red
		"#31748f", // Blue
	}

	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		// Seleccionar color basado en la posición de la línea (gradiente vertical)
		colorIndex := (i * len(colors)) / len(lines)
		if colorIndex >= len(colors) {
			colorIndex = len(colors) - 1
		}

		style := lipgloss.NewStyle().Foreground(lipgloss.Color(colors[colorIndex])).Bold(true)
		coloredBanner.WriteString(style.Render(line) + "\n")
	}

	return lipgloss.NewStyle().MarginLeft(2).MarginBottom(1).Render(coloredBanner.String())
}
