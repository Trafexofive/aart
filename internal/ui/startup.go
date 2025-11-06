package ui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mlamkadm/aart/internal/config"
	"github.com/mlamkadm/aart/internal/fileformat"
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
	
	// Animated artwork support
	artworkFrames    []string  // All frames of animated artwork
	currentArtFrame  int       // Current frame index
	artworkFPS       int       // Animation FPS
	lastFrameTime    time.Time // Last frame update time
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
	
	// Load animated artwork frames if .aa file
	artworkFrames := []string{}
	artworkFPS := 12 // Default FPS
	
	if cfg.Startup.ArtworkFile != "" {
		// Pass terminal dimensions for size calculation
		frames, fps := loadAnimatedArtwork(cfg.Startup.ArtworkFile, cfg, 0, 0) // Will be updated on first render
		if len(frames) > 0 {
			artworkFrames = frames
			artworkFPS = fps
		}
	}
	
	return StartupPage{
		theme:           theme,
		styles:          NewStyles(theme),
		config:          cfg,
		selectedOption:  0,
		selectedRecent:  0,
		focusArea:       "menu",
		options:         options,
		recentFiles:     cfg.GetRecentFiles(),
		breathing:       NewBreathingEffect(4 * time.Second),
		currentTime:     time.Now(),
		artworkFrames:   artworkFrames,
		currentArtFrame: 0,
		artworkFPS:      artworkFPS,
		lastFrameTime:   time.Now(),
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
		
		// Animate artwork frames if we have multiple
		if len(s.artworkFrames) > 1 {
			frameDuration := time.Duration(1000/s.artworkFPS) * time.Millisecond
			if time.Since(s.lastFrameTime) >= frameDuration {
				s.currentArtFrame = (s.currentArtFrame + 1) % len(s.artworkFrames)
				s.lastFrameTime = time.Now()
			}
		}
		
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
	
	// Compact header with logo inline
	b.WriteString(s.renderHeader())
	b.WriteString("\n")
	
	// Main panels side by side
	leftPanel := s.renderMenuPanel()
	rightPanel := s.renderRecentPanel()
	
	panels := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		"  ",
		rightPanel,
	)
	
	b.WriteString(panels)
	b.WriteString("\n")
	
	// Footer with complete navigation hints
	b.WriteString(s.renderFooter())
	
	// Use top alignment instead of center for better 80x24 compatibility
	// Only center horizontally
	return lipgloss.Place(
		s.width,
		s.height,
		lipgloss.Center,
		lipgloss.Top,
		b.String(),
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(s.theme.BgPrimary),
	)
}

func (s StartupPage) renderHeader() string {
	// Get custom or default ASCII art logo
	var logo string
	
	// Use animated frame if available
	if len(s.artworkFrames) > 0 {
		logo = s.artworkFrames[s.currentArtFrame]
	} else {
		logo = s.config.GetStartupArtwork()
	}
	
	// Apply breathing effect to logo color
	logoColor := s.theme.AccentPrimary
	if s.config.Startup.BreathingEffect && s.breathing.CurrentAlpha() > 0.9 {
		logoColor = s.theme.AccentSecondary
	}
	
	logoStyle := lipgloss.NewStyle().
		Foreground(logoColor)
	
	// Render logo
	logoRendered := logoStyle.Render(strings.TrimSpace(logo))
	
	// Apply border if configured
	if s.config.Startup.ArtworkBorder {
		borderStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(s.theme.Border).
			Padding(1, 2)
		logoRendered = borderStyle.Render(logoRendered)
	}
	
	// Center logo
	width := 104
	return lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(logoRendered)
}


