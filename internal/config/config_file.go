package config

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"time"
)

// DefaultFileConfig returns a FileConfig with sensible defaults (nested format).
func DefaultFileConfig() FileConfig {
	start := 9868
	end := 9968
	restartMaxRestarts := 20
	restartInitBackoffSec := 2
	restartMaxBackoffSec := 60
	restartStableAfterSec := 300
	maxTabs := 20
	allowEvaluate := false
	allowMacro := false
	allowScreencast := false
	allowDownload := false
	allowUpload := false
	maxRedirects := -1
	return FileConfig{
		Server: ServerConfig{
			Port:     "9867",
			Bind:     "127.0.0.1",
			StateDir: userConfigDir(),
		},
		Browser: BrowserConfig{
			ChromeVersion: "144.0.7559.133",
		},
		InstanceDefaults: InstanceDefaultsConfig{
			Mode:              "headless",
			MaxTabs:           &maxTabs,
			StealthLevel:      "light",
			TabEvictionPolicy: "close_lru",
		},
		Security: SecurityConfig{
			AllowEvaluate:   &allowEvaluate,
			AllowMacro:      &allowMacro,
			AllowScreencast: &allowScreencast,
			AllowDownload:   &allowDownload,
			AllowUpload:     &allowUpload,
			MaxRedirects:    &maxRedirects,
			Attach: AttachConfig{
				AllowHosts:   []string{"127.0.0.1", "localhost", "::1"},
				AllowSchemes: []string{"ws", "wss"},
			},
			IDPI: IDPIConfig{
				Enabled:        true,
				AllowedDomains: append([]string(nil), defaultLocalAllowedDomains...),
				StrictMode:     true,
				ScanContent:    true,
				WrapContent:    true,
				ScanTimeoutSec: 5,
			},
		},
		Profiles: ProfilesConfig{
			BaseDir:        filepath.Join(userConfigDir(), "profiles"),
			DefaultProfile: "default",
		},
		MultiInstance: MultiInstanceConfig{
			Strategy:          "simple",
			AllocationPolicy:  "fcfs",
			InstancePortStart: &start,
			InstancePortEnd:   &end,
			Restart: MultiInstanceRestartConfig{
				MaxRestarts:    &restartMaxRestarts,
				InitBackoffSec: &restartInitBackoffSec,
				MaxBackoffSec:  &restartMaxBackoffSec,
				StableAfterSec: &restartStableAfterSec,
			},
		},
		Timeouts: TimeoutsConfig{
			ActionSec:   30,
			NavigateSec: 60,
			ShutdownSec: 10,
			WaitNavMs:   1000,
		},
	}
}

// FileConfigFromRuntime converts the effective runtime configuration back into a
// nested file configuration shape.
func FileConfigFromRuntime(cfg *RuntimeConfig) FileConfig {
	if cfg == nil {
		return DefaultFileConfig()
	}

	noRestore := cfg.NoRestore
	blockImages := cfg.BlockImages
	blockMedia := cfg.BlockMedia
	blockAds := cfg.BlockAds
	maxTabs := cfg.MaxTabs
	maxParallelTabs := cfg.MaxParallelTabs
	noAnimations := cfg.NoAnimations
	allowEvaluate := cfg.AllowEvaluate
	allowMacro := cfg.AllowMacro
	allowScreencast := cfg.AllowScreencast
	allowDownload := cfg.AllowDownload
	allowUpload := cfg.AllowUpload
	maxRedirects := cfg.MaxRedirects
	attachEnabled := cfg.AttachEnabled
	start := cfg.InstancePortStart
	end := cfg.InstancePortEnd
	restartMaxRestarts := cfg.RestartMaxRestarts
	restartInitBackoffSec := int(cfg.RestartInitBackoff / time.Second)
	restartMaxBackoffSec := int(cfg.RestartMaxBackoff / time.Second)
	restartStableAfterSec := int(cfg.RestartStableAfter / time.Second)

	mode := "headless"
	if !cfg.Headless {
		mode = "headed"
	}

	fc := FileConfig{
		Server: ServerConfig{
			Port:     cfg.Port,
			Bind:     cfg.Bind,
			Token:    cfg.Token,
			StateDir: cfg.StateDir,
			Engine:   cfg.Engine,
		},
		Browser: BrowserConfig{
			ChromeVersion:    cfg.ChromeVersion,
			ChromeBinary:     cfg.ChromeBinary,
			ChromeExtraFlags: cfg.ChromeExtraFlags,
			ExtensionPaths:   append([]string(nil), cfg.ExtensionPaths...),
		},
		InstanceDefaults: InstanceDefaultsConfig{
			Mode:              mode,
			NoRestore:         &noRestore,
			Timezone:          cfg.Timezone,
			BlockImages:       &blockImages,
			BlockMedia:        &blockMedia,
			BlockAds:          &blockAds,
			MaxTabs:           &maxTabs,
			MaxParallelTabs:   &maxParallelTabs,
			UserAgent:         cfg.UserAgent,
			NoAnimations:      &noAnimations,
			StealthLevel:      cfg.StealthLevel,
			TabEvictionPolicy: cfg.TabEvictionPolicy,
		},
		Security: SecurityConfig{
			AllowEvaluate:   &allowEvaluate,
			AllowMacro:      &allowMacro,
			AllowScreencast: &allowScreencast,
			AllowDownload:   &allowDownload,
			AllowUpload:     &allowUpload,
			MaxRedirects:    &maxRedirects,
			Attach: AttachConfig{
				Enabled:      &attachEnabled,
				AllowHosts:   append([]string(nil), cfg.AttachAllowHosts...),
				AllowSchemes: append([]string(nil), cfg.AttachAllowSchemes...),
			},
			IDPI: cfg.IDPI,
		},
		Profiles: ProfilesConfig{
			BaseDir:        cfg.ProfilesBaseDir,
			DefaultProfile: cfg.DefaultProfile,
		},
		MultiInstance: MultiInstanceConfig{
			Strategy:          cfg.Strategy,
			AllocationPolicy:  cfg.AllocationPolicy,
			InstancePortStart: &start,
			InstancePortEnd:   &end,
			Restart: MultiInstanceRestartConfig{
				MaxRestarts:    &restartMaxRestarts,
				InitBackoffSec: &restartInitBackoffSec,
				MaxBackoffSec:  &restartMaxBackoffSec,
				StableAfterSec: &restartStableAfterSec,
			},
		},
		Timeouts: TimeoutsConfig{
			ActionSec:   int(cfg.ActionTimeout / time.Second),
			NavigateSec: int(cfg.NavigateTimeout / time.Second),
			ShutdownSec: int(cfg.ShutdownTimeout / time.Second),
			WaitNavMs:   int(cfg.WaitNavDelay / time.Millisecond),
		},
	}

	return fc
}

