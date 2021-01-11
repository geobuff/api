package database

// Quiz is the database object for a quiz entry.
type Quiz struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Code        string `json:"code"`
	MaxScore    int    `json:"maxScore"`
	Enabled     bool   `json:"enabled"`
}

// GetQuizzes returns all quizzes.
var GetQuizzes = func() ([]Quiz, error) {
	rows, err := Connection.Query("SELECT * FROM quizzes;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quizzes = []Quiz{}
	for rows.Next() {
		var quiz Quiz
		if err = rows.Scan(&quiz.ID, &quiz.Name, &quiz.Description, &quiz.Code, &quiz.MaxScore, &quiz.Enabled); err != nil {
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
	err := Connection.QueryRow(statement, id).Scan(&quiz.ID, &quiz.Name, &quiz.Description, &quiz.Code, &quiz.MaxScore, &quiz.Enabled)
	return quiz, err
}
