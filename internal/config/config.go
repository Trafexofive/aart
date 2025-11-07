package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Version   string         `yaml:"version"`
	Editor    EditorConfig   `yaml:"editor"`
	UI        UIConfig       `yaml:"ui"`
	Colors    ColorScheme    `yaml:"colors"`
	Recent    RecentFiles    `yaml:"recent"`
	Converter ConvertConfig  `yaml:"converter"`
	Startup   StartupConfig  `yaml:"startup"`
	Keybinds  KeyBindings    `yaml:"keybindings,omitempty"`
}

// EditorConfig contains editor preferences
type EditorConfig struct {
	DefaultWidth  int    `yaml:"default_width"`
	DefaultHeight int    `yaml:"default_height"`
	DefaultFPS    int    `yaml:"default_fps"`
	AutoSave      bool   `yaml:"auto_save"`
	AutoSaveInterval int `yaml:"auto_save_interval"` // seconds
	TabSize       int    `yaml:"tab_size"`
	ShowGrid      bool   `yaml:"show_grid"`
	ShowLineNumbers bool `yaml:"show_line_numbers"`
	ZenMode       bool   `yaml:"zen_mode"`
}

// UIConfig contains UI preferences
type UIConfig struct {
	Theme              string  `yaml:"theme"` // dark, light, custom
	ShowStatusBar      bool    `yaml:"show_status_bar"`
	ShowTimeline       bool    `yaml:"show_timeline"`
	ShowWheelByDefault bool    `yaml:"show_wheel_by_default"`
	CursorStyle        string  `yaml:"cursor_style"` // block, line, underline
	AnimationSmooth    bool    `yaml:"animation_smooth"`
	ProgressStyle      string  `yaml:"progress_style"` // bar, spinner, minimal
	BorderStyle        string  `yaml:"border_style"`   // rounded, thick, double, ascii
	TimelineStyle      string  `yaml:"timeline_style"` // compact, detailed, minimal
	StatusBarPosition  string  `yaml:"status_bar_position"` // top, bottom
}

// ColorScheme defines the color palette
type ColorScheme struct {
	Name       string            `yaml:"name"`
	Foreground string            `yaml:"foreground"`
	Background string            `yaml:"background"`
	Cursor     string            `yaml:"cursor"`
	Selection  string            `yaml:"selection"`
	StatusBar  string            `yaml:"status_bar"`
	Timeline   string            `yaml:"timeline"`
	Border     string            `yaml:"border"`
	Custom     map[string]string `yaml:"custom,omitempty"`
}

// RecentFiles tracks recently opened files
type RecentFiles struct {
	Files []RecentFile `yaml:"files"`
	Max   int          `yaml:"max_entries"`
}

// RecentFile represents a recently opened file
type RecentFile struct {
	Path      string    `yaml:"path"`
	Timestamp time.Time `yaml:"timestamp"`
	Frames    int       `yaml:"frames,omitempty"`
}

// ConvertConfig contains GIF conversion preferences
type ConvertConfig struct {
	DefaultMethod string `yaml:"default_method"` // luminosity, block, edge, dither
	DefaultChars  string `yaml:"default_chars,omitempty"`
	PreserveAspect bool  `yaml:"preserve_aspect"`
	Quality       string `yaml:"quality"` // low, medium, high
}

// StartupConfig contains startup screen preferences
type StartupConfig struct {
	ShowStartupPage   bool   `yaml:"show_startup_page"`    // Show startup page on launch
	ArtworkFile       string `yaml:"artwork_file"`         // Custom ASCII art file for logo
	ArtworkInline     string `yaml:"artwork_inline"`       // Inline ASCII art (multiline)
	ArtworkBorder     bool   `yaml:"artwork_border"`       // Show border around artwork
	ArtworkSize       string `yaml:"artwork_size"`         // Size as percentage (e.g. "40p" = 40% of screen height)
	ArtworkOffsetX    int    `yaml:"artwork_offset_x"`     // X offset for artwork
	ArtworkOffsetY    int    `yaml:"artwork_offset_y"`     // Y offset for artwork
	ArtworkWidth      int    `yaml:"artwork_width"`        // Max width (0 = auto, overridden by artwork_size)
	ArtworkHeight     int    `yaml:"artwork_height"`       // Max height (0 = auto, overridden by artwork_size)
	ShowRecentFiles   bool   `yaml:"show_recent_files"`    // Show recent files panel
	ShowTips          bool   `yaml:"show_tips"`            // Show rotating tips
	TipRotationSec    int    `yaml:"tip_rotation_seconds"` // Seconds before rotating tips
	BreathingEffect   bool   `yaml:"breathing_effect"`     // Enable breathing animation
}

