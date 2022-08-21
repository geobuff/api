package repo

import (
	"database/sql"
	"math/rand"
	"time"
)

type CommunityQuiz struct {
	ID          int       `json:"id"`
	UserID      int       `json:"userId"`
	StatusID    int       `json:"statusId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	MaxScore    int       `json:"maxScore"`
	Added       time.Time `json:"added"`
	Verified    bool      `json:"verified"`
	IsPublic    bool      `json:"isPublic"`
}

type CommunityQuizDto struct {
	ID          int           `json:"id"`
	UserID      int           `json:"userId"`
	Status      string        `json:"status"`
	Username    string        `json:"username"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	MaxScore    int           `json:"maxScore"`
	Added       time.Time     `json:"added"`
	Verified    bool          `json:"verified"`
	IsPublic    bool          `json:"isPublic"`
	Plays       sql.NullInt64 `json:"plays"`
}

type GetCommunityQuizDto struct {
	ID          int                           `json:"id"`
	UserID      int                           `json:"userId"`
	Status      string                        `json:"status"`
	Name        string                        `json:"name"`
	Description string                        `json:"description"`
	MaxScore    int                           `json:"maxScore"`
	IsPublic    bool                          `json:"isPublic"`
	Questions   []GetCommunityQuizQuestionDto `json:"questions"`
}

type GetCommunityQuizzesFilter struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Filter string `json:"filter"`
}

type CreateCommunityQuizDto struct {
	UserID      int                              `json:"userId"`
	Name        string                           `json:"name"`
	Description string                           `json:"description"`
	MaxScore    int                              `json:"maxScore"`
	IsPublic    bool                             `json:"isPublic"`
	IsVerified  bool                             `json:"isVerified"`
	Questions   []CreateCommunityQuizQuestionDto `json:"questions"`
}

type UpdateCommunityQuizDto struct {
	UserID      int                              `json:"userId"`
	Name        string                           `json:"name"`
	Description string                           `json:"description"`
	MaxScore    int                              `json:"maxScore"`
	IsPublic    bool                             `json:"isPublic"`
	Questions   []UpdateCommunityQuizQuestionDto `json:"questions"`
}

