package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mlamkadm/aart/internal/config"
)

// SettingsScreen for configuring preferences
type SettingsScreen struct {
	width    int
	height   int
	theme    Theme
	styles   Styles
	config   *config.Config
	returnTo tea.Model
	
	selectedIndex int
	settings      []settingItem
}

type settingItem struct {
	Name        string
	Description string
	Value       interface{}
	Type        string // "bool", "int", "string", "select"
	Options     []string
}

func NewSettingsScreen(cfg *config.Config, returnTo tea.Model) SettingsScreen {
	themeName := cfg.UI.Theme
	if themeName == "" {
		themeName = "tokyo-night"
	}
	theme := GetTheme(themeName)
	
	settings := []settingItem{
		{
			Name:        "Theme",
			Description: "Color scheme for the UI",
			Value:       cfg.UI.Theme,
			Type:        "select",
			Options:     AvailableThemes(),
		},
		{
			Name:        "Default Width",
			Description: "Default canvas width for new animations",
			Value:       cfg.Editor.DefaultWidth,
			Type:        "int",
		},
		{
			Name:        "Default Height",
			Description: "Default canvas height for new animations",
			Value:       cfg.Editor.DefaultHeight,
			Type:        "int",
		},
		{
			Name:        "Default FPS",
			Description: "Default frames per second",
			Value:       cfg.Editor.DefaultFPS,
			Type:        "int",
		},
		{
			Name:        "Show Grid",
			Description: "Show grid lines in editor",
			Value:       cfg.Editor.ShowGrid,
			Type:        "bool",
		},
		{
			Name:        "Zen Mode",
			Description: "Start in zen mode (minimal UI)",
			Value:       cfg.Editor.ZenMode,
			Type:        "bool",
		},
	}
	
	return SettingsScreen{
		theme:         theme,
		styles:        NewStyles(theme),
		config:        cfg,
		returnTo:      returnTo,
		settings:      settings,
		selectedIndex: 0,
	}
}

func (s SettingsScreen) Init() tea.Cmd {
	return nil
}

func (s SettingsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			// Save and return
			config.Save(s.config)
			if s.returnTo != nil {
				return s.returnTo, nil
			}
			return NewStartupPage(s.config), nil
		
		case "up", "k":
			if s.selectedIndex > 0 {
				s.selectedIndex--
			}
		
		case "down", "j":
			if s.selectedIndex < len(s.settings)-1 {
				s.selectedIndex++
			}
		
		case "left", "h", "-":
			s.decrementValue()
		
		case "right", "l", "+", " ", "enter":
			s.incrementValue()
		}
	
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
	}
	
	return s, nil
}

func (s *SettingsScreen) incrementValue() {
	setting := &s.settings[s.selectedIndex]
	
	switch setting.Type {
	case "bool":
		setting.Value = !setting.Value.(bool)
		s.applyToConfig(setting)
	case "int":
		setting.Value = setting.Value.(int) + 1
		s.applyToConfig(setting)
	case "select":
		currentIdx := 0
		for i, opt := range setting.Options {
			if opt == setting.Value.(string) {
				currentIdx = i
				break
			}
		}
		nextIdx := (currentIdx + 1) % len(setting.Options)
		setting.Value = setting.Options[nextIdx]
		s.applyToConfig(setting)
	}
}

func (s *SettingsScreen) decrementValue() {
	setting := &s.settings[s.selectedIndex]
	
	switch setting.Type {
	case "bool":
		setting.Value = !setting.Value.(bool)
		s.applyToConfig(setting)
	case "int":
		if setting.Value.(int) > 1 {
			setting.Value = setting.Value.(int) - 1
			s.applyToConfig(setting)
		}
	case "select":
		currentIdx := 0
		for i, opt := range setting.Options {
			if opt == setting.Value.(string) {
				currentIdx = i
				break
			}
		}
		prevIdx := (currentIdx - 1 + len(setting.Options)) % len(setting.Options)
		setting.Value = setting.Options[prevIdx]
		s.applyToConfig(setting)
	}
}

func (s *SettingsScreen) applyToConfig(setting *settingItem) {
	switch setting.Name {
	case "Theme":
		s.config.UI.Theme = setting.Value.(string)
		// Reload theme
		s.theme = GetTheme(s.config.UI.Theme)
		s.styles = NewStyles(s.theme)
	case "Default Width":
		s.config.Editor.DefaultWidth = setting.Value.(int)
	case "Default Height":
		s.config.Editor.DefaultHeight = setting.Value.(int)
	case "Default FPS":
		s.config.Editor.DefaultFPS = setting.Value.(int)
	case "Show Grid":
		s.config.Editor.ShowGrid = setting.Value.(bool)
	case "Zen Mode":
		s.config.Editor.ZenMode = setting.Value.(bool)
	}
}

