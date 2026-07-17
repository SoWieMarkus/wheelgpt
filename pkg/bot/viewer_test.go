package bot

import "testing"

func TestMentionViewer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal name",
			input:    "username",
			expected: "username",
		},
		{
			name:     "name with @",
			input:    "@username",
			expected: "username",
		},
		{
			name:     "name with multiple @",
			input:    "@@username",
			expected: "username",
		},
		{
			name:     "name with spaces",
			input:    " user name ",
			expected: "user name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MentionViewer(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