func GetCommunityQuizzes(filter GetCommunityQuizzesFilter) ([]CommunityQuizDto, error) {
	statement := "SELECT q.id, q.userid, s.name, u.username, q.name, q.description, q.maxscore, q.added, q.verified, q.ispublic, p.plays FROM communityquizzes q JOIN users u ON u.id = q.userid LEFT JOIN communityquizplays p ON p.communityQuizId = q.id JOIN communityQuizStatus s ON s.id = q.statusid WHERE (q.name ILIKE '%' || $1 || '%' OR q.description ILIKE '%' || $1 || '%') AND q.statusid != $2 AND q.ispublic LIMIT $3 OFFSET $4;"
	rows, err := Connection.Query(statement, filter.Filter, COMMUNITY_QUIZ_STATUS_PENDING, filter.Limit, filter.Page*filter.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quizzes = []CommunityQuizDto{}
	for rows.Next() {
		var quiz CommunityQuizDto
		if err = rows.Scan(&quiz.ID, &quiz.UserID, &quiz.Status, &quiz.Username, &quiz.Name, &quiz.Description, &quiz.MaxScore, &quiz.Added, &quiz.Verified, &quiz.IsPublic, &quiz.Plays); err != nil {
			return nil, err
		}
		quizzes = append(quizzes, quiz)
	}
	return quizzes, rows.Err()
}

func GetFirstCommunityQuizID(filter GetCommunityQuizzesFilter) (int, error) {
	statement := "SELECT q.id FROM communityquizzes q JOIN communityQuizStatus s ON s.id = q.statusid WHERE (q.name ILIKE '%' || $1 || '%' OR q.description ILIKE '%' || $1 || '%') AND q.statusid != $2 AND q.ispublic LIMIT 1 OFFSET $3;"
	var id int
	err := Connection.QueryRow(statement, filter.Filter, COMMUNITY_QUIZ_STATUS_PENDING, (filter.Page+1)*filter.Limit).Scan(&id)
	return id, err
}

func GetUserCommunityQuizzes(userID int) ([]CommunityQuizDto, error) {
	statement := "SELECT q.id, q.userid, s.name, u.username, q.name, q.description, q.maxscore, q.added, q.verified, q.ispublic, p.plays FROM communityquizzes q JOIN users u ON u.id = q.userid LEFT JOIN communityquizplays p ON p.communityQuizId = q.id JOIN communityQuizStatus s ON s.id = q.statusid WHERE q.userId = $1 ORDER BY q.added DESC;"
	rows, err := Connection.Query(statement, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quizzes = []CommunityQuizDto{}
	for rows.Next() {
		var quiz CommunityQuizDto
		if err = rows.Scan(&quiz.ID, &quiz.UserID, &quiz.Status, &quiz.Username, &quiz.Name, &quiz.Description, &quiz.MaxScore, &quiz.Added, &quiz.Verified, &quiz.IsPublic, &quiz.Plays); err != nil {
			return nil, err
		}
		quizzes = append(quizzes, quiz)
	}
	return quizzes, rows.Err()
}

func InsertCommunityQuiz(quiz CreateCommunityQuizDto) error {
	statement := "INSERT INTO communityquizzes (userid, statusid, name, description, maxscore, ispublic, verified, added) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;"
	var quizID int
	if err := Connection.QueryRow(statement, quiz.UserID, COMMUNITY_QUIZ_STATUS_APPROVED, quiz.Name, quiz.Description, quiz.MaxScore, quiz.IsPublic, quiz.IsVerified, time.Now()).Scan(&quizID); err != nil {
		return err
	}

	for _, question := range quiz.Questions {
		questionID, err := InsertCommunityQuizQuestion(quizID, question)
		if err != nil {
			return err
		}

		for _, answer := range question.Answers {
			if err := InsertCommunityQuizAnswer(questionID, answer); err != nil {
				return err
			}
		}
	}

	return nil
}

func UpdateCommunityQuiz(quizID int, quiz UpdateCommunityQuizDto) error {
	statement := "UPDATE communityquizzes SET name = $1, description = $2, maxscore = $3, ispublic = $4 WHERE id = $5 RETURNING id;"
	var id int
	if err := Connection.QueryRow(statement, quiz.Name, quiz.Description, quiz.MaxScore, quiz.IsPublic, quizID).Scan(&id); err != nil {
		return err
	}

	for _, question := range quiz.Questions {
		if question.ID.Valid {
			if err := DeleteCommunityQuizAnswers(int(question.ID.Int64)); err != nil {
				return err
			}

			if err := DeleteCommunityQuizQuestion(int(question.ID.Int64)); err != nil {
				return err
			}
		}

		questionID, err := InsertCommunityQuizQuestion(id, CreateCommunityQuizQuestionDto{
			TypeID:             question.TypeID,
			Question:           question.Question,
			Explainer:          question.Explainer,
			Map:                question.Map,
			Highlighted:        question.Highlighted,
			FlagCode:           question.FlagCode,
			ImageUrl:           question.ImageUrl,
			ImageAttributeName: question.ImageAttributeName,
			ImageAttributeURL:  question.ImageAttributeURL,
			ImageWidth:         question.ImageWidth,
			ImageHeight:        question.ImageHeight,
			ImageAlt:           question.ImageAlt,
		})

		if err != nil {
			return err
		}

		for _, answer := range question.Answers {
			if err := InsertCommunityQuizAnswer(questionID, answer); err != nil {
				return err
			}
		}
	}

	return nil
}

func GetCommunityQuiz(quizID int) (GetCommunityQuizDto, error) {
	statement := "SELECT q.id, q.userid, s.name, q.name, q.description, q.maxscore, q.ispublic FROM communityquizzes q JOIN communityQuizStatus s ON s.id = q.statusid WHERE q.id = $1;"
	var quiz GetCommunityQuizDto
	if err := Connection.QueryRow(statement, quizID).Scan(&quiz.ID, &quiz.UserID, &quiz.Status, &quiz.Name, &quiz.Description, &quiz.MaxScore, &quiz.IsPublic); err != nil {
		return quiz, err
	}

	questions, err := GetCommunityQuizQuestions(quizID)
	if err != nil {
		return quiz, err
	}

	rand.Shuffle(len(questions), func(i, j int) {
		questions[i], questions[j] = questions[j], questions[i]
	})

	quiz.Questions = questions

	return quiz, err
}

func DeleteCommunityQuiz(quizID int) error {
	questionIds, err := GetCommunityQuizQuestionIds(quizID)
	if err != nil {
		return err
	}

	var id int
	if err := ClearCommunityQuizPlayCommunityQuizId(quizID); err != nil && err != sql.ErrNoRows {
		return err
	}

	for _, questionId := range questionIds {
		if err = DeleteCommunityQuizAnswers(questionId); err != nil {
			return err
		}

		if err = DeleteCommunityQuizQuestion(questionId); err != nil {
			return err
		}
	}

	return Connection.QueryRow("DELETE FROM communityquizzes WHERE id = $1 RETURNING id;", quizID).Scan(&id)
}

func GetUserCommunityQuizCount(userID int) (int, error) {
	var count int
	err := Connection.QueryRow("SELECT COUNT(id) FROM communityquizzes WHERE userid = $1;", userID).Scan(&count)
	return count, err
}
