package repo

// Avatar is the database object for a avatar entry.
type Avatar struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ImageUrl   string `json:"imageUrl"`
	Background string `json:"background"`
	Border     string `json:"border"`
}

// GetAvatars returns all avatars.
var GetAvatars = func() ([]Avatar, error) {
	rows, err := Connection.Query("SELECT * FROM avatars;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var avatars = []Avatar{}
	for rows.Next() {
		var avatar Avatar
		if err = rows.Scan(&avatar.ID, &avatar.Name, &avatar.ImageUrl, &avatar.Background, &avatar.Border); err != nil {
			return nil, err
		}
		avatars = append(avatars, avatar)
	}
	return avatars, rows.Err()
}

// GetAvatar returns an avatar with the matching id.
var GetAvatar = func(id int) (Avatar, error) {
	statement := "SELECT * from avatars WHERE id = $1;"
	var avatar Avatar
	err := Connection.QueryRow(statement, id).Scan(&avatar.ID, &avatar.Name, &avatar.ImageUrl, &avatar.Background, &avatar.Border)
	return avatar, err
}
