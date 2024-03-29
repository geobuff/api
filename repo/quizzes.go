package repo

import (
	"database/sql"
	"fmt"
)

type Quiz struct {
	ID             int           `json:"id"`
	TypeID         int           `json:"typeId"`
	BadgeID        sql.NullInt64 `json:"badgeId"`
	ContinentID    sql.NullInt64 `json:"continentId"`
	Country        string        `json:"country"`
	Singular       string        `json:"singular"`
	Name           string        `json:"name"`
	MaxScore       int           `json:"maxScore"`
	Time           int           `json:"time"`
	MapSVG         string        `json:"mapSVG"`
	ImageURL       string        `json:"imageUrl"`
	Plural         string        `json:"plural"`
	APIPath        string        `json:"apiPath"`
	Route          string        `json:"route"`
	HasLeaderboard bool          `json:"hasLeaderboard"`
	HasGrouping    bool          `json:"hasGrouping"`
	HasFlags       bool          `json:"hasFlags"`
	Enabled        bool          `json:"enabled"`
}

type QuizDto struct {
	ID             int           `json:"id"`
	TypeID         int           `json:"typeId"`
	BadgeID        sql.NullInt64 `json:"badgeId"`
	ContinentID    sql.NullInt64 `json:"continentId"`
	Country        string        `json:"country"`
	Singular       string        `json:"singular"`
	Name           string        `json:"name"`
	MaxScore       int           `json:"maxScore"`
	Time           int           `json:"time"`
	MapName        string        `json:"mapName"`
	Map            MapDto        `json:"map"`
	ImageURL       string        `json:"imageUrl"`
	Plural         string        `json:"plural"`
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
	Page              int    `json:"page"`
	Limit             int    `json:"limit"`
	OrderByPopularity bool   `json:"orderByPopularity"`
	Filter            string `json:"filter"`
}

type CreateQuizDto struct {
	TypeID         int           `json:"typeId"`
	BadgeID        sql.NullInt64 `json:"badgeId"`
	ContinentID    sql.NullInt64 `json:"continentId"`
	Country        string        `json:"country"`
	Singular       string        `json:"singular"`
	Name           string        `json:"name"`
	MaxScore       int           `json:"maxScore"`
	Time           int           `json:"time"`
	MapSVG         string        `json:"mapSVG"`
	ImageURL       string        `json:"imageUrl"`
	Plural         string        `json:"plural"`
	APIPath        string        `json:"apiPath"`
	Route          string        `json:"route"`
	HasLeaderboard bool          `json:"hasLeaderboard"`
	HasGrouping    bool          `json:"hasGrouping"`
	HasFlags       bool          `json:"hasFlags"`
	Enabled        bool          `json:"enabled"`
}

type UpdateQuizDto struct {
	TypeID         int           `json:"typeId"`
	BadgeID        sql.NullInt64 `json:"badgeId"`
	ContinentID    sql.NullInt64 `json:"continentId"`
	Country        string        `json:"country"`
	Singular       string        `json:"singular"`
	Name           string        `json:"name"`
	MaxScore       int           `json:"maxScore"`
	Time           int           `json:"time"`
	MapSVG         string        `json:"mapSVG"`
	ImageURL       string        `json:"imageUrl"`
	Plural         string        `json:"plural"`
	APIPath        string        `json:"apiPath"`
	Route          string        `json:"route"`
	HasLeaderboard bool          `json:"hasLeaderboard"`
	HasGrouping    bool          `json:"hasGrouping"`
	HasFlags       bool          `json:"hasFlags"`
	Enabled        bool          `json:"enabled"`
}

