package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pinchtab/pinchtab/internal/config"
	"github.com/spf13/cobra"
)

var quickCmd = &cobra.Command{
	Use:   "quick <url>",
	Short: "Navigate + analyze page (beginner-friendly)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		runCLIWith(cfg, func(client *http.Client, base, token string) {
			cliQuick(client, base, token, args)
		})
	},
}

var navCmd = &cobra.Command{
	Use:   "nav <url>",
	Short: "Navigate to URL",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		runCLIWith(cfg, func(client *http.Client, base, token string) {
			cliNavigate(client, base, token, args)
		})
	},
}

var snapCmd = &cobra.Command{
	Use:   "snap",
	Short: "Snapshot accessibility tree",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		runCLIWith(cfg, func(client *http.Client, base, token string) {
			cliSnapshot(client, base, token, args)
		})
	},
}

var clickCmd = &cobra.Command{
	Use:   "click <ref>",
	Short: "Click element",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		runCLIWith(cfg, func(client *http.Client, base, token string) {
			cliTabOperation(client, base, token, "click", args)
		})
	},
}

var typeCmd = &cobra.Command{
	Use:   "type <ref> <text>",
	Short: "Type into element",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		runCLIWith(cfg, func(client *http.Client, base, token string) {
			cliTabOperation(client, base, token, "type", args)
		})
	},
}

var screenshotCmd = &cobra.Command{
	Use:   "screenshot",
	Short: "Take a screenshot",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		runCLIWith(cfg, func(client *http.Client, base, token string) {
			cliScreenshot(client, base, token, args)
		})
	},
}

var tabsCmd = &cobra.Command{
	Use:   "tabs",
	Short: "List or manage tabs",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		runCLIWith(cfg, func(client *http.Client, base, token string) {
			cliTabs(client, base, token, args)
		})
	},
}

var instancesCmd = &cobra.Command{
	Use:   "instances",
	Short: "List or manage instances",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		runCLIWith(cfg, func(client *http.Client, base, token string) {
			cliInstances(client, base, token)
		})
	},
}

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check server health",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		runCLIWith(cfg, func(client *http.Client, base, token string) {
			cliHealth(client, base, token)
		})
	},
}

var pressCmd = &cobra.Command{
	Use:   "press <key>",
	Short: "Press key (Enter, Tab, Escape...)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		runCLIWith(cfg, func(client *http.Client, base, token string) {
			cliAction(client, base, token, "press", args)
		})
	},
}

var fillCmd = &cobra.Command{
	Use:   "fill <ref|selector> <text>",
	Short: "Fill input directly",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		runCLIWith(cfg, func(client *http.Client, base, token string) {
			cliAction(client, base, token, "fill", args)
		})
	},
}

var hoverCmd = &cobra.Command{
	Use:   "hover <ref>",
	Short: "Hover element",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		runCLIWith(cfg, func(client *http.Client, base, token string) {
			cliAction(client, base, token, "hover", args)
		})
	},
}

var scrollCmd = &cobra.Command{
	Use:   "scroll <ref|pixels>",
	Short: "Scroll to element or by pixels",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		runCLIWith(cfg, func(client *http.Client, base, token string) {
			cliAction(client, base, token, "scroll", args)
		})
	},
}

var evalCmd = &cobra.Command{
	Use:   "eval <expression>",
	Short: "Evaluate JavaScript",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		runCLIWith(cfg, func(client *http.Client, base, token string) {
			cliEvaluate(client, base, token, args)
		})
	},
}

var pdfCmd = &cobra.Command{
	Use:   "pdf",
	Short: "Export the current page as PDF",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		runCLIWith(cfg, func(client *http.Client, base, token string) {
			cliPDF(client, base, token, args)
		})
	},
}

var textCmd = &cobra.Command{
	Use:   "text",
	Short: "Extract page text",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		runCLIWith(cfg, func(client *http.Client, base, token string) {
			cliText(client, base, token, args)
		})
	},
}

var profilesCmd = &cobra.Command{
	Use:   "profiles",
	Short: "List browser profiles",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		runCLIWith(cfg, func(client *http.Client, base, token string) {
			cliProfiles(client, base, token)
		})
	},
}

var instanceCmd = &cobra.Command{
	Use:   "instance",
	Short: "Manage browser instances",
}

func init() {
	rootCmd.AddCommand(quickCmd)
	rootCmd.AddCommand(navCmd)
	rootCmd.AddCommand(snapCmd)
	rootCmd.AddCommand(clickCmd)
	rootCmd.AddCommand(typeCmd)
	rootCmd.AddCommand(screenshotCmd)
	rootCmd.AddCommand(tabsCmd)
	rootCmd.AddCommand(instancesCmd)
	rootCmd.AddCommand(healthCmd)
	rootCmd.AddCommand(pressCmd)
	rootCmd.AddCommand(fillCmd)
	rootCmd.AddCommand(hoverCmd)
	rootCmd.AddCommand(scrollCmd)
	rootCmd.AddCommand(evalCmd)
	rootCmd.AddCommand(pdfCmd)
	rootCmd.AddCommand(textCmd)
	rootCmd.AddCommand(profilesCmd)

	instanceCmd.AddCommand(&cobra.Command{
		Use:   "start <name>",
		Short: "Start a browser instance",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cfg := config.Load()
			runCLIWith(cfg, func(client *http.Client, base, token string) {
				cliInstanceStart(client, base, token, args)
			})
		},
	})
	instanceCmd.AddCommand(&cobra.Command{
		Use:   "navigate <id> <url>",
		Short: "Navigate an instance to a URL",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			cfg := config.Load()
			runCLIWith(cfg, func(client *http.Client, base, token string) {
				cliInstanceNavigate(client, base, token, args)
			})
		},
	})
	instanceCmd.AddCommand(&cobra.Command{
		Use:   "stop <id>",
		Short: "Stop a browser instance",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cfg := config.Load()
			runCLIWith(cfg, func(client *http.Client, base, token string) {
				cliInstanceStop(client, base, token, args)
			})
		},
	})
	instanceCmd.AddCommand(&cobra.Command{
		Use:   "logs <id>",
		Short: "Get instance logs",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cfg := config.Load()
			runCLIWith(cfg, func(client *http.Client, base, token string) {
				cliInstanceLogs(client, base, token, args)
			})
		},
	})
	rootCmd.AddCommand(instanceCmd)
}

