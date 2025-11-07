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
		
		// Reload artwork frames with proper terminal dimensions if artwork_size is set
		if s.config.Startup.ArtworkFile != "" && s.config.Startup.ArtworkSize != "" && len(s.artworkFrames) == 0 {
			frames, fps := loadAnimatedArtwork(s.config.Startup.ArtworkFile, s.config, msg.Width, msg.Height)
			if len(frames) > 0 {
				s.artworkFrames = frames
				s.artworkFPS = fps
				s.lastFrameTime = time.Now()
			}
		}
	
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
	
	// Center both horizontally and vertically
	return lipgloss.Place(
		s.width,
		s.height,
		lipgloss.Center,
		lipgloss.Center,
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
	
	// Apply breathing effect with gradient transition
	alpha := s.breathing.CurrentAlpha()
	logoColor := s.theme.AccentPrimary
	if s.config.Startup.BreathingEffect && alpha > 0.85 {
		logoColor = s.theme.AccentSecondary
	}
	
	logoStyle := lipgloss.NewStyle().
		Foreground(logoColor).
		Bold(true)
	
	// Render logo with glow effect
	logoRendered := logoStyle.Render(strings.TrimSpace(logo))
	
	// Enhanced border with shadow effect
	if s.config.Startup.ArtworkBorder {
		borderStyle := lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(s.theme.AccentPrimary).
			Padding(1, 3).
			Margin(0, 2)
		logoRendered = borderStyle.Render(logoRendered)
	}
	
	// Add subtle shadow line below
	shadow := lipgloss.NewStyle().
		Foreground(s.theme.BgSecondary).
		Render("▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓")
	
	combined := lipgloss.JoinVertical(lipgloss.Center, logoRendered, shadow)
	
	// Center logo
	return lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(combined)
}


func (s StartupPage) renderMenuPanel() string {
	var menuItems []string
	
	// Animated title bar with pulsing effect
	var titleBar string
	pulseChar := "✦"
	if s.breathing.CurrentAlpha() > 0.9 {
		pulseChar = "✧"
	}
	
	if s.focusArea == "menu" {
		titleStyle := lipgloss.NewStyle().
			Foreground(s.theme.FgBright).
			Background(s.theme.AccentPrimary).
			Bold(true).
			Width(54).
			Align(lipgloss.Center).
			Padding(0, 2)
		titleBar = titleStyle.Render(fmt.Sprintf("%s QUICK START %s", pulseChar, pulseChar))
	} else {
		titleStyle := lipgloss.NewStyle().
			Foreground(s.theme.FgMuted).
			Background(s.theme.BgSecondary).
			Width(54).
			Align(lipgloss.Center).
			Padding(0, 2)
		titleBar = titleStyle.Render("Quick Start")
	}
	
	menuItems = append(menuItems, titleBar)
	menuItems = append(menuItems, "")
	
	// Enhanced menu items with icons and better spacing
	for i, opt := range s.options {
		isSelected := i == s.selectedOption && s.focusArea == "menu"
		
		var line string
		var lineStyle lipgloss.Style
		
		// Enhanced formatting with icon + key + title
		firstLetter := strings.ToUpper(string(opt.Shortcut))
		
		// Animated cursor for selected item
		cursor := "  "
		if isSelected {
			if s.breathing.CurrentAlpha() > 0.85 {
				cursor = "▸ "
			} else {
				cursor = "▹ "
			}
		}
		
		if isSelected {
			// Selected: glassmorphic effect with gradient
			lineStyle = lipgloss.NewStyle().
				Foreground(s.theme.FgBright).
				Background(s.theme.AccentPrimary).
				Bold(true).
				Width(54).
				Padding(0, 2)
			
			iconStyle := lipgloss.NewStyle().
				Foreground(s.theme.AccentSecondary)
			
			keyStyle := lipgloss.NewStyle().
				Foreground(s.theme.FgBright).
				Background(s.theme.AccentSecondary).
				Bold(true)
			
			titleStyle := lipgloss.NewStyle().
				Foreground(s.theme.FgBright).
				Bold(true)
			
			line = fmt.Sprintf("%s%s %s %s",
				cursor,
				iconStyle.Render(opt.Icon),
				keyStyle.Render(fmt.Sprintf(" %s ", firstLetter)),
				titleStyle.Render(opt.Title),
			)
			
			menuItems = append(menuItems, lineStyle.Render(line))
			
			// Add description on separate line for selected item
			descLineStyle := lipgloss.NewStyle().
				Foreground(s.theme.BgPrimary).
				Background(s.theme.AccentPrimary).
				Width(54).
				Padding(0, 2).
				Italic(true)
			menuItems = append(menuItems, descLineStyle.Render("  └─ "+opt.Description))
		} else {
			// Not selected: clean subtle look
			lineStyle = lipgloss.NewStyle().
				Foreground(s.theme.FgSecondary).
				Width(54).
				Padding(0, 2)
			
			iconStyle := lipgloss.NewStyle().
				Foreground(s.theme.FgMuted)
			
			keyStyle := lipgloss.NewStyle().
				Foreground(s.theme.AccentPrimary).
				Bold(true)
			
			line = fmt.Sprintf("%s%s %s %s",
				cursor,
				iconStyle.Render(opt.Icon),
				keyStyle.Render(fmt.Sprintf("[%s]", firstLetter)),
				opt.Title,
			)
			
			menuItems = append(menuItems, lineStyle.Render(line))
		}
	}
	
	menuItems = append(menuItems, "")
	
	// Enhanced border with double-line effect and shadow
	borderStyle := s.theme.Border
	borderType := lipgloss.RoundedBorder()
	
	if s.focusArea == "menu" {
		borderStyle = s.theme.AccentPrimary
		borderType = lipgloss.DoubleBorder()
	}
	
	box := lipgloss.NewStyle().
		Border(borderType).
		BorderForeground(borderStyle).
		Padding(1, 1).
		Width(58).
		Margin(0, 1)
	
	return box.Render(strings.Join(menuItems, "\n"))
}

