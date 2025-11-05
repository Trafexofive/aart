package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mlamkadm/aart/internal/config"
)

// StartupPage is the initial welcome screen
type StartupPage struct {
	width           int
	height          int
	theme           Theme
	styles          Styles
	config          *config.Config
	selectedOption  int
	options         []StartupOption
	recentFiles     []config.RecentFile
	breathing       *BreathingEffect
	currentTime     time.Time
}

// StartupOption represents a menu option
type StartupOption struct {
	Icon        string
	Title       string
	Description string
	Action      string
	Shortcut    string
}

// NewStartupPage creates the startup screen
func NewStartupPage(cfg *config.Config) StartupPage {
	themeName := cfg.UI.Theme
	if themeName == "" {
		themeName = "tokyo-night"
	}
	theme := GetTheme(themeName)
	
	options := []StartupOption{
		{
			Icon:        "🎨",
			Title:       "New Animation",
			Description: "Create a new ASCII art animation",
			Action:      "new",
			Shortcut:    "n",
		},
		{
			Icon:        "📂",
			Title:       "Open File",
			Description: "Open an existing .aart file",
			Action:      "open",
			Shortcut:    "o",
		},
		{
			Icon:        "🎬",
			Title:       "Import GIF",
			Description: "Convert GIF to ASCII animation",
			Action:      "import",
			Shortcut:    "i",
		},
		{
			Icon:        "🌈",
			Title:       "Change Theme",
			Description: fmt.Sprintf("Current: %s", themeName),
			Action:      "theme",
			Shortcut:    "t",
		},
		{
			Icon:        "⚙️",
			Title:       "Settings",
			Description: "Configure editor preferences",
			Action:      "settings",
			Shortcut:    "s",
		},
		{
			Icon:        "📚",
			Title:       "Examples",
			Description: "Load example animations",
			Action:      "examples",
			Shortcut:    "e",
		},
		{
			Icon:        "❓",
			Title:       "Help",
			Description: "View keyboard shortcuts and guide",
			Action:      "help",
			Shortcut:    "?",
		},
		{
			Icon:        "🚪",
			Title:       "Quit",
			Description: "Exit aart",
			Action:      "quit",
			Shortcut:    "q",
		},
	}
	
	return StartupPage{
		theme:          theme,
		styles:         NewStyles(theme),
		config:         cfg,
		selectedOption: 0,
		options:        options,
		recentFiles:    cfg.GetRecentFiles(),
		breathing:      NewBreathingEffect(4 * time.Second),
		currentTime:    time.Now(),
	}
}

func (s StartupPage) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		tickCmd(),
	)
}

func (s StartupPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return s, tea.Quit
		case "up", "k":
			s.selectedOption = (s.selectedOption - 1 + len(s.options)) % len(s.options)
		case "down", "j":
			s.selectedOption = (s.selectedOption + 1) % len(s.options)
		case "enter", " ":
			return s.handleAction()
		case "n":
			s.selectedOption = 0
			return s.handleAction()
		case "o":
			s.selectedOption = 1
			return s.handleAction()
		case "i":
			s.selectedOption = 2
			return s.handleAction()
		case "t":
			s.selectedOption = 3
			return s.handleAction()
		case "s":
			s.selectedOption = 4
			return s.handleAction()
		case "e":
			s.selectedOption = 5
			return s.handleAction()
		case "?":
			s.selectedOption = 6
			return s.handleAction()
		case "1", "2", "3", "4", "5":
			// Quick access to recent files
			idx := int(msg.String()[0] - '1')
			if idx < len(s.recentFiles) {
				// TODO: Open recent file
				return s, tea.Quit
			}
		}
	
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
	
	case tickMsg:
		s.currentTime = time.Time(msg)
		return s, tickCmd()
	}
	
	return s, nil
}

