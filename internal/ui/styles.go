package ui

import (
	"fmt"
	"math"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// Styles creates beautiful, themed styles
type Styles struct {
	Theme Theme
	
	// Status bar
	StatusBar        lipgloss.Style
	StatusBarTitle   lipgloss.Style
	StatusBarSection lipgloss.Style
	StatusBarDivider lipgloss.Style
	
	// Canvas
	CanvasBorder       lipgloss.Style
	CanvasBorderActive lipgloss.Style
	Canvas             lipgloss.Style
	Cursor             lipgloss.Style
	
	// Timeline
	Timeline         lipgloss.Style
	TimelineFrame    lipgloss.Style
	TimelineActive   lipgloss.Style
	TimelinePlayhead lipgloss.Style
	TimelineInfo     lipgloss.Style
	
	// Bottom status
	BottomStatus lipgloss.Style
	CommandMode  lipgloss.Style
	InsertMode   lipgloss.Style
	NormalMode   lipgloss.Style
	
	// Wheel
	WheelSection       lipgloss.Style
	WheelSectionActive lipgloss.Style
	WheelItem          lipgloss.Style
	WheelItemActive    lipgloss.Style
	
	// Zen mode
	ZenCanvas lipgloss.Style
	ZenStatus lipgloss.Style
	
	// Accents
	Success lipgloss.Style
	Warning lipgloss.Style
	Error   lipgloss.Style
	Info    lipgloss.Style
	Muted   lipgloss.Style
	Bright  lipgloss.Style
}

// NewStyles creates a beautiful style set from theme
func NewStyles(theme Theme) Styles {
	s := Styles{Theme: theme}
	
	// Status Bar - Sleek with gradient effect
	s.StatusBar = lipgloss.NewStyle().
		Background(theme.StatusBg).
		Foreground(theme.StatusFg).
		Padding(0, 1)
	
	s.StatusBarTitle = lipgloss.NewStyle().
		Background(theme.AccentPrimary).
		Foreground(theme.BgPrimary).
		Bold(true).
		Padding(0, 2)
	
	s.StatusBarSection = lipgloss.NewStyle().
		Foreground(theme.FgSecondary).
		Padding(0, 1)
	
	s.StatusBarDivider = lipgloss.NewStyle().
		Foreground(theme.Border).
		SetString("│")
	
	// Canvas - Clean borders with active state
	s.CanvasBorder = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Border).
		Padding(0, 1)
	
	s.CanvasBorderActive = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.BorderActive).
		Padding(0, 1)
	
	s.Canvas = lipgloss.NewStyle().
		Background(theme.BgCanvas)
	
	s.Cursor = lipgloss.NewStyle().
		Foreground(theme.Cursor).
		Bold(true)
	
	// Timeline - Beautiful frame indicators
	s.Timeline = lipgloss.NewStyle().
		Background(theme.BgSecondary).
		Foreground(theme.FgSecondary).
		Padding(0, 1)
	
	s.TimelineFrame = lipgloss.NewStyle().
		Foreground(theme.TimelineInactive).
		Padding(0, 1)
	
	s.TimelineActive = lipgloss.NewStyle().
		Foreground(theme.TimelineActive).
		Bold(true).
		Padding(0, 1)
	
	s.TimelinePlayhead = lipgloss.NewStyle().
		Background(theme.PlayheadColor).
		Foreground(theme.BgPrimary).
		Bold(true).
		Padding(0, 1)
	
	s.TimelineInfo = lipgloss.NewStyle().
		Foreground(theme.FgMuted).
		Padding(0, 1)
	
	// Bottom Status - Mode-aware colors
	s.BottomStatus = lipgloss.NewStyle().
		Background(theme.BgSecondary).
		Foreground(theme.FgSecondary).
		Padding(0, 1)
	
	s.CommandMode = lipgloss.NewStyle().
		Background(theme.AccentWarning).
		Foreground(theme.BgPrimary).
		Bold(true).
		Padding(0, 1)
	
	s.InsertMode = lipgloss.NewStyle().
		Background(theme.AccentSuccess).
		Foreground(theme.BgPrimary).
		Bold(true).
		Padding(0, 1)
	
	s.NormalMode = lipgloss.NewStyle().
		Background(theme.AccentInfo).
		Foreground(theme.BgPrimary).
		Bold(true).
		Padding(0, 1)
	
	// Wheel - Radial menu styling
	s.WheelSection = lipgloss.NewStyle().
		Foreground(theme.FgSecondary).
		Padding(0, 1)
	
	s.WheelSectionActive = lipgloss.NewStyle().
		Foreground(theme.AccentPrimary).
		Bold(true).
		Padding(0, 1)
	
	s.WheelItem = lipgloss.NewStyle().
		Foreground(theme.FgMuted).
		Padding(0, 1)
	
	s.WheelItemActive = lipgloss.NewStyle().
		Foreground(theme.AccentSecondary).
		Bold(true).
		Padding(0, 1)
	
	// Zen Mode - Minimal and beautiful
	s.ZenCanvas = lipgloss.NewStyle().
		Background(theme.BgCanvas)
	
	s.ZenStatus = lipgloss.NewStyle().
		Foreground(theme.FgMuted).
		Faint(true)
	
	// Accent colors
	s.Success = lipgloss.NewStyle().
		Foreground(theme.AccentSuccess).
		Bold(true)
	
	s.Warning = lipgloss.NewStyle().
		Foreground(theme.AccentWarning).
		Bold(true)
	
	s.Error = lipgloss.NewStyle().
		Foreground(theme.AccentError).
		Bold(true)
	
	s.Info = lipgloss.NewStyle().
		Foreground(theme.AccentInfo)
	
	s.Muted = lipgloss.NewStyle().
		Foreground(theme.FgMuted)
	
	s.Bright = lipgloss.NewStyle().
		Foreground(theme.FgBright).
		Bold(true)
	
	return s
}

