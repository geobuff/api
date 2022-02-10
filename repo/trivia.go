package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/geobuff/mapping"
)

type Trivia struct {
	ID   int       `json:"id"`
	Name string    `json:"name"`
	Date time.Time `json:"date"`
}

type TriviaQuestion struct {
	ID          int    `json:"id"`
	TriviaId    int    `json:"triviaId"`
	TypeID      int    `json:"typeID"`
	Question    string `json:"question"`
	Map         string `json:"map"`
	Highlighted string `json:"hightlighted"`
	FlagCode    string `json:"flagCode"`
	ImageURL    string `json:"imageUrl"`
}

type TriviaAnswer struct {
	ID               int    `json:"id"`
	TriviaQuestionID int    `json:"triviaQuestionId"`
	Text             string `json:"text"`
	IsCorrect        bool   `json:"isCorrect"`
	FlagCode         string `json:"flagCode"`
}

type TriviaDto struct {
	ID        int           `json:"id"`
	Name      string        `json:"name"`
	Questions []QuestionDto `json:"questions"`
}

type QuestionDto struct {
	ID          int         `json:"id"`
	Type        string      `json:"type"`
	Question    string      `json:"question"`
	Map         string      `json:"map"`
	Highlighted string      `json:"highlighted"`
	FlagCode    string      `json:"flagCode"`
	ImageURL    string      `json:"imageUrl"`
	Answers     []AnswerDto `json:"answers"`
}

