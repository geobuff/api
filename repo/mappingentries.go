package repo

import (
	"database/sql"
	"strings"

	"github.com/lib/pq"
)

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
	FlagUrl          sql.NullString  `json:"flagUrl"`
	SVGName          string          `json:"svgName"`
	AlternativeNames *pq.StringArray `json:"alternativeNames"`
	Prefixes         *pq.StringArray `json:"prefixes"`
	Grouping         string          `json:"grouping"`
}

type CreateMappingEntryDto struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type UpdateMappingEntryDto struct {
	ID               int             `json:"id"`
	Name             string          `json:"name"`
	Code             string          `json:"code"`
	SVGName          string          `json:"svgName"`
	AlternativeNames *pq.StringArray `json:"alternativeNames"`
	Prefixes         *pq.StringArray `json:"prefixes"`
	Grouping         string          `json:"grouping"`
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

func CreateMappingEntry(groupId int, entry CreateMappingEntryDto) error {
	var id int
	statement := "INSERT INTO mappingentries (groupid, name, code, svgname, alternativenames, prefixes, grouping) values ($1, $2, $3, $4, $5, $6, $7) RETURNING id;"
	return Connection.QueryRow(statement, groupId, strings.ToLower(entry.Name), entry.Code, entry.Name, pq.Array([]string{}), pq.Array([]string{}), "").Scan(&id)
}

func DeleteMappingEntries(groupId int) error {
	var id int
	return Connection.QueryRow("DELETE FROM mappingentries where groupId = $1 RETURNING id;", groupId).Scan(&id)
}

func UpdateMappingEntry(entry UpdateMappingEntryDto) error {
	statement := "UPDATE mappingentries SET name = $2, code = $3, svgname = $4, alternativenames = $5, prefixes = $6, grouping = $7 WHERE id = $1 RETURNING id;"
	var id int
	return Connection.QueryRow(statement, entry.ID, entry.Name, entry.Code, entry.SVGName, pq.Array(entry.AlternativeNames), pq.Array(entry.Prefixes), entry.Grouping).Scan(&id)
}