func runCLIWith(cfg *config.RuntimeConfig, fn func(client *http.Client, base, token string)) {
	client := &http.Client{Timeout: 60 * time.Second}
	dashPort := cfg.Port
	if dashPort == "" {
		dashPort = "9870"
	}
	base := fmt.Sprintf("http://localhost:%s", dashPort)
	token := cfg.Token

	fn(client, base, token)
}

func printHelp() {
	fmt.Printf(`pinchtab %s - Browser control for AI agents

Usage:
  pinchtab [command] [flags]

Primary commands (most users start here):
  pinchtab                  Start the full server (default: headed or headless via env)
                            -> Launches Chrome instances + HTTP API on http://localhost:9867
  pinchtab onboard          Guided interactive setup (recommended first run)
                            -> Configures profiles, daemon, stealth, and security defaults

OTHER COMMANDS:
  pinchtab server           Start full server explicitly
  pinchtab bridge           Start single-instance bridge-only server
<<<<<<< HEAD
  pinchtab mcp              Start MCP (Model Context Protocol) server on stdio
=======
  pinchtab daemon           Manage the background service
>>>>>>> 97fee9e (refator: splitting big files based on responsibility)
  pinchtab connect <name>   Get URL for a running profile
  pinchtab security         Review runtime security posture

QUICK START (requires running server):
  pinchtab quick <url>                  Navigate + analyze page (beginner-friendly)
  pinchtab onboard --install-daemon     Guided setup + install a user daemon

CLI COMMANDS:
  pinchtab nav <url>                    Navigate to URL
  pinchtab snap [-i] [-c] [-d]         Snapshot accessibility tree
  pinchtab click <ref>                  Click element
  pinchtab type <ref> <text>            Type into element
  pinchtab press <key>                  Press key (Enter, Tab, Escape...)
  pinchtab fill <ref|selector> <text>   Fill input directly
  pinchtab hover <ref>                  Hover element
  pinchtab scroll <ref|pixels>          Scroll to element or by pixels
  pinchtab select <ref> <value>         Select dropdown option
  pinchtab focus <ref>                  Focus element
  pinchtab text [--raw]                 Extract readable text
  pinchtab tabs [new <url>|close <id>]  Manage tabs
  pinchtab ss [-o file] [-q 80]         Screenshot
  pinchtab eval <expression>            Run JavaScript
  pinchtab pdf [-o file] [--tab <id>] [options]  Export page as PDF (see PDF FLAGS)
  pinchtab instances                    List running instances
  pinchtab profiles                     List available profiles
  pinchtab health                       Check server status

SNAPSHOT FLAGS:
  -i, --interactive    Interactive elements only (buttons, links, inputs)
  -c, --compact        Compact format (most token-efficient)
  -d, --diff           Only changes since last snapshot
  -s, --selector CSS   Scope to CSS selector
  --max-tokens N       Truncate to ~N tokens
  --depth N            Max tree depth
  --tab ID             Target specific tab

PDF FLAGS:
  -o, --output FILE          Output filename (default: page-{timestamp}.pdf)
  --landscape                Landscape orientation
  --paper-width N            Paper width in inches (default: 8.5)
  --paper-height N           Paper height in inches (default: 11)
  --margin-top N             Top margin in inches (default: 0.4)
  --margin-bottom N          Bottom margin in inches (default: 0.4)
  --margin-left N            Left margin in inches (default: 0.4)
  --margin-right N           Right margin in inches (default: 0.4)
  --scale N                  Print scale 0.1-2.0 (default: 1.0)
  --page-ranges RANGE        Pages to export (e.g., "1-3,5")
  --prefer-css-page-size     Honor CSS @page size
  --display-header-footer    Show header and footer
  --header-template HTML     HTML template for header
  --footer-template HTML     HTML template for footer
  --generate-tagged-pdf      Generate accessible/tagged PDF
  --generate-document-outline  Embed document outline
  --file-output              Save to disk (server-side)
  --path PATH                Custom file path (with --file-output)
  --tab ID                   Target specific tab

ENVIRONMENT:
  PINCHTAB_URL         Server URL (default: http://127.0.0.1:9867)
  PINCHTAB_TOKEN       Auth token
  PINCHTAB_PORT        Server port (default: 9867)

FLAGS (global, place before or after command):
  --instance <id>      Target a specific instance by ID (e.g., pinchtab nav --instance abc123 https://...)
  -I <id>              Alias for --instance

Examples:
  pinchtab nav https://pinchtab.com
  pinchtab snap -i -c
  pinchtab click e5
  pinchtab type e12 hello world
  pinchtab press Enter
  pinchtab text | jq .text
  pinchtab eval "document.title"
`, version)
}

var cliCommands = map[string]bool{
	"nav": true, "navigate": true,
	"snap": true, "snapshot": true,
	"click": true, "type": true, "press": true, "fill": true,
	"hover": true, "scroll": true, "select": true, "focus": true,
	"text": true, "tabs": true, "tab": true,
	"screenshot": true, "ss": true,
	"eval": true, "evaluate": true,
	"pdf": true, "health": true,
	"help": true, "quick": true,
	"instance": true, "instances": true,
	"profiles": true,
}

func isCLICommand(cmd string) bool {
	return cliCommands[cmd]
}

