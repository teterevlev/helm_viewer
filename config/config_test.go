package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name     string
		envPort  string
		expected string
	}{
		{
			name:     "Default port when PORT env is not set",
			envPort:  "",
			expected: "8080",
		},
		{
			name:     "Custom port from PORT env",
			envPort:  "3000",
			expected: "3000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save current env and restore after test
			originalPort := os.Getenv("PORT")
			defer os.Setenv("PORT", originalPort)

			// Set test env
			os.Setenv("PORT", tt.envPort)

			// Create new config
			config := NewConfig()

			// Check if port matches expected
			require.Equal(t, tt.expected, config.Port)
		})
	}
}
