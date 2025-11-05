package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Beautiful themed render functions

func (m Model) renderStatusBar() string {
	// Beautiful gradient title
	title := m.styles.StatusBarTitle.Render("✨ aart")
	
	// File info with icon
	fileIcon := "📄"
	fileStyle := m.styles.StatusBarSection
	if m.modified {
		fileIcon = "✏️"
		fileStyle = m.styles.Warning
	}
	fileInfo := fileStyle.Render(fmt.Sprintf("%s %s", fileIcon, m.filename))
	
	// Canvas size with breathing effect
	sizeStyle := m.styles.StatusBarSection
	alpha := m.breathing.CurrentAlpha()
	if alpha > 0.9 {
		sizeStyle = m.styles.Info
	}
	canvasInfo := sizeStyle.Render(fmt.Sprintf("📐 %dx%d", 
		m.frames[m.currentFrame].Width, 
		m.frames[m.currentFrame].Height))
	
	// Frame indicator with color
	frameStyle := lipgloss.NewStyle().
		Foreground(m.theme.AccentPrimary).
		Bold(true)
	frameInfo := frameStyle.Render(fmt.Sprintf("🎬 %d/%d", 
		m.currentFrame+1, len(m.frames)))
	
	// FPS with playback icon
	fpsIcon := "⏸"
	fpsStyle := m.styles.Muted
	if m.playing {
		fpsIcon = "▶"
		fpsStyle = lipgloss.NewStyle().
			Foreground(m.theme.AccentSuccess).
			Bold(true)
	}
	fpsInfo := fpsStyle.Render(fmt.Sprintf("%s %dfps", fpsIcon, m.fps))
	
	// Tool info with icon
	toolIcons := map[Tool]string{
		ToolPencil:     "✏️",
		ToolFill:       "🪣",
		ToolSelect:     "⬚",
		ToolLine:       "╱",
		ToolBox:        "▢",
		ToolText:       "𝑻",
		ToolEyedropper: "💧",
		ToolMove:       "✥",
	}
	toolIcon := toolIcons[m.selectedTool]
	toolName := toolNames[m.selectedTool]
	toolStyle := lipgloss.NewStyle().
		Foreground(m.theme.AccentSecondary)
	toolInfo := toolStyle.Render(fmt.Sprintf("%s %s", toolIcon, toolName))
	
	// Colors with visual preview
	fgPreview := lipgloss.NewStyle().
		Foreground(m.theme.AccentPrimary).
		Bold(true).
		Render(string(m.fgChar))
	bgPreview := lipgloss.NewStyle().
		Foreground(m.theme.FgMuted).
		Render("░")
	colorInfo := m.styles.StatusBarSection.Render(fmt.Sprintf("fg:%s bg:%s", fgPreview, bgPreview))
	
	// Layer info
	layerStyle := lipgloss.NewStyle().
		Foreground(m.theme.FgSecondary)
	layerInfo := layerStyle.Render(fmt.Sprintf("📑 %d/%d", 
		m.currentLayer+1, len(m.layers)))
	
	// Divider with theme color
	div := lipgloss.NewStyle().
		Foreground(m.theme.Border).
		Render(" │ ")
	
	// Assemble with spacing
	sections := []string{
		title,
		fileInfo,
		canvasInfo,
		frameInfo,
		fpsInfo,
		toolInfo,
		colorInfo,
		layerInfo,
	}
	
	content := strings.Join(sections, div)
	
	return m.styles.StatusBar.Render(content)
}