func runCLI(cfg *config.RuntimeConfig) {
	cmd := os.Args[1]
	rawArgs := os.Args[2:]

	// Extract --instance/-I flag and strip it from args so sub-commands don't see it.
	var instanceID string
	args := make([]string, 0, len(rawArgs))
	for i := 0; i < len(rawArgs); i++ {
		if (rawArgs[i] == "--instance" || rawArgs[i] == "-I") && i+1 < len(rawArgs) {
			instanceID = rawArgs[i+1]
			i++ // skip the value
		} else {
			args = append(args, rawArgs[i])
		}
	}

	orchBase := fmt.Sprintf("http://%s:%s", cfg.Bind, cfg.Port)
	if envURL := os.Getenv("PINCHTAB_URL"); envURL != "" {
		orchBase = strings.TrimRight(envURL, "/")
	}

	token := cfg.Token
	if envToken := os.Getenv("PINCHTAB_TOKEN"); envToken != "" {
		token = envToken
	}

	// --instance resolves the target base URL from the named instance's port.
	base := orchBase
	if instanceID != "" {
		base = resolveInstanceBase(orchBase, token, instanceID, cfg.Bind)
	}

	client := &http.Client{Timeout: 30 * time.Second}

	// Check if server is running (except for help)
	if cmd != "help" {
		if !checkServerAndGuide(client, base, token) {
			return
		}
	}

	switch cmd {
	case "nav", "navigate":
		cliNavigate(client, base, token, args)
	case "snap", "snapshot":
		cliSnapshot(client, base, token, args)
	case "click", "type", "press", "fill", "hover", "scroll", "select", "focus":
		cliAction(client, base, token, cmd, args)
	case "text":
		cliText(client, base, token, args)
	case "tabs", "tab":
		cliTabs(client, base, token, args)
	case "screenshot", "ss":
		cliScreenshot(client, base, token, args)
	case "eval", "evaluate":
		cliEvaluate(client, base, token, args)
	case "pdf":
		cliPDF(client, base, token, args)
	case "health":
		cliHealth(client, base, token)
	case "instance":
		cliInstance(client, base, token, args)
	case "instances":
		cliInstances(client, base, token)
	case "profiles":
		cliProfiles(client, base, token)
	case "quick":
		cliQuick(client, base, token, args)
	case "help":
		cliHelp()
	}
}

func cliHelp() {
	fmt.Print(`Pinchtab CLI - browser control from the command line

Usage: pinchtab <command> [args] [flags]

QUICK START:
  pinchtab quick <url>    Navigate and show page structure (combines nav + snap)

WORKFLOW:
  1. Start server:        pinchtab                  (or: pinchtab server)
  2. Navigate:           pinchtab nav https://pinchtab.com
  3. See page:           pinchtab snap             (shows clickable refs)
  4. Interact:           pinchtab click e5         (click element)
  5. Check result:       pinchtab snap             (see changes)

Commands:
  quick <url>             Navigate and analyze page (beginner-friendly)

  INSTANCE MANAGEMENT:
  instance launch         Create new instance (--mode headed, --port 9999)
  instance logs <id>      Get instance logs (for debugging)
  instance stop <id>      Stop instance
  instances               List all running instances

  BROWSER CONTROL:
  nav, navigate <url>     Navigate to URL (--new-tab, --block-images, --block-ads)
  snap, snapshot          Accessibility tree snapshot (-i, -c, -d, --max-tokens N)
  click <ref>             Click element by ref
  type <ref> <text>       Type text into element
  fill <ref> <text>       Set input value (no key events)
  press <key>             Press a key (Enter, Tab, Escape, ...)
  hover <ref>             Hover over element
  scroll <direction>      Scroll page (up, down, left, right)
  select <ref> <value>    Select dropdown option
  focus <ref>             Focus element
  text                    Extract page text (--raw for innerText)
  tabs                    List open tabs
  tabs new <url>          Open new tab
  tabs close <tabId>      Close tab
  ss, screenshot          Take screenshot (-o file, -q quality)
  eval <expression>       Evaluate JavaScript
  pdf                     Export page as PDF (-o file, --landscape, --scale N)

  OTHER:
  health                  Server health check
  help                    Show this help

Environment:
  PINCHTAB_URL            Server URL (default: http://localhost:9867)
  PINCHTAB_TOKEN          Auth token (sent as Bearer)

Flags (global):
  --instance <id>, -I <id>  Target a specific instance (e.g., pinchtab snap --instance abc123)

Pipe with jq:
  pinchtab snap -i | jq '.nodes[] | select(.role=="link")'
`)
	os.Exit(0)
}

// --- navigate ---

func cliNavigate(client *http.Client, base, token string, args []string) {
	if len(args) < 1 {
		fatal("Usage: pinchtab nav <url> [--new-tab] [--block-images] [--block-ads]")
	}
	body := map[string]any{"url": args[0]}
	for _, a := range args[1:] {
		switch a {
		case "--new-tab":
			body["newTab"] = true
		case "--block-images":
			body["blockImages"] = true
		case "--block-ads":
			body["blockAds"] = true
		}
	}
	result := doPost(client, base, token, "/navigate", body)
	suggestNextAction("navigate", result)
}

// --- snapshot ---

func cliSnapshot(client *http.Client, base, token string, args []string) {
	params := url.Values{}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--interactive", "-i":
			params.Set("filter", "interactive")
		case "--compact", "-c":
			params.Set("format", "compact")
		case "--text":
			params.Set("format", "text")
		case "--diff", "-d":
			params.Set("diff", "true")
		case "--selector", "-s":
			if i+1 < len(args) {
				i++
				params.Set("selector", args[i])
			}
		case "--max-tokens":
			if i+1 < len(args) {
				i++
				params.Set("maxTokens", args[i])
			}
		case "--depth":
			if i+1 < len(args) {
				i++
				params.Set("depth", args[i])
			}
		case "--tab":
			if i+1 < len(args) {
				i++
				params.Set("tabId", args[i])
			}
		}
	}
	result := doGet(client, base, token, "/snapshot", params)
	suggestNextAction("snapshot", result)
}

// --- element actions ---

func cliAction(client *http.Client, base, token, kind string, args []string) {
	body := map[string]any{"kind": kind}

	switch kind {
	case "click", "hover", "focus":
		var cssSelector string
		var refArg string
		for i := 0; i < len(args); i++ {
			switch args[i] {
			case "--css":
				if i+1 < len(args) {
					i++
					cssSelector = args[i]
				}
			case "--wait-nav":
				body["waitNav"] = true
			default:
				if refArg == "" {
					refArg = args[i]
				}
			}
		}
		if cssSelector != "" {
			body["selector"] = cssSelector
		} else if refArg != "" {
			body["ref"] = refArg
		} else {
			fatal("Usage: pinchtab %s <ref> [--css <selector>] [--wait-nav]", kind)
		}
	case "type":
		if len(args) < 2 {
			fatal("Usage: pinchtab type <ref> <text>")
		}
		body["ref"] = args[0]
		body["text"] = strings.Join(args[1:], " ")
	case "fill":
		if len(args) < 2 {
			fatal("Usage: pinchtab fill <ref|selector> <text>")
		}
		if strings.HasPrefix(args[0], "e") {
			body["ref"] = args[0]
		} else {
			body["selector"] = args[0]
		}
		body["text"] = strings.Join(args[1:], " ")
	case "press":
		if len(args) < 1 {
			fatal("Usage: pinchtab press <key>  (e.g. Enter, Tab, Escape)")
		}
		body["key"] = args[0]
	case "scroll":
		if len(args) < 1 {
			fatal("Usage: pinchtab scroll <ref|pixels|direction>  (e.g. e5, 800, down, up)")
		}
		if strings.HasPrefix(args[0], "e") {
			body["ref"] = args[0]
		} else if v, err := strconv.Atoi(args[0]); err == nil {
			body["scrollY"] = v
		} else {
			// direction: down, up, etc.
			body["direction"] = args[0]
		}
	case "select":
		if len(args) < 2 {
			fatal("Usage: pinchtab select <ref> <value>")
		}
		body["ref"] = args[0]
		body["value"] = args[1]
	}

	doPost(client, base, token, "/action", body)
}