// KeyBindings contains custom keybindings
type KeyBindings struct {
	Play        string `yaml:"play,omitempty"`
	Quit        string `yaml:"quit,omitempty"`
	Save        string `yaml:"save,omitempty"`
	ZenMode     string `yaml:"zen_mode,omitempty"`
	NextFrame   string `yaml:"next_frame,omitempty"`
	PrevFrame   string `yaml:"prev_frame,omitempty"`
}

var (
	// DefaultConfig provides sensible defaults
	DefaultConfig = Config{
		Version: "0.1.0",
		Editor: EditorConfig{
			DefaultWidth:     80,
			DefaultHeight:    24,
			DefaultFPS:       12,
			AutoSave:         false,
			AutoSaveInterval: 300,
			TabSize:          4,
			ShowGrid:         false,
			ShowLineNumbers:  false,
			ZenMode:          false,
		},
		UI: UIConfig{
			Theme:              "tokyo-night",
			ShowStatusBar:      true,
			ShowTimeline:       true,
			ShowWheelByDefault: false,
			CursorStyle:        "line",
			AnimationSmooth:    true,
			ProgressStyle:      "bar",
			BorderStyle:        "rounded",
			TimelineStyle:      "detailed",
			StatusBarPosition:  "bottom",
		},
		Colors: ColorScheme{
			Name:       "default",
			Foreground: "#FFFFFF",
			Background: "#000000",
			Cursor:     "#FFFF00",
			Selection:  "#444444",
			StatusBar:  "#333333",
			Timeline:   "#222222",
			Border:     "#666666",
		},
		Recent: RecentFiles{
			Files: []RecentFile{},
			Max:   10,
		},
		Converter: ConvertConfig{
			DefaultMethod:  "luminosity",
			DefaultChars:   "",
			PreserveAspect: true,
			Quality:        "high",
		},
		Startup: StartupConfig{
			ShowStartupPage:   true,
			ArtworkFile:       "",
			ArtworkInline:     "",
			ArtworkBorder:     true,
			ArtworkOffsetX:    0,
			ArtworkOffsetY:    0,
			ArtworkWidth:      0,
			ArtworkHeight:     0,
			ShowRecentFiles:   true,
			ShowTips:          true,
			TipRotationSec:    5,
			BreathingEffect:   true,
		},
	}
)

// ConfigDir returns the configuration directory path
func ConfigDir() (string, error) {
	// Try XDG_CONFIG_HOME first
	if configHome := os.Getenv("XDG_CONFIG_HOME"); configHome != "" {
		return filepath.Join(configHome, "aart"), nil
	}
	
	// Fall back to ~/.config
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	
	return filepath.Join(home, ".config", "aart"), nil
}

// ConfigPath returns the full path to the config file
func ConfigPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yml"), nil
}

// Load loads the configuration from file
func Load() (*Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	// If config doesn't exist, return default
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &DefaultConfig, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Merge with defaults for any missing fields
	mergeDefaults(&config)

	return &config, nil
}

// Save saves the configuration to file
func Save(config *Config) error {
	dir, err := ConfigDir()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	path := filepath.Join(dir, "config.yml")

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// Init initializes the configuration directory and creates default config
func Init() error {
	dir, err := ConfigDir()
	if err != nil {
		return err
	}

	// Create config directory
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create subdirectories
	dirs := []string{"themes", "templates", "recent", "backups"}
	for _, d := range dirs {
		subdir := filepath.Join(dir, d)
		if err := os.MkdirAll(subdir, 0755); err != nil {
			return fmt.Errorf("failed to create %s directory: %w", d, err)
		}
	}

	// Create default config if it doesn't exist
	path := filepath.Join(dir, "config.yml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := Save(&DefaultConfig); err != nil {
			return err
		}
	}

	return nil
}

// mergeDefaults merges default values for missing fields
func mergeDefaults(config *Config) {
	if config.Editor.DefaultWidth == 0 {
		config.Editor.DefaultWidth = DefaultConfig.Editor.DefaultWidth
	}
	if config.Editor.DefaultHeight == 0 {
		config.Editor.DefaultHeight = DefaultConfig.Editor.DefaultHeight
	}
	if config.Editor.DefaultFPS == 0 {
		config.Editor.DefaultFPS = DefaultConfig.Editor.DefaultFPS
	}
	if config.Editor.AutoSaveInterval == 0 {
		config.Editor.AutoSaveInterval = DefaultConfig.Editor.AutoSaveInterval
	}
	if config.Recent.Max == 0 {
		config.Recent.Max = DefaultConfig.Recent.Max
	}
	if config.UI.Theme == "" {
		config.UI.Theme = DefaultConfig.UI.Theme
	}
	if config.UI.CursorStyle == "" {
		config.UI.CursorStyle = DefaultConfig.UI.CursorStyle
	}
	if config.Converter.DefaultMethod == "" {
		config.Converter.DefaultMethod = DefaultConfig.Converter.DefaultMethod
	}
}