func (m Model) renderTimeline() string {
	var b strings.Builder
	
	// Top border with title - more zen, breathing
	titleStyle := lipgloss.NewStyle().
		Foreground(m.theme.AccentPrimary).
		Bold(true)
	
	// Breathing effect on title
	alpha := m.breathing.CurrentAlpha()
	spacer := " "
	if alpha > 0.9 {
		spacer = " ✦ "
	}
	title := titleStyle.Render(spacer + "TIMELINE" + spacer)
	
	borderStyle := lipgloss.NewStyle().
		Foreground(m.theme.Border)
	
	lineWidth := 110
	titleLen := lipgloss.Width(title)
	leftPad := (lineWidth - titleLen) / 2
	rightPad := lineWidth - titleLen - leftPad
	
	b.WriteString(borderStyle.Render("├" + strings.Repeat("─", leftPad)))
	b.WriteString(title)
	b.WriteString(borderStyle.Render(strings.Repeat("─", rightPad) + "┤\n│ "))
	
	// Frame indicators with beautiful box drawing
	maxFrames := 30
	displayFrames := len(m.frames)
	if displayFrames > maxFrames {
		displayFrames = maxFrames
	}
	
	for i := 0; i < displayFrames; i++ {
		var frameStyle lipgloss.Style
		frameText := fmt.Sprintf("%2d", i+1)
		
		if i == m.currentFrame {
			// Current frame - bold block
			frameStyle = lipgloss.NewStyle().
				Background(m.theme.PlayheadColor).
				Foreground(m.theme.BgPrimary).
				Bold(true)
			frameText = "▓▓"
		} else if i < m.currentFrame {
			// Past frames - subtle
			frameStyle = lipgloss.NewStyle().
				Foreground(m.theme.TimelineInactive)
		} else {
			// Future frames - normal brightness
			frameStyle = lipgloss.NewStyle().
				Foreground(m.theme.TimelineActive)
		}
		
		b.WriteString(frameStyle.Render(frameText))
		b.WriteString(" ")
	}
	
	if len(m.frames) > maxFrames {
		moreStyle := lipgloss.NewStyle().
			Foreground(m.theme.FgMuted)
		b.WriteString(moreStyle.Render(fmt.Sprintf("…+%d", len(m.frames)-maxFrames)))
	}
	
	// Pad to full width
	currentWidth := displayFrames*3 + 5
	if len(m.frames) > maxFrames {
		currentWidth += 5
	}
	b.WriteString(strings.Repeat(" ", max(0, lineWidth-currentWidth)))
	
	b.WriteString(borderStyle.Render("\n│ "))
	
	// Playback status line - zen layout
	divStyle := lipgloss.NewStyle().
		Foreground(m.theme.Border)
	div := divStyle.Render(" │ ")
	
	// Status icon with animation
	statusIcon := "⏹"
	statusText := "stopped"
	statusStyle := lipgloss.NewStyle().
		Foreground(m.theme.FgMuted)
	
	if m.playing {
		statusIcon = "▶"
		statusText = "playing"
		statusStyle = lipgloss.NewStyle().
			Foreground(m.theme.AccentSuccess).
			Bold(true)
		
		// Add subtle breathing pulse when playing
		if alpha > 0.95 {
			statusIcon = "▷"
		}
	}
	
	b.WriteString(statusStyle.Render(fmt.Sprintf("%s %s", statusIcon, statusText)))
	b.WriteString(div)
	
	// Frame timing info
	if len(m.frames) > m.currentFrame {
		duration := 1000 / m.fps
		timingStyle := lipgloss.NewStyle().
			Foreground(m.theme.FgSecondary)
		b.WriteString(timingStyle.Render(fmt.Sprintf("%dms/frame", duration)))
	}
	
	b.WriteString(div)
	
	// Loop indicator with toggle hint
	loopStyle := lipgloss.NewStyle().
		Foreground(m.theme.AccentInfo)
	b.WriteString(loopStyle.Render("loop: "))
	loopState := lipgloss.NewStyle().
		Foreground(m.theme.AccentSuccess).
		Bold(true).
		Render("on")
	b.WriteString(loopState)
	
	b.WriteString(div)
	
	// Shortcuts - zen grouping
	hintStyle := lipgloss.NewStyle().
		Foreground(m.theme.FgMuted)
	keyStyle := lipgloss.NewStyle().
		Foreground(m.theme.AccentSecondary)
	
	b.WriteString(hintStyle.Render("wheel: "))
	b.WriteString(keyStyle.Render("^j/k"))
	
	b.WriteString(div)
	
	// Playback control
	playHint := "[space]"
	playAction := "play"
	if m.playing {
		playAction = "pause"
	}
	b.WriteString(keyStyle.Render(playHint))
	b.WriteString(hintStyle.Render(" " + playAction))
	
	// Frame navigation
	b.WriteString(div)
	b.WriteString(keyStyle.Render(",/."))
	b.WriteString(hintStyle.Render(" frame"))
	
	// Breathing indicator when playing
	if m.playing && alpha > 0.95 {
		breatheStyle := lipgloss.NewStyle().
			Foreground(m.theme.AccentSuccess)
		b.WriteString(breatheStyle.Render(" ◉"))
	}
	
	b.WriteString(strings.Repeat(" ", max(0, lineWidth-lipgloss.Width(b.String())+2)))
	
	b.WriteString(borderStyle.Render("\n│ "))
	
	// Command hints - cleaner layout
	commandStyle := lipgloss.NewStyle().
		Foreground(m.theme.AccentWarning).
		Bold(false)
	b.WriteString(commandStyle.Render(":"))
	
	hintCommands := []string{
		"export", "import", "new", "save", "quit",
	}
	
	for i, cmd := range hintCommands {
		if i > 0 {
			b.WriteString(hintStyle.Render(" "))
		}
		b.WriteString(hintStyle.Render(cmd))
	}
	
	b.WriteString(divStyle.Render(" │ "))
	
	// Navigation hints
	navHints := []struct {
		key    string
		action string
	}{
		{"hjkl", "move"},
		{"+/-", "zoom"},
		{"g", "grid"},
		{"z", "zen"},
		{"?", "help"},
		{"q", "quit"},
	}
	
	for i, hint := range navHints {
		if i > 0 {
			b.WriteString(hintStyle.Render(" "))
		}
		b.WriteString(keyStyle.Render(hint.key))
		b.WriteString(hintStyle.Render(":" + hint.action))
	}
	
	return m.styles.Timeline.Render(b.String())
}

func (m Model) renderBottomStatus() string {
	var b strings.Builder
	
	modeStyle := m.styles.NormalMode
	modeText := "NORMAL"
	
	switch m.mode {
	case ModeInsert:
		modeStyle = m.styles.InsertMode
		modeText = "INSERT"
	case ModeCommand:
		modeStyle = m.styles.CommandMode
		modeText = "COMMAND"
		b.WriteString(modeStyle.Render(fmt.Sprintf(" %s ", modeText)))
		b.WriteString(" ")
		
		// Command input with cursor
		cmdStyle := lipgloss.NewStyle().
			Foreground(m.theme.AccentWarning).
			Bold(true)
		b.WriteString(cmdStyle.Render(":"))
		b.WriteString(lipgloss.NewStyle().Foreground(m.theme.FgPrimary).Render(m.command))
		b.WriteString(lipgloss.NewStyle().
			Foreground(m.theme.Cursor).
			Bold(true).
			Render("▌"))
		return b.String()
	}
	
	b.WriteString(modeStyle.Render(fmt.Sprintf(" %s ", modeText)))
	
	return b.String()
}
