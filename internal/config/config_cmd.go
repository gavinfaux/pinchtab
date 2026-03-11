package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pinchtab/pinchtab/internal/cliui"
	"github.com/spf13/cobra"
)

// ConfigCmd is the root command for configuration management.
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
}

func init() {
	ConfigCmd.AddCommand(&cobra.Command{
		Use:   "show",
		Short: "Display current configuration",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := Load()
			handleConfigShow(cfg)
		},
	})
	ConfigCmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "Initialize a new config file",
		Run: func(cmd *cobra.Command, args []string) {
			handleConfigInit()
		},
	})
	ConfigCmd.AddCommand(&cobra.Command{
		Use:   "path",
		Short: "Show config file path",
		Run: func(cmd *cobra.Command, args []string) {
			handleConfigPath()
		},
	})
	ConfigCmd.AddCommand(&cobra.Command{
		Use:   "validate",
		Short: "Validate config file",
		Run: func(cmd *cobra.Command, args []string) {
			handleConfigValidate()
		},
	})
	ConfigCmd.AddCommand(&cobra.Command{
		Use:   "get <path>",
		Short: "Get a config value (e.g., server.port)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			handleConfigGet(args[0])
		},
	})
	ConfigCmd.AddCommand(&cobra.Command{
		Use:   "set <path> <val>",
		Short: "Set a config value (e.g., server.port 8080)",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			handleConfigSet(args[0], args[1])
		},
	})
	ConfigCmd.AddCommand(&cobra.Command{
		Use:   "patch <json>",
		Short: "Merge JSON into config",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			handleConfigPatch(args[0])
		},
	})
}

func handleConfigShow(cfg *RuntimeConfig) {
	fmt.Println(cliui.RenderStdout(cliui.HeadingStyle, "Current configuration (env > file > defaults):"))
	fmt.Println()
	fmt.Println(cliui.RenderStdout(cliui.HeadingStyle, "Server"))
	fmt.Printf("  Port:           %s\n", cfg.Port)
	fmt.Printf("  Bind:           %s\n", cfg.Bind)
	fmt.Printf("  Token:          %s\n", MaskToken(cfg.Token))
	fmt.Printf("  State Dir:      %s\n", cfg.StateDir)
	fmt.Printf("  Instance Ports: %d-%d\n", cfg.InstancePortStart, cfg.InstancePortEnd)
	fmt.Println()
	fmt.Println(cliui.RenderStdout(cliui.HeadingStyle, "Security"))
	fmt.Printf("  Evaluate:       %v\n", cfg.AllowEvaluate)
	fmt.Printf("  Macro:          %v\n", cfg.AllowMacro)
	fmt.Printf("  Screencast:     %v\n", cfg.AllowScreencast)
	fmt.Printf("  Download:       %v\n", cfg.AllowDownload)
	fmt.Printf("  Upload:         %v\n", cfg.AllowUpload)
	fmt.Println()
	fmt.Println(cliui.RenderStdout(cliui.HeadingStyle, "Browser / Instance Defaults"))
	fmt.Printf("  Headless:       %v\n", cfg.Headless)
	fmt.Printf("  No Restore:     %v\n", cfg.NoRestore)
	fmt.Printf("  Profile Dir:    %s\n", cfg.ProfileDir)
	fmt.Printf("  Profiles Dir:   %s\n", cfg.ProfilesBaseDir)
	fmt.Printf("  Default Profile: %s\n", cfg.DefaultProfile)
	fmt.Printf("  Max Tabs:       %d\n", cfg.MaxTabs)
	fmt.Printf("  Stealth:        %s\n", cfg.StealthLevel)
	fmt.Printf("  Tab Eviction:   %s\n", cfg.TabEvictionPolicy)
	fmt.Printf("  Extensions:     %v\n", cfg.ExtensionPaths)
	fmt.Println()
	fmt.Println(cliui.RenderStdout(cliui.HeadingStyle, "Multi-Instance"))
	fmt.Printf("  Strategy:       %s\n", cfg.Strategy)
	fmt.Printf("  Allocation:     %s\n", cfg.AllocationPolicy)
	fmt.Printf("  Max Restarts:   %d\n", cfg.RestartMaxRestarts)
	fmt.Printf("  Init Backoff:   %v\n", cfg.RestartInitBackoff)
	fmt.Printf("  Max Backoff:    %v\n", cfg.RestartMaxBackoff)
	fmt.Printf("  Stable After:   %v\n", cfg.RestartStableAfter)
	fmt.Println()
	fmt.Println(cliui.RenderStdout(cliui.HeadingStyle, "Attach"))
	fmt.Printf("  Enabled:        %v\n", cfg.AttachEnabled)
	fmt.Printf("  Allow Hosts:    %v\n", cfg.AttachAllowHosts)
	fmt.Printf("  Allow Schemes:  %v\n", cfg.AttachAllowSchemes)
	fmt.Println()
	fmt.Println(cliui.RenderStdout(cliui.HeadingStyle, "Timeouts"))
	fmt.Printf("  Action:         %v\n", cfg.ActionTimeout)
	fmt.Printf("  Navigate:       %v\n", cfg.NavigateTimeout)
	fmt.Printf("  Shutdown:       %v\n", cfg.ShutdownTimeout)
}

