package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/pinchtab/pinchtab/internal/cliui"
	"github.com/pinchtab/pinchtab/internal/config"
)

type startupBannerOptions struct {
	Mode       string
	ListenAddr string
	PublicURL  string
	Strategy   string
	Allocation string
	ProfileDir string
}

func printStartupBanner(cfg *config.RuntimeConfig, opts startupBannerOptions) {
	writeBannerLine(renderStartupLogo(blankIfEmpty(opts.Mode, "server")))
	writeBannerf("  %s  %s\n", styleLabel("listen"), styleValue(blankIfEmpty(opts.ListenAddr, cfg.ListenAddr())))
	if opts.PublicURL != "" {
		writeBannerf("  %s  %s\n", styleLabel("url"), styleValue(opts.PublicURL))
	}
	strat := blankIfEmpty(opts.Strategy, "manual")
	alloc := blankIfEmpty(opts.Allocation, "none")
	writeBannerf("  %s  %s\n", styleLabel("str,plc"), styleValue(fmt.Sprintf("%s,%s", strat, alloc)))

	daemonStatus := styleStdout(cliWarningStyle, "not installed")
	if IsDaemonInstalled() {
		daemonStatus = styleStdout(cliSuccessStyle, "ok")
	}
	writeBannerf("  %s  %s\n", styleLabel("daemon"), daemonStatus)

	if opts.ProfileDir != "" {
		writeBannerf("  %s  %s\n", styleLabel("profile"), styleValue(opts.ProfileDir))
	}
	printSecuritySummary(os.Stdout, cfg, "  ")
	writeBannerLine("")
}

func printSecuritySummary(w io.Writer, cfg *config.RuntimeConfig, prefix string) {
	posture := assessSecurityPosture(cfg)

	writeSummaryf(
		w,
		"%s%s  %s  %s\n",
		prefix,
		styleLabel("security"),
		styleSecurityLevel(posture.Level),
		styleSecurityBar(posture.Level, renderPostureBar(posture.Passed, posture.Total)),
	)
	for _, check := range posture.Checks {
		writeSummaryf(
			w,
			"%s  [%s] %s %s\n",
			prefix,
			styleMarker(check.Passed),
			styleCheckLabel(check.Label),
			styleCheckDetail(check.Passed, check.Detail),
		)
	}
}

func renderPostureBar(passed, total int) string {
	if total <= 0 {
		return "[--------]"
	}
	var b strings.Builder
	b.WriteByte('[')
	for i := 0; i < total; i++ {
		if i < passed {
			b.WriteByte('#')
		} else {
			b.WriteByte('-')
		}
	}
	b.WriteByte(']')
	return b.String()
}

func blankIfEmpty(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func renderStartupLogo(mode string) string {
	return styleLogo(startupLogo) + "  " + styleMode(mode)
}

func writeBannerLine(line string) {
	_, _ = fmt.Fprintln(os.Stdout, line)
}

func writeBannerf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stdout, format, args...)
}

func writeSummaryf(w io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(w, format, args...)
}

func styleLogo(text string) string {
	return applyStyle(text, lipgloss.NewStyle().Foreground(cliui.ColorAccent).Bold(true))
}

func styleMode(text string) string {
	return applyStyle(text, lipgloss.NewStyle().Foreground(cliui.ColorTextMuted))
}

func styleLabel(text string) string {
	return applyStyle(fmt.Sprintf("%-8s", text), lipgloss.NewStyle().Foreground(cliui.ColorTextMuted))
}

func styleValue(text string) string {
	return applyStyle(text, lipgloss.NewStyle().Foreground(cliui.ColorTextPrimary).Bold(true))
}

func styleCheckLabel(text string) string {
	return applyStyle(fmt.Sprintf("%-20s", text), lipgloss.NewStyle().Foreground(cliui.ColorTextMuted))
}

func styleCheckDetail(passed bool, text string) string {
	if passed {
		return applyStyle(text, lipgloss.NewStyle().Foreground(cliui.ColorSuccess))
	}
	return applyStyle(text, lipgloss.NewStyle().Foreground(cliui.ColorWarning))
}

func styleMarker(passed bool) string {
	if passed {
		return applyStyle("ok", lipgloss.NewStyle().Foreground(cliui.ColorSuccess).Bold(true))
	}
	return applyStyle("!!", lipgloss.NewStyle().Foreground(cliui.ColorDanger).Bold(true))
}

func styleSecurityLevel(level string) string {
	return applyStyle(level, lipgloss.NewStyle().Foreground(lipgloss.Color(securityLevelColor(level))).Bold(true))
}

func styleSecurityBar(level, bar string) string {
	return applyStyle(bar, lipgloss.NewStyle().Foreground(lipgloss.Color(securityLevelColor(level))).Bold(true))
}

func securityLevelColor(level string) string {
	switch level {
	case "LOCKED":
		return string(cliui.ColorSuccess)
	case "GUARDED":
		return string(cliui.ColorWarning)
	case "ELEVATED":
		return string(cliui.ColorDanger)
	default:
		return string(cliui.ColorDanger)
	}
}

func applyStyle(text string, style lipgloss.Style) string {
	return cliui.RenderStdout(style, text)
}

const startupLogo = `   ____  _            _     _____     _
  |  _ \(_)_ __   ___| |__ |_   _|_ _| |__
  | |_) | | '_ \ / __| '_ \  | |/ _  | '_ \
  |  __/| | | | | (__| | | | | | (_| | |_) |
  |_|   |_|_| |_|\___|_| |_| |_|\__,_|_.__/`