// --- text ---

func cliText(client *http.Client, base, token string, args []string) {
	params := url.Values{}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--raw":
			params.Set("mode", "raw")
		case "--tab":
			if i+1 < len(args) {
				i++
				params.Set("tabId", args[i])
			}
		}
	}
	doGet(client, base, token, "/text", params)
}

// --- tabs ---

func cliTabs(client *http.Client, base, token string, args []string) {
	if len(args) == 0 {
		// List all tabs
		doGet(client, base, token, "/tabs", nil)
		return
	}

	cmd := args[0]
	subArgs := args[1:]

	// Check if this is a tab operation (navigate, snapshot, click, etc.)
	// Pattern: pinchtab tab <operation> <tabId> [args...]
	if isTabOperation(cmd) {
		cliTabOperation(client, base, token, cmd, subArgs)
		return
	}

	// Legacy: pinchtab tab new/close
	switch cmd {
	case "new":
		url := ""
		if len(subArgs) > 0 {
			url = subArgs[0]
		}

		// Check if any instances are running
		instances := getInstances(client, base, token)
		if len(instances) == 0 {
			fmt.Fprintln(os.Stderr, styleStderr(cliWarningStyle, "No instances running, launching default..."))
			launchInstance(client, base, token, "default")
			fmt.Fprintln(os.Stderr, styleStderr(cliSuccessStyle, "Instance launched"))
		}

		body := map[string]any{"action": "new"}
		if url != "" {
			body["url"] = url
		}
		doPost(client, base, token, "/tab", body)

	case "close":
		if len(subArgs) < 1 {
			fatal("Usage: pinchtab tab close <tabId>")
		}
		doPost(client, base, token, "/tab", map[string]any{
			"action": "close",
			"tabId":  subArgs[0],
		})

	default:
		cliTabOperation(client, base, token, cmd, subArgs)
	}
}

func isTabOperation(op string) bool {
	ops := map[string]bool{
		"navigate": true, "snapshot": true, "screenshot": true,
		"click": true, "type": true, "press": true, "fill": true,
		"hover": true, "scroll": true, "select": true, "focus": true,
		"text": true, "eval": true, "evaluate": true, "pdf": true,
		"cookies": true, "lock": true, "unlock": true, "locks": true,
		"fingerprint": true, "info": true,
	}
	return ops[op]
}

