package repo

type TriviaQuestionCategory struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	IsActive  bool   `json:"isActive"`
	ImageOnly bool   `json:"imageOnly"`
}

func GetTriviaQuestionCategories(onlyActive bool) ([]TriviaQuestionCategory, error) {
	statement := "SELECT * from triviaquestioncategory"
	if onlyActive {
		statement += " WHERE isactive"
	}
	statement += ";"

	rows, err := Connection.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories = []TriviaQuestionCategory{}
	for rows.Next() {
		var category TriviaQuestionCategory
		if err = rows.Scan(&category.ID, &category.Name, &category.IsActive, &category.ImageOnly); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, rows.Err()
}
