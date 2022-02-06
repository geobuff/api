package repo

const (
	QUIZ_TYPE_MAP int = iota + 1
	QUIZ_TYPE_FLAG
)

type QuizType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetQuizTypes() ([]QuizType, error) {
	rows, err := Connection.Query("SELECT * from quiztype;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types = []QuizType{}
	for rows.Next() {
		var quizType QuizType
		if err = rows.Scan(&quizType.ID, &quizType.Name); err != nil {
			return nil, err
		}
		types = append(types, quizType)
	}
	return types, rows.Err()
}