func (s StartupPage) renderRecentPanel() string {
	var items []string
	
	// Enhanced title bar with animated icon
	var titleBar string
	pulseChar := "★"
	if s.breathing.CurrentAlpha() > 0.9 {
		pulseChar = "☆"
	}
	
	if s.focusArea == "recent" {
		titleStyle := lipgloss.NewStyle().
			Foreground(s.theme.FgBright).
			Background(s.theme.AccentPrimary).
			Bold(true).
			Width(54).
			Align(lipgloss.Center).
			Padding(0, 2)
		titleBar = titleStyle.Render(fmt.Sprintf("%s RECENT FILES %s", pulseChar, pulseChar))
	} else {
		titleStyle := lipgloss.NewStyle().
			Foreground(s.theme.FgMuted).
			Background(s.theme.BgSecondary).
			Width(54).
			Align(lipgloss.Center).
			Padding(0, 2)
		
		tabHintStyle := lipgloss.NewStyle().
			Foreground(s.theme.AccentPrimary).
			Bold(true)
		titleBar = titleStyle.Render("Recent Files") + " " + tabHintStyle.Render("[Tab]")
	}
	
	items = append(items, titleBar)
	items = append(items, "")
	
	if len(s.recentFiles) == 0 {
		emptyBox := lipgloss.NewStyle().
			Foreground(s.theme.FgMuted).
			Background(s.theme.BgSecondary).
			Align(lipgloss.Center).
			Italic(true).
			Padding(2, 4).
			Width(50).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(s.theme.Border)
		
		emptyContent := lipgloss.JoinVertical(
			lipgloss.Center,
			lipgloss.NewStyle().Foreground(s.theme.FgMuted).Bold(true).Render("📭"),
			"",
			"No recent files yet",
			"",
			lipgloss.NewStyle().Foreground(s.theme.AccentSecondary).Render("Press 'n' to create"),
			lipgloss.NewStyle().Foreground(s.theme.AccentSecondary).Render("Press 'i' to import GIF"),
		)
		
		items = append(items, emptyBox.Render(emptyContent))
	} else {
		// Enhanced recent file list with card design
		for i, rf := range s.recentFiles {
			if i >= 8 {
				break
			}
			
			isSelected := s.focusArea == "recent" && i == s.selectedRecent
			
			filename := filepath.Base(rf.Path)
			if len(filename) > 30 {
				filename = filename[:27] + "..."
			}
			
			timeAgo := formatTimeAgo(rf.Timestamp)
			timeAgo = strings.Replace(timeAgo, " ago", "", 1)
			timeAgo = strings.Replace(timeAgo, "minutes", "m", 1)
			timeAgo = strings.Replace(timeAgo, "minute", "m", 1)
			timeAgo = strings.Replace(timeAgo, "hours", "h", 1)
			timeAgo = strings.Replace(timeAgo, "hour", "h", 1)
			timeAgo = strings.Replace(timeAgo, "days", "d", 1)
			timeAgo = strings.Replace(timeAgo, "day", "d", 1)
			
			// Animated cursor for selected item
			cursor := "  "
			if isSelected {
				if s.breathing.CurrentAlpha() > 0.85 {
					cursor = "▸ "
				} else {
					cursor = "▹ "
				}
			}
			
			var line string
			var lineStyle lipgloss.Style
			
			if isSelected {
				// Selected: card-style highlight
				lineStyle = lipgloss.NewStyle().
					Foreground(s.theme.FgBright).
					Background(s.theme.AccentPrimary).
					Bold(true).
					Width(54).
					Padding(0, 2)
				
				numStyle := lipgloss.NewStyle().
					Foreground(s.theme.AccentSecondary).
					Background(s.theme.FgBright).
					Bold(true)
				
				nameStyle := lipgloss.NewStyle().
					Foreground(s.theme.FgBright).
					Bold(true)
				
				metaStyle := lipgloss.NewStyle().
					Foreground(s.theme.BgPrimary)
				
				line = fmt.Sprintf("%s%s %-30s %s",
					cursor,
					numStyle.Render(fmt.Sprintf(" %d ", i+1)),
					nameStyle.Render(filename),
					metaStyle.Render(fmt.Sprintf("%df·%s", rf.Frames, timeAgo)),
				)
			} else {
				// Not selected: subtle card
				lineStyle = lipgloss.NewStyle().
					Foreground(s.theme.FgSecondary).
					Width(54).
					Padding(0, 2)
				
				numStyle := lipgloss.NewStyle().
					Foreground(s.theme.AccentPrimary).
					Bold(true)
				
				metaStyle := lipgloss.NewStyle().
					Foreground(s.theme.FgMuted)
				
				line = fmt.Sprintf("%s%s %-30s %s",
					cursor,
					numStyle.Render(fmt.Sprintf("[%d]", i+1)),
					filename,
					metaStyle.Render(fmt.Sprintf("%df·%s", rf.Frames, timeAgo)),
				)
			}
			
			items = append(items, lineStyle.Render(line))
		}
		
		items = append(items, "")
		
		// Enhanced action hints
		hintStyle := lipgloss.NewStyle().
			Foreground(s.theme.FgMuted).
			Background(s.theme.BgSecondary).
			Align(lipgloss.Center).
			Width(54).
			Padding(0, 2)
		items = append(items, hintStyle.Render("↵ open  │  ⌫ remove  │  1-8 quick select"))
	}
	
	items = append(items, "")
	items = append(items, "")
	
	// Enhanced statistics panel with card design
	divider := lipgloss.NewStyle().
		Foreground(s.theme.AccentPrimary).
		Align(lipgloss.Center).
		Bold(true).
		Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	items = append(items, divider)
	items = append(items, "")
	
	// Stats with icons
	statsStyle := lipgloss.NewStyle().
		Foreground(s.theme.FgSecondary).
		Padding(0, 2).
		Width(54)
	
	iconStyle := lipgloss.NewStyle().
		Foreground(s.theme.AccentPrimary)
	
	valueStyle := lipgloss.NewStyle().
		Foreground(s.theme.AccentSecondary).
		Bold(true)
	
	stats1 := statsStyle.Render(fmt.Sprintf("%s Canvas: %s",
		iconStyle.Render("📐"),
		valueStyle.Render(fmt.Sprintf("%dx%d @ %dfps",
			s.config.Editor.DefaultWidth,
			s.config.Editor.DefaultHeight,
			s.config.Editor.DefaultFPS)),
	))
	items = append(items, stats1)
	
	stats2 := statsStyle.Render(fmt.Sprintf("%s Method: %s",
		iconStyle.Render("🎨"),
		valueStyle.Render(s.config.Converter.DefaultMethod),
	))
	items = append(items, stats2)
	
	items = append(items, "")
	
	// Animated rotating tips with better styling
	tipStyle := lipgloss.NewStyle().
		Foreground(s.theme.FgMuted).
		Background(s.theme.BgSecondary).
		Italic(true).
		Align(lipgloss.Center).
		Width(54).
		Padding(1, 2)
	
	tips := []string{
		"💡 Luminosity method = best quality",
		"🎨 Press 't' to switch themes instantly",
		"🌐 Import GIFs directly from URLs",
		"📏 Use --ratio fit to preserve aspect",
		"⚙️  Edit config.yml for customization",
		"🖼️  Edge method for wireframe style",
		"🧱 Block method for solid appearance",
		"⚡ Number keys quick-open recent files",
	}
	
	tipIdx := int(s.currentTime.Unix()/5) % len(tips)
	
	// Breathing effect on tip
	rotateIcon := "⟳"
	if s.breathing.CurrentAlpha() > 0.95 {
		rotateIcon = "⟲"
	}
	
	tipLine := fmt.Sprintf("%s  %s",
		tips[tipIdx],
		lipgloss.NewStyle().Foreground(s.theme.AccentPrimary).Render(rotateIcon),
	)
	items = append(items, tipStyle.Render(tipLine))
	
	// Enhanced border with glow effect
	borderStyle := s.theme.Border
	borderType := lipgloss.RoundedBorder()
	
	if s.focusArea == "recent" {
		borderStyle = s.theme.AccentPrimary
		borderType = lipgloss.DoubleBorder()
	}
	
	box := lipgloss.NewStyle().
		Border(borderType).
		BorderForeground(borderStyle).
		Padding(1, 1).
		Width(58).
		Margin(0, 1)
	
	return box.Render(strings.Join(items, "\n"))
}

