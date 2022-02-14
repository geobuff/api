package repo

const (
	QUESTION_TYPE_TEXT int = iota + 1
	QUESTION_TYPE_IMAGE
	QUESTION_TYPE_FLAG
	QUESTION_TYPE_MAP
)

type TriviaQuestionType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetTriviaQuestionTypes() ([]TriviaQuestionType, error) {
	rows, err := Connection.Query("SELECT * from triviaquestiontype;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types = []TriviaQuestionType{}
	for rows.Next() {
		var quizType TriviaQuestionType
		if err = rows.Scan(&quizType.ID, &quizType.Name); err != nil {
			return nil, err
		}
		types = append(types, quizType)
	}
	return types, rows.Err()
}