var GetQuizzes = func(filter QuizzesFilterDto) ([]Quiz, error) {
	statement := "SELECT q.id, q.typeid, q.badgeid, q.continentid, q.country, q.singular, q.name, q.maxscore, q.time, q.mapsvg, q.imageurl, q.plural, q.apipath, q.route, q.hasleaderboard, q.hasgrouping, q.hasflags, q.enabled FROM quizzes q JOIN quizType t ON t.id = q.typeId LEFT JOIN quizPlays p ON q.id = p.quizId WHERE q.name ILIKE '%' || $1 || '%' OR t.name ILIKE '%' || $1 || '%' OR q.country ILIKE '%' || $1 || '%' "

	if filter.OrderByPopularity {
		statement = statement + "ORDER BY p.plays, q.country NULLS FIRST LIMIT $2 OFFSET $3;"
	} else {
		statement = statement + "ORDER BY q.country NULLS FIRST, q.maxscore DESC LIMIT $2 OFFSET $3;"
	}

	rows, err := Connection.Query(statement, filter.Filter, filter.Limit, filter.Page*filter.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quizzes = []Quiz{}
	for rows.Next() {
		var quiz Quiz
		if err = rows.Scan(&quiz.ID, &quiz.TypeID, &quiz.BadgeID, &quiz.ContinentID, &quiz.Country, &quiz.Singular, &quiz.Name, &quiz.MaxScore, &quiz.Time, &quiz.MapSVG, &quiz.ImageURL, &quiz.Plural, &quiz.APIPath, &quiz.Route, &quiz.HasLeaderboard, &quiz.HasGrouping, &quiz.HasFlags, &quiz.Enabled); err != nil {
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

var GetQuiz = func(id int) (Quiz, error) {
	statement := "SELECT * FROM quizzes WHERE id = $1;"
	var quiz Quiz
	err := Connection.QueryRow(statement, id).Scan(&quiz.ID, &quiz.TypeID, &quiz.BadgeID, &quiz.ContinentID, &quiz.Country, &quiz.Singular, &quiz.Name, &quiz.MaxScore, &quiz.Time, &quiz.MapSVG, &quiz.ImageURL, &quiz.Plural, &quiz.APIPath, &quiz.Route, &quiz.HasLeaderboard, &quiz.HasGrouping, &quiz.HasFlags, &quiz.Enabled)
	return quiz, err
}

var GetQuizByRoute = func(route string) (QuizDto, error) {
	statement := "SELECT * FROM quizzes WHERE route = $1;"
	var quiz QuizDto
	err := Connection.QueryRow(statement, route).Scan(&quiz.ID, &quiz.TypeID, &quiz.BadgeID, &quiz.ContinentID, &quiz.Country, &quiz.Singular, &quiz.Name, &quiz.MaxScore, &quiz.Time, &quiz.MapName, &quiz.ImageURL, &quiz.Plural, &quiz.APIPath, &quiz.Route, &quiz.HasLeaderboard, &quiz.HasGrouping, &quiz.HasFlags, &quiz.Enabled)

	if quiz.MapName != "" {
		svgMap, err := GetMap(quiz.MapName)
		if err != nil {
			return QuizDto{}, err
		}
		quiz.Map = svgMap
	}

	return quiz, err
}

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

func CreateQuiz(newQuiz CreateQuizDto) (Quiz, error) {
	statement := "INSERT INTO quizzes (typeId, badgeId, continentId, country, singular, name, maxScore, time, mapSVG, imageUrl, plural, apiPath, route, hasLeaderboard, hasGrouping, hasFlags, enabled) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17) RETURNING *;"
	var quiz Quiz
	err := Connection.QueryRow(statement, newQuiz.TypeID, newQuiz.BadgeID, newQuiz.ContinentID, newQuiz.Country, newQuiz.Singular, newQuiz.Name, newQuiz.MaxScore, newQuiz.Time, newQuiz.MapSVG, newQuiz.ImageURL, newQuiz.Plural, newQuiz.APIPath, newQuiz.Route, newQuiz.HasLeaderboard, newQuiz.HasGrouping, newQuiz.HasFlags, newQuiz.Enabled).Scan(&quiz.ID, &quiz.TypeID, &quiz.BadgeID, &quiz.ContinentID, &quiz.Country, &quiz.Singular, &quiz.Name, &quiz.MaxScore, &quiz.Time, &quiz.MapSVG, &quiz.ImageURL, &quiz.Plural, &quiz.APIPath, &quiz.Route, &quiz.HasLeaderboard, &quiz.HasGrouping, &quiz.HasFlags, &quiz.Enabled)
	return quiz, err
}

func UpdateQuiz(quizID int, quiz UpdateQuizDto) error {
	statement := "UPDATE quizzes SET typeId = $1, badgeId = $2, continentId = $3, country = $4, singular = $5, name = $6, maxScore = $7, time = $8, mapSVG = $9, imageUrl = $10, plural = $11, apiPath = $12, route = $13, hasLeaderboard = $14, hasGrouping = $15, hasFlags = $16, enabled = $17 WHERE id = $18 RETURNING id;"
	var id int
	return Connection.QueryRow(statement, quiz.TypeID, quiz.BadgeID, quiz.ContinentID, quiz.Country, quiz.Singular, quiz.Name, quiz.MaxScore, quiz.Time, quiz.MapSVG, quiz.ImageURL, quiz.Plural, quiz.APIPath, quiz.Route, quiz.HasLeaderboard, quiz.HasGrouping, quiz.HasFlags, quiz.Enabled, quizID).Scan(&id)
}

func DeleteQuiz(quizID int) error {
	// Delete quiz plays.
	var id int
	err := Connection.QueryRow("DELETE FROM quizplays where quizid = $1 RETURNING id;", quizID).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// Delete leaderboard entries.
	err = Connection.QueryRow("DELETE FROM leaderboard where quizid = $1 RETURNING id;", quizID).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// Delete quiz.
	var key string
	var typeId int
	err = Connection.QueryRow("DELETE FROM quizzes where id = $1 RETURNING typeId, apiPath;", quizID).Scan(&typeId, &key)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if typeId == QUIZ_TYPE_MAP {
		// Delete svg map.
		mapId, err := GetMapId(key)
		if err != nil {
			return err
		}

		if err = DeleteMapElements(mapId); err != nil {
			return err
		}

		if err = DeleteMap(mapId); err != nil {
			return err
		}
	}

	return err
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

func GetQuizRoutes() ([]string, error) {
	rows, err := Connection.Query("SELECT route FROM quizzes;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routes = []string{}
	for rows.Next() {
		var route string
		if err = rows.Scan(&route); err != nil {
			return nil, err
		}
		routes = append(routes, fmt.Sprintf("quiz/%s", route))
	}
	return routes, rows.Err()
}
