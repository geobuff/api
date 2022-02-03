package repo

import (
	"database/sql"
)

// Quiz is the database object for a quiz entry.
type Quiz struct {
	ID             int           `json:"id"`
	TypeID         int           `json:"typeId"`
	BadgeID        int           `json:"badgeId"`
	ContinentID    sql.NullInt64 `json:"continentId"`
	Country        string        `json:"country"`
	Singular       string        `json:"singular"`
	Name           string        `json:"name"`
	MaxScore       int           `json:"maxScore"`
	Time           int           `json:"time"`
	MapSVG         string        `json:"mapSVG"`
	ImageURL       string        `json:"imageUrl"`
	Verb           string        `json:"verb"`
	APIPath        string        `json:"apiPath"`
	Route          string        `json:"route"`
	HasLeaderboard bool          `json:"hasLeaderboard"`
	HasGrouping    bool          `json:"hasGrouping"`
	HasFlags       bool          `json:"hasFlags"`
	Enabled        bool          `json:"enabled"`
}

type TriviaQuizDto struct {
	Name     string `json:"name"`
	MapSVG   string `json:"mapSVG"`
	APIPath  string `json:"apiPath"`
	Singular string `json:"singular"`
	Country  string `json:"country"`
}

type QuizzesFilterDto struct {
	Filter string `json:"filter"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}

// GetQuizzes returns all quizzes.
var GetQuizzes = func(filter QuizzesFilterDto) ([]Quiz, error) {
	statement := "SELECT q.id, q.typeid, q.badgeid, q.continentid, q.country, q.singular, q.name, q.maxscore, q.time, q.mapsvg, q.imageurl, q.verb, q.apipath, q.route, q.hasleaderboard, q.hasgrouping, q.hasflags, q.enabled FROM quizzes q JOIN continents c ON c.id = q.continentId WHERE q.name ILIKE '%' || $1 || '%' OR c.name ILIKE '%' || $1 || '%' ORDER BY q.id LIMIT $2 OFFSET $3;"
	rows, err := Connection.Query(statement, filter.Filter, filter.Limit, filter.Page*filter.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quizzes = []Quiz{}
	for rows.Next() {
		var quiz Quiz
		if err = rows.Scan(&quiz.ID, &quiz.TypeID, &quiz.BadgeID, &quiz.ContinentID, &quiz.Country, &quiz.Singular, &quiz.Name, &quiz.MaxScore, &quiz.Time, &quiz.MapSVG, &quiz.ImageURL, &quiz.Verb, &quiz.APIPath, &quiz.Route, &quiz.HasLeaderboard, &quiz.HasGrouping, &quiz.HasFlags, &quiz.Enabled); err != nil {
			return nil, err
		}
		quizzes = append(quizzes, quiz)
	}
	return quizzes, rows.Err()
}

var GetFirstQuizID = func(offset int) (int, error) {
	statement := "SELECT id FROM quizzes LIMIT 1 OFFSET $1;"
	var id int
	err := Connection.QueryRow(statement, offset).Scan(&id)
	return id, err
}

// GetQuiz return a quiz with the matching id.
var GetQuiz = func(id int) (Quiz, error) {
	statement := "SELECT * FROM quizzes WHERE id = $1;"
	var quiz Quiz
	err := Connection.QueryRow(statement, id).Scan(&quiz.ID, &quiz.TypeID, &quiz.BadgeID, &quiz.ContinentID, &quiz.Country, &quiz.Singular, &quiz.Name, &quiz.MaxScore, &quiz.Time, &quiz.MapSVG, &quiz.ImageURL, &quiz.Verb, &quiz.APIPath, &quiz.Route, &quiz.HasLeaderboard, &quiz.HasGrouping, &quiz.HasFlags, &quiz.Enabled)
	return quiz, err
}

// GetQuizID gets the quiz id based on name.
func GetQuizID(name string) (int, error) {
	statement := "SELECT id FROM quizzes WHERE name ILIKE '%' || $1 || '%';"
	var id int
	err := Connection.QueryRow(statement, name).Scan(&id)
	return id, err
}

var ScoreExceedsMax = func(quizID, score int) (bool, error) {
	quiz, err := GetQuiz(quizID)
	if err != nil {
		return false, err
	}
	return score > quiz.MaxScore, nil
}

func ToggleQuizEnabled(quizID int) (int, error) {
	statement := "UPDATE quizzes set enabled = NOT enabled WHERE id = $1 RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, quizID).Scan(&id)
	return id, err
}

func getCountryRegionQuizzes() ([]TriviaQuizDto, error) {
	statement := "SELECT country, singular, name, mapsvg, apipath FROM quizzes WHERE typeid = $1 AND name NOT LIKE '%World%' AND name NOT LIKE '%Countries%' AND name NOT LIKE '%US States%';"
	rows, err := Connection.Query(statement, QUIZ_TYPE_MAP)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quizzes = []TriviaQuizDto{}
	for rows.Next() {
		var quiz TriviaQuizDto
		if err = rows.Scan(&quiz.Country, &quiz.Singular, &quiz.Name, &quiz.MapSVG, &quiz.APIPath); err != nil {
			return nil, err
		}
		quizzes = append(quizzes, quiz)
	}
	return quizzes, rows.Err()
}

func getFlagRegionQuizzes() ([]TriviaQuizDto, error) {
	statement := "SELECT country, singular, name, mapsvg, apipath FROM quizzes WHERE typeId = $1 AND name NOT LIKE '%World%';"
	rows, err := Connection.Query(statement, QUIZ_TYPE_FLAG)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quizzes = []TriviaQuizDto{}
	for rows.Next() {
		var quiz TriviaQuizDto
		if err = rows.Scan(&quiz.Country, &quiz.Singular, &quiz.Name, &quiz.MapSVG, &quiz.APIPath); err != nil {
			return nil, err
		}
		quizzes = append(quizzes, quiz)
	}
	return quizzes, rows.Err()
}

func getWorldQuizCount(badgeID int) (int, error) {
	statement := "SELECT COUNT(id) FROM quizzes WHERE badgeid = $1;"
	var count int
	err := Connection.QueryRow(statement, badgeID).Scan(&count)
	return count, err
}

func getContinentQuizCount(continentID int) (int, error) {
	statement := "SELECT COUNT(id) FROM quizzes WHERE continentid = $1;"
	var count int
	err := Connection.QueryRow(statement, continentID).Scan(&count)
	return count, err
}