func (s StartupPage) renderFooter() string {
	// Enhanced tagline with gradient-style separators
	taglineStyle := lipgloss.NewStyle().
		Foreground(s.theme.FgBright).
		Bold(true).
		Align(lipgloss.Center)
	
	tagline := taglineStyle.Render("✦ ASCII ART ANIMATION EDITOR ✦")
	
	// Subtitle with better styling
	accentStyle := lipgloss.NewStyle().
		Foreground(s.theme.AccentPrimary).
		Bold(true)
	
	mutedStyle := lipgloss.NewStyle().
		Foreground(s.theme.FgMuted)
	
	subtitle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Render(
			accentStyle.Render("Convert") + 
			mutedStyle.Render(" • ") + 
			accentStyle.Render("Create") + 
			mutedStyle.Render(" • ") + 
			accentStyle.Render("Animate") + 
			mutedStyle.Render("         ") + 
			lipgloss.NewStyle().Foreground(s.theme.FgSecondary).Render("v0.1.0"),
		)
	
	// Enhanced navigation hints with better visual hierarchy
	keyStyle := lipgloss.NewStyle().
		Foreground(s.theme.AccentPrimary).
		Bold(true)
	
	labelStyle := lipgloss.NewStyle().
		Foreground(s.theme.FgSecondary)
	
	divStyle := lipgloss.NewStyle().
		Foreground(s.theme.Border).
		Bold(true)
	
	// Build navigation with styled components
	hints := keyStyle.Render("↑↓") + labelStyle.Render("navigate") +
		divStyle.Render(" ┃ ") +
		keyStyle.Render("Tab") + labelStyle.Render("switch") +
		divStyle.Render(" ┃ ") +
		keyStyle.Render("↵") + labelStyle.Render("select") +
		divStyle.Render(" ┃ ") +
		keyStyle.Render("Esc") + labelStyle.Render("cancel") +
		divStyle.Render(" ┃ ") +
		keyStyle.Render("Q") + labelStyle.Render("quit")
	
	centeredHints := lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(hints)
	
	// Add breathing effect indicator
	breathingIndicator := ""
	if s.config.Startup.BreathingEffect {
		alpha := s.breathing.CurrentAlpha()
		if alpha > 0.9 {
			breathingIndicator = lipgloss.NewStyle().
				Foreground(s.theme.AccentPrimary).
				Render("◉")
		} else {
			breathingIndicator = lipgloss.NewStyle().
				Foreground(s.theme.FgMuted).
				Render("○")
		}
		breathingIndicator = lipgloss.NewStyle().
			Width(s.width).
			Align(lipgloss.Center).
			Render(breathingIndicator)
	}
	
	// Combine all footer elements
	return lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		tagline,
		subtitle,
		"",
		centeredHints,
		breathingIndicator,
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
