package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mlamkadm/aart/internal/config"
)

// ThemeSelector allows users to preview and select themes
type ThemeSelector struct {
	width      int
	height     int
	theme      Theme
	styles     Styles
	config     *config.Config
	returnTo   tea.Model
	
	selectedIndex int
	themes        []string
	previews      map[string]string
}

func NewThemeSelector(cfg *config.Config, returnTo tea.Model) ThemeSelector {
	themeName := cfg.UI.Theme
	if themeName == "" {
		themeName = "tokyo-night"
	}
	theme := GetTheme(themeName)
	
	themes := AvailableThemes()
	selectedIdx := 0
	for i, t := range themes {
		if t == cfg.UI.Theme {
			selectedIdx = i
			break
		}
	}
	
	return ThemeSelector{
		theme:         theme,
		styles:        NewStyles(theme),
		config:        cfg,
		returnTo:      returnTo,
		themes:        themes,
		selectedIndex: selectedIdx,
		previews:      generateThemePreviews(),
	}
}

func (t ThemeSelector) Init() tea.Cmd {
	return nil
}

func (t ThemeSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			// Return without saving
			if t.returnTo != nil {
				return t.returnTo, nil
			}
			return NewStartupPage(t.config), nil
		
		case "up", "k":
			if t.selectedIndex > 0 {
				t.selectedIndex--
				t.updatePreview()
			}
		
		case "down", "j":
			if t.selectedIndex < len(t.themes)-1 {
				t.selectedIndex++
				t.updatePreview()
			}
		
		case "enter", " ":
			// Apply theme
			t.config.UI.Theme = t.themes[t.selectedIndex]
			config.Save(t.config)
			
			// Return to previous screen with new theme
			if t.returnTo != nil {
				return NewStartupPage(t.config), nil
			}
			return NewStartupPage(t.config), nil
		}
	
	case tea.WindowSizeMsg:
		t.width = msg.Width
		t.height = msg.Height
	}
	
	return t, nil
}

func (t *ThemeSelector) updatePreview() {
	// Update theme for live preview
	t.theme = GetTheme(t.themes[t.selectedIndex])
	t.styles = NewStyles(t.theme)
}

func (t ThemeSelector) View() string {
	width := t.width
	height := t.height
	if width == 0 {
		width = 120
		height = 40
	}
	
	var b strings.Builder
	
	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(t.theme.AccentPrimary).
		Bold(true).
		Align(lipgloss.Center).
		Width(100)
	
	b.WriteString(titleStyle.Render("🌈 Theme Selector"))
	b.WriteString("\n\n")
	
	// Theme list on left, preview on right
	leftPanel := t.renderThemeList()
	rightPanel := t.renderPreview()
	
	combined := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		"  ",
		rightPanel,
	)
	
	b.WriteString(combined)
	b.WriteString("\n\n")
	
	// Instructions
	hintStyle := lipgloss.NewStyle().
		Foreground(t.theme.FgMuted).
		Italic(true).
		Align(lipgloss.Center).
		Width(100)
	
	hints := "hjkl/arrows: navigate │ enter/space: apply theme │ esc/q: cancel"
	b.WriteString(hintStyle.Render(hints))
	
	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		b.String(),
	)
}

func (t ThemeSelector) renderThemeList() string {
	var items []string
	
	for i, themeName := range t.themes {
		var style lipgloss.Style
		prefix := "  "
		
		if i == t.selectedIndex {
			style = lipgloss.NewStyle().
				Background(t.theme.Selection).
				Foreground(t.theme.AccentPrimary).
				Bold(true).
				Width(25).
				Padding(0, 2)
			prefix = "▶ "
		} else {
			style = lipgloss.NewStyle().
				Foreground(t.theme.FgSecondary).
				Width(25).
				Padding(0, 2)
		}
		
		displayName := strings.ReplaceAll(themeName, "-", " ")
		displayName = strings.Title(displayName)
		
		items = append(items, style.Render(prefix+displayName))
	}
	
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.theme.BorderActive).
		Padding(1, 2).
		Width(30)
	
	return box.Render(strings.Join(items, "\n"))
}

func (t ThemeSelector) renderPreview() string {
	preview := t.previews[t.themes[t.selectedIndex]]
	if preview == "" {
		preview = t.generatePreview()
	}
	
	// Apply current theme colors to preview
	previewStyle := lipgloss.NewStyle().
		Foreground(t.theme.FgPrimary)
	
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.theme.BorderActive).
		Padding(1, 2).
		Width(60)
	
	return box.Render(previewStyle.Render(preview))
}

func (t ThemeSelector) generatePreview() string {
	var b strings.Builder
	
	// Show color samples
	titleStyle := lipgloss.NewStyle().
		Foreground(t.theme.AccentPrimary).
		Bold(true)
	
	b.WriteString(titleStyle.Render(fmt.Sprintf("Preview: %s", t.themes[t.selectedIndex])))
	b.WriteString("\n\n")
	
	// Primary accent
	accentStyle := lipgloss.NewStyle().
		Foreground(t.theme.AccentPrimary).
		Bold(true)
	b.WriteString(accentStyle.Render("█████ Primary Accent"))
	b.WriteString("\n")
	
	// Secondary accent
	accent2Style := lipgloss.NewStyle().
		Foreground(t.theme.AccentSecondary).
		Bold(true)
	b.WriteString(accent2Style.Render("█████ Secondary Accent"))
	b.WriteString("\n")
	
	// Info accent
	infoStyle := lipgloss.NewStyle().
		Foreground(t.theme.AccentInfo)
	b.WriteString(infoStyle.Render("█████ Info"))
	b.WriteString("\n")
	
	// Success
	successStyle := lipgloss.NewStyle().
		Foreground(t.theme.AccentSuccess)
	b.WriteString(successStyle.Render("█████ Success"))
	b.WriteString("\n")
	
	// Warning
	warnStyle := lipgloss.NewStyle().
		Foreground(t.theme.AccentWarning)
	b.WriteString(warnStyle.Render("█████ Warning"))
	b.WriteString("\n")
	
	// Error
	errorStyle := lipgloss.NewStyle().
		Foreground(t.theme.AccentError)
	b.WriteString(errorStyle.Render("█████ Error"))
	b.WriteString("\n\n")
	
	// Text samples
	fgStyle := lipgloss.NewStyle().Foreground(t.theme.FgPrimary)
	b.WriteString(fgStyle.Render("Primary text"))
	b.WriteString("\n")
	
	fg2Style := lipgloss.NewStyle().Foreground(t.theme.FgSecondary)
	b.WriteString(fg2Style.Render("Secondary text"))
	b.WriteString("\n")
	
	mutedStyle := lipgloss.NewStyle().Foreground(t.theme.FgMuted)
	b.WriteString(mutedStyle.Render("Muted text"))
	b.WriteString("\n\n")
	
	// Selection sample
	selStyle := lipgloss.NewStyle().
		Background(t.theme.Selection).
		Foreground(t.theme.AccentPrimary).
		Padding(0, 2)
	b.WriteString(selStyle.Render("Selected item"))
	b.WriteString("\n")
	
	return b.String()
}

func generateThemePreviews() map[string]string {
	// Pre-generate previews for all themes
	// For now, return empty map - previews will be generated on demand
	return make(map[string]string)
}