func cliTabOperation(client *http.Client, base, token string, op string, args []string) {
	if len(args) < 1 {
		fatal("Usage: pinchtab tab %s <tabId> [args...]", op)
	}

	tabID := args[0]
	restArgs := args[1:]

	switch op {
	case "navigate":
		if len(restArgs) < 1 {
			fatal("Usage: pinchtab tab navigate <tabId> <url> [--timeout N] [--block-images]")
		}
		body := map[string]any{"url": restArgs[0]}
		for i := 1; i < len(restArgs); i++ {
			switch restArgs[i] {
			case "--timeout":
				if i+1 < len(restArgs) {
					body["timeout"] = restArgs[i+1]
					i++
				}
			case "--block-images":
				body["blockImages"] = true
			case "--block-ads":
				body["blockAds"] = true
			}
		}
		doPost(client, base, token, fmt.Sprintf("/tabs/%s/navigate", tabID), body)

	case "snapshot":
		params := url.Values{}
		for _, arg := range restArgs {
			switch arg {
			case "-i", "--interactive":
				params.Set("interactive", "true")
			case "-c", "--compact":
				params.Set("compact", "true")
			case "-d", "--diff":
				params.Set("diff", "true")
			}
		}
		doGet(client, base, token, fmt.Sprintf("/tabs/%s/snapshot", tabID), params)

	case "screenshot", "ss":
		params := url.Values{}
		outFile := ""
		for i := 0; i < len(restArgs); i++ {
			switch restArgs[i] {
			case "-o", "--output":
				if i+1 < len(restArgs) {
					outFile = restArgs[i+1]
					i++
				}
			case "-q", "--quality":
				if i+1 < len(restArgs) {
					params.Set("quality", restArgs[i+1])
					i++
				}
			}
		}
		params.Set("raw", "true")
		data := doGetRaw(client, base, token, fmt.Sprintf("/tabs/%s/screenshot", tabID), params)
		if outFile == "" {
			outFile = fmt.Sprintf("screenshot-%s.png", time.Now().Format("20060102-150405"))
		}
		if data != nil {
			if err := os.WriteFile(outFile, data, 0600); err == nil {
				fmt.Println(styleStdout(cliSuccessStyle, fmt.Sprintf("Saved %s (%d bytes)", outFile, len(data))))
			}
		}

	case "click", "hover", "focus":
		if len(restArgs) < 1 {
			fatal("Usage: pinchtab tab %s <tabId> <ref>", op)
		}
		body := map[string]any{"kind": op, "ref": restArgs[0]}
		doPost(client, base, token, fmt.Sprintf("/tabs/%s/action", tabID), body)

	case "type":
		if len(restArgs) < 2 {
			fatal("Usage: pinchtab tab type <tabId> <ref> <text>")
		}
		body := map[string]any{"kind": "type", "ref": restArgs[0], "text": strings.Join(restArgs[1:], " ")}
		doPost(client, base, token, fmt.Sprintf("/tabs/%s/action", tabID), body)

	case "fill":
		if len(restArgs) < 2 {
			fatal("Usage: pinchtab tab fill <tabId> <ref> <text>")
		}
		body := map[string]any{"kind": "fill", "ref": restArgs[0], "text": strings.Join(restArgs[1:], " ")}
		doPost(client, base, token, fmt.Sprintf("/tabs/%s/action", tabID), body)

	case "press":
		if len(restArgs) < 1 {
			fatal("Usage: pinchtab tab press <tabId> <key>")
		}
		body := map[string]any{"kind": "press", "key": restArgs[0]}
		doPost(client, base, token, fmt.Sprintf("/tabs/%s/action", tabID), body)

	case "scroll":
		if len(restArgs) < 1 {
			fatal("Usage: pinchtab tab scroll <tabId> <direction|pixels>")
		}
		body := map[string]any{}
		if v, err := strconv.Atoi(restArgs[0]); err == nil {
			body["kind"] = "scroll"
			body["scrollY"] = v
		} else {
			body["kind"] = "scroll"
			body["direction"] = restArgs[0]
		}
		doPost(client, base, token, fmt.Sprintf("/tabs/%s/action", tabID), body)

	case "select":
		if len(restArgs) < 2 {
			fatal("Usage: pinchtab tab select <tabId> <ref> <value>")
		}
		body := map[string]any{"kind": "select", "ref": restArgs[0], "value": restArgs[1]}
		doPost(client, base, token, fmt.Sprintf("/tabs/%s/action", tabID), body)

	case "text":
		params := url.Values{}
		for _, arg := range restArgs {
			if arg == "--raw" {
				params.Set("raw", "true")
			}
		}
		doGet(client, base, token, fmt.Sprintf("/tabs/%s/text", tabID), params)

	case "eval", "evaluate":
		if len(restArgs) < 1 {
			fatal("Usage: pinchtab tab eval <tabId> <expression>")
		}
		body := map[string]any{"expression": strings.Join(restArgs, " ")}
		doPost(client, base, token, fmt.Sprintf("/tabs/%s/evaluate", tabID), body)

	case "pdf":
		params := url.Values{}
		outFile := ""
		for i := 0; i < len(restArgs); i++ {
			switch restArgs[i] {
			case "-o", "--output":
				if i+1 < len(restArgs) {
					outFile = restArgs[i+1]
					i++
				}
			case "--landscape":
				params.Set("landscape", "true")
			case "--scale":
				if i+1 < len(restArgs) {
					params.Set("scale", restArgs[i+1])
					i++
				}
			}
		}
		params.Set("raw", "true")
		data := doGetRaw(client, base, token, fmt.Sprintf("/tabs/%s/pdf", tabID), params)
		if outFile == "" {
			outFile = fmt.Sprintf("page-%s.pdf", time.Now().Format("20060102-150405"))
		}
		if data != nil {
			if err := os.WriteFile(outFile, data, 0600); err == nil {
				fmt.Printf("Saved %s (%d bytes)\n", outFile, len(data))
			}
		}

	case "cookies":
		doGet(client, base, token, fmt.Sprintf("/tabs/%s/cookies", tabID), nil)

	case "lock":
		body := map[string]any{}
		for i := 0; i < len(restArgs); i++ {
			switch restArgs[i] {
			case "--owner":
				if i+1 < len(restArgs) {
					body["owner"] = restArgs[i+1]
					i++
				}
			case "--ttl":
				if i+1 < len(restArgs) {
					if ttl, err := strconv.Atoi(restArgs[i+1]); err == nil {
						body["ttl"] = ttl
					}
					i++
				}
			}
		}
		doPost(client, base, token, fmt.Sprintf("/tabs/%s/lock", tabID), body)

	case "unlock":
		body := map[string]any{}
		for i := 0; i < len(restArgs); i++ {
			switch restArgs[i] {
			case "--owner":
				if i+1 < len(restArgs) {
					body["owner"] = restArgs[i+1]
					i++
				}
			}
		}
		doPost(client, base, token, fmt.Sprintf("/tabs/%s/unlock", tabID), body)

	case "locks":
		doGet(client, base, token, fmt.Sprintf("/tabs/%s/locks", tabID), nil)

	case "info":
		doGet(client, base, token, fmt.Sprintf("/tabs/%s", tabID), nil)

	default:
		fatal("Unknown tab operation: %s", op)
	}
}

// getInstances fetches the list of running instances
func getInstances(client *http.Client, base, token string) []map[string]any {
	resp, err := http.NewRequest("GET", base+"/instances", nil)
	if err != nil {
		return nil
	}
	if token != "" {
		resp.Header.Set("Authorization", "Bearer "+token)
	}

	result, err := client.Do(resp)
	if err != nil || result.StatusCode >= 400 {
		return nil
	}
	defer func() { _ = result.Body.Close() }()

	var data map[string]any
	if err := json.NewDecoder(result.Body).Decode(&data); err != nil {
		log.Printf("warning: error decoding instances response: %v", err)
	}

	if instances, ok := data["instances"].([]interface{}); ok {
		converted := make([]map[string]any, len(instances))
		for i, inst := range instances {
			if m, ok := inst.(map[string]any); ok {
				converted[i] = m
			}
		}
		return converted
	}
	return nil
}

// launchInstance launches a default instance
func launchInstance(client *http.Client, base, token string, profile string) {
	body := map[string]any{"profile": profile}
	doPost(client, base, token, "/instances/launch", body)
}

// --- screenshot ---

func cliScreenshot(client *http.Client, base, token string, args []string) {
	params := url.Values{}
	params.Set("raw", "true")
	outFile := ""
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-o", "--output":
			if i+1 < len(args) {
				i++
				outFile = args[i]
			}
		case "--quality", "-q":
			if i+1 < len(args) {
				i++
				params.Set("quality", args[i])
			}
		case "--tab":
			if i+1 < len(args) {
				i++
				params.Set("tabId", args[i])
			}
		}
	}

	if outFile == "" {
		outFile = fmt.Sprintf("screenshot-%s.jpg", time.Now().Format("20060102-150405"))
	}

	data := doGetRaw(client, base, token, "/screenshot", params)
	if data == nil {
		return
	}
	if err := os.WriteFile(outFile, data, 0600); err != nil {
		fatal("Write failed: %v", err)
	}
	fmt.Println(styleStdout(cliSuccessStyle, fmt.Sprintf("Saved %s (%d bytes)", outFile, len(data))))
}

// --- evaluate ---

