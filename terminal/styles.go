package terminal

import "github.com/charmbracelet/lipgloss"

var styleDebug = lipgloss.NewStyle().
	Faint(true).
	Foreground(lipgloss.Color("11"))

var styleFaint = lipgloss.NewStyle().Faint(true)

var styleCode = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))

var styleErr = lipgloss.NewStyle().Foreground(lipgloss.Color("#ed4747"))

var styleNestedContainer = lipgloss.NewStyle().PaddingTop(2).PaddingLeft(2)

var styleSuggestion = lipgloss.NewStyle().PaddingLeft(1).PaddingTop(1).Foreground(lipgloss.Color("2"))
var styleSelectedSuggestion = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Underline(true)
var styleSelectedPrefix = styleSelectedSuggestion.Underline(false)
