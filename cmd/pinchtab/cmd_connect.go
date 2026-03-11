package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/pinchtab/pinchtab/internal/config"
	"github.com/spf13/cobra"
)

type profileInstanceStatus struct {
	Name    string `json:"name"`
	Running bool   `json:"running"`
	Status  string `json:"status"`
	Port    string `json:"port"`
	ID      string `json:"id"`
	Error   string `json:"error"`
}

var connectCmd = &cobra.Command{
	Use:   "connect <profile>",
	Short: "Get URL for a running profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		jsonOut, _ := cmd.Flags().GetBool("json")
		dashboardURL, _ := cmd.Flags().GetString("dashboard")
		handleConnectCommand(cfg, args[0], dashboardURL, jsonOut)
	},
}

func init() {
	connectCmd.Flags().Bool("json", false, "Output as JSON")
	connectCmd.Flags().String("dashboard", "", "Dashboard URL (e.g. http://localhost:9867)")
	rootCmd.AddCommand(connectCmd)
}

func handleConnectCommand(cfg *config.RuntimeConfig, profile, dashboardURL string, jsonOut bool) {
	if dashboardURL == "" {
		dashboardURL = os.Getenv("PINCHTAB_DASHBOARD_URL")
	}
	if dashboardURL == "" {
		dashPort := cfg.Port
		if dashPort == "" {
			dashPort = "9870"
		}
		dashboardURL = "http://localhost:" + dashPort
	}

	reqURL := stringsTrimRightSlash(dashboardURL) + "/profiles/" + url.PathEscape(profile) + "/instance"
	res, err := http.Get(reqURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect failed: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(res.Body, 8<<10))
		fmt.Fprintf(os.Stderr, "connect failed: dashboard returned %d: %s\n", res.StatusCode, string(body))
		os.Exit(1)
	}

	var st profileInstanceStatus
	if err := json.NewDecoder(res.Body).Decode(&st); err != nil {
		fmt.Fprintf(os.Stderr, "connect failed: invalid response: %v\n", err)
		os.Exit(1)
	}
	if !st.Running || st.Port == "" {
		errMsg := st.Error
		if errMsg == "" {
			errMsg = st.Status
		}
		fmt.Fprintf(os.Stderr, "profile %q not running (%s)\n", profile, errMsg)
		os.Exit(1)
	}

	instanceURL := "http://localhost:" + st.Port
	if jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(map[string]string{
			"profile": st.Name,
			"id":      st.ID,
			"status":  st.Status,
			"port":    st.Port,
			"url":     instanceURL,
		})
		return
	}

	fmt.Println(instanceURL)
}

func stringsTrimRightSlash(s string) string {
	for len(s) > 0 && s[len(s)-1] == '/' {
		s = s[:len(s)-1]
	}
	return s
}