func cliEvaluate(client *http.Client, base, token string, args []string) {
	if len(args) < 1 {
		fatal("Usage: pinchtab eval <expression>")
	}
	expr := strings.Join(args, " ")
	doPost(client, base, token, "/evaluate", map[string]any{
		"expression": expr,
	})
}

// --- pdf ---

func cliPDF(client *http.Client, base, token string, args []string) {
	params := url.Values{}
	params.Set("raw", "true")
	outFile := ""
	tabID := ""
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-o", "--output":
			if i+1 < len(args) {
				i++
				outFile = args[i]
			}
		case "--landscape":
			params.Set("landscape", "true")
		case "--scale":
			if i+1 < len(args) {
				i++
				params.Set("scale", args[i])
			}
		case "--tab":
			if i+1 < len(args) {
				i++
				tabID = args[i]
			}
		// Paper dimensions
		case "--paper-width":
			if i+1 < len(args) {
				i++
				params.Set("paperWidth", args[i])
			}
		case "--paper-height":
			if i+1 < len(args) {
				i++
				params.Set("paperHeight", args[i])
			}
		// Margins
		case "--margin-top":
			if i+1 < len(args) {
				i++
				params.Set("marginTop", args[i])
			}
		case "--margin-bottom":
			if i+1 < len(args) {
				i++
				params.Set("marginBottom", args[i])
			}
		case "--margin-left":
			if i+1 < len(args) {
				i++
				params.Set("marginLeft", args[i])
			}
		case "--margin-right":
			if i+1 < len(args) {
				i++
				params.Set("marginRight", args[i])
			}
		// Content options
		case "--page-ranges":
			if i+1 < len(args) {
				i++
				params.Set("pageRanges", args[i])
			}
		case "--prefer-css-page-size":
			params.Set("preferCSSPageSize", "true")
		// Header/Footer
		case "--display-header-footer":
			params.Set("displayHeaderFooter", "true")
		case "--header-template":
			if i+1 < len(args) {
				i++
				params.Set("headerTemplate", args[i])
			}
		case "--footer-template":
			if i+1 < len(args) {
				i++
				params.Set("footerTemplate", args[i])
			}
		// Accessibility
		case "--generate-tagged-pdf":
			params.Set("generateTaggedPDF", "true")
		case "--generate-document-outline":
			params.Set("generateDocumentOutline", "true")
		// Output options
		case "--file-output":
			params.Del("raw")
			params.Set("output", "file")
		case "--path":
			if i+1 < len(args) {
				i++
				params.Set("path", args[i])
			}
		case "--raw":
			params.Set("raw", "true")
		}
	}

	if outFile == "" {
		outFile = fmt.Sprintf("page-%s.pdf", time.Now().Format("20060102-150405"))
	}

	var data []byte
	if tabID != "" {
		data = doGetRaw(client, base, token, fmt.Sprintf("/tabs/%s/pdf", tabID), params)
	} else {
		data = doGetRaw(client, base, token, "/pdf", params)
	}
	if data == nil {
		return
	}
	if err := os.WriteFile(outFile, data, 0600); err != nil {
		fatal("Write failed: %v", err)
	}
	fmt.Println(styleStdout(cliSuccessStyle, fmt.Sprintf("Saved %s (%d bytes)", outFile, len(data))))
}

// --- quick command ---

func cliQuick(client *http.Client, base, token string, args []string) {
	if len(args) < 1 {
		fatal("Usage: pinchtab quick <url>")
	}

	fmt.Println(styleStdout(cliHeadingStyle, fmt.Sprintf("Navigating to %s...", args[0])))

	// Navigate
	navBody := map[string]any{"url": args[0]}
	navResult := doPost(client, base, token, "/navigate", navBody)

	// Small delay for page to stabilize
	time.Sleep(1 * time.Second)

	fmt.Println()
	fmt.Println(styleStdout(cliHeadingStyle, "Page structure"))

	// Snapshot with interactive filter
	snapParams := url.Values{}
	snapParams.Set("filter", "interactive")
	snapParams.Set("compact", "true")
	doGet(client, base, token, "/snapshot", snapParams)

	// Extract info from navigation result
	if title, ok := navResult["title"].(string); ok {
		fmt.Println()
		fmt.Printf("%s %s\n", styleStdout(cliMutedStyle, "Title:"), styleStdout(cliValueStyle, title))
	}
	if urlStr, ok := navResult["url"].(string); ok {
		fmt.Printf("%s %s\n", styleStdout(cliMutedStyle, "URL:"), styleStdout(cliValueStyle, urlStr))
	}

	fmt.Println()
	fmt.Println(styleStdout(cliHeadingStyle, "Quick actions"))
	fmt.Printf("  %s %s\n", styleStdout(cliCommandStyle, "pinchtab click <ref>"), styleStdout(cliMutedStyle, "# Click an element (use refs from above)"))
	fmt.Printf("  %s %s\n", styleStdout(cliCommandStyle, "pinchtab type <ref> <text>"), styleStdout(cliMutedStyle, "# Type into input field"))
	fmt.Printf("  %s %s\n", styleStdout(cliCommandStyle, "pinchtab screenshot"), styleStdout(cliMutedStyle, "# Take a screenshot"))
	fmt.Printf("  %s %s\n", styleStdout(cliCommandStyle, "pinchtab pdf --tab <id> -o output.pdf"), styleStdout(cliMutedStyle, "# Save tab as PDF"))
}

// --- health ---

func cliHealth(client *http.Client, base, token string) {
	doGet(client, base, token, "/health", nil)
}

// --- instance ---

func cliInstance(client *http.Client, base, token string, args []string) {
	if len(args) < 1 {
		fatal("Usage: pinchtab instance <subcommand> [options]\nSubcommands: start, launch (alias), navigate, logs, stop")
	}

	subCmd := args[0]
	subArgs := args[1:]

	switch subCmd {
	case "start", "launch": // "start" is new Phase 2 API, "launch" is legacy
		cliInstanceStart(client, base, token, subArgs)
	case "navigate":
		cliInstanceNavigate(client, base, token, subArgs)
	case "logs":
		cliInstanceLogs(client, base, token, subArgs)
	case "stop":
		cliInstanceStop(client, base, token, subArgs)
	default:
		fatal("Unknown subcommand: %s", subCmd)
	}
}

