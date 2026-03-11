package config

import (
	"testing"
)

func TestMaskToken(t *testing.T) {
	tests := []struct {
		token string
		want  string
	}{
		{"", "(none)"},
		{"short", "***"},
		{"very-long-token-secret", "very...cret"},
	}

	for _, tt := range tests {
		if got := MaskToken(tt.token); got != tt.want {
			t.Errorf("MaskToken(%q) = %v, want %v", tt.token, got, tt.want)
		}
	}
}

func TestGenerateAuthToken(t *testing.T) {
	token, err := GenerateAuthToken()
	if err != nil {
		t.Fatalf("GenerateAuthToken() error = %v", err)
	}
	if len(token) != 48 {
		t.Fatalf("GenerateAuthToken() len = %d, want 48", len(token))
	}
	for _, r := range token {
		switch {
		case r >= '0' && r <= '9':
		case r >= 'a' && r <= 'f':
		default:
			t.Fatalf("GenerateAuthToken() produced non-hex rune %q", r)
		}
	}
}

func TestListenAddr(t *testing.T) {
	cfg := &RuntimeConfig{Bind: "127.0.0.1", Port: "9867"}
	if got := cfg.ListenAddr(); got != "127.0.0.1:9867" {
		t.Errorf("expected 127.0.0.1:9867, got %s", got)
	}

	cfg = &RuntimeConfig{Bind: "0.0.0.0", Port: "8080"}
	if got := cfg.ListenAddr(); got != "0.0.0.0:8080" {
		t.Errorf("expected 0.0.0.0:8080, got %s", got)
	}
}
