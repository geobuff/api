package repo

type FlagGroup struct {
	ID    int    `json:"id"`
	Key   string `json:"key"`
	Label string `json:"label"`
}

type FlagEntry struct {
	ID      int    `json:"id"`
	GroupID int    `json:"groupId"`
	Code    string `json:"code"`
	Url     string `json:"url"`
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

func GetFlagEntries(key string) ([]FlagEntry, error) {
	rows, err := Connection.Query("SELECT e.id, e.groupId, e.code, e.url from flagEntries e JOIN flagGroups g ON g.id = e.groupId WHERE g.key = $1;", key)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries = []FlagEntry{}
	for rows.Next() {
		var entry FlagEntry
		if err = rows.Scan(&entry.ID, &entry.GroupID, &entry.Code, &entry.Url); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}

func GetFlagUrl(code string) (string, error) {
	statement := "SELECT url from flagEntries where code = $1;"
	var url string
	err := Connection.QueryRow(statement, code).Scan(&url)
	return url, err
}
