package repo

// Badge is the database object for a badge entry.
type Badge struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Total       int    `json:"total"`
	ImageUrl    string `json:"imageUrl"`
	Background  string `json:"background"`
	Border      string `json:"border"`
}

// GetBadges returns all badges.
var GetBadges = func() ([]Badge, error) {
	rows, err := Connection.Query("SELECT * FROM badges;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var badges = []Badge{}
	for rows.Next() {
		var badge Badge
		if err = rows.Scan(&badge.ID, &badge.Name, &badge.Description, &badge.Total, &badge.ImageUrl, &badge.Background, &badge.Border); err != nil {
			return nil, err
		}
		badges = append(badges, badge)
	}
	return badges, rows.Err()
}
