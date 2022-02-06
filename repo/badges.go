package repo

import (
	"database/sql"
	"errors"
)

// Badge is the database object for a badge entry.
type Badge struct {
	ID          int           `json:"id"`
	TypeID      int           `json:"typeId"`
	ContinentID sql.NullInt64 `json:"continentId"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	ImageUrl    string        `json:"imageUrl"`
	Background  string        `json:"background"`
	Border      string        `json:"border"`
}

type BadgeDto struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageUrl    string `json:"imageUrl"`
	Background  string `json:"background"`
	Border      string `json:"border"`
	Progress    int    `json:"progress"`
	Total       int    `json:"total"`
}

type CreateQuizBadgeDto struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetBadges() ([]CreateQuizBadgeDto, error) {
	rows, err := Connection.Query("SELECT id, name from badges;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var badges = []CreateQuizBadgeDto{}
	for rows.Next() {
		var badge CreateQuizBadgeDto
		if err = rows.Scan(&badge.ID, &badge.Name); err != nil {
			return nil, err
		}
		badges = append(badges, badge)
	}
	return badges, rows.Err()
}

// GetUserBadges returns all badges for a user.
var GetUserBadges = func(userId int) ([]BadgeDto, error) {
	leaderboardEntries, err := GetUserLeaderboardEntries(userId)
	if err != nil {
		return nil, err
	}

	rows, err := Connection.Query("SELECT * FROM badges;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var badges = []BadgeDto{}
	for rows.Next() {
		var badge Badge
		if err = rows.Scan(&badge.ID, &badge.TypeID, &badge.ContinentID, &badge.Name, &badge.Description, &badge.ImageUrl, &badge.Background, &badge.Border); err != nil {
			return nil, err
		}

		total, err := getTotal(badge.TypeID, badge.ID, badge.ContinentID)
		if err != nil {
			return nil, err
		}

		dto := BadgeDto{
			ID:          badge.ID,
			Name:        badge.Name,
			Description: badge.Description,
			ImageUrl:    badge.ImageUrl,
			Background:  badge.Background,
			Border:      badge.Border,
			Progress:    getProgress(leaderboardEntries, badge.ID, badge.TypeID),
			Total:       total,
		}

		badges = append(badges, dto)
	}
	return badges, rows.Err()
}

func getTotal(typeID, badgeID int, continentID sql.NullInt64) (int, error) {
	switch typeID {
	case BADGE_TYPE_ONE_OFF:
		return 1, nil
	case BADGE_TYPE_WORLD:
		return getWorldQuizCount(badgeID)
	case BADGE_TYPE_CONTINENT:
		return getContinentQuizCount(int(continentID.Int64))
	default:
		return 0, errors.New("invalid type id passed to getTotal")
	}
}

func getProgress(entries []UserLeaderboardEntryDto, badgeId, typeID int) int {
	if typeID == BADGE_TYPE_ONE_OFF {
		if len(entries) > 0 {
			return 1
		}
		return 0
	}

	var count int
	for _, val := range entries {
		if val.BadgeID == badgeId {
			count = count + 1
		}
	}

	return count
}