func cliInstanceStart(client *http.Client, base, token string, args []string) {
	body := map[string]any{}

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--profileId":
			if i+1 < len(args) {
				body["profileId"] = args[i+1]
				i++
			}
		case "--mode":
			if i+1 < len(args) {
				body["mode"] = args[i+1]
				i++
			}
		case "--port":
			if i+1 < len(args) {
				body["port"] = args[i+1]
				i++
			}
		}
	}

	// Use new /instances/start endpoint if available, fall back to /instances/launch for backward compat
	endpoint := "/instances/start"
	doPost(client, base, token, endpoint, body)
}

func cliInstanceNavigate(client *http.Client, base, token string, args []string) {
	if len(args) < 2 {
		fatal("Usage: pinchtab instance navigate <instance-id> <url>")
	}

	instID := args[0]
	targetURL := args[1]

	// Instance navigate now works via tab-scoped navigation:
	// open a tab on the instance, then navigate that tab.
	openResp := doPost(client, base, token, fmt.Sprintf("/instances/%s/tabs/open", instID), map[string]any{
		"url": "about:blank",
	})
	tabID, _ := openResp["tabId"].(string)
	if tabID == "" {
		fatal("failed to open tab for instance %s", instID)
	}

	// doPost auto-prints JSON response.
	doPost(client, base, token, fmt.Sprintf("/tabs/%s/navigate", tabID), map[string]any{
		"url": targetURL,
	})
}

func cliInstanceLogs(client *http.Client, base, token string, args []string) {
	var instID string

	// Support both positional argument and --id flag
	if len(args) == 0 {
		fatal("Usage: pinchtab instance logs <instance-id> OR pinchtab instance logs --id <instance-id>")
	}

	// Check if first arg is --id flag
	if args[0] == "--id" {
		if len(args) < 2 {
			fatal("Usage: --id requires instance ID")
		}
		instID = args[1]
	} else {
		// Positional argument (backward compat)
		instID = args[0]
	}

	logs := doGetRaw(client, base, token, fmt.Sprintf("/instances/%s/logs", instID), nil)
	fmt.Println(string(logs))
}

func cliInstanceStop(client *http.Client, base, token string, args []string) {
	var instID string

	// Support both positional argument and --id flag
	if len(args) == 0 {
		fatal("Usage: pinchtab instance stop <instance-id> OR pinchtab instance stop --id <instance-id>")
	}

	// Check if first arg is --id flag
	if args[0] == "--id" {
		if len(args) < 2 {
			fatal("Usage: --id requires instance ID")
		}
		instID = args[1]
	} else {
		// Positional argument (backward compat)
		instID = args[0]
	}

	// doPost auto-prints JSON response
	doPost(client, base, token, fmt.Sprintf("/instances/%s/stop", instID), nil)
}

// --- instances ---

func cliInstances(client *http.Client, base, token string) {
	body := doGetRaw(client, base, token, "/instances", nil)

	// Parse and format as JSON
	var instances []map[string]any
	if err := json.Unmarshal(body, &instances); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse instances: %v\n", err)
		os.Exit(1)
	}

	// Transform to cleaner output format
	output := make([]map[string]any, len(instances))
	for i, inst := range instances {
		id, _ := inst["id"].(string)
		port, _ := inst["port"].(string)
		headless, _ := inst["headless"].(bool)
		status, _ := inst["status"].(string)

		mode := "headless"
		if !headless {
			mode = "headed"
		}

		output[i] = map[string]any{
			"id":     id,
			"port":   port,
			"mode":   mode,
			"status": status,
		}
	}

	// Output as JSON
	data, _ := json.MarshalIndent(output, "", "  ")
	fmt.Println(string(data))
}

// --- profiles ---

func cliProfiles(client *http.Client, base, token string) {
	result := doGet(client, base, token, "/profiles", nil)

	// Display profiles in a friendly format
	if profiles, ok := result["profiles"].([]interface{}); ok && len(profiles) > 0 {
		fmt.Println()
		for _, prof := range profiles {
			if m, ok := prof.(map[string]any); ok {
				name, _ := m["name"].(string)

				fmt.Printf("👤 %s\n", name)
			}
		}
		fmt.Println()
	} else {
		fmt.Println("No profiles available")
	}
}

// --- helpers ---

func doGet(client *http.Client, base, token, path string, params url.Values) map[string]any {
	u := base + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}
	req, _ := http.NewRequest("GET", u, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		fatal("Request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		fmt.Fprintln(os.Stderr, styleStderr(cliErrorStyle, fmt.Sprintf("Error %d: %s", resp.StatusCode, string(body))))
		os.Exit(1)
	}

	// Pretty-print JSON if possible
	var buf bytes.Buffer
	if json.Indent(&buf, body, "", "  ") == nil {
		fmt.Println(buf.String())
	} else {
		fmt.Println(string(body))
	}

	// Parse and return result
	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("warning: error unmarshaling response: %v", err)
	}
	return result
}

func doGetRaw(client *http.Client, base, token, path string, params url.Values) []byte {
	u := base + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}
	req, _ := http.NewRequest("GET", u, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		fatal("Request failed: %v", err)
		return nil
	}
	defer func() { _ = resp.Body.Close() }()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		fmt.Fprintln(os.Stderr, styleStderr(cliErrorStyle, fmt.Sprintf("Error %d: %s", resp.StatusCode, string(body))))
		os.Exit(1)
	}
	return body
}

func doPost(client *http.Client, base, token, path string, body map[string]any) map[string]any {
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", base+path, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		fatal("Request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		fmt.Fprintln(os.Stderr, styleStderr(cliErrorStyle, fmt.Sprintf("Error %d: %s", resp.StatusCode, string(respBody))))
		os.Exit(1)
	}

	var buf bytes.Buffer
	if json.Indent(&buf, respBody, "", "  ") == nil {
		fmt.Println(buf.String())
	} else {
		fmt.Println(string(respBody))
	}

	// Parse and return result for suggestions
	var result map[string]any
	if err := json.Unmarshal(respBody, &result); err != nil {
		log.Printf("warning: error unmarshaling response: %v", err)
	}
	return result
}

