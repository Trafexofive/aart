package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mlamkadm/aart/internal/config"
)

// FilePickerScreen allows browsing and selecting files
type FilePickerScreen struct {
	width         int
	height        int
	theme         Theme
	styles        Styles
	config        *config.Config
	currentDir    string
	files         []os.FileInfo
	selectedIndex int
	showHidden    bool
	filter        string // ".aart" or "" for all
	title         string
	returnTo      tea.Model
}

// NewFilePicker creates a file browser
func NewFilePicker(cfg *config.Config, title, filter string, returnTo tea.Model) FilePickerScreen {
	themeName := cfg.UI.Theme
	if themeName == "" {
		themeName = "tokyo-night"
	}
	theme := GetTheme(themeName)
	
	// Start in home directory
	homeDir, _ := os.UserHomeDir()
	
	fp := FilePickerScreen{
		theme:      theme,
		styles:     NewStyles(theme),
		config:     cfg,
		currentDir: homeDir,
		filter:     filter,
		title:      title,
		returnTo:   returnTo,
	}
	
	fp.loadDirectory()
	return fp
}

func (f *FilePickerScreen) loadDirectory() {
	entries, err := os.ReadDir(f.currentDir)
	if err != nil {
		return
	}
	
	f.files = make([]os.FileInfo, 0)
	
	// Add parent directory option
	if f.currentDir != "/" {
		info, _ := os.Stat(filepath.Dir(f.currentDir))
		if info != nil {
			f.files = append(f.files, info)
		}
	}
	
	// Add files and directories
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		
		// Skip hidden files unless enabled
		if !f.showHidden && strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		
		// Apply filter
		if f.filter != "" && !info.IsDir() {
			if !strings.HasSuffix(entry.Name(), f.filter) {
				continue
			}
		}
		
		f.files = append(f.files, info)
	}
	
	f.selectedIndex = 0
}

func (f FilePickerScreen) Init() tea.Cmd {
	return nil
}

func (f FilePickerScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			// Return to previous screen
			if f.returnTo != nil {
				return f.returnTo, nil
			}
			return NewStartupPage(f.config), nil
		
		case "up", "k":
			if f.selectedIndex > 0 {
				f.selectedIndex--
			}
		
		case "down", "j":
			if f.selectedIndex < len(f.files)-1 {
				f.selectedIndex++
			}
		
		case "enter", " ":
			return f.handleSelect()
		
		case ".":
			f.showHidden = !f.showHidden
			f.loadDirectory()
		
		case "h":
			// Go to home directory
			homeDir, _ := os.UserHomeDir()
			f.currentDir = homeDir
			f.loadDirectory()
		}
	
	case tea.WindowSizeMsg:
		f.width = msg.Width
		f.height = msg.Height
	}
	
	return f, nil
}

func (f FilePickerScreen) handleSelect() (tea.Model, tea.Cmd) {
	if f.selectedIndex >= len(f.files) {
		return f, nil
	}
	
	selected := f.files[f.selectedIndex]
	
	if selected.IsDir() {
		// Navigate into directory
		if selected.Name() == ".." {
			f.currentDir = filepath.Dir(f.currentDir)
		} else {
			f.currentDir = filepath.Join(f.currentDir, selected.Name())
		}
		f.loadDirectory()
		return f, nil
	}
	
	// File selected - load it
	fullPath := filepath.Join(f.currentDir, selected.Name())
	
	// Load the file and open editor
	frames, err := LoadFile(fullPath)
	if err != nil {
		// TODO: Show error message
		return f, nil
	}
	
	// Add to recent files
	f.config.AddRecentFile(fullPath, len(frames))
	config.Save(f.config)
	
	// Open editor with loaded frames
	return newModelWithConfigAndFrames(frames, fullPath, f.config), nil
}

func (f FilePickerScreen) View() string {
	if f.width == 0 {
		return "Loading..."
	}
	
	var b strings.Builder
	
	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(f.theme.AccentPrimary).
		Bold(true).
		Align(lipgloss.Center).
		Width(80)
	
	b.WriteString(titleStyle.Render(f.title))
	b.WriteString("\n\n")
	
	// Current directory
	pathStyle := lipgloss.NewStyle().
		Foreground(f.theme.AccentSecondary).
		Bold(true)
	
	b.WriteString(pathStyle.Render(fmt.Sprintf("📁 %s", f.currentDir)))
	b.WriteString("\n\n")
	
	// File list
	for i, file := range f.files {
		var icon string
		var style lipgloss.Style
		
		if file.IsDir() {
			icon = "📁"
			style = lipgloss.NewStyle().Foreground(f.theme.AccentInfo)
		} else {
			icon = "📄"
			style = lipgloss.NewStyle().Foreground(f.theme.FgSecondary)
		}
		
		prefix := "  "
		if i == f.selectedIndex {
			style = lipgloss.NewStyle().
				Background(f.theme.Selection).
				Foreground(f.theme.AccentPrimary).
				Bold(true).
				Width(60)
			prefix = "▶ "
		}
		
		name := file.Name()
		if file.IsDir() && name != ".." {
			name += "/"
		}
		
		line := fmt.Sprintf("%s%s %s", prefix, icon, name)
		b.WriteString(style.Render(line))
		b.WriteString("\n")
		
		// Only show first 20 files
		if i >= 19 {
			moreStyle := lipgloss.NewStyle().
				Foreground(f.theme.FgMuted).
				Italic(true)
			b.WriteString(moreStyle.Render(fmt.Sprintf("  ... and %d more", len(f.files)-20)))
			break
		}
	}
	
	// Instructions
	b.WriteString("\n\n")
	hintStyle := lipgloss.NewStyle().
		Foreground(f.theme.FgMuted).
		Italic(true)
	
	hints := []string{
		"hjkl/arrows: navigate",
		"enter: select",
		"h: home",
		".: toggle hidden",
		"q/esc: back",
	}
	b.WriteString(hintStyle.Render(strings.Join(hints, " │ ")))
	
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(f.theme.BorderActive).
		Padding(2, 4).
		Width(70)
	
	return lipgloss.Place(
		f.width,
		f.height,
		lipgloss.Center,
		lipgloss.Center,
		box.Render(b.String()),
	)
}
