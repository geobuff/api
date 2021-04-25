package config

import "testing"

func TestLoad(t *testing.T) {
	tt := []struct {
		name     string
		filename string
		err      bool
	}{
		{
			name:     "valid filename",
			filename: "../config.json",
			err:      false,
		},
		{
			name:     "invalid filename",
			filename: "test.json",
			err:      true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := Load(tc.filename)
			if err != nil && tc.err {
				t.Fatal("expected error; got none")
			}
		})
	}
}
