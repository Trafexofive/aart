package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mlamkadm/aart/internal/config"
)

type editorFinishedMsg struct{}

type loadFileMsg struct {
	path string
}

func LoadFileCmd(path string) tea.Cmd {
	return func() tea.Msg {
		return loadFileMsg{path: path}
	}
}

// StartupPage is the initial welcome screen
type StartupPage struct {
	width            int
	height           int
	theme            Theme
	styles           Styles
	config           *config.Config
	selectedOption   int
	selectedRecent   int
	focusArea        string // "menu" or "recent"
	options          []StartupOption
	recentFiles      []config.RecentFile
	breathing        *BreathingEffect
	currentTime      time.Time
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
			Icon:        "📝",
			Title:       "Edit Config",
			Description: "Edit config.yml with $EDITOR",
			Action:      "editconfig",
			Shortcut:    "c",
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
		selectedRecent: 0,
		focusArea:      "menu",
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
		
		case "tab":
			// Switch focus between menu and recent files
			if s.focusArea == "menu" {
				s.focusArea = "recent"
			} else {
				s.focusArea = "menu"
			}
		
		case "up", "k":
			if s.focusArea == "menu" {
				s.selectedOption = (s.selectedOption - 1 + len(s.options)) % len(s.options)
			} else if len(s.recentFiles) > 0 {
				s.selectedRecent = (s.selectedRecent - 1 + len(s.recentFiles)) % len(s.recentFiles)
			}
		
		case "down", "j":
			if s.focusArea == "menu" {
				s.selectedOption = (s.selectedOption + 1) % len(s.options)
			} else if len(s.recentFiles) > 0 {
				s.selectedRecent = (s.selectedRecent + 1) % len(s.recentFiles)
			}
		
		case "enter", " ":
			if s.focusArea == "recent" && len(s.recentFiles) > 0 {
				// Open selected recent file
				return s.openRecentFile(s.selectedRecent)
			}
			return s.handleAction()
		
		case "n":
			s.selectedOption = 0
			s.focusArea = "menu"
			return s.handleAction()
		case "o":
			s.selectedOption = 1
			s.focusArea = "menu"
			return s.handleAction()
		case "i":
			s.selectedOption = 2
			s.focusArea = "menu"
			return s.handleAction()
		case "t":
			s.selectedOption = 3
			s.focusArea = "menu"
			return s.handleAction()
		case "s":
			s.selectedOption = 4
			s.focusArea = "menu"
			return s.handleAction()
		case "c":
			s.selectedOption = 5
			s.focusArea = "menu"
			return s.handleAction()
		case "e":
			s.selectedOption = 6
			s.focusArea = "menu"
			return s.handleAction()
		case "?":
			s.selectedOption = 7
			s.focusArea = "menu"
			return s.handleAction()
		case "1", "2", "3", "4", "5":
			// Quick access to recent files
			idx := int(msg.String()[0] - '1')
			if idx < len(s.recentFiles) {
				return s.openRecentFile(idx)
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
		// Create new animation - transition to editor with proper initialization
		model := NewWithConfig(s.config)
		return model, model.Init()
	case "open":
		// Open file picker
		picker := NewFilePicker(s.config, "📂 Open File", ".aart", s)
		return picker, picker.Init()
	case "import":
		// Show GIF import dialog
		importer := NewImportGIFScreen(s.config, s)
		return importer, importer.Init()
	case "quit":
		return s, tea.Quit
	case "theme":
		// Show theme selector instead of cycling
		return NewThemeSelector(s.config, s), nil
	case "help":
		// Show help screen
		help := NewHelpScreen(s.config)
		return help, help.Init()
	case "settings":
		// Show settings screen
		settings := NewSettingsScreen(s.config, s)
		return settings, settings.Init()
	case "editconfig":
		// Edit config with $EDITOR
		return s, s.editConfigCmd()
	case "examples":
		// Show examples gallery
		return NewExamplesScreen(s.config, s), nil
	default:
		return s, nil
	}
}

func (s StartupPage) openRecentFile(idx int) (tea.Model, tea.Cmd) {
	if idx < 0 || idx >= len(s.recentFiles) {
		return s, nil
	}
	
	rf := s.recentFiles[idx]
	// Load the file and transition to editor
	model := NewWithConfig(s.config)
	return model, LoadFileCmd(rf.Path)
}