func (s StartupPage) handleAction() (tea.Model, tea.Cmd) {
	action := s.options[s.selectedOption].Action
	
	switch action {
	case "new":
		// Create new animation - transition to editor
		return NewWithConfig(s.config), nil
	case "quit":
		return s, tea.Quit
	case "theme":
		// Cycle through themes
		themes := AvailableThemes()
		currentIdx := 0
		for i, t := range themes {
			if t == s.config.UI.Theme {
				currentIdx = i
				break
			}
		}
		nextIdx := (currentIdx + 1) % len(themes)
		s.config.UI.Theme = themes[nextIdx]
		config.Save(s.config)
		
		// Reload with new theme
		return NewStartupPage(s.config), nil
	case "help":
		// Show help screen
		return NewHelpScreen(s.config), nil
	default:
		// Other actions not implemented yet
		return s, nil
	}
}

func (s StartupPage) View() string {
	if s.width == 0 {
		return "Loading..."
	}
	
	var b strings.Builder
	
	// Header with ASCII art logo
	b.WriteString(s.renderHeader())
	b.WriteString("\n\n")
	
	// Main menu and recent files side by side
	leftPanel := s.renderMenuPanel()
	rightPanel := s.renderRecentPanel()
	
	// Combine panels
	b.WriteString(lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		"  ",
		rightPanel,
	))
	
	b.WriteString("\n\n")
	
	// Footer with tips
	b.WriteString(s.renderFooter())
	
	return lipgloss.Place(
		s.width,
		s.height,
		lipgloss.Center,
		lipgloss.Center,
		b.String(),
	)
}

func (s StartupPage) renderHeader() string {
	// ASCII art logo with breathing effect
	logo := `
    ▄▄▄        ▄▄▄       ██▀███  ▄▄▄█████▓
   ▒████▄     ▒████▄    ▓██ ▒ ██▒▓  ██▒ ▓▒
   ▒██  ▀█▄   ▒██  ▀█▄  ▓██ ░▄█ ▒▒ ▓██░ ▒░
   ░██▄▄▄▄██  ░██▄▄▄▄██ ▒██▀▀█▄  ░ ▓██▓ ░ 
    ▓█   ▓██▒  ▓█   ▓██▒░██▓ ▒██▒  ▒██▒ ░ 
    ▒▒   ▓▒█░  ▒▒   ▓▒█░░ ▒▓ ░▒▓░  ▒ ░░   
     ▒   ▒▒ ░   ▒   ▒▒ ░  ░▒ ░ ▒░    ░    
     ░   ▒      ░   ▒     ░░   ░   ░      
         ░  ░       ░  ░   ░              
`
	
	logoStyle := lipgloss.NewStyle().
		Foreground(s.theme.AccentPrimary).
		Bold(true)
	
	if s.breathing.CurrentAlpha() > 0.9 {
		logoStyle = logoStyle.Foreground(s.theme.AccentSecondary)
	}
	
	tagline := lipgloss.NewStyle().
		Foreground(s.theme.FgSecondary).
		Italic(true).
		Render("ASCII Art Animation Editor")
	
	version := lipgloss.NewStyle().
		Foreground(s.theme.FgMuted).
		Render("v0.1.0")
	
	return lipgloss.JoinVertical(
		lipgloss.Center,
		logoStyle.Render(logo),
		tagline,
		version,
	)
}

func (s StartupPage) renderMenuPanel() string {
	var menuItems []string
	
	titleStyle := lipgloss.NewStyle().
		Foreground(s.theme.AccentPrimary).
		Bold(true).
		Padding(0, 2)
	
	menuItems = append(menuItems, titleStyle.Render("✨ Quick Start"))
	menuItems = append(menuItems, "")
	
	for i, opt := range s.options {
		var style lipgloss.Style
		prefix := "  "
		
		if i == s.selectedOption {
			// Selected item - highlighted
			style = lipgloss.NewStyle().
				Background(s.theme.Selection).
				Foreground(s.theme.AccentPrimary).
				Bold(true).
				Width(45).
				Padding(0, 2)
			prefix = "▶ "
		} else {
			// Normal item
			style = lipgloss.NewStyle().
				Foreground(s.theme.FgSecondary).
				Width(45).
				Padding(0, 2)
		}
		
		shortcut := lipgloss.NewStyle().
			Foreground(s.theme.FgMuted).
			Render(fmt.Sprintf("[%s]", opt.Shortcut))
		
		line := fmt.Sprintf("%s%s %s  %s", 
			prefix,
			opt.Icon,
			opt.Title,
			shortcut,
		)
		
		desc := lipgloss.NewStyle().
			Foreground(s.theme.FgMuted).
			Italic(true).
			Render(opt.Description)
		
		menuItems = append(menuItems, style.Render(line))
		if i == s.selectedOption {
			menuItems = append(menuItems, "     "+desc)
		}
	}
	
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(s.theme.Border).
		Padding(1, 2).
		Width(50)
	
	return box.Render(strings.Join(menuItems, "\n"))
}

