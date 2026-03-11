package main

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/pinchtab/pinchtab/internal/cliui"
)

var (
	pinchtabColorBorder      = cliui.ColorBorder
	pinchtabColorTextPrimary = cliui.ColorTextPrimary
	pinchtabColorTextMuted   = cliui.ColorTextMuted
	pinchtabColorAccent      = cliui.ColorAccent
	pinchtabColorAccentLight = cliui.ColorAccentLight
	pinchtabColorSuccess     = cliui.ColorSuccess
	pinchtabColorWarning     = cliui.ColorWarning
	pinchtabColorDanger      = cliui.ColorDanger
)

var (
	cliHeadingStyle = cliui.HeadingStyle
	cliCommandStyle = cliui.CommandStyle
	cliMutedStyle   = cliui.MutedStyle
	cliSuccessStyle = cliui.SuccessStyle
	cliWarningStyle = cliui.WarningStyle
	cliErrorStyle   = cliui.ErrorStyle
	cliValueStyle   = cliui.ValueStyle
)

func renderToWriter(w *os.File, style lipgloss.Style, text string) string {
	if w == os.Stdout {
		return cliui.RenderStdout(style, text)
	}
	if w == os.Stderr {
		return cliui.RenderStderr(style, text)
	}
	return style.Render(text)
}

func styleStdout(style lipgloss.Style, text string) string {
	return renderToWriter(os.Stdout, style, text)
}

func styleStderr(style lipgloss.Style, text string) string {
	return renderToWriter(os.Stderr, style, text)
}
