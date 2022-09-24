package repo

type MappingGroup struct {
	ID    int    `json:"id"`
	Key   string `json:"key"`
	Label string `json:"label"`
}

type CreateMappingsDto struct {
	Key     string                  `json:"key"`
	Label   string                  `json:"label"`
	Entries []CreateMappingEntryDto `json:"entries"`
}

type UpdateMappingDto struct {
	Label   string                  `json:"label"`
	Entries []UpdateMappingEntryDto `json:"entries"`
}

type MappingsWithoutFlagDto struct {
	Key     string      `json:"key"`
	Entries []FlagEntry `json:"entries"`
}

func GetMappingGroups() ([]MappingGroup, error) {
	rows, err := Connection.Query("SELECT * from mappingGroups ORDER BY key ASC;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups = []MappingGroup{}
	for rows.Next() {
		var group MappingGroup
		if err = rows.Scan(&group.ID, &group.Key, &group.Label); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, rows.Err()
}

func CreateMappings(mappings CreateMappingsDto) error {
	var id int
	statement := "INSERT INTO mappinggroups (key, label) values ($1, $2) RETURNING id;"
	if err := Connection.QueryRow(statement, mappings.Key, mappings.Label).Scan(&id); err != nil {
		return err
	}

	for _, entry := range mappings.Entries {
		if err := CreateMappingEntry(id, entry); err != nil {
			return err
		}
	}
	return nil
}

func GetMappingsWithoutFlags() ([]MappingsWithoutFlagDto, error) {
	rows, err := Connection.Query("SELECT m.key FROM mappinggroups m LEFT JOIN flaggroups f ON f.key = m.key WHERE f.id IS NULL;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results = []MappingsWithoutFlagDto{}
	for rows.Next() {
		var result MappingsWithoutFlagDto
		if err = rows.Scan(&result.Key); err != nil {
			return nil, err
		}

		rows, err := Connection.Query("SELECT e.code FROM mappingentries e JOIN mappinggroups g ON g.id = e.groupid WHERE g.key = $1;", result.Key)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var entries = []FlagEntry{}
		for rows.Next() {
			var entry FlagEntry
			if err = rows.Scan(&entry.Code); err != nil {
				return nil, err
			}
			entries = append(entries, entry)
		}

		result.Entries = entries
		results = append(results, result)
	}

	return results, rows.Err()
}

func UpdateMapping(key string, update UpdateMappingDto) error {
	existingMappingEntries, err := GetMappingEntries(key)
	if err != nil {
		return err
	}

	svgMap, err := GetMapUsingKey(key)
	if err != nil {
		return err
	}

	for _, mapEntry := range svgMap.Elements {
		for _, mappingEntry := range existingMappingEntries {
			if mappingEntry.SVGName == mapEntry.Name {
				var updatedEntry UpdateMapElementDto
				for _, val := range update.Entries {
					if val.ID == mappingEntry.ID {
						updatedEntry = UpdateMapElementDto{
							Name:      val.Name,
							ElementID: val.Code,
						}
					}
				}

				if err := UpdateMapElement(mapEntry.EntryID, updatedEntry); err != nil {
					return err
				}
			}
		}
	}

	for _, entry := range update.Entries {
		if err := UpdateMappingEntry(entry); err != nil {
			return err
		}
	}

	var id int
	return Connection.QueryRow("UPDATE mappingGroups set label = $1 RETURNING id;", update.Label).Scan(&id)
}

func GetMappingGroupId(key string) (int, error) {
	var id int
	err := Connection.QueryRow("SELECT id from mappingGroups WHERE key = $1;", key).Scan(&id)
	return id, err
}

func DeleteMappingGroup(groupId int) error {
	var id int
	return Connection.QueryRow("DELETE FROM mappingGroups where id = $1 RETURNING id;", groupId).Scan(&id)
}

func DeleteMapping(key string) error {
	mappingGroupId, err := GetMappingGroupId(key)
	if err != nil {
		return err
	}

	if err = DeleteMappingEntries(mappingGroupId); err != nil {
		return err
	}

	return DeleteMappingGroup(mappingGroupId)
}
