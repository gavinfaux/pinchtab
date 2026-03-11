package main

import (
	"fmt"
	"os"
	"strings"

	prompt "github.com/c-bata/go-prompt"
)

var demoRefs = []prompt.Suggest{
	{Text: "e12", Description: "button Login"},
	{Text: "e18", Description: "input Email"},
	{Text: "e24", Description: "link Forgot password"},
}

var rootSuggestions = []prompt.Suggest{
	{Text: "help", Description: "Show the available commands"},
	{Text: "nav", Description: "Navigate to a URL"},
	{Text: "snap", Description: "Capture a page snapshot"},
	{Text: "click", Description: "Click an element from the last snapshot"},
	{Text: "type", Description: "Type into an element from the last snapshot"},
	{Text: "pdf", Description: "Export the current page as PDF"},
	{Text: "exit", Description: "Quit the demo"},
}

func main() {
	fmt.Println("PinchTab prompt demo")
	fmt.Println("Tab completes commands and a few cached refs from a fake snapshot.")
	fmt.Println("Try: nav https://example.com, snap -i, click e12, type e18 hello")

	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix("pinchtab> "),
		prompt.OptionTitle("PinchTab Prompt Demo"),
	)
	p.Run()
}

func executor(in string) {
	line := strings.TrimSpace(in)
	if line == "" {
		return
	}

	switch {
	case line == "exit" || line == "quit":
		fmt.Println("bye")
		os.Exit(0)
	case line == "help":
		fmt.Println("commands: nav, snap, click, type, pdf, exit")
	case strings.HasPrefix(line, "nav "):
		fmt.Printf("navigate -> %s\n", strings.TrimSpace(strings.TrimPrefix(line, "nav ")))
	case strings.HasPrefix(line, "snap"):
		fmt.Println("snapshot -> cached refs: e12(Login), e18(Email), e24(Forgot password)")
	case strings.HasPrefix(line, "click "):
		fmt.Printf("click -> %s\n", strings.TrimSpace(strings.TrimPrefix(line, "click ")))
	case strings.HasPrefix(line, "type "):
		fmt.Printf("type -> %s\n", strings.TrimSpace(strings.TrimPrefix(line, "type ")))
	case strings.HasPrefix(line, "pdf"):
		fmt.Println("pdf -> /tmp/example.pdf")
	default:
		fmt.Printf("unknown command: %s\n", line)
	}
}

func completer(d prompt.Document) []prompt.Suggest {
	words := strings.Fields(d.TextBeforeCursor())
	if len(words) == 0 {
		return prompt.FilterHasPrefix(rootSuggestions, d.GetWordBeforeCursor(), true)
	}

	switch words[0] {
	case "nav":
		return prompt.FilterHasPrefix([]prompt.Suggest{
			{Text: "https://example.com", Description: "demo URL"},
			{Text: "https://pinchtab.com", Description: "PinchTab website"},
		}, d.GetWordBeforeCursor(), true)
	case "snap":
		return prompt.FilterHasPrefix([]prompt.Suggest{
			{Text: "-i", Description: "interactive elements only"},
			{Text: "-c", Description: "compact format"},
			{Text: "--diff", Description: "only changes since last snapshot"},
		}, d.GetWordBeforeCursor(), true)
	case "click", "type":
		return prompt.FilterHasPrefix(demoRefs, d.GetWordBeforeCursor(), true)
	case "pdf":
		return prompt.FilterHasPrefix([]prompt.Suggest{
			{Text: "-o", Description: "output file"},
			{Text: "--tab", Description: "target tab"},
		}, d.GetWordBeforeCursor(), true)
	default:
		return prompt.FilterHasPrefix(rootSuggestions, d.GetWordBeforeCursor(), true)
	}
}
