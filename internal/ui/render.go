package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Beautiful themed render functions

func (m Model) renderStatusBar() string {
	// Compact status bar without emoji clutter
	
	// File info - only show icon if modified
	fileStyle := m.styles.StatusBarSection
	fileInfo := m.filename
	if m.modified {
		fileStyle = m.styles.Warning
		fileInfo = m.filename + " *"
	}
	
	// Canvas dimensions
	canvasInfo := fmt.Sprintf("%dx%d", 
		m.frames[m.currentFrame].Width, 
		m.frames[m.currentFrame].Height)
	
	// Frame counter
	frameInfo := fmt.Sprintf("frame %d/%d", 
		m.currentFrame+1, len(m.frames))
	
	// FPS with playback state
	fpsInfo := fmt.Sprintf("%dfps", m.fps)
	if m.playing {
		fpsStyle := lipgloss.NewStyle().
			Foreground(m.theme.AccentSuccess).
			Bold(true)
		fpsInfo = fpsStyle.Render("▸ " + fpsInfo)
	}
	
	// Tool info - clean without icon
	toolName := toolNames[m.selectedTool]
	toolInfo := fmt.Sprintf("%s", toolName)
	
	// Layer info
	layerInfo := fmt.Sprintf("layer %d/%d", 
		m.currentLayer+1, len(m.layers))
	
	// Divider
	div := lipgloss.NewStyle().
		Foreground(m.theme.Border).
		Render(" │ ")
	
	// Assemble compact status
	sections := []string{
		fileStyle.Render(fileInfo),
		canvasInfo,
		frameInfo,
		fpsInfo,
		toolInfo,
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
	
	// Frame indicators - scrolling window centered on current frame
	maxFrames := 30
	totalFrames := len(m.frames)
	
	// Calculate visible frame window
	var startFrame, endFrame int
	if totalFrames <= maxFrames {
		// Show all frames
		startFrame = 0
		endFrame = totalFrames
	} else {
		// Scroll window centered on current frame
		halfWindow := maxFrames / 2
		startFrame = m.currentFrame - halfWindow
		endFrame = m.currentFrame + halfWindow
		
		// Clamp to valid range
		if startFrame < 0 {
			startFrame = 0
			endFrame = maxFrames
		} else if endFrame > totalFrames {
			endFrame = totalFrames
			startFrame = totalFrames - maxFrames
		}
	}
	
	// Show window indicator if scrolled
	if startFrame > 0 {
		leftStyle := lipgloss.NewStyle().
			Foreground(m.theme.FgMuted)
		b.WriteString(leftStyle.Render("‹ "))
	}
	
	for i := startFrame; i < endFrame; i++ {
		var frameStyle lipgloss.Style
		var frameText string
		
		if i == m.currentFrame {
			// Current frame - prominent indicator
			frameStyle = lipgloss.NewStyle().
				Foreground(m.theme.AccentPrimary).
				Bold(true)
			frameText = "●"
			// Add breathing pulse
			if alpha > 0.9 {
				frameText = "◉"
			}
		} else if m.frames[i].Modified {
			// Modified frame marker
			frameStyle = lipgloss.NewStyle().
				Foreground(m.theme.AccentWarning)
			frameText = "◉"
		} else if i < m.currentFrame {
			// Past frames
			frameStyle = lipgloss.NewStyle().
				Foreground(m.theme.TimelineInactive)
			frameText = "·"
		} else {
			// Future frames
			frameStyle = lipgloss.NewStyle().
				Foreground(m.theme.FgMuted)
			frameText = "·"
		}
		
		b.WriteString(frameStyle.Render(frameText))
		if i < endFrame-1 {
			b.WriteString(" ")
		}
	}
	
	// Show right indicator if more frames
	if endFrame < totalFrames {
		rightStyle := lipgloss.NewStyle().
			Foreground(m.theme.FgMuted)
		b.WriteString(rightStyle.Render(" ›"))
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
