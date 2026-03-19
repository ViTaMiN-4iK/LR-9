package main

import "testing"

func TestProcessMessage(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "any message returns greeting",
			input:    "ping",
			expected: "Hello from Go",
		},
		{
			name:     "empty message returns greeting",
			input:    "",
			expected: "Hello from Go",
		},
		{
			name:     "long message returns greeting",
			input:    "very long test message here",
			expected: "Hello from Go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessMessage(tt.input)
			if result != tt.expected {
				t.Errorf("ProcessMessage(%q) = %q; want %q",
					tt.input, result, tt.expected)
			}
		})
	}
}

func TestCleanMessage(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "simple message",
			input:    []byte("ping"),
			expected: "ping",
		},
		{
			name:     "message with newline",
			input:    []byte("ping\n"),
			expected: "ping",
		},
		{
			name:     "message with spaces",
			input:    []byte("  ping  "),
			expected: "ping",
		},
		{
			name:     "empty bytes",
			input:    []byte{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanMessage(tt.input)
			if result != tt.expected {
				t.Errorf("CleanMessage(%q) = %q; want %q",
					tt.input, result, tt.expected)
			}
		})
	}
}