func (s StartupPage) renderMenuPanel() string {
	var menuItems []string
	
	// Title with visual active indicator
	var titleBar string
	if s.focusArea == "menu" {
		// Active panel - inverse title bar
		titleStyle := lipgloss.NewStyle().
			Foreground(s.theme.BgPrimary).
			Background(s.theme.AccentPrimary).
			Bold(true).
			Width(46).
			Padding(0, 1)
		titleBar = titleStyle.Render("▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓")
	} else {
		// Inactive panel
		titleStyle := lipgloss.NewStyle().
			Foreground(s.theme.FgSecondary).
			Width(46).
			Padding(0, 1)
		titleBar = titleStyle.Render("Quick Start")
	}
	
	menuItems = append(menuItems, titleBar)
	menuItems = append(menuItems, "")
	
	for i, opt := range s.options {
		isSelected := i == s.selectedOption && s.focusArea == "menu"
		
		// Build line without emoji - format: "▸ [N]ew Animation"
		var line string
		var lineStyle lipgloss.Style
		
		// Highlight first letter of title
		firstLetter := strings.ToUpper(string(opt.Shortcut))
		restOfTitle := opt.Title[1:]
		
		if isSelected {
			// Selected: inverse video with arrow
			lineStyle = lipgloss.NewStyle().
				Foreground(s.theme.BgPrimary).
				Background(s.theme.AccentPrimary).
				Bold(true).
				Width(46).
				Padding(0, 1)
			
			keyStyle := lipgloss.NewStyle().
				Foreground(s.theme.BgPrimary).
				Background(s.theme.AccentSecondary).
				Bold(true)
			
			line = fmt.Sprintf("▸ %s%s",
				keyStyle.Render(fmt.Sprintf("[%s]", firstLetter)),
				restOfTitle,
			)
		} else {
			// Not selected: muted
			lineStyle = lipgloss.NewStyle().
				Foreground(s.theme.FgSecondary).
				Width(46).
				Padding(0, 1)
			
			keyStyle := lipgloss.NewStyle().
				Foreground(s.theme.AccentSecondary)
			
			line = fmt.Sprintf("  %s%s",
				keyStyle.Render(fmt.Sprintf("[%s]", firstLetter)),
				restOfTitle,
			)
		}
		
		menuItems = append(menuItems, lineStyle.Render(line))
	}
	
	menuItems = append(menuItems, "")
	
	// Border style based on focus
	borderStyle := s.theme.Border
	borderType := lipgloss.RoundedBorder()
	
	if s.focusArea == "menu" {
		borderStyle = s.theme.AccentPrimary
		borderType = lipgloss.ThickBorder()
	}
	
	box := lipgloss.NewStyle().
		Border(borderType).
		BorderForeground(borderStyle).
		Padding(0, 1).
		Width(50)
	
	return box.Render(strings.Join(menuItems, "\n"))
}

