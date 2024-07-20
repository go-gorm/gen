package main

import (
	"testing"

	"gorm.io/gen"
)

func TestParseMode(t *testing.T) {
	tests := []struct {
		name     string
		modes    []string
		expected gen.GenerateMode
	}{
		{
			name:     "Empty modes",
			modes:    []string{},
			expected: gen.GenerateMode(0),
		},
		{
			name:     "Single mode",
			modes:    []string{"WithDefaultQuery"},
			expected: gen.WithDefaultQuery,
		},
		{
			name:     "Multiple modes",
			modes:    []string{"WithDefaultQuery", "WithoutContext", "WithQueryInterface"},
			expected: gen.WithDefaultQuery | gen.WithoutContext | gen.WithQueryInterface,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseMode(tt.modes...)
			if result != tt.expected {
				t.Errorf("parseMode(%v) = %v, want %v", tt.modes, result, tt.expected)
			}
		})
	}
}
