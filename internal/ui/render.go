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
	
	// Breathing effect
	alpha := m.breathing.CurrentAlpha()
	
	borderStyle := lipgloss.NewStyle().
		Foreground(m.theme.Border)
	
	lineWidth := 110
	
	// Simplified top border - no title, just clean line
	b.WriteString(borderStyle.Render("├" + strings.Repeat("─", lineWidth) + "┤\n│ "))
	
	// Frame indicators - sleek minimal
	maxFrames := 30
	displayFrames := len(m.frames)
	if displayFrames > maxFrames {
		displayFrames = maxFrames
	}
	
	for i := 0; i < displayFrames; i++ {
		var frameStyle lipgloss.Style
		frameText := fmt.Sprintf("%2d", i+1)
		
		if i == m.currentFrame {
			// Current frame - sleek indicator
			frameStyle = lipgloss.NewStyle().
				Foreground(m.theme.AccentPrimary).
				Bold(true)
			frameText = "●"
			// Add breathing pulse
			if alpha > 0.9 {
				frameText = "◉"
			}
		} else if i < m.currentFrame {
			// Past frames - done
			frameStyle = lipgloss.NewStyle().
				Foreground(m.theme.TimelineInactive)
			frameText = "·"
		} else {
			// Future frames - todo
			frameStyle = lipgloss.NewStyle().
				Foreground(m.theme.FgMuted)
			frameText = "·"
		}
		
		b.WriteString(frameStyle.Render(frameText))
		if i < displayFrames-1 {
			b.WriteString(" ")
		}
	}
	
	if len(m.frames) > maxFrames {
		moreStyle := lipgloss.NewStyle().
			Foreground(m.theme.FgMuted)
		b.WriteString(moreStyle.Render(fmt.Sprintf(" …+%d", len(m.frames)-maxFrames)))
	}
	
	// Frame count indicator
	countStyle := lipgloss.NewStyle().
		Foreground(m.theme.FgSecondary).
		MarginLeft(2)
	b.WriteString(countStyle.Render(fmt.Sprintf("  %d/%d", m.currentFrame+1, len(m.frames))))
	
	// Pad to full width
	currentLine := b.String()
	lines := strings.Split(currentLine, "\n")
	lastLine := lines[len(lines)-1]
	// Remove border prefix for width calculation
	contentAfterBorder := strings.TrimPrefix(lastLine, "│ ")
	currentWidth := lipgloss.Width(contentAfterBorder)
	b.WriteString(strings.Repeat(" ", max(0, lineWidth-currentWidth)))
	
	b.WriteString(borderStyle.Render("\n│ "))
	
	// Playback status - zen and minimal
	divStyle := lipgloss.NewStyle().
		Foreground(m.theme.Border)
	div := divStyle.Render("  ")
	
	// Status with icon
	statusIcon := "■"
	statusStyle := lipgloss.NewStyle().
		Foreground(m.theme.FgMuted)
	
	if m.playing {
		statusIcon = "▸"
		statusStyle = lipgloss.NewStyle().
			Foreground(m.theme.AccentSuccess)
		
		// Subtle breathing pulse when playing
		if alpha > 0.95 {
			statusIcon = "▹"
		}
	}
	
	b.WriteString(statusStyle.Render(statusIcon))
	b.WriteString(div)
	
	// Frame timing
	if len(m.frames) > m.currentFrame {
		duration := 1000 / m.fps
		timingStyle := lipgloss.NewStyle().
			Foreground(m.theme.FgSecondary)
		b.WriteString(timingStyle.Render(fmt.Sprintf("%dms", duration)))
	}
	
	b.WriteString(div)
	
	// FPS
	fpsStyle := lipgloss.NewStyle().
		Foreground(m.theme.FgSecondary)
	b.WriteString(fpsStyle.Render(fmt.Sprintf("%dfps", m.fps)))
	
	b.WriteString(div)
	
	// Loop state
	loopStyle := lipgloss.NewStyle().
		Foreground(m.theme.AccentInfo)
	b.WriteString(loopStyle.Render("↻"))
	
	b.WriteString(strings.Repeat(" ", max(0, lineWidth-lipgloss.Width(b.String())+2)))
	
	b.WriteString(borderStyle.Render("\n│ "))
	
	// Commands and shortcuts - clean minimal layout
	hintStyle := lipgloss.NewStyle().
		Foreground(m.theme.FgMuted)
	keyStyle := lipgloss.NewStyle().
		Foreground(m.theme.AccentSecondary)
	
	// Playback shortcuts
	var shortcuts []string
	if m.playing {
		shortcuts = append(shortcuts, keyStyle.Render("space")+hintStyle.Render(":pause"))
	} else {
		shortcuts = append(shortcuts, keyStyle.Render("space")+hintStyle.Render(":play"))
	}
	shortcuts = append(shortcuts, keyStyle.Render(",")+hintStyle.Render("/")+keyStyle.Render(".")+hintStyle.Render(":frame"))
	shortcuts = append(shortcuts, keyStyle.Render("^j/k")+hintStyle.Render(":scroll"))
	
	// Mode shortcuts
	shortcuts = append(shortcuts, keyStyle.Render("i")+hintStyle.Render(":insert"))
	shortcuts = append(shortcuts, keyStyle.Render(":")+hintStyle.Render("cmd"))
	shortcuts = append(shortcuts, keyStyle.Render("?")+hintStyle.Render(":help"))
	shortcuts = append(shortcuts, keyStyle.Render("q")+hintStyle.Render(":quit"))
	
	b.WriteString(strings.Join(shortcuts, hintStyle.Render("  ")))
	
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
