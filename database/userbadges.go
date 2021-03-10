package database

// UserBadge is the database object for a user badge entry.
type UserBadge struct {
	ID       int `json:"id"`
	UserID   int `json:"userId"`
	BadgeID  int `json:"badgeId"`
	Progress int `json:"progress"`
}

// GetUserBadges gets all badge entries for a user.
func GetUserBadges(userID int) ([]UserBadge, error) {
	query := "SELECT * FROM userbadges WHERE userId = $1;"
	rows, err := Connection.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userBadges = []UserBadge{}
	for rows.Next() {
		var userBadge UserBadge
		if err = rows.Scan(&userBadge.ID, &userBadge.UserID, &userBadge.BadgeID, &userBadge.Progress); err != nil {
			return nil, err
		}
		userBadges = append(userBadges, userBadge)
	}
	return userBadges, rows.Err()
}