func (s StartupPage) renderRecentPanel() string {
	var items []string
	
	// Title with Tab hint
	var titleBar string
	if s.focusArea == "recent" {
		// Active panel - inverse title bar
		titleStyle := lipgloss.NewStyle().
			Foreground(s.theme.BgPrimary).
			Background(s.theme.AccentPrimary).
			Bold(true).
			Width(46).
			Padding(0, 1)
		titleBar = titleStyle.Render("▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓")
	} else {
		titleStyle := lipgloss.NewStyle().
			Foreground(s.theme.FgSecondary).
			Width(46).
			Padding(0, 1)
		tabHintStyle := lipgloss.NewStyle().
			Foreground(s.theme.FgMuted)
		titleBar = titleStyle.Render("Recent Files") + tabHintStyle.Render("  [Tab]")
	}
	
	items = append(items, titleBar)
	items = append(items, "")
	
	if len(s.recentFiles) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(s.theme.FgMuted).
			Italic(true).
			Padding(0, 1)
		items = append(items, emptyStyle.Render("No recent files"))
		items = append(items, "")
		items = append(items, emptyStyle.Render("Press 'n' to create new animation"))
		items = append(items, emptyStyle.Render("Press 'i' to import a GIF"))
	} else {
		// Compact single-line format
		for i, rf := range s.recentFiles {
			if i >= 8 { // Show more files now that we're compact
				break
			}
			
			isSelected := s.focusArea == "recent" && i == s.selectedRecent
			
			// Format: "1. filename.aart         150f • 13m"
			filename := filepath.Base(rf.Path)
			if len(filename) > 22 {
				filename = filename[:19] + "..."
			}
			
			timeAgo := formatTimeAgo(rf.Timestamp)
			// Compact time format
			timeAgo = strings.Replace(timeAgo, " ago", "", 1)
			timeAgo = strings.Replace(timeAgo, "minutes", "m", 1)
			timeAgo = strings.Replace(timeAgo, "minute", "m", 1)
			timeAgo = strings.Replace(timeAgo, "hours", "h", 1)
			timeAgo = strings.Replace(timeAgo, "hour", "h", 1)
			timeAgo = strings.Replace(timeAgo, "days", "d", 1)
			timeAgo = strings.Replace(timeAgo, "day", "d", 1)
			
			var line string
			var lineStyle lipgloss.Style
			
			if isSelected {
				lineStyle = lipgloss.NewStyle().
					Foreground(s.theme.BgPrimary).
					Background(s.theme.AccentPrimary).
					Bold(true).
					Width(46).
					Padding(0, 1)
				
				line = fmt.Sprintf("▸ %d. %-22s %3df • %s",
					i+1,
					filename,
					rf.Frames,
					timeAgo,
				)
				items = append(items, lineStyle.Render(line))
			} else {
				// Build line without embedded styles to avoid escape sequence issues
				line = fmt.Sprintf("  %d. %-22s %3df • %s",
					i+1,
					filename,
					rf.Frames,
					timeAgo,
				)
				
				lineStyle = lipgloss.NewStyle().
					Foreground(s.theme.FgSecondary).
					Width(46).
					Padding(0, 1)
				
				items = append(items, lineStyle.Render(line))
			}
		}
		
		items = append(items, "")
		
		// Action hints
		hintStyle := lipgloss.NewStyle().
			Foreground(s.theme.FgMuted).
			Padding(0, 1)
		items = append(items, hintStyle.Render("Enter: open │ Del: remove"))
	}
	
	items = append(items, "")
	
	// Compact statistics
	statsStyle := lipgloss.NewStyle().
		Foreground(s.theme.FgMuted).
		Padding(0, 1)
	
	divider := lipgloss.NewStyle().
		Foreground(s.theme.Border).
		Render("─────────────────────────────────────────────")
	items = append(items, divider)
	
	stats := fmt.Sprintf("Default: %dx%d @ %dfps",
		s.config.Editor.DefaultWidth,
		s.config.Editor.DefaultHeight,
		s.config.Editor.DefaultFPS,
	)
	items = append(items, statsStyle.Render(stats))
	
	methodStats := fmt.Sprintf("Method: %s", s.config.Converter.DefaultMethod)
	items = append(items, statsStyle.Render(methodStats))
	
	items = append(items, "")
	
	// Tip with counter
	tipStyle := lipgloss.NewStyle().
		Foreground(s.theme.FgMuted).
		Italic(true).
		Padding(0, 1)
	
	tips := []string{
		"Use luminosity method for best quality",
		"Press 't' to change theme instantly",
		"Import GIFs from URLs with 'i' key",
		"Use --ratio fit to preserve aspect",
		"Edit config.yml for custom settings",
		"Try edge method for wireframe style",
		"Block method gives solid appearance",
		"Number keys quick-open recent files",
	}
	
	tipIdx := int(s.currentTime.Unix()/5) % len(tips)
	
	// Breathing effect on tip rotation indicator
	rotateIcon := "⟳"
	if s.breathing.CurrentAlpha() > 0.95 {
		rotateIcon = "⟲"
	}
	
	tipLine := fmt.Sprintf("Tip %d/%d: %s %s",
		tipIdx+1,
		len(tips),
		tips[tipIdx],
		rotateIcon,
	)
	items = append(items, tipStyle.Render(tipLine))
	
	// Border style based on focus
	borderStyle := s.theme.Border
	borderType := lipgloss.RoundedBorder()
	
	if s.focusArea == "recent" {
		borderStyle = s.theme.AccentPrimary
		borderType = lipgloss.ThickBorder()
	}
	
	box := lipgloss.NewStyle().
		Border(borderType).
		BorderForeground(borderStyle).
		Padding(0, 1).
		Width(50)
	
	return box.Render(strings.Join(items, "\n"))
}

