package fileformat

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ExportFormat represents supported export formats
type ExportFormat string

const (
	FormatJSON ExportFormat = "json"
	FormatCSV  ExportFormat = "csv"
	FormatANSI ExportFormat = "ansi"
	FormatTXT  ExportFormat = "txt"
	FormatHTML ExportFormat = "html"
	FormatSVG  ExportFormat = "svg"
	FormatGIF  ExportFormat = "gif"
)

// ExportOptions contains export configuration
type ExportOptions struct {
	Format      ExportFormat
	FrameIndex  int  // -1 for all frames
	IncludeMeta bool
	Compact     bool // For JSON
	Colors      bool // For ANSI/TXT
}

// Export exports to the specified format
func Export(aart *AartFile, path string, opts ExportOptions) error {
	switch opts.Format {
	case FormatJSON:
		return exportJSON(aart, path, opts)
	case FormatCSV:
		return exportCSV(aart, path, opts)
	case FormatANSI:
		return exportANSI(aart, path, opts)
	case FormatTXT:
		return exportTXT(aart, path, opts)
	case FormatHTML:
		return exportHTML(aart, path, opts)
	case FormatSVG:
		return exportSVG(aart, path, opts)
	default:
		return fmt.Errorf("unsupported export format: %s", opts.Format)
	}
}

