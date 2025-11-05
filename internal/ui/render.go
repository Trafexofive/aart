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
	
	// Top border with title
	titleStyle := lipgloss.NewStyle().
		Foreground(m.theme.AccentPrimary).
		Bold(true)
	title := titleStyle.Render(" TIMELINE ")
	
	borderStyle := lipgloss.NewStyle().
		Foreground(m.theme.Border)
	
	b.WriteString(borderStyle.Render("├─"))
	b.WriteString(title)
	b.WriteString(borderStyle.Render(strings.Repeat("─", 100)))
	b.WriteString(borderStyle.Render("┤\n│ "))
	
	// Frame indicators with beautiful styling
	for i := 0; i < len(m.frames) && i < 40; i++ {
		var frameStyle lipgloss.Style
		frameNum := fmt.Sprintf("%2d", i+1)
		
		if i == m.currentFrame {
			// Current frame - highlighted with playhead
			frameStyle = lipgloss.NewStyle().
				Background(m.theme.PlayheadColor).
				Foreground(m.theme.BgPrimary).
				Bold(true)
			frameNum = "▓▓"
		} else if i < m.currentFrame {
			// Past frames - muted
			frameStyle = lipgloss.NewStyle().
				Foreground(m.theme.TimelineInactive)
		} else {
			// Future frames - normal
			frameStyle = lipgloss.NewStyle().
				Foreground(m.theme.TimelineActive)
		}
		
		b.WriteString(frameStyle.Render(frameNum))
		b.WriteString(" ")
	}
	
	if len(m.frames) > 40 {
		moreStyle := lipgloss.NewStyle().
			Foreground(m.theme.FgMuted)
		b.WriteString(moreStyle.Render(fmt.Sprintf("... +%d", len(m.frames)-40)))
	}
	
	b.WriteString(borderStyle.Render("\n│ "))
	
	// Playback info with icons
	statusIcon := "⏹"
	statusText := "stopped"
	statusStyle := m.styles.Muted
	
	if m.playing {
		statusIcon = "▶"
		statusText = "playing"
		statusStyle = lipgloss.NewStyle().
			Foreground(m.theme.AccentSuccess).
			Bold(true)
	}
	
	b.WriteString(statusStyle.Render(fmt.Sprintf("%s %s", statusIcon, statusText)))
	b.WriteString(lipgloss.NewStyle().Foreground(m.theme.Border).Render(" │ "))
	
	// Frame duration
	if len(m.frames) > m.currentFrame {
		duration := 1000 / m.fps
		durationStyle := lipgloss.NewStyle().Foreground(m.theme.FgSecondary)
		b.WriteString(durationStyle.Render(fmt.Sprintf("%dms/frame", duration)))
	}
	
	b.WriteString(lipgloss.NewStyle().Foreground(m.theme.Border).Render(" │ "))
	
	// Loop indicator
	loopStyle := lipgloss.NewStyle().Foreground(m.theme.AccentInfo)
	b.WriteString(loopStyle.Render("loop: on"))
	
	b.WriteString(lipgloss.NewStyle().Foreground(m.theme.Border).Render(" │ "))
	
	// Shortcuts hint
	hintStyle := lipgloss.NewStyle().Foreground(m.theme.FgMuted)
	b.WriteString(hintStyle.Render("ctrl-j/k: wheel"))
	
	b.WriteString(lipgloss.NewStyle().Foreground(m.theme.Border).Render(" │ "))
	
	controlStyle := lipgloss.NewStyle().Foreground(m.theme.AccentSecondary)
	b.WriteString(controlStyle.Render("[space] pause"))
	
	// Breathing effect on timeline
	if m.breathing.CurrentAlpha() > 0.95 && m.playing {
		b.WriteString(lipgloss.NewStyle().
			Foreground(m.theme.AccentSuccess).
			Render(" ●"))
	}
	
	b.WriteString(borderStyle.Render("\n│ "))
	
	// Bottom hints
	commandStyle := lipgloss.NewStyle().
		Foreground(m.theme.AccentWarning)
	b.WriteString(commandStyle.Render(":export out.ans"))
	
	b.WriteString(lipgloss.NewStyle().Foreground(m.theme.Border).Render(" | "))
	
	b.WriteString(hintStyle.Render("hjkl:move ctrl-j/k:wheel +/-:zoom g:grid z:zen ?:help q:quit"))
	
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
