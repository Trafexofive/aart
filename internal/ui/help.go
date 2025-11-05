package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mlamkadm/aart/internal/config"
)

// HelpScreen shows keyboard shortcuts and guide
type HelpScreen struct {
	width  int
	height int
	theme  Theme
	styles Styles
	config *config.Config
}

func NewHelpScreen(cfg *config.Config) HelpScreen {
	themeName := cfg.UI.Theme
	if themeName == "" {
		themeName = "tokyo-night"
	}
	theme := GetTheme(themeName)
	
	return HelpScreen{
		theme:  theme,
		styles: NewStyles(theme),
		config: cfg,
	}
}

func (h HelpScreen) Init() tea.Cmd {
	return nil
}

func (h HelpScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "enter":
			// Return to startup
			return NewStartupPage(h.config), nil
		}
	case tea.WindowSizeMsg:
		h.width = msg.Width
		h.height = msg.Height
	}
	
	return h, nil
}

func (h HelpScreen) View() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(h.theme.AccentPrimary).
		Bold(true).
		Align(lipgloss.Center).
		Width(80)
	
	sectionStyle := lipgloss.NewStyle().
		Foreground(h.theme.AccentSecondary).
		Bold(true).
		MarginTop(1)
	
	keyStyle := lipgloss.NewStyle().
		Foreground(h.theme.AccentInfo).
		Bold(true).
		Width(15)
	
	descStyle := lipgloss.NewStyle().
		Foreground(h.theme.FgSecondary)
	
	var content strings.Builder
	
	content.WriteString(titleStyle.Render("⌨️  Keyboard Shortcuts"))
	content.WriteString("\n\n")
	
	// Navigation
	content.WriteString(sectionStyle.Render("Navigation"))
	content.WriteString("\n")
	content.WriteString(keyStyle.Render("hjkl / arrows") + descStyle.Render("Move cursor"))
	content.WriteString("\n")
	content.WriteString(keyStyle.Render(", .") + descStyle.Render("Previous/next frame"))
	content.WriteString("\n")
	content.WriteString(keyStyle.Render("space") + descStyle.Render("Play/pause animation"))
	content.WriteString("\n\n")
	
	// Modes
	content.WriteString(sectionStyle.Render("Modes"))
	content.WriteString("\n")
	content.WriteString(keyStyle.Render("i") + descStyle.Render("Insert mode"))
	content.WriteString("\n")
	content.WriteString(keyStyle.Render("d") + descStyle.Render("Draw mode"))
	content.WriteString("\n")
	content.WriteString(keyStyle.Render(":") + descStyle.Render("Command mode"))
	content.WriteString("\n")
	content.WriteString(keyStyle.Render("esc") + descStyle.Render("Normal mode"))
	content.WriteString("\n\n")
	
	// View
	content.WriteString(sectionStyle.Render("View"))
	content.WriteString("\n")
	content.WriteString(keyStyle.Render("z") + descStyle.Render("Toggle zen mode"))
	content.WriteString("\n")
	content.WriteString(keyStyle.Render("g") + descStyle.Render("Toggle grid"))
	content.WriteString("\n")
	content.WriteString(keyStyle.Render("+ -") + descStyle.Render("Zoom in/out"))
	content.WriteString("\n\n")
	
	// Tools
	content.WriteString(sectionStyle.Render("Tools"))
	content.WriteString("\n")
	content.WriteString(keyStyle.Render("ctrl-j/k") + descStyle.Render("Cycle wheel menu"))
	content.WriteString("\n\n")
	
	// Commands
	content.WriteString(sectionStyle.Render("Commands"))
	content.WriteString("\n")
	content.WriteString(keyStyle.Render(":w") + descStyle.Render("Save"))
	content.WriteString("\n")
	content.WriteString(keyStyle.Render(":q") + descStyle.Render("Quit"))
	content.WriteString("\n")
	content.WriteString(keyStyle.Render(":export") + descStyle.Render("Export to file"))
	content.WriteString("\n\n")
	
	// Footer
	footerStyle := lipgloss.NewStyle().
		Foreground(h.theme.FgMuted).
		Italic(true).
		Align(lipgloss.Center).
		Width(80).
		MarginTop(2)
	
	content.WriteString(footerStyle.Render("Press any key to return"))
	
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(h.theme.BorderActive).
		Padding(2, 4).
		Width(82)
	
	return lipgloss.Place(
		h.width,
		h.height,
		lipgloss.Center,
		lipgloss.Center,
		box.Render(content.String()),
	)
}
