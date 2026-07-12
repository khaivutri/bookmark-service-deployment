package logger

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestSetLogLevel(t *testing.T) {
	originalLevel := zerolog.GlobalLevel()
	t.Cleanup(func() {
		zerolog.SetGlobalLevel(originalLevel)
	})

	tests := []struct {
		name          string
		levelStr      string
		expectedLevel zerolog.Level
	}{
		{
			name:          "sets debug level",
			levelStr:      "debug",
			expectedLevel: zerolog.DebugLevel,
		},
		{
			name:          "sets info level",
			levelStr:      "info",
			expectedLevel: zerolog.InfoLevel,
		},
		{
			name:          "sets warn level",
			levelStr:      "warn",
			expectedLevel: zerolog.WarnLevel,
		},
		{
			name:          "sets error level",
			levelStr:      "error",
			expectedLevel: zerolog.ErrorLevel,
		},
		{
			name:          "sets no level for invalid input",
			levelStr:      "invalid",
			expectedLevel: zerolog.NoLevel,
		},
		{
			name:          "sets no level for empty input",
			levelStr:      "",
			expectedLevel: zerolog.NoLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetLogLevel(tt.levelStr)

			assert.Equal(t, tt.expectedLevel, zerolog.GlobalLevel())
		})
	}
}
