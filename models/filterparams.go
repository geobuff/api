package models

// GetEntriesFilterParams contain the filter fields for the leaderboard.
type GetEntriesFilterParams struct {
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Range string `json:"range"`
	User  string `json:"user"`
}
