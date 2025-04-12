package terminal

import "github.com/charmbracelet/lipgloss"

var styleDebug = lipgloss.NewStyle().
	Faint(true).
	Foreground(lipgloss.Color("11"))

var styleFaint = lipgloss.NewStyle().Faint(true)

var styleCode = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))

var styleErr = lipgloss.NewStyle().Foreground(lipgloss.Color("#ed4747"))

var styleNestedContainer = lipgloss.NewStyle().PaddingTop(2).PaddingLeft(2)
