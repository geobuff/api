package utils

import "testing"

func TestUsernameValid(t *testing.T) {
	tt := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "length small than min",
			input:    "as",
			expected: false,
		},
		{
			name:     "length greater than max",
			input:    "aaaaaaaaaaaaaaaaaaaaa",
			expected: false,
		},
		{
			name:     "contains space",
			input:    "test test",
			expected: false,
		},
		{
			name:     "valid username",
			input:    "mrscrub",
			expected: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := UsernameValid(tc.input)
			if result != tc.expected {
				t.Errorf("expected %v; got %v", tc.expected, result)
			}
		})
	}
}

func TestPasswordValid(t *testing.T) {
	tt := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "length shorter than min",
			input:    "test",
			expected: false,
		},
		{
			name:     "correct length, missing upper, lower and number",
			input:    "!!!!!!!!!!!",
			expected: false,
		},
		{
			name:     "correct length, missing upper",
			input:    "testingtest1",
			expected: false,
		},
		{
			name:     "correct length, missing lower",
			input:    "TESTINGTEST1",
			expected: false,
		},
		{
			name:     "correct length, missing number",
			input:    "Testingtest",
			expected: false,
		},
		{
			name:     "happy path",
			input:    "Password1!",
			expected: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := PasswordValid(tc.input)
			if result != tc.expected {
				t.Errorf("expected %v; got %v", result, tc.expected)
			}
		})
	}
}
