package repo

type FlagGroup struct {
	ID    int    `json:"id"`
	Key   string `json:"key"`
	Label string `json:"label"`
}

type CreateFlagsDto struct {
	Key     string               `json:"key"`
	Label   string               `json:"label"`
	Entries []CreateFlagEntryDto `json:"entries"`
}

func GetFlagGroups() ([]FlagGroup, error) {
	rows, err := Connection.Query("SELECT * from flagGroups;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups = []FlagGroup{}
	for rows.Next() {
		var group FlagGroup
		if err = rows.Scan(&group.ID, &group.Key, &group.Label); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, rows.Err()
}
