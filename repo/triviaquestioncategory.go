package repo

type TriviaQuestionCategory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetTriviaQuestionCategories() ([]TriviaQuestionCategory, error) {
	rows, err := Connection.Query("SELECT * from triviaquestioncategory;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories = []TriviaQuestionCategory{}
	for rows.Next() {
		var category TriviaQuestionCategory
		if err = rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, rows.Err()
}