// checkServerAndGuide checks if pinchtab server is running and provides guidance
func checkServerAndGuide(client *http.Client, base, token string) bool {
	req, _ := http.NewRequest("GET", base+"/health", nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "dial tcp") {
			fmt.Fprintln(os.Stderr, styleStderr(cliErrorStyle, fmt.Sprintf("Pinchtab server is not running on %s", base)))
			fmt.Fprintln(os.Stderr)
			fmt.Fprintln(os.Stderr, styleStderr(cliHeadingStyle, "To start the server"))
			fmt.Fprintf(os.Stderr, "  %s %s\n", styleStderr(cliCommandStyle, "pinchtab"), styleStderr(cliMutedStyle, "# Run in foreground (recommended for beginners)"))
			fmt.Fprintf(os.Stderr, "  %s %s\n", styleStderr(cliCommandStyle, "pinchtab &"), styleStderr(cliMutedStyle, "# Run in background"))
			fmt.Fprintf(os.Stderr, "  %s %s\n", styleStderr(cliCommandStyle, "PINCHTAB_PORT=9868 pinchtab"), styleStderr(cliMutedStyle, "# Use different port"))
			fmt.Fprintln(os.Stderr)
			fmt.Fprintln(os.Stderr, styleStderr(cliHeadingStyle, "Then try your command again"))
			fmt.Fprintf(os.Stderr, "  %s\n", styleStderr(cliCommandStyle, strings.Join(os.Args, " ")))
			fmt.Fprintln(os.Stderr)
			fmt.Fprintf(os.Stderr, "%s %s\n", styleStderr(cliMutedStyle, "Learn more:"), styleStderr(cliCommandStyle, "https://github.com/pinchtab/pinchtab#quick-start"))
			return false
		}
		// Other connection errors
		fmt.Fprintln(os.Stderr, styleStderr(cliErrorStyle, fmt.Sprintf("Cannot connect to Pinchtab server: %v", err)))
		return false
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == 401 {
		fmt.Fprintln(os.Stderr, styleStderr(cliErrorStyle, "Authentication required. Set PINCHTAB_TOKEN."))
		return false
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		fmt.Fprintln(os.Stderr, styleStderr(cliErrorStyle, fmt.Sprintf("Server error %d: %s", resp.StatusCode, string(body))))
		return false
	}

	return true
}

// suggestNextAction provides helpful suggestions based on the current command and state
func suggestNextAction(cmd string, result map[string]any) {
	switch cmd {
	case "nav", "navigate":
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, styleStderr(cliHeadingStyle, "Next steps"))
		fmt.Fprintf(os.Stderr, "  %s %s\n", styleStderr(cliCommandStyle, "pinchtab snap"), styleStderr(cliMutedStyle, "# See page structure"))
		fmt.Fprintf(os.Stderr, "  %s %s\n", styleStderr(cliCommandStyle, "pinchtab screenshot"), styleStderr(cliMutedStyle, "# Capture visual"))
		fmt.Fprintf(os.Stderr, "  %s %s\n", styleStderr(cliCommandStyle, "pinchtab click <ref>"), styleStderr(cliMutedStyle, "# Click an element"))
		fmt.Fprintf(os.Stderr, "  %s %s\n", styleStderr(cliCommandStyle, "pinchtab pdf --tab <id> -o output.pdf"), styleStderr(cliMutedStyle, "# Save tab as PDF"))

	case "snap", "snapshot":
		refs := extractRefs(result)
		if len(refs) > 0 {
			fmt.Fprintln(os.Stderr)
			fmt.Fprintln(os.Stderr, styleStderr(cliHeadingStyle, fmt.Sprintf("Found %d interactive elements", len(refs))))
			for i, ref := range refs[:min(3, len(refs))] {
				fmt.Fprintf(os.Stderr, "  %s %s\n", styleStderr(cliCommandStyle, fmt.Sprintf("pinchtab click %s", ref.id)), styleStderr(cliMutedStyle, "# "+ref.desc))
				if i >= 2 {
					break
				}
			}
			if len(refs) > 3 {
				fmt.Fprintf(os.Stderr, "  %s\n", styleStderr(cliMutedStyle, fmt.Sprintf("... and %d more", len(refs)-3)))
			}
		}

	case "click", "type", "fill":
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, styleStderr(cliHeadingStyle, "Action completed"))
		fmt.Fprintf(os.Stderr, "  %s %s\n", styleStderr(cliCommandStyle, "pinchtab snap"), styleStderr(cliMutedStyle, "# See updated page"))
		fmt.Fprintf(os.Stderr, "  %s %s\n", styleStderr(cliCommandStyle, "pinchtab screenshot"), styleStderr(cliMutedStyle, "# Visual confirmation"))
	}
}

type refInfo struct {
	id   string
	desc string
}

func extractRefs(data map[string]any) []refInfo {
	var refs []refInfo

	// Handle different snapshot formats
	if elements, ok := data["elements"].([]any); ok {
		for _, elem := range elements {
			if m, ok := elem.(map[string]any); ok {
				if ref, ok := m["ref"].(string); ok && ref != "" {
					desc := ""
					if role, ok := m["role"].(string); ok {
						desc = role
					}
					if name, ok := m["name"].(string); ok && name != "" {
						desc += ": " + name
					}
					// Only include interactive elements
					if role, ok := m["role"].(string); ok {
						if role == "button" || role == "link" || role == "textbox" ||
							role == "checkbox" || role == "radio" || role == "combobox" {
							refs = append(refs, refInfo{id: ref, desc: desc})
						}
					}
				}
			}
		}
	}

	return refs
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// resolveInstanceBase fetches the named instance from the orchestrator and returns
// a base URL pointing directly at that instance's API port.
func resolveInstanceBase(orchBase, token, instanceID, bind string) string {
	c := &http.Client{Timeout: 10 * time.Second}
	body := doGetRaw(c, orchBase, token, fmt.Sprintf("/instances/%s", instanceID), nil)

	var inst struct {
		Port string `json:"port"`
	}
	if err := json.Unmarshal(body, &inst); err != nil {
		fatal("failed to parse instance %q: %v", instanceID, err)
	}
	if inst.Port == "" {
		fatal("instance %q has no port assigned (is it still starting?)", instanceID)
	}
	return fmt.Sprintf("http://%s:%s", bind, inst.Port)
}

func fatal(format string, args ...any) {
	fmt.Fprintln(os.Stderr, styleStderr(cliErrorStyle, fmt.Sprintf(format, args...)))
	os.Exit(1)
}