// legacyFileConfig is the old flat structure for backward compatibility.
type legacyFileConfig struct {
	Port              string `json:"port"`
	InstancePortStart *int   `json:"instancePortStart,omitempty"`
	InstancePortEnd   *int   `json:"instancePortEnd,omitempty"`
	Token             string `json:"token,omitempty"`
	AllowEvaluate     *bool  `json:"allowEvaluate,omitempty"`
	AllowMacro        *bool  `json:"allowMacro,omitempty"`
	AllowScreencast   *bool  `json:"allowScreencast,omitempty"`
	AllowDownload     *bool  `json:"allowDownload,omitempty"`
	AllowUpload       *bool  `json:"allowUpload,omitempty"`
	StateDir          string `json:"stateDir"`
	ProfileDir        string `json:"profileDir"`
	Headless          *bool  `json:"headless,omitempty"`
	NoRestore         bool   `json:"noRestore"`
	MaxTabs           *int   `json:"maxTabs,omitempty"`
	TimeoutSec        int    `json:"timeoutSec,omitempty"`
	NavigateSec       int    `json:"navigateSec,omitempty"`
}

// convertLegacyConfig converts flat config to nested structure.
func convertLegacyConfig(lc *legacyFileConfig) *FileConfig {
	fc := &FileConfig{}

	// Server
	fc.Server.Port = lc.Port
	fc.Server.Token = lc.Token
	fc.Server.StateDir = lc.StateDir

	// Browser / instance defaults
	if lc.Headless != nil {
		if *lc.Headless {
			fc.InstanceDefaults.Mode = "headless"
		} else {
			fc.InstanceDefaults.Mode = "headed"
		}
	}
	fc.InstanceDefaults.MaxTabs = lc.MaxTabs
	if lc.NoRestore {
		b := true
		fc.InstanceDefaults.NoRestore = &b
	}

	// Profiles
	if lc.ProfileDir != "" {
		fc.Profiles.BaseDir = filepath.Dir(lc.ProfileDir)
		fc.Profiles.DefaultProfile = filepath.Base(lc.ProfileDir)
	}

	// Security
	fc.Security.AllowEvaluate = lc.AllowEvaluate
	fc.Security.AllowMacro = lc.AllowMacro
	fc.Security.AllowScreencast = lc.AllowScreencast
	fc.Security.AllowDownload = lc.AllowDownload
	fc.Security.AllowUpload = lc.AllowUpload

	// Timeouts
	fc.Timeouts.ActionSec = lc.TimeoutSec
	fc.Timeouts.NavigateSec = lc.NavigateSec

	// Multi-instance
	fc.MultiInstance.InstancePortStart = lc.InstancePortStart
	fc.MultiInstance.InstancePortEnd = lc.InstancePortEnd

	return fc
}

// isLegacyConfig detects if JSON is flat (legacy) or nested (new).
func isLegacyConfig(data []byte) bool {
	var probe map[string]json.RawMessage
	if err := json.Unmarshal(data, &probe); err != nil {
		return false
	}

	// If any new nested keys exist, it's new format
	newKeys := []string{"server", "browser", "instanceDefaults", "profiles", "multiInstance", "security", "attach", "timeouts"}
	for _, key := range newKeys {
		if _, has := probe[key]; has {
			return false
		}
	}

	// If "port" or "headless" exist at top level, it's legacy
	if _, hasPort := probe["port"]; hasPort {
		return true
	}
	if _, hasHeadless := probe["headless"]; hasHeadless {
		return true
	}

	return false
}

func modeToHeadless(mode string, fallback bool) bool {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "":
		return fallback
	case "headless":
		return true
	case "headed":
		return false
	default:
		return fallback
	}
}
