package repo

type AdminDashboardData struct {
	UserCount   int        `json:"userCount"`
	Discounts   []Discount `json:"discounts"`
	QuizPlays   []PlaysDto `json:"quizPlays"`
	TriviaPlays []PlaysDto `json:"triviaPlays"`
}

func GetAdminDashboardData() (*AdminDashboardData, error) {
	var userCount int
	err := Connection.QueryRow("SELECT COUNT(id) FROM users;").Scan(&userCount)
	if err != nil {
		return nil, err
	}

	discounts, err := GetDiscounts()
	if err != nil {
		return nil, err
	}

	quizPlays, err := GetTopFiveQuizPlays()
	if err != nil {
		return nil, err
	}

	triviaPlays, err := GetLastFiveTriviaPlays()
	if err != nil {
		return nil, err
	}

	result := AdminDashboardData{
		UserCount:   userCount,
		Discounts:   discounts,
		QuizPlays:   quizPlays,
		TriviaPlays: triviaPlays,
	}
	return &result, nil
}