func (s SettingsScreen) View() string {
	// Use defaults if not set yet
	width := s.width
	height := s.height
	if width == 0 {
		width = 120
		height = 40
	}
	
	var b strings.Builder
	
	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(s.theme.AccentPrimary).
		Bold(true).
		Align(lipgloss.Center).
		Width(80)
	
	b.WriteString(titleStyle.Render("⚙️  Settings"))
	b.WriteString("\n\n")
	
	// Settings list
	for i, setting := range s.settings {
		var style lipgloss.Style
		prefix := "  "
		
		if i == s.selectedIndex {
			style = lipgloss.NewStyle().
				Background(s.theme.Selection).
				Foreground(s.theme.AccentPrimary).
				Bold(true).
				Width(70).
				Padding(0, 2)
			prefix = "▶ "
		} else {
			style = lipgloss.NewStyle().
				Foreground(s.theme.FgSecondary).
				Width(70).
				Padding(0, 2)
		}
		
		valueStr := fmt.Sprintf("%v", setting.Value)
		if setting.Type == "bool" {
			if setting.Value.(bool) {
				valueStr = "✓ ON"
			} else {
				valueStr = "✗ OFF"
			}
		}
		
		valueStyle := lipgloss.NewStyle().
			Foreground(s.theme.AccentSecondary).
			Bold(true)
		
		line := fmt.Sprintf("%s%-20s  %s", prefix, setting.Name, valueStyle.Render(valueStr))
		
		b.WriteString(style.Render(line))
		b.WriteString("\n")
		
		if i == s.selectedIndex {
			descStyle := lipgloss.NewStyle().
				Foreground(s.theme.FgMuted).
				Italic(true).
				Padding(0, 4)
			b.WriteString(descStyle.Render(setting.Description))
			b.WriteString("\n")
		}
	}
	
	// Instructions
	b.WriteString("\n\n")
	hintStyle := lipgloss.NewStyle().
		Foreground(s.theme.FgMuted).
		Italic(true)
	
	hints := "hjkl/arrows: navigate │ space/enter/+: increase │ -/h: decrease │ q/esc: save & exit"
	b.WriteString(hintStyle.Render(hints))
	
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(s.theme.BorderActive).
		Padding(2, 4).
		Width(80)
	
	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		box.Render(b.String()),
	)
}

// ExamplesScreen shows example animations
type ExamplesScreen struct {
	width    int
	height   int
	theme    Theme
	styles   Styles
	config   *config.Config
	returnTo tea.Model
}

func NewExamplesScreen(cfg *config.Config, returnTo tea.Model) ExamplesScreen {
	themeName := cfg.UI.Theme
	if themeName == "" {
		themeName = "tokyo-night"
	}
	theme := GetTheme(themeName)
	
	return ExamplesScreen{
		theme:    theme,
		styles:   NewStyles(theme),
		config:   cfg,
		returnTo: returnTo,
	}
}

func (e ExamplesScreen) Init() tea.Cmd {
	return nil
}

func (e ExamplesScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if e.returnTo != nil {
			return e.returnTo, nil
		}
		return NewStartupPage(e.config), nil
	
	case tea.WindowSizeMsg:
		e.width = msg.Width
		e.height = msg.Height
	}
	
	return e, nil
}

func (e ExamplesScreen) View() string {
	// Use defaults if not set yet
	width := e.width
	height := e.height
	if width == 0 {
		width = 120
		height = 40
	}
	
	titleStyle := lipgloss.NewStyle().
		Foreground(e.theme.AccentPrimary).
		Bold(true).
		Align(lipgloss.Center).
		Width(80)
	
	content := titleStyle.Render("📚 Examples") + "\n\n"
	
	msgStyle := lipgloss.NewStyle().
		Foreground(e.theme.FgMuted).
		Italic(true).
		Align(lipgloss.Center).
		Width(80)
	
	content += msgStyle.Render("Example gallery coming soon!\n\n")
	content += msgStyle.Render("Press any key to return")
	
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(e.theme.Border).
		Padding(2, 4)
	
	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		box.Render(content),
	)
}
