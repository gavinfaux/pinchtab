package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

var defaultLocalAllowedDomains = []string{"127.0.0.1", "localhost", "::1"}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	h, _ := os.UserHomeDir()
	return h
}

func userConfigDir() string {
	home := homeDir()
	legacyPath := filepath.Join(home, ".pinchtab")

	configDir, err := os.UserConfigDir()
	if err != nil {
		return legacyPath
	}

	newPath := filepath.Join(configDir, "pinchtab")

	legacyExists := dirExists(legacyPath)
	newExists := dirExists(newPath)

	if legacyExists && !newExists {
		return legacyPath
	}

	return newPath
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func (c *RuntimeConfig) ListenAddr() string {
	return c.Bind + ":" + c.Port
}

func GenerateAuthToken() (string, error) {
	buf := make([]byte, 24)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate auth token: %w", err)
	}
	return hex.EncodeToString(buf), nil
}

func MaskToken(t string) string {
	if t == "" {
		return "(none)"
	}
	if len(t) <= 8 {
		return "***"
	}
	return t[:4] + "..." + t[len(t)-4:]
}
