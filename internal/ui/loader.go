package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// LoaderModel shows a loading screen with progress
type LoaderModel struct {
	spinner  spinner.Model
	progress progress.Model
	message  string
	current  int
	total    int
	done     bool
	err      error
}

type progressMsg struct {
	current int
	total   int
	message string
}

type doneMsg struct {
	err error
}

func NewLoader(message string) LoaderModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	p := progress.New(progress.WithDefaultGradient())
	
	return LoaderModel{
		spinner: s,
		progress: p,
		message: message,
	}
}

func (m LoaderModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m LoaderModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		
	case progressMsg:
		m.current = msg.current
		m.total = msg.total
		if msg.message != "" {
			m.message = msg.message
		}
		return m, nil
		
	case doneMsg:
		m.done = true
		m.err = msg.err
		return m, tea.Quit
		
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	
	return m, nil
}

func (m LoaderModel) View() string {
	if m.done {
		if m.err != nil {
			return lipgloss.NewStyle().
				Foreground(lipgloss.Color("196")).
				Render(fmt.Sprintf("\n✗ Error: %v\n", m.err))
		}
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Render("\n✓ Done!\n")
	}

	var s strings.Builder
	
	s.WriteString("\n")
	s.WriteString(m.spinner.View())
	s.WriteString(" ")
	s.WriteString(m.message)
	s.WriteString("\n\n")
	
	if m.total > 0 {
		percent := float64(m.current) / float64(m.total)
		s.WriteString(m.progress.ViewAs(percent))
		s.WriteString(fmt.Sprintf(" %d/%d", m.current, m.total))
		s.WriteString("\n")
	}
	
	s.WriteString("\n")
	
	return s.String()
}

// Progress messages for GIF conversion
func SendProgress(current, total int, message string) tea.Cmd {
	return func() tea.Msg {
		return progressMsg{current: current, total: total, message: message}
	}
}

func SendDone(err error) tea.Cmd {
	return func() tea.Msg {
		return doneMsg{err: err}
	}
}

// ConversionProgress wraps a conversion with progress reporting
type ConversionProgress struct {
	Current int
	Total   int
	Message string
}

// ProgressReporter is a callback for progress updates
type ProgressReporter func(current, total int, message string)