type AnswerDto struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"isCorrect"`
	FlagCode  string `json:"flagCode"`
}

var topLandmass = []string{
	"Russia",
	"Canada",
	"China",
	"United States",
	"Brazil",
	"Australia",
	"India",
	"Argentina",
	"Kazakhstan",
	"Algeria",
	"Democratic Republic of the Congo",
	"Denmark",
	"Saudi Arabia",
	"Mexico",
	"Indonesia",
	"Sudan",
	"Libya",
	"Iran",
	"Mongolia",
	"Peru",
	"Chad",
	"Niger",
	"Angola",
	"Mali",
	"South Africa",
	"Colombia",
	"Ethiopia",
	"Bolivia",
	"Mauritania",
	"Egypt",
	"Tanzania",
	"Nigeria",
	"Venezuela",
	"Pakistan",
	"Namibia",
	"Mozambique",
	"Turkey",
	"Chile",
	"Zambia",
	"Myanmar",
	"Afghanistan",
	"Somalia",
	"Central African Republic",
	"South Sudan",
	"Ukraine",
	"Madagascar",
	"Botswana",
	"Kenya",
	"France",
	"Yemen",
	"New Zealand",
}

func CreateTrivia() error {
	date := time.Now().AddDate(0, 0, 1)
	var id int
	if err := Connection.QueryRow("SELECT id FROM trivia WHERE date = $1", date).Scan(&id); err != sql.ErrNoRows {
		return errors.New("trivia for current date already created")
	}

	_, month, day := date.Date()
	weekday := date.Weekday().String()
	statement := "INSERT INTO trivia (name, date) VALUES ($1, $2) RETURNING id;"
	err := Connection.QueryRow(statement, fmt.Sprintf("%s, %s %d", weekday, month, day), date).Scan(&id)
	if err != nil {
		return err
	}
	return generateQuestions(id)
}

func generateQuestions(triviaId int) error {
	err := whatCountry(triviaId)
	if err != nil {
		return err
	}

	err = whatCapital(triviaId)
	if err != nil {
		return err
	}

	err = whatUSState(triviaId)
	if err != nil {
		return err
	}

	err = whatFlag(triviaId)
	if err != nil {
		return err
	}

	err = whatRegionInCountry(triviaId)
	if err != nil {
		return err
	}

	err = whatFlagInCountry(triviaId)
	if err != nil {
		return err
	}

	err = trueFalseCountryInContinent(triviaId, randomBool())
	if err != nil {
		return err
	}

	err = trueFalseCapitalOfCountry(triviaId, randomBool())
	if err != nil {
		return err
	}

	err = trueFalseRegionInCountry(triviaId, randomBool())
	if err != nil {
		return err
	}

	err = trueFalseFlagForCountry(triviaId, randomBool())
	if err != nil {
		return err
	}

	return nil
}

func randomBool() bool {
	return rand.Float32() < 0.5
}

func createQuestion(question TriviaQuestion) (int, error) {
	statement := "INSERT INTO triviaQuestions (triviaId, typeId, question, map, highlighted, flagCode, imageUrl) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, question.TriviaId, question.TypeID, question.Question, question.Map, question.Highlighted, question.FlagCode, question.ImageURL).Scan(&id)
	return id, err
}

func createAnswer(answer TriviaAnswer) error {
	statement := "INSERT INTO triviaAnswers (triviaQuestionId, text, isCorrect, flagCode) VALUES ($1, $2, $3, $4) RETURNING id;"
	var id int
	return Connection.QueryRow(statement, answer.TriviaQuestionID, answer.Text, answer.IsCorrect, answer.FlagCode).Scan(&id)
}

func whatCountry(triviaId int) error {
	max := len(topLandmass)
	index := rand.Intn(max)
	country := topLandmass[index]

	question := TriviaQuestion{
		TriviaId:    triviaId,
		TypeID:      QUESTION_TYPE_MAP,
		Question:    "Which country is highlighted above?",
		Map:         "WorldCountries",
		Highlighted: country,
	}

	questionId, err := createQuestion(question)
	if err != nil {
		return err
	}

	answer := TriviaAnswer{
		TriviaQuestionID: questionId,
		Text:             country,
		IsCorrect:        true,
	}
	err = createAnswer(answer)
	if err != nil {
		return err
	}

	countries := copyMapping(mapping.Mappings["world-countries"])
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

		err = createAnswer(answer)
		if err != nil {
			return err
		}

		countries = append(countries[:index], countries[index+1:]...)
		max = max - 1
	}

	return nil
}

func whatCapital(triviaId int) error {
	max := len(topLandmass)
	index := rand.Intn(max)
	country := topLandmass[index]
	var code string
	for _, val := range mapping.Mappings["world-countries"] {
		if val.SVGName == country {
			code = val.Code
			break
		}
	}
	capitalName := getCapitalName(code)

	question := TriviaQuestion{
		TriviaId:    triviaId,
		TypeID:      QUESTION_TYPE_MAP,
		Question:    fmt.Sprintf("What is the capital city of %s?", country),
		Map:         "WorldCapitals",
		Highlighted: capitalName,
	}

	questionId, err := createQuestion(question)
	if err != nil {
		return err
	}

	answer := TriviaAnswer{
		TriviaQuestionID: questionId,
		Text:             capitalName,
		IsCorrect:        true,
	}
	err = createAnswer(answer)
	if err != nil {
		return err
	}

	capitals := copyMapping(mapping.Mappings["world-capitals"])
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

		err = createAnswer(answer)
		if err != nil {
			return err
		}

		capitals = append(capitals[:index], capitals[index+1:]...)
		max = max - 1
	}

	return nil
}

func getCountryName(code string) string {
	for _, value := range mapping.Mappings["world-countries"] {
		if value.Code == code {
			return value.SVGName
		}
	}
	return ""
}

func whatUSState(triviaId int) error {
	states := copyMapping(mapping.Mappings["us-states"])
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

	questionId, err := createQuestion(question)
	if err != nil {
		return err
	}

	answer := TriviaAnswer{
		TriviaQuestionID: questionId,
		Text:             state.SVGName,
		IsCorrect:        true,
	}

	err = createAnswer(answer)
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

		err = createAnswer(answer)
		if err != nil {
			return err
		}

		states = append(states[:index], states[index+1:]...)
		max = max - 1
	}

	return nil
}

func whatFlag(triviaId int) error {
	countries := copyMapping(mapping.Mappings["world-countries"])
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

	questionId, err := createQuestion(question)
	if err != nil {
		return err
	}

	answer := TriviaAnswer{
		TriviaQuestionID: questionId,
		Text:             country.SVGName,
		IsCorrect:        true,
	}
	err = createAnswer(answer)
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

		err = createAnswer(answer)
		if err != nil {
			return err
		}

		countries = append(countries[:index], countries[index+1:]...)
		max = max - 1
	}

	return nil
}

func whatRegionInCountry(triviaId int) error {
	quizzes, err := getCountryRegionQuizzes()
	if err != nil {
		return err
	}

	max := len(quizzes)
	index := rand.Intn(max)
	quiz := quizzes[index]
	mapping := copyMapping(mapping.Mappings[quiz.APIPath])

	max = len(mapping)
	index = rand.Intn(max)
	region := mapping[index]

	question := TriviaQuestion{
		TriviaId:    triviaId,
		TypeID:      QUESTION_TYPE_MAP,
		Question:    fmt.Sprintf("Which %s of %s is highlighted above?", quiz.Singular, quiz.Country),
		Map:         quiz.MapSVG,
		Highlighted: region.SVGName,
	}
	questionId, err := createQuestion(question)
	if err != nil {
		return err
	}

	answer := TriviaAnswer{
		TriviaQuestionID: questionId,
		Text:             region.SVGName,
		IsCorrect:        true,
	}
	err = createAnswer(answer)
	if err != nil {
		return err
	}

	mapping = append(mapping[:index], mapping[index+1:]...)
	max = max - 1
	for i := 0; i < 3; i++ {
		index := rand.Intn(max)
		region = mapping[index]
		answer := TriviaAnswer{
			TriviaQuestionID: questionId,
			Text:             region.SVGName,
			IsCorrect:        false,
		}

		err = createAnswer(answer)
		if err != nil {
			return err
		}

		mapping = append(mapping[:index], mapping[index+1:]...)
		max = max - 1
	}

	return nil
}

func whatFlagInCountry(triviaId int) error {
	quizzes, err := getFlagRegionQuizzes()
	if err != nil {
		return err
	}

	max := len(quizzes)
	index := rand.Intn(max)
	quiz := quizzes[index]
	mapping := copyMapping(mapping.Mappings[quiz.APIPath])

	max = len(mapping)
	index = rand.Intn(max)
	region := mapping[index]

	quizNameSplit := strings.Split(quiz.Name, " ")
	question := TriviaQuestion{
		TriviaId: triviaId,
		TypeID:   QUESTION_TYPE_FLAG,
		Question: fmt.Sprintf("Which of the flags of %s is shown above?", quizNameSplit[len(quizNameSplit)-1]),
		FlagCode: region.Code,
	}
	questionId, err := createQuestion(question)
	if err != nil {
		return err
	}

	answer := TriviaAnswer{
		TriviaQuestionID: questionId,
		Text:             region.SVGName,
		IsCorrect:        true,
	}
	err = createAnswer(answer)
	if err != nil {
		return err
	}

	mapping = append(mapping[:index], mapping[index+1:]...)
	max = max - 1
	for i := 0; i < 3; i++ {
		index := rand.Intn(max)
		region = mapping[index]
		answer := TriviaAnswer{
			TriviaQuestionID: questionId,
			Text:             region.SVGName,
			IsCorrect:        false,
		}

		err = createAnswer(answer)
		if err != nil {
			return err
		}

		mapping = append(mapping[:index], mapping[index+1:]...)
		max = max - 1
	}

	return nil
}

func trueFalseCountryInContinent(triviaId int, answer bool) error {
	countries := copyMapping(mapping.Mappings["world-countries"])
	max := len(countries)
	index := rand.Intn(max)
	country := countries[index]

	var continent string
	if answer {
		continent = strings.Title(country.Group)
	} else {
		group := strings.Title(country.Group)

		continents, err := GetContinents()
		if err != nil {
			return err
		}

		var index int
		for i, val := range continents {
			if strings.ToLower(val.Name) == group {
				index = i
				break
			}
		}

		continents[index] = continents[len(continents)-1]
		continents = continents[:len(continents)-1]
		continent = continents[rand.Intn(len(continents))].Name
	}

	question := TriviaQuestion{
		TriviaId: triviaId,
		TypeID:   QUESTION_TYPE_TEXT,
		Question: fmt.Sprintf("%s is in %s", country.SVGName, continent),
	}

	questionId, err := createQuestion(question)
	if err != nil {
		return err
	}

	answers := []TriviaAnswer{
		{
			TriviaQuestionID: questionId,
			Text:             "True",
			IsCorrect:        answer,
		},
		{
			TriviaQuestionID: questionId,
			Text:             "False",
			IsCorrect:        !answer,
		},
	}

	for _, answer := range answers {
		err = createAnswer(answer)
		if err != nil {
			return err
		}
	}

	return nil
}

func trueFalseCapitalOfCountry(triviaId int, answer bool) error {
	countries := copyMapping(mapping.Mappings["world-countries"])
	max := len(countries)
	index := rand.Intn(max)
	country := countries[index]

	var capitalName string
	if answer {
		capitalName = getCapitalName(country.Code)
	} else {
		capitals := copyMapping(mapping.Mappings["world-capitals"])
		var index int
		for i, val := range capitals {
			if val.Code == country.Code {
				index = i
				break
			}
		}

		capitals[index] = capitals[len(capitals)-1]
		capitals = capitals[:len(capitals)-1]
		max := len(capitals)
		index = rand.Intn(max)
		capital := capitals[index]
		capitalName = capital.SVGName
	}

	question := TriviaQuestion{
		TriviaId: triviaId,
		TypeID:   QUESTION_TYPE_TEXT,
		Question: fmt.Sprintf("%s is the capital city of %s", capitalName, country.SVGName),
	}

	questionId, err := createQuestion(question)
	if err != nil {
		return err
	}

	answers := []TriviaAnswer{
		{
			TriviaQuestionID: questionId,
			Text:             "True",
			IsCorrect:        answer,
		},
		{
			TriviaQuestionID: questionId,
			Text:             "False",
			IsCorrect:        !answer,
		},
	}

	for _, answer := range answers {
		err = createAnswer(answer)
		if err != nil {
			return err
		}
	}

	return nil
}

func getCapitalName(code string) string {
	for _, value := range mapping.Mappings["world-capitals"] {
		if value.Code == code {
			return value.SVGName
		}
	}
	return ""
}

func trueFalseRegionInCountry(triviaId int, answer bool) error {
	quizzes, err := getCountryRegionQuizzes()
	if err != nil {
		return err
	}

	max := len(quizzes)
	index := rand.Intn(max)
	quiz := quizzes[index]
	mappingOne := copyMapping(mapping.Mappings[quiz.APIPath])

	var region string
	if answer {
		max = len(mappingOne)
		index = rand.Intn(max)
		region = mappingOne[index].SVGName
	} else {
		var index int
		for i, val := range quizzes {
			if val.Name == quiz.Name {
				index = i
				break
			}
		}

		quizzes[index] = quizzes[len(quizzes)-1]
		quizzes = quizzes[:len(quizzes)-1]
		max = max - 1

		index = rand.Intn(max)
		quizTwo := quizzes[index]
		mappingTwo := copyMapping(mapping.Mappings[quizTwo.APIPath])

		max = len(mappingTwo)
		index = rand.Intn(max)
		region = mappingTwo[index].SVGName
	}

	question := TriviaQuestion{
		TriviaId: triviaId,
		TypeID:   QUESTION_TYPE_TEXT,
		Question: fmt.Sprintf("%s is a %s of %s", region, quiz.Singular, quiz.Country),
	}

	questionId, err := createQuestion(question)
	if err != nil {
		return err
	}

	answers := []TriviaAnswer{
		{
			TriviaQuestionID: questionId,
			Text:             "True",
			IsCorrect:        answer,
		},
		{
			TriviaQuestionID: questionId,
			Text:             "False",
			IsCorrect:        !answer,
		},
	}

	for _, answer := range answers {
		err = createAnswer(answer)
		if err != nil {
			return err
		}
	}

	return nil
}

func trueFalseFlagForCountry(triviaId int, answer bool) error {
	quizzes, err := getFlagRegionQuizzes()
	if err != nil {
		return err
	}

	max := len(quizzes)
	index := rand.Intn(max)
	quiz := quizzes[index]
	mappingOne := copyMapping(mapping.Mappings[quiz.APIPath])

	var flagCode string
	if answer {
		max = len(mappingOne)
		index = rand.Intn(max)
		flagCode = mappingOne[index].Code
	} else {
		var index int
		for i, val := range quizzes {
			if val.Name == quiz.Name {
				index = i
				break
			}
		}

		quizzes[index] = quizzes[len(quizzes)-1]
		quizzes = quizzes[:len(quizzes)-1]
		max = max - 1

		index = rand.Intn(max)
		quizTwo := quizzes[index]
		mappingTwo := copyMapping(mapping.Mappings[quizTwo.APIPath])

		max = len(mappingTwo)
		index = rand.Intn(max)
		flagCode = mappingTwo[index].Code
	}

	quizNameSplit := strings.Split(quiz.Name, " ")
	question := TriviaQuestion{
		TriviaId: triviaId,
		TypeID:   QUESTION_TYPE_FLAG,
		Question: fmt.Sprintf("This is a flag of %s", quizNameSplit[len(quizNameSplit)-1]),
		FlagCode: flagCode,
	}

	questionId, err := createQuestion(question)
	if err != nil {
		return err
	}

	answers := []TriviaAnswer{
		{
			TriviaQuestionID: questionId,
			Text:             "True",
			IsCorrect:        answer,
		},
		{
			TriviaQuestionID: questionId,
			Text:             "False",
			IsCorrect:        !answer,
		},
	}

	for _, answer := range answers {
		err = createAnswer(answer)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetAllTrivia() ([]Trivia, error) {
	rows, err := Connection.Query("SELECT * FROM trivia ORDER BY date DESC LIMIT 20;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trivia = []Trivia{}
	for rows.Next() {
		var quiz Trivia
		if err = rows.Scan(&quiz.ID, &quiz.Name, &quiz.Date); err != nil {
			return nil, err
		}
		trivia = append(trivia, quiz)
	}
	return trivia, rows.Err()
}

func GetTrivia(date string) (*TriviaDto, error) {
	var result TriviaDto
	err := Connection.QueryRow("SELECT id, name from trivia WHERE date = $1;", date).Scan(&result.ID, &result.Name)
	if err != nil {
		return nil, err
	}

	rows, err := Connection.Query("SELECT q.id, t.name, q.question, q.map, q.highlighted, q.flagCode, q.imageUrl FROM triviaQuestions q JOIN triviaQuestionType t ON t.id = q.typeId WHERE q.triviaId = $1;", result.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions = []QuestionDto{}
	for rows.Next() {
		var question QuestionDto
		if err = rows.Scan(&question.ID, &question.Type, &question.Question, &question.Map, &question.Highlighted, &question.FlagCode, &question.ImageURL); err != nil {
			return nil, err
		}

		rows, err := Connection.Query("SELECT text, isCorrect, flagCode FROM triviaAnswers WHERE triviaQuestionId = $1;", question.ID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var answers = []AnswerDto{}
		for rows.Next() {
			var answer AnswerDto
			if err = rows.Scan(&answer.Text, &answer.IsCorrect, &answer.FlagCode); err != nil {
				return nil, err
			}
			answers = append(answers, answer)
		}

		if len(answers) > 2 {
			shuffleAnswers(answers)
		}

		question.Answers = answers
		questions = append(questions, question)
	}

	rand.Shuffle(len(questions), func(i, j int) {
		questions[i], questions[j] = questions[j], questions[i]
	})

	result.Questions = questions
	return &result, rows.Err()
}

func shuffleAnswers(slice []AnswerDto) {
	for i := range slice {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func copyMapping(orig []mapping.Mapping) []mapping.Mapping {
	cpy := make([]mapping.Mapping, len(orig))
	copy(cpy, orig)
	return cpy
}
