package config

import "testing"

func TestLoad(t *testing.T) {
	tt := []struct {
		name        string
		filename    string
		shouldError bool
	}{
		{
			name:        "valid filename",
			filename:    "../config.json",
			shouldError: false,
		},
		{
			name:        "invalid filename",
			filename:    "test.json",
			shouldError: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := Load(tc.filename)
			if err != nil && !tc.shouldError {
				t.Fatalf("unexpected error %v", err)
			}
		})
	}
}