// Breathing effect - pulsing opacity/brightness
type BreathingEffect struct {
	StartTime time.Time
	Duration  time.Duration
	MinAlpha  float64
	MaxAlpha  float64
}

// NewBreathingEffect creates a new breathing animation
func NewBreathingEffect(duration time.Duration) *BreathingEffect {
	return &BreathingEffect{
		StartTime: time.Now(),
		Duration:  duration,
		MinAlpha:  0.5,
		MaxAlpha:  1.0,
	}
}

// CurrentAlpha returns current alpha value in breathing cycle
func (b *BreathingEffect) CurrentAlpha() float64 {
	elapsed := time.Since(b.StartTime)
	progress := float64(elapsed%b.Duration) / float64(b.Duration)
	
	// Sine wave for smooth breathing
	alpha := b.MinAlpha + (b.MaxAlpha-b.MinAlpha)*(1+math.Sin(progress*2*math.Pi))/2
	return alpha
}

// ApplyBreathing applies breathing effect to a color
func (b *BreathingEffect) ApplyBreathing(style lipgloss.Style) lipgloss.Style {
	alpha := b.CurrentAlpha()
	// Lipgloss doesn't support alpha directly, but we can simulate with brightness
	if alpha < 0.8 {
		return style.Faint(true)
	}
	return style
}

// Gradient creates a horizontal gradient effect
func Gradient(text string, startColor, endColor lipgloss.Color) string {
	if len(text) == 0 {
		return ""
	}
	
	result := ""
	for i, char := range text {
		progress := float64(i) / float64(len(text)-1)
		color := interpolateColor(startColor, endColor, progress)
		style := lipgloss.NewStyle().Foreground(color)
		result += style.Render(string(char))
	}
	return result
}

// interpolateColor interpolates between two colors
func interpolateColor(start, end lipgloss.Color, progress float64) lipgloss.Color {
	// Simple linear interpolation
	// This is a simplified version - for production, you'd parse hex properly
	return start // Placeholder - lipgloss colors are complex
}

// ProgressBar creates a beautiful progress bar
func ProgressBar(current, total int, width int, theme Theme) string {
	if total == 0 {
		return ""
	}
	
	percent := float64(current) / float64(total)
	filled := int(float64(width) * percent)
	
	filledStyle := lipgloss.NewStyle().Foreground(theme.AccentPrimary)
	emptyStyle := lipgloss.NewStyle().Foreground(theme.FgMuted)
	
	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += filledStyle.Render("━")
		} else {
			bar += emptyStyle.Render("━")
		}
	}
	
	return bar
}

// BoxWithTitle creates a beautiful titled box
func BoxWithTitle(title, content string, theme Theme, active bool) string {
	borderColor := theme.Border
	titleColor := theme.FgSecondary
	
	if active {
		borderColor = theme.BorderActive
		titleColor = theme.AccentPrimary
	}
	
	titleStyle := lipgloss.NewStyle().
		Foreground(titleColor).
		Bold(true)
	
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1)
	
	// Add title to top border
	titleStr := titleStyle.Render(fmt.Sprintf(" %s ", title))
	box := boxStyle.Render(content)
	
	// Simple approach: render title and box separately
	return titleStr + "\n" + box
}

// Sparkle adds sparkle effect (random stars)
func Sparkle(text string, intensity float64, theme Theme) string {
	if intensity <= 0 {
		return text
	}
	
	sparkleStyle := lipgloss.NewStyle().
		Foreground(theme.AccentPrimary).
		Bold(true)
	
	result := ""
	for _, char := range text {
		// Random sparkle based on intensity
		if math.Mod(float64(char), 1.0/intensity) < 0.1 {
			result += sparkleStyle.Render("✨")
		} else {
			result += string(char)
		}
	}
	return result
}

// Glow adds a glow effect around text
func Glow(text string, theme Theme) string {
	glowStyle := lipgloss.NewStyle().
		Foreground(theme.AccentPrimary).
		Faint(true)
	
	brightStyle := lipgloss.NewStyle().
		Foreground(theme.FgBright).
		Bold(true)
	
	return glowStyle.Render("▓") + brightStyle.Render(text) + glowStyle.Render("▓")
}
