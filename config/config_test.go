package config

import "testing"

func TestLoad(t *testing.T) {
	tt := []struct {
		name     string
		filename string
		err      string
	}{
		{
			name:     "valid filename",
			filename: "../config.json",
			err:      "",
		},
		{
			name:     "invalid filename",
			filename: "test.json",
			err:      "open test.json: The system cannot find the file specified.",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := Load(tc.filename)
			if err != nil && err.Error() != tc.err {
				t.Fatalf("expected error %v; got %v", tc.err, err.Error())
			}
		})
	}
}
