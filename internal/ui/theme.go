package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme defines a complete color scheme
type Theme struct {
	Name string
	
	// Backgrounds
	BgPrimary   lipgloss.Color
	BgSecondary lipgloss.Color
	BgTertiary  lipgloss.Color
	BgCanvas    lipgloss.Color
	
	// Foregrounds
	FgPrimary   lipgloss.Color
	FgSecondary lipgloss.Color
	FgMuted     lipgloss.Color
	FgBright    lipgloss.Color
	
	// Accents
	AccentPrimary   lipgloss.Color
	AccentSecondary lipgloss.Color
	AccentSuccess   lipgloss.Color
	AccentWarning   lipgloss.Color
	AccentError     lipgloss.Color
	AccentInfo      lipgloss.Color
	
	// UI Elements
	Border       lipgloss.Color
	BorderActive lipgloss.Color
	Selection    lipgloss.Color
	Cursor       lipgloss.Color
	
	// Timeline
	TimelineActive   lipgloss.Color
	TimelineInactive lipgloss.Color
	PlayheadColor    lipgloss.Color
	
	// Status
	StatusBg lipgloss.Color
	StatusFg lipgloss.Color
}

// Predefined themes
var (
	// Nord theme - Cool, professional blues
	NordTheme = Theme{
		Name:             "nord",
		BgPrimary:        lipgloss.Color("#2E3440"),
		BgSecondary:      lipgloss.Color("#3B4252"),
		BgTertiary:       lipgloss.Color("#434C5E"),
		BgCanvas:         lipgloss.Color("#2E3440"),
		FgPrimary:        lipgloss.Color("#ECEFF4"),
		FgSecondary:      lipgloss.Color("#D8DEE9"),
		FgMuted:          lipgloss.Color("#4C566A"),
		FgBright:         lipgloss.Color("#FFFFFF"),
		AccentPrimary:    lipgloss.Color("#88C0D0"),
		AccentSecondary:  lipgloss.Color("#81A1C1"),
		AccentSuccess:    lipgloss.Color("#A3BE8C"),
		AccentWarning:    lipgloss.Color("#EBCB8B"),
		AccentError:      lipgloss.Color("#BF616A"),
		AccentInfo:       lipgloss.Color("#5E81AC"),
		Border:           lipgloss.Color("#4C566A"),
		BorderActive:     lipgloss.Color("#88C0D0"),
		Selection:        lipgloss.Color("#434C5E"),
		Cursor:           lipgloss.Color("#88C0D0"),
		TimelineActive:   lipgloss.Color("#88C0D0"),
		TimelineInactive: lipgloss.Color("#4C566A"),
		PlayheadColor:    lipgloss.Color("#A3BE8C"),
		StatusBg:         lipgloss.Color("#3B4252"),
		StatusFg:         lipgloss.Color("#ECEFF4"),
	}

	// Dracula theme - Purple and pink vibes
	DraculaTheme = Theme{
		Name:             "dracula",
		BgPrimary:        lipgloss.Color("#282A36"),
		BgSecondary:      lipgloss.Color("#44475A"),
		BgTertiary:       lipgloss.Color("#6272A4"),
		BgCanvas:         lipgloss.Color("#282A36"),
		FgPrimary:        lipgloss.Color("#F8F8F2"),
		FgSecondary:      lipgloss.Color("#E6E6E6"),
		FgMuted:          lipgloss.Color("#6272A4"),
		FgBright:         lipgloss.Color("#FFFFFF"),
		AccentPrimary:    lipgloss.Color("#BD93F9"),
		AccentSecondary:  lipgloss.Color("#FF79C6"),
		AccentSuccess:    lipgloss.Color("#50FA7B"),
		AccentWarning:    lipgloss.Color("#F1FA8C"),
		AccentError:      lipgloss.Color("#FF5555"),
		AccentInfo:       lipgloss.Color("#8BE9FD"),
		Border:           lipgloss.Color("#6272A4"),
		BorderActive:     lipgloss.Color("#BD93F9"),
		Selection:        lipgloss.Color("#44475A"),
		Cursor:           lipgloss.Color("#FF79C6"),
		TimelineActive:   lipgloss.Color("#BD93F9"),
		TimelineInactive: lipgloss.Color("#44475A"),
		PlayheadColor:    lipgloss.Color("#50FA7B"),
		StatusBg:         lipgloss.Color("#44475A"),
		StatusFg:         lipgloss.Color("#F8F8F2"),
	}

	// Tokyo Night theme - Modern, vibrant
	TokyoNightTheme = Theme{
		Name:             "tokyo-night",
		BgPrimary:        lipgloss.Color("#1A1B26"),
		BgSecondary:      lipgloss.Color("#24283B"),
		BgTertiary:       lipgloss.Color("#414868"),
		BgCanvas:         lipgloss.Color("#1A1B26"),
		FgPrimary:        lipgloss.Color("#C0CAF5"),
		FgSecondary:      lipgloss.Color("#A9B1D6"),
		FgMuted:          lipgloss.Color("#565F89"),
		FgBright:         lipgloss.Color("#FFFFFF"),
		AccentPrimary:    lipgloss.Color("#7AA2F7"),
		AccentSecondary:  lipgloss.Color("#BB9AF7"),
		AccentSuccess:    lipgloss.Color("#9ECE6A"),
		AccentWarning:    lipgloss.Color("#E0AF68"),
		AccentError:      lipgloss.Color("#F7768E"),
		AccentInfo:       lipgloss.Color("#7DCFFF"),
		Border:           lipgloss.Color("#414868"),
		BorderActive:     lipgloss.Color("#7AA2F7"),
		Selection:        lipgloss.Color("#283457"),
		Cursor:           lipgloss.Color("#BB9AF7"),
		TimelineActive:   lipgloss.Color("#7AA2F7"),
		TimelineInactive: lipgloss.Color("#414868"),
		PlayheadColor:    lipgloss.Color("#9ECE6A"),
		StatusBg:         lipgloss.Color("#24283B"),
		StatusFg:         lipgloss.Color("#C0CAF5"),
	}

	// Gruvbox theme - Warm, retro
	GruvboxTheme = Theme{
		Name:             "gruvbox",
		BgPrimary:        lipgloss.Color("#282828"),
		BgSecondary:      lipgloss.Color("#3C3836"),
		BgTertiary:       lipgloss.Color("#504945"),
		BgCanvas:         lipgloss.Color("#282828"),
		FgPrimary:        lipgloss.Color("#EBDBB2"),
		FgSecondary:      lipgloss.Color("#D5C4A1"),
		FgMuted:          lipgloss.Color("#665C54"),
		FgBright:         lipgloss.Color("#FBF1C7"),
		AccentPrimary:    lipgloss.Color("#83A598"),
		AccentSecondary:  lipgloss.Color("#D3869B"),
		AccentSuccess:    lipgloss.Color("#B8BB26"),
		AccentWarning:    lipgloss.Color("#FABD2F"),
		AccentError:      lipgloss.Color("#FB4934"),
		AccentInfo:       lipgloss.Color("#8EC07C"),
		Border:           lipgloss.Color("#504945"),
		BorderActive:     lipgloss.Color("#83A598"),
		Selection:        lipgloss.Color("#3C3836"),
		Cursor:           lipgloss.Color("#FE8019"),
		TimelineActive:   lipgloss.Color("#83A598"),
		TimelineInactive: lipgloss.Color("#504945"),
		PlayheadColor:    lipgloss.Color("#B8BB26"),
		StatusBg:         lipgloss.Color("#3C3836"),
		StatusFg:         lipgloss.Color("#EBDBB2"),
	}

	// Catppuccin Mocha - Pastel perfection
	CatppuccinTheme = Theme{
		Name:             "catppuccin",
		BgPrimary:        lipgloss.Color("#1E1E2E"),
		BgSecondary:      lipgloss.Color("#313244"),
		BgTertiary:       lipgloss.Color("#45475A"),
		BgCanvas:         lipgloss.Color("#1E1E2E"),
		FgPrimary:        lipgloss.Color("#CDD6F4"),
		FgSecondary:      lipgloss.Color("#BAC2DE"),
		FgMuted:          lipgloss.Color("#6C7086"),
		FgBright:         lipgloss.Color("#FFFFFF"),
		AccentPrimary:    lipgloss.Color("#89B4FA"),
		AccentSecondary:  lipgloss.Color("#CBA6F7"),
		AccentSuccess:    lipgloss.Color("#A6E3A1"),
		AccentWarning:    lipgloss.Color("#F9E2AF"),
		AccentError:      lipgloss.Color("#F38BA8"),
		AccentInfo:       lipgloss.Color("#94E2D5"),
		Border:           lipgloss.Color("#45475A"),
		BorderActive:     lipgloss.Color("#89B4FA"),
		Selection:        lipgloss.Color("#313244"),
		Cursor:           lipgloss.Color("#F5C2E7"),
		TimelineActive:   lipgloss.Color("#89B4FA"),
		TimelineInactive: lipgloss.Color("#45475A"),
		PlayheadColor:    lipgloss.Color("#A6E3A1"),
		StatusBg:         lipgloss.Color("#313244"),
		StatusFg:         lipgloss.Color("#CDD6F4"),
	}

	// Oceanic - Blue depths
	OceanicTheme = Theme{
		Name:             "oceanic",
		BgPrimary:        lipgloss.Color("#1B2B34"),
		BgSecondary:      lipgloss.Color("#343D46"),
		BgTertiary:       lipgloss.Color("#4F5B66"),
		BgCanvas:         lipgloss.Color("#1B2B34"),
		FgPrimary:        lipgloss.Color("#C0C5CE"),
		FgSecondary:      lipgloss.Color("#A7ADBA"),
		FgMuted:          lipgloss.Color("#65737E"),
		FgBright:         lipgloss.Color("#D8DEE9"),
		AccentPrimary:    lipgloss.Color("#6699CC"),
		AccentSecondary:  lipgloss.Color("#C594C5"),
		AccentSuccess:    lipgloss.Color("#99C794"),
		AccentWarning:    lipgloss.Color("#FAC863"),
		AccentError:      lipgloss.Color("#EC5f67"),
		AccentInfo:       lipgloss.Color("#5FB3B3"),
		Border:           lipgloss.Color("#4F5B66"),
		BorderActive:     lipgloss.Color("#6699CC"),
		Selection:        lipgloss.Color("#343D46"),
		Cursor:           lipgloss.Color("#5FB3B3"),
		TimelineActive:   lipgloss.Color("#6699CC"),
		TimelineInactive: lipgloss.Color("#4F5B66"),
		PlayheadColor:    lipgloss.Color("#99C794"),
		StatusBg:         lipgloss.Color("#343D46"),
		StatusFg:         lipgloss.Color("#C0C5CE"),
	}
)

// GetTheme returns a theme by name
func GetTheme(name string) Theme {
	switch name {
	case "nord":
		return NordTheme
	case "dracula":
		return DraculaTheme
	case "tokyo-night":
		return TokyoNightTheme
	case "gruvbox":
		return GruvboxTheme
	case "catppuccin":
		return CatppuccinTheme
	case "oceanic":
		return OceanicTheme
	default:
		return TokyoNightTheme // Default to Tokyo Night
	}
}

// AvailableThemes returns list of theme names
func AvailableThemes() []string {
	return []string{
		"nord",
		"dracula",
		"tokyo-night",
		"gruvbox",
		"catppuccin",
		"oceanic",
	}
}