func (s StartupPage) renderFooter() string {
	// Tagline and version first
	taglineStyle := lipgloss.NewStyle().
		Foreground(s.theme.FgSecondary).
		Align(lipgloss.Center)
	
	tagline := taglineStyle.Render("ASCII Art Animation Editor")
	subtitle := lipgloss.NewStyle().
		Foreground(s.theme.FgMuted).
		Align(lipgloss.Center).
		Render("Convert • Create • Animate         v0.1.0")
	
	// Then navigation hints below
	hintStyle := lipgloss.NewStyle().
		Foreground(s.theme.FgSecondary).
		Bold(true)
	
	divStyle := lipgloss.NewStyle().
		Foreground(s.theme.FgMuted)
	
	// Complete navigation hints
	hints := hintStyle.Render("hjkl:navigate") +
		divStyle.Render(" │ ") +
		hintStyle.Render("Tab:switch") +
		divStyle.Render(" │ ") +
		hintStyle.Render("Enter:select") +
		divStyle.Render(" │ ") +
		hintStyle.Render("Esc:cancel") +
		divStyle.Render(" │ ") +
		hintStyle.Render("q:quit")
	
	centeredHints := lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(hints)
	
	// Combine tagline, subtitle, and hints
	return lipgloss.JoinVertical(
		lipgloss.Center,
		tagline,
		subtitle,
		centeredHints,
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

// loadAnimatedArtwork loads all frames from a .aa file for animated startup logo
func loadAnimatedArtwork(path string, cfg *config.Config, termWidth, termHeight int) ([]string, int) {
	// Try absolute path first
	aartFile, err := fileformat.Load(path)
	if err != nil {
		// Try relative to config directory
		dir, err := config.ConfigDir()
		if err == nil {
			artPath := filepath.Join(dir, path)
			aartFile, err = fileformat.Load(artPath)
		}
		if err != nil {
			// Not a .aa file, return empty
			return []string{}, 12
		}
	}
	
	// Calculate max dimensions based on artwork_size
	maxWidth := 80
	maxHeight := 20
	
	if cfg.Startup.ArtworkSize != "" && termHeight > 0 {
		// Parse percentage (e.g. "40p" = 40% of screen)
		sizeStr := strings.TrimSuffix(cfg.Startup.ArtworkSize, "p")
		if percent, err := strconv.ParseInt(sizeStr, 10, 64); err == nil && percent > 0 && percent <= 100 {
			// Calculate height as percentage of terminal
			maxHeight = int(float64(termHeight) * float64(percent) / 100.0)
			if maxHeight < 10 {
				maxHeight = 10
			}
			if maxHeight > 40 {
				maxHeight = 40
			}
			
			// Calculate width proportionally (keep it reasonable)
			maxWidth = maxHeight * 3 // Roughly 3:1 ratio for ASCII art
			if maxWidth > termWidth-10 {
				maxWidth = termWidth - 10
			}
		}
	}
	
	// Account for border padding
	if cfg.Startup.ArtworkBorder {
		maxWidth -= 6  // Border + padding
		maxHeight -= 4 // Border + padding
	}
	
	// Override with explicit width/height if artwork_size not set
	if cfg.Startup.ArtworkSize == "" {
		if cfg.Startup.ArtworkWidth > 0 {
			maxWidth = cfg.Startup.ArtworkWidth
		}
		if cfg.Startup.ArtworkHeight > 0 {
			maxHeight = cfg.Startup.ArtworkHeight
		}
	}
	
	// Extract ASCII art from each frame (TEXT ONLY, no colors)
	frames := make([]string, len(aartFile.Frames))
	
	for i, frame := range aartFile.Frames {
		var lines []string
		
		// Limit to max height
		frameCells := frame.Cells
		if len(frameCells) > maxHeight {
			// Take center portion
			startY := (len(frameCells) - maxHeight) / 2
			frameCells = frameCells[startY : startY+maxHeight]
		}
		
		for _, row := range frameCells {
			var line strings.Builder
			charCount := 0
			
			for _, cell := range row {
				if charCount >= maxWidth {
					break
				}
				
				// Extract ONLY the character, ignore colors
				if cell.Char == "" || cell.Char == " " {
					line.WriteString(" ")
				} else {
					line.WriteString(cell.Char)
				}
				charCount++
			}
			
			// Trim trailing spaces
			lineStr := strings.TrimRight(line.String(), " ")
			lines = append(lines, lineStr)
		}
		
		// Remove trailing empty lines
		for len(lines) > 0 && lines[len(lines)-1] == "" {
			lines = lines[:len(lines)-1]
		}
		
		// Remove leading empty lines
		for len(lines) > 0 && lines[0] == "" {
			lines = lines[1:]
		}
		
		frames[i] = strings.Join(lines, "\n")
	}
	
	// Calculate FPS from first frame duration
	fps := 12
	if len(aartFile.Frames) > 0 && aartFile.Frames[0].Duration > 0 {
		fps = 1000 / aartFile.Frames[0].Duration
		if fps < 1 {
			fps = 1
		}
		if fps > 60 {
			fps = 60
		}
	}
	
	return frames, fps
}