func (s StartupPage) renderRecentPanel() string {
	var items []string
	
	titleStyle := lipgloss.NewStyle().
		Foreground(s.theme.AccentSecondary).
		Bold(true).
		Padding(0, 2)
	
	items = append(items, titleStyle.Render("📁 Recent Files"))
	items = append(items, "")
	
	if len(s.recentFiles) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(s.theme.FgMuted).
			Italic(true).
			Padding(0, 2)
		items = append(items, emptyStyle.Render("No recent files"))
		items = append(items, "")
		items = append(items, emptyStyle.Render("Start by creating a new"))
		items = append(items, emptyStyle.Render("animation or importing a GIF!"))
	} else {
		for i, rf := range s.recentFiles {
			if i >= 5 {
				break
			}
			
			fileStyle := lipgloss.NewStyle().
				Foreground(s.theme.FgSecondary).
				Width(40).
				Padding(0, 1)
			
			numberStyle := lipgloss.NewStyle().
				Foreground(s.theme.AccentInfo).
				Bold(true)
			
			timeAgo := formatTimeAgo(rf.Timestamp)
			timeStyle := lipgloss.NewStyle().
				Foreground(s.theme.FgMuted)
			
			items = append(items, fileStyle.Render(fmt.Sprintf(
				"%s %s",
				numberStyle.Render(fmt.Sprintf("[%d]", i+1)),
				truncate(rf.Path, 35),
			)))
			
			items = append(items, fileStyle.Render(fmt.Sprintf(
				"    %d frames • %s",
				rf.Frames,
				timeStyle.Render(timeAgo),
			)))
			items = append(items, "")
		}
	}
	
	// Statistics
	items = append(items, "")
	statsStyle := lipgloss.NewStyle().
		Foreground(s.theme.FgMuted).
		Padding(0, 2)
	items = append(items, statsStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━"))
	
	stats := fmt.Sprintf("Theme: %s", s.config.UI.Theme)
	items = append(items, statsStyle.Render(stats))
	
	stats = fmt.Sprintf("Default: %dx%d @ %dfps",
		s.config.Editor.DefaultWidth,
		s.config.Editor.DefaultHeight,
		s.config.Editor.DefaultFPS,
	)
	items = append(items, statsStyle.Render(stats))
	
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(s.theme.Border).
		Padding(1, 2).
		Width(50)
	
	return box.Render(strings.Join(items, "\n"))
}

func (s StartupPage) renderFooter() string {
	tipStyle := lipgloss.NewStyle().
		Foreground(s.theme.FgMuted).
		Italic(true)
	
	tips := []string{
		"💡 Pro tip: Press the number key to quickly open recent files",
		"💡 Pro tip: Use hjkl or arrow keys to navigate",
		"💡 Pro tip: Press 't' to cycle through beautiful themes",
		"💡 Pro tip: Import any GIF with the 'i' key",
	}
	
	// Rotate tips based on time
	tipIdx := int(s.currentTime.Unix()/5) % len(tips)
	tip := tips[tipIdx]
	
	// Breathing effect on tip
	if s.breathing.CurrentAlpha() > 0.95 {
		tipStyle = tipStyle.Foreground(s.theme.AccentInfo)
	}
	
	return lipgloss.NewStyle().
		Width(100).
		Align(lipgloss.Center).
		Render(tipStyle.Render(tip))
}

// Helper functions

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func formatTimeAgo(t time.Time) string {
	duration := time.Since(t)
	
	if duration < time.Minute {
		return "just now"
	} else if duration < time.Hour {
		mins := int(duration.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", mins)
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else {
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}
}