// AddRecentFile adds a file to the recent files list
func (c *Config) AddRecentFile(path string, frames int) {
	// Remove if already exists
	for i, rf := range c.Recent.Files {
		if rf.Path == path {
			c.Recent.Files = append(c.Recent.Files[:i], c.Recent.Files[i+1:]...)
			break
		}
	}

	// Add to front
	c.Recent.Files = append([]RecentFile{
		{
			Path:      path,
			Timestamp: time.Now(),
			Frames:    frames,
		},
	}, c.Recent.Files...)

	// Trim to max
	if len(c.Recent.Files) > c.Recent.Max {
		c.Recent.Files = c.Recent.Files[:c.Recent.Max]
	}
}

// GetRecentFiles returns the list of recent files
func (c *Config) GetRecentFiles() []RecentFile {
	return c.Recent.Files
}

// GetStartupArtwork returns the custom startup artwork or default
func (c *Config) GetStartupArtwork() string {
	// Try loading from file first
	if c.Startup.ArtworkFile != "" {
		artwork := c.loadArtworkFile(c.Startup.ArtworkFile)
		if artwork != "" {
			return artwork
		}
	}
	
	// Use inline artwork if specified
	if c.Startup.ArtworkInline != "" {
		return c.Startup.ArtworkInline
	}
	
	// Return default logo
	return DefaultStartupArtwork
}

// loadArtworkFile loads artwork from a file, handling different formats
func (c *Config) loadArtworkFile(path string) string {
	// Try absolute path first
	data, err := os.ReadFile(path)
	if err != nil {
		// Try relative to config directory
		dir, err := ConfigDir()
		if err == nil {
			artPath := filepath.Join(dir, path)
			data, err = os.ReadFile(artPath)
		}
		if err != nil {
			return ""
		}
	}
	
	// Check if it's a .aa file (JSON format)
	if strings.HasSuffix(path, ".aa") || strings.HasSuffix(path, ".aart") {
		return extractArtworkFromAA(data)
	}
	
	// Otherwise treat as plain text
	return string(data)
}

// extractArtworkFromAA extracts ASCII art from a .aa file
func extractArtworkFromAA(data []byte) string {
	// Parse the .aa file JSON
	var aart struct {
		Canvas struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"canvas"`
		Frames []struct {
			Cells [][]struct {
				Char string `json:"char"`
			} `json:"cells"`
		} `json:"frames"`
	}
	
	if err := json.Unmarshal(data, &aart); err != nil {
		// If parsing fails, return empty (will fallback to default)
		return ""
	}
	
	// Get first frame
	if len(aart.Frames) == 0 || len(aart.Frames[0].Cells) == 0 {
		return ""
	}
	
	// Build ASCII art from cells
	var lines []string
	for _, row := range aart.Frames[0].Cells {
		var line strings.Builder
		for _, cell := range row {
			if cell.Char == "" {
				line.WriteString(" ")
			} else {
				line.WriteString(cell.Char)
			}
		}
		// Trim trailing spaces
		lineStr := strings.TrimRight(line.String(), " ")
		lines = append(lines, lineStr)
	}
	
	// Remove trailing empty lines
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	
	return strings.Join(lines, "\n")
}

// DefaultStartupArtwork is the default ASCII art logo (enhanced modern version)
const DefaultStartupArtwork = `
    ▄▄▄▄▄▄▄▄▄        ▄▄▄▄▄▄▄▄▄       ██▀███   ▄▄▄█████▓
   ▐████████▌      ▐████████▌     ▓██ ▒ ██▒▓  ██▒ ▓▒
    ▀██████▀        ▀██████▀      ▓██ ░▄█ ▒▒ ▓██░ ▒░
     ▐████▌          ▐████▌       ▒██▀▀█▄  ░ ▓██▓ ░ 
      ▀██▀            ▀██▀        ░██▓ ▒██▒  ▒██▒ ░ 
       ▀▀              ▀▀         ░ ▒▓ ░▒▓░  ▒ ░░   
                                                     
         ┌─────────────────────────────────┐         
         │  Your ASCII Animation Studio    │         
         └─────────────────────────────────┘         `

