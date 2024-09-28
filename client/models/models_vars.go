package models

import (
	"errors"

	"github.com/charmbracelet/lipgloss"
)

var (
	// раскрашиваем текст
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#aaffaa"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#111222"))
	cursorStyle  = focusedStyle
	noStyle      = lipgloss.NewStyle()
	helpStyle    = blurredStyle

	styleBlue  = lipgloss.NewStyle().Foreground(lipgloss.Color("#1d3cdb"))
	styleRed   = lipgloss.NewStyle().Foreground(lipgloss.Color("#c40202"))
	styleError = lipgloss.NewStyle().Foreground(lipgloss.Color("#870707"))

	// стандартные ошибки
	ErrTokenUndefined = errors.New("токен не определен")
)
