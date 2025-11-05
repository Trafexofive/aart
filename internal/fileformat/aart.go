package fileformat

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// AartFile represents the native .aart file format
type AartFile struct {
	Version    string     `json:"version"`
	Metadata   Metadata   `json:"metadata"`
	Canvas     Canvas     `json:"canvas"`
	Frames     []Frame    `json:"frames"`
	Layers     []Layer    `json:"layers,omitempty"`
	Palette    []Color    `json:"palette,omitempty"`
	Audio      *Audio     `json:"audio,omitempty"`
}

// Metadata contains file information
type Metadata struct {
	Title       string    `json:"title"`
	Author      string    `json:"author,omitempty"`
	Description string    `json:"description,omitempty"`
	Created     time.Time `json:"created"`
	Modified    time.Time `json:"modified"`
	Tags        []string  `json:"tags,omitempty"`
	Source      string    `json:"source,omitempty"` // Original GIF URL, etc.
}

// Canvas defines the canvas properties
type Canvas struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Frame represents a single animation frame
type Frame struct {
	Index    int      `json:"index"`
	Duration int      `json:"duration"` // milliseconds
	Cells    [][]Cell `json:"cells"`
	Name     string   `json:"name,omitempty"`
}

// Cell represents a single character cell
type Cell struct {
	Char       string `json:"char"`             // UTF-8 character
	Foreground string `json:"fg"`               // Hex color
	Background string `json:"bg"`               // Hex color
	Bold       bool   `json:"bold,omitempty"`
	Italic     bool   `json:"italic,omitempty"`
	Underline  bool   `json:"underline,omitempty"`
}

// Layer represents a drawing layer
type Layer struct {
	Name      string  `json:"name"`
	Visible   bool    `json:"visible"`
	Opacity   float64 `json:"opacity"`
	BlendMode string  `json:"blend_mode"`
}

// Color represents a palette color
type Color struct {
	Name string `json:"name"`
	Hex  string `json:"hex"`
}

// Audio represents optional audio track
type Audio struct {
	Path   string `json:"path,omitempty"`
	Offset int    `json:"offset,omitempty"` // milliseconds
	Loop   bool   `json:"loop"`
}

// Load loads a .aart file
func Load(path string) (*AartFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var aart AartFile
	if err := json.Unmarshal(data, &aart); err != nil {
		return nil, fmt.Errorf("failed to parse .aart file: %w", err)
	}

	return &aart, nil
}

// Save saves a .aart file
func Save(path string, aart *AartFile) error {
	// Update modified timestamp
	aart.Metadata.Modified = time.Now()

	// Pretty print JSON
	data, err := json.MarshalIndent(aart, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal .aart file: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// NewAartFile creates a new .aart file structure
func NewAartFile(width, height int, title string) *AartFile {
	now := time.Now()
	return &AartFile{
		Version: "1.0",
		Metadata: Metadata{
			Title:    title,
			Created:  now,
			Modified: now,
		},
		Canvas: Canvas{
			Width:  width,
			Height: height,
		},
		Frames: []Frame{},
		Layers: []Layer{
			{Name: "Background", Visible: true, Opacity: 1.0, BlendMode: "normal"},
			{Name: "Foreground", Visible: true, Opacity: 1.0, BlendMode: "normal"},
		},
	}
}

// AddFrame adds a frame to the file
func (a *AartFile) AddFrame(cells [][]Cell, duration int) {
	frame := Frame{
		Index:    len(a.Frames),
		Duration: duration,
		Cells:    cells,
	}
	a.Frames = append(a.Frames, frame)
}

// GetFrame returns a frame by index
func (a *AartFile) GetFrame(index int) (*Frame, error) {
	if index < 0 || index >= len(a.Frames) {
		return nil, fmt.Errorf("frame index out of range: %d", index)
	}
	return &a.Frames[index], nil
}

// FrameCount returns the number of frames
func (a *AartFile) FrameCount() int {
	return len(a.Frames)
}

// TotalDuration returns total animation duration in milliseconds
func (a *AartFile) TotalDuration() int {
	total := 0
	for _, frame := range a.Frames {
		total += frame.Duration
	}
	return total
}

// Validate checks if the file is valid
func (a *AartFile) Validate() error {
	if a.Canvas.Width <= 0 || a.Canvas.Height <= 0 {
		return fmt.Errorf("invalid canvas dimensions: %dx%d", a.Canvas.Width, a.Canvas.Height)
	}

	for i, frame := range a.Frames {
		if len(frame.Cells) != a.Canvas.Height {
			return fmt.Errorf("frame %d: wrong height %d, expected %d", i, len(frame.Cells), a.Canvas.Height)
		}
		for y, row := range frame.Cells {
			if len(row) != a.Canvas.Width {
				return fmt.Errorf("frame %d, row %d: wrong width %d, expected %d", i, y, len(row), a.Canvas.Width)
			}
		}
	}

	return nil
}
