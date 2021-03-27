package database

// Quiz is the database object for a quiz entry.
type Quiz struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	MaxScore       int    `json:"maxScore"`
	Time           int    `json:"time"`
	MapSVG         string `json:"mapSVG"`
	ImageURL       string `json:"imageUrl"`
	Verb           string `json:"verb"`
	APIPath        string `json:"apiPath"`
	Route          string `json:"route"`
	HasLeaderboard bool   `json:"hasLeaderboard"`
	HasGrouping    bool   `json:"hasGrouping"`
	HasFlags       bool   `json:"hasFlags"`
	Enabled        bool   `json:"enabled"`
}

// GetQuizzes returns all quizzes.
var GetQuizzes = func(filter string) ([]Quiz, error) {
	statement := "SELECT * FROM quizzes WHERE name ILIKE '%' || $1 || '%';"
	rows, err := Connection.Query(statement, filter)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quizzes = []Quiz{}
	for rows.Next() {
		var quiz Quiz
		if err = rows.Scan(&quiz.ID, &quiz.Name, &quiz.MaxScore, &quiz.Time, &quiz.MapSVG, &quiz.ImageURL, &quiz.Verb, &quiz.APIPath, &quiz.Route, &quiz.HasLeaderboard, &quiz.HasGrouping, &quiz.HasFlags, &quiz.Enabled); err != nil {
			return nil, err
		}
		quizzes = append(quizzes, quiz)
	}
	return quizzes, rows.Err()
}

// GetQuiz return a quiz with the matching id.
var GetQuiz = func(id int) (Quiz, error) {
	statement := "SELECT * FROM quizzes WHERE id = $1;"
	var quiz Quiz
	err := Connection.QueryRow(statement, id).Scan(&quiz.ID, &quiz.Name, &quiz.MaxScore, &quiz.Time, &quiz.MapSVG, &quiz.ImageURL, &quiz.Verb, &quiz.APIPath, &quiz.Route, &quiz.HasLeaderboard, &quiz.HasGrouping, &quiz.HasFlags, &quiz.Enabled)
	return quiz, err
}

// GetQuizID gets the quiz id based on name.
func GetQuizID(name string) (int, error) {
	statement := "SELECT id FROM quizzes WHERE name ILIKE '%' || $1 || '%';"
	var id int
	err := Connection.QueryRow(statement, name).Scan(&id)
	return id, err
}
