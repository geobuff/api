package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/geobuff/api/types"
)

type Trivia struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Date     time.Time `json:"date"`
	MaxScore int       `json:"maxScore"`
}

type TriviaDto struct {
	ID        int           `json:"id"`
	Name      string        `json:"name"`
	MaxScore  int           `json:"maxScore"`
	Questions []QuestionDto `json:"questions"`
}

type GetTriviaFilter struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Filter string `json:"filter"`
}

func CreateTrivia() error {
	date := time.Now().AddDate(0, 0, 1)
	return CreateTriviaForDate(date)
}

func CreateTriviaForDate(date time.Time) error {
	var id int
	if err := Connection.QueryRow("SELECT id FROM trivia WHERE date = $1", date).Scan(&id); err != sql.ErrNoRows {
		return errors.New("trivia for current date already created")
	}

	_, month, day := date.Date()
	weekday := date.Weekday().String()
	statement := "INSERT INTO trivia (name, date, maxscore) VALUES ($1, $2, $3) RETURNING id;"
	if err := Connection.QueryRow(statement, fmt.Sprintf("%s, %s %d", weekday, month, day), date, 0).Scan(&id); err != nil {
		return err
	}

	count, err := generateQuestions(id, 10)
	if err != nil {
		return err
	}

	return setTriviaMaxScore(id, count)
}

func RegenerateTrivia(dateString string) error {
	trivia, err := GetTrivia(dateString)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err != sql.ErrNoRows {
		if err = deleteTrivia(trivia); err != nil {
			return err
		}
	}

	date, err := time.Parse("2006-01-02", dateString)
	return CreateTriviaForDate(date)
}

func generateQuestions(triviaId, max int) (int, error) {
	countries, err := GetMappingEntries("world-countries")
	if err != nil {
		return 0, err
	}

	capitals, err := GetMappingEntries("world-capitals")
	if err != nil {
		return 0, err
	}

	states, err := GetMappingEntries("us-states")
	if err != nil {
		return 0, err
	}

	count := 0
	err = whatCountry(triviaId, countries)
	if err != nil {
		return count, err
	}

	err = whatCapital(triviaId, countries, capitals)
	if err != nil {
		return count, err
	}

	err = whatUSState(triviaId, states)
	if err != nil {
		return count, err
	}

	err = whatFlag(triviaId, countries)
	if err != nil {
		return count, err
	}
	count = count + 4

	questions, err := GetTodaysManualTriviaQuestions()
	if err != nil && err != sql.ErrNoRows {
		return count, err
	}

	todaysQuestionCount, err := createQuestionsAndAnswers(questions, triviaId, max-count)
	if err != nil {
		return count, err
	}

	count = count + todaysQuestionCount
	if count < max {
		categories, err := GetTriviaQuestionCategories(true)
		if err != nil {
			return count, err
		}

		remainder := max - count
		var allowedCategories []TriviaQuestionCategory
		for i := 0; i < remainder; i++ {
			index := rand.Intn(len(categories))
			allowedCategories = append(allowedCategories, categories[index])
			categories = append(categories[:index], categories[index+1:]...)
		}

		maxTextCount := remainder / 2
		var textCategories []int
		var imageCategories []int
		for i := 0; i < len(allowedCategories); i++ {
			curr := allowedCategories[i]
			if !curr.ImageOnly && len(textCategories) < maxTextCount {
				textCategories = append(textCategories, curr.ID)
			} else {
				imageCategories = append(imageCategories, curr.ID)
			}
		}

		textQuestionCount, err := setRandomManualTriviaQuestions(triviaId, QUESTION_TYPE_TEXT, maxTextCount, textCategories)
		if err != nil && err != sql.ErrNoRows {
			return count, err
		}
		count = count + textQuestionCount

		imageQuestionCount, err := setRandomManualTriviaQuestions(triviaId, QUESTION_TYPE_IMAGE, max-count, imageCategories)
		if err != nil && err != sql.ErrNoRows {
			return count, err
		}
		count = count + imageQuestionCount
	}

	return count, nil
}

func whatCountry(triviaId int, countries []MappingEntryDto) error {
	max := len(types.TopLandmass)
	index := rand.Intn(max)
	country := types.TopLandmass[index]

	question := TriviaQuestion{
		TriviaId:    triviaId,
		TypeID:      QUESTION_TYPE_MAP,
		Question:    "Which country is highlighted above?",
		Map:         "WorldCountries",
		Highlighted: country,
	}

	questionId, err := CreateTriviaQuestion(question)
	if err != nil {
		return err
	}

	answer := TriviaAnswer{
		TriviaQuestionID: questionId,
		Text:             country,
		IsCorrect:        true,
	}
	err = CreateTriviaAnswer(answer)
	if err != nil {
		return err
	}

	for i, val := range countries {
		if val.SVGName == country {
			index = i
			break
		}
	}
	countries = append(countries[:index], countries[index+1:]...)
	max = len(countries)

	for i := 0; i < 3; i++ {
		index := rand.Intn(max)
		country = countries[index].SVGName
		answer := TriviaAnswer{
			TriviaQuestionID: questionId,
			Text:             country,
			IsCorrect:        false,
		}

		err = CreateTriviaAnswer(answer)
		if err != nil {
			return err
		}

		countries = append(countries[:index], countries[index+1:]...)
		max = max - 1
	}

	return nil
}

