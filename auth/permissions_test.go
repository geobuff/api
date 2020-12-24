package auth

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHasPermission(t *testing.T) {
	tt := []struct {
		name       string
		token      string
		permission string
		want       bool
		err        string
	}{
		{
			name:       "valid token, valid permission",
			token:      "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InpuNmd6ZnBxZDRZbU5hMUZabzFTSiJ9.eyJodHRwOi8vZ2VvYnVmZi5jb20vdXNlcm5hbWUiOiJtcnNjcnViIiwiaXNzIjoiaHR0cHM6Ly9kZXYtZ2VvYnVmZi5hdS5hdXRoMC5jb20vIiwic3ViIjoiYXV0aDB8NWZjNmM5YWRiMTBkNmQwMDZmNjEwNTk2IiwiYXVkIjpbImh0dHBzOi8vZGV2LWdlb2J1ZmYtYXBpLmNvbSIsImh0dHBzOi8vZGV2LWdlb2J1ZmYuYXUuYXV0aDAuY29tL3VzZXJpbmZvIl0sImlhdCI6MTYwODc5Njk5NiwiZXhwIjoxNjA4ODgzMzk2LCJhenAiOiJpUXRrTmc2cHYwZzhsQXFFQUprUkpxMjFONzNzaG9UMSIsInNjb3BlIjoib3BlbmlkIHByb2ZpbGUgZW1haWwiLCJwZXJtaXNzaW9ucyI6WyJyZWFkOnNjb3JlcyIsInJlYWQ6dXNlcnMiLCJ3cml0ZTpsZWFkZXJib2FyZCIsIndyaXRlOnNjb3JlcyIsIndyaXRlOnVzZXJzIl19.YyQUeX1ttxPF2CUGOlE9m3dBBbi8kDbFXq7-2dld3YPctHJb_wgRxc10RFcyryLkSOdbWFdELoF0sDRh2miEYchdrit7Ac0rTufWkn7pfv4SV8e9yYmlCV-cVBtDjKPDQP2e7TDSuoXOOVBcbdlhyro524bEqAibzWPqfh5apn1HeP_tvbHqdZRdrRZOVu1W89f7uaoua9-DmBASSgeO3KixLYREHAPAI2kFgt2YmWZ-Y_474A2KoIAJgwrbhG094sbDiBuuhRRl2GUURLqbxRB0gIRS1ut4J2DGvkr59t5ev789XIWc4zI3rHvBLJ0hGDyi66Ip8TZW3D-pkNCl6Q",
			permission: "read:scores",
			want:       true,
			err:        "",
		},
		{
			name:       "valid token, missing permission",
			token:      "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InpuNmd6ZnBxZDRZbU5hMUZabzFTSiJ9.eyJodHRwOi8vZ2VvYnVmZi5jb20vdXNlcm5hbWUiOiJtcnNjcnViIiwiaXNzIjoiaHR0cHM6Ly9kZXYtZ2VvYnVmZi5hdS5hdXRoMC5jb20vIiwic3ViIjoiYXV0aDB8NWZjNmM5YWRiMTBkNmQwMDZmNjEwNTk2IiwiYXVkIjpbImh0dHBzOi8vZGV2LWdlb2J1ZmYtYXBpLmNvbSIsImh0dHBzOi8vZGV2LWdlb2J1ZmYuYXUuYXV0aDAuY29tL3VzZXJpbmZvIl0sImlhdCI6MTYwODc5NzM5MSwiZXhwIjoxNjA4ODgzNzkxLCJhenAiOiJpUXRrTmc2cHYwZzhsQXFFQUprUkpxMjFONzNzaG9UMSIsInNjb3BlIjoib3BlbmlkIHByb2ZpbGUgZW1haWwiLCJwZXJtaXNzaW9ucyI6WyJyZWFkOnVzZXJzIiwid3JpdGU6bGVhZGVyYm9hcmQiLCJ3cml0ZTpzY29yZXMiLCJ3cml0ZTp1c2VycyJdfQ.McI54hJHhMKBoonfCmX4lXScRkK2fxPa4Ej_QkkV_12bzTGeFmtvQsGYnYsnOOZSwRsdsG5JvA2c39l4IQU_7DjY4OFkgSuJnrLfQsInfm4799uefYkpx_ELdbClLSitNAZMObM85iYc_tUZwpI5Yas5stnmUAgf6YJMG_cQVWiVnQZ-9lRhXXdrLjxwdJgtK5MS0YHMSsgZE8QmB3SPZysGx5ICqFshqcKEOmre3LdEwk1jFw0JHLROAPcwQF7D1Ba45XHnZnODxS1xLmpWevdMISHyuXaHDDDVRy_JQCk5C8A1bFsULBilclCDWnGcyGNzNLXBZkVvun9-5JTm6g",
			permission: "read:scores",
			want:       false,
			err:        "",
		},
		{
			name:       "invalid token",
			token:      "",
			permission: "read:scores",
			want:       false,
			err:        "token missing or invalid length",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			request, err := http.NewRequest("GET", "http://localhost:8080/api/scores", nil)
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}
			request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tc.token))

			result, err := HasPermission(request, tc.permission)
			if err != nil {
				if err.Error() != tc.err {
					t.Fatalf("expected error %v; got %v", tc.err, err.Error())
				}
			} else if tc.err != "" {
				t.Fatalf("expected error %v; got nil", tc.err)
			}

			if result != tc.want {
				t.Fatalf("got %v; want %v", result, tc.want)
			}
		})
	}
}