func handleConfigInit() {
	configPath := filepath.Join(userConfigDir(), "config.json")

	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("Config file already exists at %s\n", configPath)
		fmt.Print("Overwrite? (y/N): ")
		var response string
		_, _ = fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			return
		}
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		os.Exit(1)
	}

	fc := DefaultFileConfig()
	token, err := GenerateAuthToken()
	if err != nil {
		fmt.Printf("Error generating auth token: %v\n", err)
		os.Exit(1)
	}
	fc.Server.Token = token
	data, _ := json.MarshalIndent(fc, "", "  ")
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		fmt.Printf("Error writing config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Config file created at %s\n", configPath)
}

func handleConfigPath() {
	configPath := envOr("PINCHTAB_CONFIG", filepath.Join(userConfigDir(), "config.json"))
	fmt.Println(configPath)
}

func handleConfigGet(path string) {
	fc, _, err := LoadFileConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	value, err := GetConfigValue(fc, path)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(value)
}

func handleConfigSet(path, value string) {
	fc, configPath, err := LoadFileConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	if err := SetConfigValue(fc, path, value); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if errs := ValidateFileConfig(fc); len(errs) > 0 {
		fmt.Printf("Warning: new value causes validation error(s):\n")
		for _, e := range errs {
			fmt.Printf("  - %v\n", e)
		}
		fmt.Print("Save anyway? (y/N): ")
		var response string
		_, _ = fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			return
		}
	}

	if err := SaveFileConfig(fc, configPath); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Set %s = %s\n", path, value)
}

func handleConfigPatch(jsonPatch string) {
	fc, configPath, err := LoadFileConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	if err := PatchConfigJSON(fc, jsonPatch); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if errs := ValidateFileConfig(fc); len(errs) > 0 {
		fmt.Printf("Warning: patch causes validation error(s):\n")
		for _, e := range errs {
			fmt.Printf("  - %v\n", e)
		}
		fmt.Print("Save anyway? (y/N): ")
		var response string
		_, _ = fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			return
		}
	}

	if err := SaveFileConfig(fc, configPath); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Config patched successfully")
}

func handleConfigValidate() {
	configPath := envOr("PINCHTAB_CONFIG", filepath.Join(userConfigDir(), "config.json"))
	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		os.Exit(1)
	}

	var fc *FileConfig
	if isLegacyConfig(data) {
		var lc legacyFileConfig
		if err := json.Unmarshal(data, &lc); err != nil {
			fmt.Printf("Error parsing legacy config: %v\n", err)
			os.Exit(1)
		}
		fc = convertLegacyConfig(&lc)
	} else {
		fc = &FileConfig{}
		if err := json.Unmarshal(data, fc); err != nil {
			fmt.Printf("Error parsing config: %v\n", err)
			os.Exit(1)
		}
	}

	if errs := ValidateFileConfig(fc); len(errs) > 0 {
		fmt.Printf("✗ Config file has %d error(s):\n", len(errs))
		for _, e := range errs {
			fmt.Printf("  - %v\n", e)
		}
		os.Exit(1)
	}
	fmt.Printf("✓ Config file is valid: %s\n", configPath)
}
