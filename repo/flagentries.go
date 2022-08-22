package repo

type FlagEntry struct {
	ID      int    `json:"id"`
	GroupID int    `json:"groupId"`
	Code    string `json:"code"`
	Url     string `json:"url"`
}

type CreateFlagEntryDto struct {
	Code string `json:"code"`
	Url  string `json:"url"`
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

func CreateFlagEntry(groupId int, entry CreateFlagEntryDto) error {
	statement := "INSERT INTO flagEntries (groupId, code, url) VALUES ($1, $2, $3) RETURNING id;"
	var id string
	return Connection.QueryRow(statement, groupId, entry.Code, entry.Url).Scan(&id)
}