func whatCapital(triviaId int, countries []MappingEntryDto, capitals []MappingEntryDto) error {
	max := len(types.TopLandmass)
	index := rand.Intn(max)
	country := types.TopLandmass[index]
	var code string
	for _, val := range countries {
		if val.SVGName == country {
			code = val.Code
			break
		}
	}

	var capitalName string
	for _, value := range capitals {
		if value.Code == code {
			capitalName = value.SVGName
		}
	}

	question := TriviaQuestion{
		TriviaId:    triviaId,
		TypeID:      QUESTION_TYPE_MAP,
		Question:    fmt.Sprintf("What is the capital city of %s?", country),
		Map:         "WorldCapitals",
		Highlighted: capitalName,
	}

	questionId, err := CreateTriviaQuestion(question)
	if err != nil {
		return err
	}

	answer := TriviaAnswer{
		TriviaQuestionID: questionId,
		Text:             capitalName,
		IsCorrect:        true,
	}
	err = CreateTriviaAnswer(answer)
	if err != nil {
		return err
	}

	for i, val := range capitals {
		if val.SVGName == capitalName {
			index = i
			break
		}
	}

	capitals = append(capitals[:index], capitals[index+1:]...)
	max = len(capitals)
	for i := 0; i < 3; i++ {
		index := rand.Intn(max)
		capital := capitals[index]
		answer := TriviaAnswer{
			TriviaQuestionID: questionId,
			Text:             capital.SVGName,
			IsCorrect:        false,
		}

		err = CreateTriviaAnswer(answer)
		if err != nil {
			return err
		}

		capitals = append(capitals[:index], capitals[index+1:]...)
		max = max - 1
	}

	return nil
}

func whatUSState(triviaId int, states []MappingEntryDto) error {
	max := len(states)
	index := rand.Intn(max)
	state := states[index]
	states = append(states[:index], states[index+1:]...)
	max = max - 1

	question := TriviaQuestion{
		TriviaId:    triviaId,
		TypeID:      QUESTION_TYPE_MAP,
		Question:    "Which US state is highlighted above?",
		Map:         "UsStates",
		Highlighted: state.SVGName,
	}

	questionId, err := CreateTriviaQuestion(question)
	if err != nil {
		return err
	}

	answer := TriviaAnswer{
		TriviaQuestionID: questionId,
		Text:             state.SVGName,
		IsCorrect:        true,
	}

	err = CreateTriviaAnswer(answer)
	if err != nil {
		return err
	}

	for i := 0; i < 3; i++ {
		index := rand.Intn(max)
		state = states[index]
		answer := TriviaAnswer{
			TriviaQuestionID: questionId,
			Text:             state.SVGName,
			IsCorrect:        false,
		}

		err = CreateTriviaAnswer(answer)
		if err != nil {
			return err
		}

		states = append(states[:index], states[index+1:]...)
		max = max - 1
	}

	return nil
}

func whatFlag(triviaId int, countries []MappingEntryDto) error {
	max := len(countries)
	index := rand.Intn(max)
	country := countries[index]
	countries = append(countries[:index], countries[index+1:]...)
	max = max - 1

	question := TriviaQuestion{
		TriviaId: triviaId,
		TypeID:   QUESTION_TYPE_FLAG,
		Question: "Which country has this flag?",
		FlagCode: country.Code,
	}

	questionId, err := CreateTriviaQuestion(question)
	if err != nil {
		return err
	}

	answer := TriviaAnswer{
		TriviaQuestionID: questionId,
		Text:             country.SVGName,
		IsCorrect:        true,
	}
	err = CreateTriviaAnswer(answer)
	if err != nil {
		return err
	}

	for i := 0; i < 3; i++ {
		index := rand.Intn(max)
		country = countries[index]
		answer := TriviaAnswer{
			TriviaQuestionID: questionId,
			Text:             country.SVGName,
			IsCorrect:        false,
		}

		err = CreateTriviaAnswer(answer)
		if err != nil {
			return err
		}

		countries = append(countries[:index], countries[index+1:]...)
		max = max - 1
	}

	return nil
}

func setRandomManualTriviaQuestions(triviaID, typeID, quantity int, allowedCategories []int) (int, error) {
	lastUsedMax := time.Now().AddDate(0, 0, -7)
	questions, err := GetManualTriviaQuestions(typeID, lastUsedMax.Format("2006-01-02"), allowedCategories)
	if err != nil {
		return 0, err
	}

	return createQuestionsAndAnswers(questions, triviaID, quantity)
}