func (s StartupPage) editConfigCmd() tea.Cmd {
	return tea.ExecProcess(exec.Command(
		getEditor(),
		getConfigFilePath(),
	), func(err error) tea.Msg {
		if err != nil {
			return editorFinishedMsg{}
		}
		// Reload config after editing
		newCfg, _ := config.Load()
		if newCfg != nil {
			*s.config = *newCfg
		}
		return editorFinishedMsg{}
	})
}

func getConfigFilePath() string {
	path, err := config.ConfigPath()
	if err != nil {
		return ""
	}
	return path
}

func getEditor() string {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor
	}
	// Default fallbacks
	for _, e := range []string{"vim", "nano", "vi"} {
		if _, err := exec.LookPath(e); err == nil {
			return e
		}
	}
	return "vi"
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
	
	if s.focusArea == "menu" {
		titleStyle = titleStyle.Underline(true)
	}
	
	menuItems = append(menuItems, titleStyle.Render("✨ Quick Start"))
	menuItems = append(menuItems, "")
	
	for i, opt := range s.options {
		var style lipgloss.Style
		prefix := "  "
		
		if i == s.selectedOption && s.focusArea == "menu" {
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
		if i == s.selectedOption && s.focusArea == "menu" {
			menuItems = append(menuItems, "     "+desc)
		}
	}
	
	borderStyle := s.theme.Border
	if s.focusArea == "menu" {
		borderStyle = s.theme.BorderActive
	}
	
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderStyle).
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
	
	if s.focusArea == "recent" {
		titleStyle = titleStyle.Foreground(s.theme.AccentPrimary).Underline(true)
	}
	
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
			
			isSelected := s.focusArea == "recent" && i == s.selectedRecent
			
			var fileStyle lipgloss.Style
			prefix := "  "
			
			if isSelected {
				fileStyle = lipgloss.NewStyle().
					Background(s.theme.Selection).
					Foreground(s.theme.AccentPrimary).
					Bold(true).
					Width(45).
					Padding(0, 1)
				prefix = "▶ "
			} else {
				fileStyle = lipgloss.NewStyle().
					Foreground(s.theme.FgSecondary).
					Width(45).
					Padding(0, 1)
			}
			
			numberStyle := lipgloss.NewStyle().
				Foreground(s.theme.AccentInfo).
				Bold(true)
			
			timeAgo := formatTimeAgo(rf.Timestamp)
			timeStyle := lipgloss.NewStyle().
				Foreground(s.theme.FgMuted)
			
			items = append(items, fileStyle.Render(fmt.Sprintf(
				"%s%s %s",
				prefix,
				numberStyle.Render(fmt.Sprintf("[%d]", i+1)),
				truncate(rf.Path, 35),
			)))
			
			detailStyle := fileStyle
			if !isSelected {
				detailStyle = lipgloss.NewStyle().
					Foreground(s.theme.FgMuted).
					Width(45).
					Padding(0, 1)
			}
			
			items = append(items, detailStyle.Render(fmt.Sprintf(
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
	
	borderStyle := s.theme.Border
	if s.focusArea == "recent" {
		borderStyle = s.theme.BorderActive
	}
	
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderStyle).
		Padding(1, 2).
		Width(50)
	
	return box.Render(strings.Join(items, "\n"))
}

func (s StartupPage) renderFooter() string {
	tipStyle := lipgloss.NewStyle().
		Foreground(s.theme.FgMuted).
		Italic(true)
	
	tips := []string{
		"💡 Pro tip: Press Tab to switch between menu and recent files",
		"💡 Pro tip: Press the number key to quickly open recent files",
		"💡 Pro tip: Use hjkl or arrow keys to navigate",
		"💡 Pro tip: Press 't' to change themes",
		"💡 Pro tip: Import any GIF with the 'i' key",
		"💡 Pro tip: Press 'c' to edit config.yml with $EDITOR",
	}
	
	// Rotate tips based on time
	tipIdx := int(s.currentTime.Unix()/5) % len(tips)
	tip := tips[tipIdx]
	
	// Breathing effect on tip
	if s.breathing.CurrentAlpha() > 0.95 {
		tipStyle = tipStyle.Foreground(s.theme.AccentInfo)
	}
	
	hintStyle := lipgloss.NewStyle().
		Foreground(s.theme.FgSecondary).
		Bold(true)
	
	hints := "hjkl/↑↓: navigate │ Tab: switch panel │ Enter: select │ q: quit"
	
	return lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.NewStyle().Width(100).Align(lipgloss.Center).Render(hintStyle.Render(hints)),
		lipgloss.NewStyle().Width(100).Align(lipgloss.Center).Render(tipStyle.Render(tip)),
	)
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
