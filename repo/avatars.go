package repo

type Avatar struct {
	ID                int    `json:"id"`
	TypeID            int    `json:"typeId"`
	CountryCode       string `json:"countryCode"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	PrimaryImageUrl   string `json:"primaryImageUrl"`
	SecondaryImageUrl string `json:"secondaryImageUrl"`
	GridPlacement     int    `json:"gridPlacement"`
}

type AvatarDto struct {
	ID                int    `json:"id"`
	Type              string `json:"type"`
	CountryCode       string `json:"countryCode"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	PrimaryImageUrl   string `json:"primaryImageUrl"`
	SecondaryImageUrl string `json:"secondaryImageUrl"`
	GridPlacement     int    `json:"gridPlacement"`
}

var GetAvatars = func() ([]AvatarDto, error) {
	rows, err := Connection.Query("SELECT a.id, t.name, a.countrycode, a.name, a.description, a.primaryImageUrl, a.secondaryImageUrl, a.gridplacement FROM avatars a JOIN avatarTypes t ON t.id = a.typeid ORDER BY a.gridplacement;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var avatars = []AvatarDto{}
	for rows.Next() {
		var avatar AvatarDto
		if err = rows.Scan(&avatar.ID, &avatar.Type, &avatar.CountryCode, &avatar.Name, &avatar.Description, &avatar.PrimaryImageUrl, &avatar.SecondaryImageUrl, &avatar.GridPlacement); err != nil {
			return nil, err
		}
		avatars = append(avatars, avatar)
	}
	return avatars, rows.Err()
}

var GetAvatar = func(id int) (AvatarDto, error) {
	statement := "SELECT a.id, t.name, a.countrycode, a.name, a.description, a.primaryImageUrl, a.secondaryImageUrl, a.gridplacement FROM avatars a JOIN avatarTypes t ON t.id = a.typeid WHERE a.id = $1;"
	var avatar AvatarDto
	err := Connection.QueryRow(statement, id).Scan(&avatar.ID, &avatar.Type, &avatar.CountryCode, &avatar.Name, &avatar.Description, &avatar.PrimaryImageUrl, &avatar.SecondaryImageUrl, &avatar.GridPlacement)
	return avatar, err
}