func createQuestionsAndAnswers(questions []ManualTriviaQuestion, triviaID, quantity int) (int, error) {
	count := 0
	for i := 0; i < quantity; i++ {
		if len(questions) == 0 {
			break
		}

		max := len(questions)
		index := rand.Intn(max)
		manualQuestion := questions[index]

		question := TriviaQuestion{
			TriviaId:           triviaID,
			TypeID:             manualQuestion.TypeID,
			Question:           manualQuestion.Question,
			Explainer:          manualQuestion.Explainer,
			Map:                manualQuestion.Map,
			Highlighted:        manualQuestion.Highlighted,
			FlagCode:           manualQuestion.FlagCode,
			ImageURL:           manualQuestion.ImageURL,
			ImageAttributeName: manualQuestion.ImageAttributeName,
			ImageAttributeURL:  manualQuestion.ImageAttributeURL,
			ImageWidth:         manualQuestion.ImageWidth,
			ImageHeight:        manualQuestion.ImageHeight,
			ImageAlt:           manualQuestion.ImageAlt,
		}

		questionID, err := CreateTriviaQuestion(question)
		if err != nil {
			return count, err
		}

		answers, err := GetManualTriviaAnswers(manualQuestion.ID)
		if err != nil {
			return count, err
		}

		for _, answer := range answers {
			newAnswer := TriviaAnswer{
				TriviaQuestionID: questionID,
				Text:             answer.Text,
				IsCorrect:        answer.IsCorrect,
				FlagCode:         answer.FlagCode,
			}

			if err := CreateTriviaAnswer(newAnswer); err != nil {
				return count, err
			}
		}

		questions = append(questions[:index], questions[index+1:]...)
		if err := UpdateManualTriviaQuestionLastUsed(manualQuestion.ID); err != nil {
			return count, err
		}
		count = count + 1
	}

	return count, nil
}

func GetAllTrivia(filter GetTriviaFilter) ([]Trivia, error) {
	statement := "SELECT * FROM trivia WHERE name ILIKE '%' || $1 || '%' ORDER BY date DESC LIMIT $2 OFFSET $3;"
	rows, err := Connection.Query(statement, filter.Filter, filter.Limit, filter.Page*filter.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trivia = []Trivia{}
	for rows.Next() {
		var quiz Trivia
		if err = rows.Scan(&quiz.ID, &quiz.Name, &quiz.Date, &quiz.MaxScore); err != nil {
			return nil, err
		}
		trivia = append(trivia, quiz)
	}
	return trivia, rows.Err()
}

func GetFirstTriviaID(filter GetTriviaFilter) (int, error) {
	statement := "SELECT id FROM trivia WHERE name ILIKE '%' || $1 || '%' ORDER BY date DESC LIMIT 1 OFFSET $2;"
	var id int
	err := Connection.QueryRow(statement, filter.Filter, (filter.Page+1)*filter.Page).Scan(&id)
	return id, err
}

func GetTrivia(date string) (*TriviaDto, error) {
	var result TriviaDto
	err := Connection.QueryRow("SELECT id, name, maxscore from trivia WHERE date = $1;", date).Scan(&result.ID, &result.Name, &result.MaxScore)
	if err != nil {
		return nil, err
	}

	questions, err := GetTriviaQuestions(result.ID)
	if err != nil {
		return nil, err
	}

	result.Questions = questions
	return &result, nil
}

func setTriviaMaxScore(triviaID, maxScore int) error {
	statement := "UPDATE trivia SET maxScore = $1 WHERE id = $2 RETURNING id;"
	var id int
	return Connection.QueryRow(statement, maxScore, triviaID).Scan(&id)
}

func DeleteTriviaByDate(dateString string) error {
	trivia, err := GetTrivia(dateString)
	if err != nil {
		return err
	}

	return deleteTrivia(trivia)
}

func deleteTrivia(trivia *TriviaDto) error {
	if err := ClearTriviaPlayTriviaId(trivia.ID); err != nil && err != sql.ErrNoRows {
		return err
	}

	for _, question := range trivia.Questions {
		if err := DeleteTriviaAnswers(question.ID); err != nil && err != sql.ErrNoRows {
			return err
		}

		if err := DeleteTriviaQuestion(question.ID); err != nil && err != sql.ErrNoRows {
			return err
		}
	}

	var id int
	return Connection.QueryRow("DELETE FROM trivia WHERE id = $1 RETURNING id;", trivia.ID).Scan(&id)
}

func DeleteOldTrivia(newTriviaCount int) error {
	dates, err := getOldTriviaDates(newTriviaCount)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return err
	}

	for _, date := range dates {
		formattedDate := date.Format("2006-01-02")
		trivia, err := GetTrivia(formattedDate)
		if err != nil {
			return err
		}

		if err = deleteTrivia(trivia); err != nil {
			return err
		}
	}
	return err
}

func getOldTriviaDates(newTriviaCount int) ([]time.Time, error) {
	rows, err := Connection.Query("SELECT date FROM trivia WHERE date < $1;", time.Now().AddDate(0, 0, 0-newTriviaCount))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dates = []time.Time{}
	for rows.Next() {
		var date time.Time
		if err = rows.Scan(&date); err != nil {
			return nil, err
		}
		dates = append(dates, date)
	}

	return dates, rows.Err()
}