// exportJSON exports to JSON format
func exportJSON(aart *AartFile, path string, opts ExportOptions) error {
	var data []byte
	var err error

	if opts.Compact {
		data, err = json.Marshal(aart)
	} else {
		data, err = json.MarshalIndent(aart, "", "  ")
	}

	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// exportCSV exports frames to CSV format
func exportCSV(aart *AartFile, path string, opts ExportOptions) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Header
	headers := []string{"frame", "x", "y", "char", "fg", "bg"}
	if err := writer.Write(headers); err != nil {
		return err
	}

	// Determine which frames to export
	startFrame := 0
	endFrame := len(aart.Frames)
	if opts.FrameIndex >= 0 && opts.FrameIndex < len(aart.Frames) {
		startFrame = opts.FrameIndex
		endFrame = opts.FrameIndex + 1
	}

	// Write data
	for frameIdx := startFrame; frameIdx < endFrame; frameIdx++ {
		frame := aart.Frames[frameIdx]
		for y, row := range frame.Cells {
			for x, cell := range row {
				record := []string{
					fmt.Sprintf("%d", frameIdx),
					fmt.Sprintf("%d", x),
					fmt.Sprintf("%d", y),
					cell.Char,
					cell.Foreground,
					cell.Background,
				}
				if err := writer.Write(record); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// exportANSI exports to ANSI art format
func exportANSI(aart *AartFile, path string, opts ExportOptions) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Export single frame or first frame
	frameIdx := 0
	if opts.FrameIndex >= 0 && opts.FrameIndex < len(aart.Frames) {
		frameIdx = opts.FrameIndex
	}

	frame := aart.Frames[frameIdx]

	var output strings.Builder

	for _, row := range frame.Cells {
		for _, cell := range row {
			if opts.Colors {
				// ANSI color codes
				fg := hexToANSI(cell.Foreground)
				bg := hexToANSI(cell.Background)
				output.WriteString(fmt.Sprintf("\x1b[38;5;%dm\x1b[48;5;%dm%s", fg, bg, cell.Char))
			} else {
				output.WriteString(cell.Char)
			}
		}
		if opts.Colors {
			output.WriteString("\x1b[0m") // Reset
		}
		output.WriteString("\n")
	}

	_, err = file.WriteString(output.String())
	return err
}

// exportTXT exports to plain text
func exportTXT(aart *AartFile, path string, opts ExportOptions) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Export single frame or first frame
	frameIdx := 0
	if opts.FrameIndex >= 0 && opts.FrameIndex < len(aart.Frames) {
		frameIdx = opts.FrameIndex
	}

	frame := aart.Frames[frameIdx]

	var output strings.Builder

	// Add metadata if requested
	if opts.IncludeMeta {
		output.WriteString(fmt.Sprintf("Title: %s\n", aart.Metadata.Title))
		output.WriteString(fmt.Sprintf("Size: %dx%d\n", aart.Canvas.Width, aart.Canvas.Height))
		output.WriteString(fmt.Sprintf("Frames: %d\n", len(aart.Frames)))
		output.WriteString("\n")
	}

	for _, row := range frame.Cells {
		for _, cell := range row {
			output.WriteString(cell.Char)
		}
		output.WriteString("\n")
	}

	_, err = file.WriteString(output.String())
	return err
}

// exportHTML exports to HTML format
func exportHTML(aart *AartFile, path string, opts ExportOptions) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	frameIdx := 0
	if opts.FrameIndex >= 0 && opts.FrameIndex < len(aart.Frames) {
		frameIdx = opts.FrameIndex
	}

	frame := aart.Frames[frameIdx]

	var html strings.Builder
	
	html.WriteString("<!DOCTYPE html>\n<html>\n<head>\n")
	html.WriteString(fmt.Sprintf("<title>%s</title>\n", aart.Metadata.Title))
	html.WriteString("<style>\n")
	html.WriteString("body { background: #000; margin: 20px; }\n")
	html.WriteString("pre { font-family: monospace; line-height: 1; margin: 0; }\n")
	html.WriteString(".cell { display: inline-block; }\n")
	html.WriteString("</style>\n")
	html.WriteString("</head>\n<body>\n<pre>\n")

	for _, row := range frame.Cells {
		for _, cell := range row {
			if opts.Colors {
				html.WriteString(fmt.Sprintf(
					"<span class=\"cell\" style=\"color:%s;background:%s\">%s</span>",
					cell.Foreground,
					cell.Background,
					escapeHTML(cell.Char),
				))
			} else {
				html.WriteString(escapeHTML(cell.Char))
			}
		}
		html.WriteString("\n")
	}

	html.WriteString("</pre>\n</body>\n</html>")

	_, err = file.WriteString(html.String())
	return err
}

// exportSVG exports to SVG format
func exportSVG(aart *AartFile, path string, opts ExportOptions) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	frameIdx := 0
	if opts.FrameIndex >= 0 && opts.FrameIndex < len(aart.Frames) {
		frameIdx = opts.FrameIndex
	}

	frame := aart.Frames[frameIdx]

	cellWidth := 10
	cellHeight := 16
	width := aart.Canvas.Width * cellWidth
	height := aart.Canvas.Height * cellHeight

	var svg strings.Builder
	svg.WriteString(fmt.Sprintf("<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"%d\" height=\"%d\">\n", width, height))
	svg.WriteString("<style>text { font-family: monospace; font-size: 14px; }</style>\n")

	for y, row := range frame.Cells {
		for x, cell := range row {
			px := x * cellWidth
			py := (y + 1) * cellHeight
			
			// Background rect
			if cell.Background != "#000000" && cell.Background != "" {
				svg.WriteString(fmt.Sprintf(
					"<rect x=\"%d\" y=\"%d\" width=\"%d\" height=\"%d\" fill=\"%s\"/>\n",
					px, py-cellHeight, cellWidth, cellHeight, cell.Background,
				))
			}
			
			// Text
			svg.WriteString(fmt.Sprintf(
				"<text x=\"%d\" y=\"%d\" fill=\"%s\">%s</text>\n",
				px, py, cell.Foreground, escapeHTML(cell.Char),
			))
		}
	}

	svg.WriteString("</svg>")

	_, err = file.WriteString(svg.String())
	return err
}

// Helper functions

func hexToANSI(hex string) int {
	// Simplified hex to ANSI 256 color conversion
	// For a full implementation, parse hex and map to closest ANSI color
	if hex == "#FFFFFF" || hex == "#ffffff" {
		return 15 // White
	}
	if hex == "#000000" || hex == "#000000" {
		return 0 // Black
	}
	// Default to white
	return 15
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}
