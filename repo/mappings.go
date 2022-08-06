package repo

import (
	"github.com/lib/pq"
)

type MappingGroup struct {
	ID    int    `json:"id"`
	Key   string `json:"key"`
	Label string `json:"label"`
}

type MappingEntry struct {
	ID               int             `json:"id"`
	GroupID          int             `json:"groupId"`
	Name             string          `json:"name"`
	Code             string          `json:"code"`
	SVGName          string          `json:"svgName"`
	AlternativeNames *pq.StringArray `json:"alternativeNames"`
	Prefixes         *pq.StringArray `json:"prefixes"`
	Grouping         string          `json:"grouping"`
}

type MappingEntryDto struct {
	ID               int             `json:"id"`
	GroupID          int             `json:"groupId"`
	Name             string          `json:"name"`
	Code             string          `json:"code"`
	FlagUrl          string          `json:"flagUrl"`
	SVGName          string          `json:"svgName"`
	AlternativeNames *pq.StringArray `json:"alternativeNames"`
	Prefixes         *pq.StringArray `json:"prefixes"`
	Grouping         string          `json:"grouping"`
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

func GetMappingEntries(key string) ([]MappingEntryDto, error) {
	rows, err := Connection.Query("SELECT m.id, m.groupid, m.name, m.code, f.url, m.svgname, m.alternativenames, m.prefixes, m.grouping from mappingEntries m JOIN mappingGroups g ON g.id = m.groupId LEFT JOIN flagEntries f ON f.code = m.code WHERE g.key = $1;", key)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries = []MappingEntryDto{}
	for rows.Next() {
		var entry MappingEntryDto
		if err = rows.Scan(&entry.ID, &entry.GroupID, &entry.Name, &entry.Code, &entry.FlagUrl, &entry.SVGName, &entry.AlternativeNames, &entry.Prefixes, &entry.Grouping); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}
