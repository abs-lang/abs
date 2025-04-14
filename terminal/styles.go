package terminal

import "github.com/charmbracelet/lipgloss"

// ANSI color codes
// https://raw.githubusercontent.com/fidian/ansi/master/images/color-codes.png

var styleDebug = lipgloss.NewStyle().
	Faint(true).
	Foreground(lipgloss.Color("11"))

var styleFaint = lipgloss.NewStyle().Faint(true)

var styleCode = lipgloss.NewStyle().Foreground(lipgloss.Color("178"))

var styleErr = lipgloss.NewStyle().Foreground(lipgloss.Color("#ed4747"))

var styleNestedContainer = lipgloss.NewStyle().PaddingTop(2).PaddingLeft(2)

var styleSuggestion = lipgloss.NewStyle().PaddingLeft(1).PaddingTop(1)
var styleSuggestions = map[suggestionType]lipgloss.Style{
	SUGGESTION_FUNCTION:   lipgloss.NewStyle().Foreground(lipgloss.Color("4")),
	SUGGESTION_IDENTIFIER: lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
	SUGGESTION_PROPERTY:   lipgloss.NewStyle().Foreground(lipgloss.Color("6")),
}
var styleSelectedSuggestion = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Underline(true)
var styleSelectedPrefix = styleSelectedSuggestion.Underline(false)

var styleSearch = styleSuggestion
var styleSearchPrompt = lipgloss.NewStyle().Foreground(lipgloss.Color("178")).Faint(true)
var styleSearchText = styleCode
