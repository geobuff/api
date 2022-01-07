package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/geobuff/mapping"
	pluralize "github.com/gertd/go-pluralize"
)

type DailyTrivia struct {
	ID    int       `json:"id"`
	Name  string    `json:"name"`
	Date  time.Time `json:"date"`
	Plays int       `json:"plays"`
}

type DailyTriviaQuestion struct {
	ID            int    `json:"id"`
	DailyTriviaId int    `json:"dailyTriviaId"`
	Type          string `json:"type"`
	Question      string `json:"question"`
	Map           string `json:"map"`
	Highlighted   string `json:"hightlighted"`
	FlagCode      string `json:"flagCode"`
	ImageURL      string `json:"imageUrl"`
}

type DailyTriviaAnswer struct {
	ID                    int    `json:"id"`
	DailyTriviaQuestionID int    `json:"dailyTriviaQuestionId"`
	Text                  string `json:"text"`
	IsCorrect             bool   `json:"isCorrect"`
	FlagCode              string `json:"flagCode"`
}

type DailyTriviaDto struct {
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

var continents = []string{
	"Africa",
	"Asia",
	"Europe",
	"South America",
	"North America",
	"Oceania",
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
}

var p = pluralize.NewClient()

func CreateDailyTrivia() error {
	date := time.Now().AddDate(0, 0, 1)
	var id int
	if err := Connection.QueryRow("SELECT id FROM dailyTrivia WHERE date = $1", date).Scan(&id); err != sql.ErrNoRows {
		return errors.New("daily trivia for current date already created")
	}

	_, month, day := date.Date()
	weekday := date.Weekday().String()
	statement := "INSERT INTO dailyTrivia (name, date, plays) VALUES ($1, $2, $3) RETURNING id;"
	err := Connection.QueryRow(statement, fmt.Sprintf("%s, %s %d", weekday, month, day), date, 0).Scan(&id)
	if err != nil {
		return err
	}
	return generateQuestions(id)
}

func generateQuestions(dailyTriviaId int) error {
	err := whatCountry(dailyTriviaId)
	if err != nil {
		return err
	}

	err = whatCapital(dailyTriviaId)
	if err != nil {
		return err
	}

	err = whatUSState(dailyTriviaId)
	if err != nil {
		return err
	}

	err = whatFlag(dailyTriviaId)
	if err != nil {
		return err
	}

	err = whatRegionInCountry(dailyTriviaId)
	if err != nil {
		return err
	}

	err = whatFlagInCountry(dailyTriviaId)
	if err != nil {
		return err
	}

	err = trueFalseCountryInContinent(dailyTriviaId, randomBool())
	if err != nil {
		return err
	}

	err = trueFalseCapitalOfCountry(dailyTriviaId, randomBool())
	if err != nil {
		return err
	}

	err = trueFalseRegionInCountry(dailyTriviaId, randomBool())
	if err != nil {
		return err
	}

	err = trueFalseFlagForCountry(dailyTriviaId, randomBool())
	if err != nil {
		return err
	}

	return nil
}

func randomBool() bool {
	return rand.Float32() < 0.5
}

func createQuestion(question DailyTriviaQuestion) (int, error) {
	statement := "INSERT INTO dailyTriviaQuestions (dailyTriviaId, type, question, map, highlighted, flagCode, imageUrl) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, question.DailyTriviaId, question.Type, question.Question, question.Map, question.Highlighted, question.FlagCode, question.ImageURL).Scan(&id)
	return id, err
}

func createAnswer(answer DailyTriviaAnswer) error {
	statement := "INSERT INTO dailyTriviaAnswers (dailyTriviaQuestionId, text, isCorrect, flagCode) VALUES ($1, $2, $3, $4) RETURNING id;"
	var id int
	return Connection.QueryRow(statement, answer.DailyTriviaQuestionID, answer.Text, answer.IsCorrect, answer.FlagCode).Scan(&id)
}

func whatCountry(dailyTriviaId int) error {
	max := len(topLandmass)
	index := rand.Intn(max)
	country := topLandmass[index]

	question := DailyTriviaQuestion{
		DailyTriviaId: dailyTriviaId,
		Type:          "map",
		Question:      "Which country is highlighted above?",
		Map:           "WorldCountries",
		Highlighted:   country,
	}

	questionId, err := createQuestion(question)
	if err != nil {
		return err
	}

	answer := DailyTriviaAnswer{
		DailyTriviaQuestionID: questionId,
		Text:                  country,
		IsCorrect:             true,
	}
	err = createAnswer(answer)
	if err != nil {
		return err
	}

	countries := mapping.Mappings["world-countries"]
	for i, val := range countries {
		if val.SVGName == country {
			index = i
			break
		}
	}
	countries = append(countries[:index], countries[index+1:]...)
	max = len(countries)

	for i := 0; i < 3; i++ {
		index := rand.Intn(max + 1)
		country = countries[index].SVGName
		answer := DailyTriviaAnswer{
			DailyTriviaQuestionID: questionId,
			Text:                  country,
			IsCorrect:             false,
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

func whatCapital(dailyTriviaId int) error {
	capitals := mapping.Mappings["world-capitals"]
	max := len(capitals)
	index := rand.Intn(max)
	capital := capitals[index]
	capitals = append(capitals[:index], capitals[index+1:]...)
	max = max - 1
	countryName := getCountryName(capital.Code)

	question := DailyTriviaQuestion{
		DailyTriviaId: dailyTriviaId,
		Type:          "map",
		Question:      fmt.Sprintf("What is the capital city of %s?", countryName),
		Map:           "WorldCapitals",
		Highlighted:   capital.SVGName,
	}

	questionId, err := createQuestion(question)
	if err != nil {
		return err
	}

	answer := DailyTriviaAnswer{
		DailyTriviaQuestionID: questionId,
		Text:                  capital.SVGName,
		IsCorrect:             true,
	}
	err = createAnswer(answer)
	if err != nil {
		return err
	}

	for i := 0; i < 3; i++ {
		index := rand.Intn(max)
		capital := capitals[index]
		answer := DailyTriviaAnswer{
			DailyTriviaQuestionID: questionId,
			Text:                  capital.SVGName,
			IsCorrect:             false,
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

func whatUSState(dailyTriviaId int) error {
	states := mapping.Mappings["us-states"]
	max := len(states)
	index := rand.Intn(max)
	state := states[index]
	states = append(states[:index], states[index+1:]...)
	max = max - 1

	question := DailyTriviaQuestion{
		DailyTriviaId: dailyTriviaId,
		Type:          "map",
		Question:      "Which US state is highlighted above?",
		Map:           "UsStates",
		Highlighted:   state.SVGName,
	}

	questionId, err := createQuestion(question)
	if err != nil {
		return err
	}

	answer := DailyTriviaAnswer{
		DailyTriviaQuestionID: questionId,
		Text:                  state.SVGName,
		IsCorrect:             true,
	}

	err = createAnswer(answer)
	if err != nil {
		return err
	}

	for i := 0; i < 3; i++ {
		index := rand.Intn(max + 1)
		state = states[index]
		answer := DailyTriviaAnswer{
			DailyTriviaQuestionID: questionId,
			Text:                  state.SVGName,
			IsCorrect:             false,
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

func whatFlag(dailyTriviaId int) error {
	countries := mapping.Mappings["world-countries"]
	max := len(countries)
	index := rand.Intn(max)
	country := countries[index]
	countries = append(countries[:index], countries[index+1:]...)
	max = max - 1

	question := DailyTriviaQuestion{
		DailyTriviaId: dailyTriviaId,
		Type:          "flag",
		Question:      "Which country has this flag?",
		FlagCode:      country.Code,
	}

	questionId, err := createQuestion(question)
	if err != nil {
		return err
	}

	answer := DailyTriviaAnswer{
		DailyTriviaQuestionID: questionId,
		Text:                  country.SVGName,
		IsCorrect:             true,
	}
	err = createAnswer(answer)
	if err != nil {
		return err
	}

	for i := 0; i < 3; i++ {
		index := rand.Intn(max + 1)
		country = countries[index]
		answer := DailyTriviaAnswer{
			DailyTriviaQuestionID: questionId,
			Text:                  country.SVGName,
			IsCorrect:             false,
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

func whatRegionInCountry(dailyTriviaId int) error {
	quizzes, err := getCountryRegionQuizzes()
	if err != nil {
		return err
	}

	max := len(quizzes)
	index := rand.Intn(max)
	quiz := quizzes[index]
	mapping := mapping.Mappings[quiz.APIPath]

	max = len(mapping)
	index = rand.Intn(max)
	region := mapping[index]

	quizNameSplit := strings.Split(quiz.Name, " ")
	question := DailyTriviaQuestion{
		DailyTriviaId: dailyTriviaId,
		Type:          "map",
		Question:      fmt.Sprintf("Which %s of %s is highlighted above?", p.Singular(strings.ToLower(quizNameSplit[0])), quizNameSplit[len(quizNameSplit)-1]),
		Map:           quiz.MapSVG,
		Highlighted:   region.SVGName,
	}
	questionId, err := createQuestion(question)
	if err != nil {
		return err
	}

	answer := DailyTriviaAnswer{
		DailyTriviaQuestionID: questionId,
		Text:                  region.SVGName,
		IsCorrect:             true,
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
		answer := DailyTriviaAnswer{
			DailyTriviaQuestionID: questionId,
			Text:                  region.SVGName,
			IsCorrect:             false,
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

func whatFlagInCountry(dailyTriviaId int) error {
	quizzes, err := getFlagRegionQuizzes()
	if err != nil {
		return err
	}

	max := len(quizzes)
	index := rand.Intn(max)
	quiz := quizzes[index]
	mapping := mapping.Mappings[quiz.APIPath]

	max = len(mapping)
	index = rand.Intn(max)
	region := mapping[index]

	quizNameSplit := strings.Split(quiz.Name, " ")
	question := DailyTriviaQuestion{
		DailyTriviaId: dailyTriviaId,
		Type:          "flag",
		Question:      fmt.Sprintf("Which of the flags of %s is shown above?", quizNameSplit[len(quizNameSplit)-1]),
		FlagCode:      region.Code,
	}
	questionId, err := createQuestion(question)
	if err != nil {
		return err
	}

	answer := DailyTriviaAnswer{
		DailyTriviaQuestionID: questionId,
		Text:                  region.SVGName,
		IsCorrect:             true,
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
		answer := DailyTriviaAnswer{
			DailyTriviaQuestionID: questionId,
			Text:                  region.SVGName,
			IsCorrect:             false,
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

func trueFalseCountryInContinent(dailyTriviaId int, answer bool) error {
	countries := mapping.Mappings["world-countries"]
	max := len(countries)
	index := rand.Intn(max)
	country := countries[index]

	var continent string
	if answer {
		continent = strings.Title(country.Group)
	} else {
		group := strings.Title(country.Group)
		temp := make([]string, len(continents))
		copy(temp, continents)
		var index int
		for i, val := range temp {
			if strings.ToLower(val) == group {
				index = i
				break
			}
		}

		temp[index] = temp[len(temp)-1]
		temp = temp[:len(temp)-1]
		continent = temp[rand.Intn(len(temp))-1]
	}

	question := DailyTriviaQuestion{
		DailyTriviaId: dailyTriviaId,
		Type:          "text",
		Question:      fmt.Sprintf("%s is in %s", country.SVGName, continent),
	}

	questionId, err := createQuestion(question)
	if err != nil {
		return err
	}

	answers := []DailyTriviaAnswer{
		{
			DailyTriviaQuestionID: questionId,
			Text:                  "True",
			IsCorrect:             answer,
		},
		{
			DailyTriviaQuestionID: questionId,
			Text:                  "False",
			IsCorrect:             !answer,
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

func trueFalseCapitalOfCountry(dailyTriviaId int, answer bool) error {
	countries := mapping.Mappings["world-countries"]
	max := len(countries)
	index := rand.Intn(max)
	country := countries[index]

	var capitalName string
	if answer {
		capitalName = getCapitalName(country.Code)
	} else {
		capitals := mapping.Mappings["world-capitals"]
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

	question := DailyTriviaQuestion{
		DailyTriviaId: dailyTriviaId,
		Type:          "text",
		Question:      fmt.Sprintf("%s is the capital city of %s", capitalName, country.SVGName),
	}

	questionId, err := createQuestion(question)
	if err != nil {
		return err
	}

	answers := []DailyTriviaAnswer{
		{
			DailyTriviaQuestionID: questionId,
			Text:                  "True",
			IsCorrect:             answer,
		},
		{
			DailyTriviaQuestionID: questionId,
			Text:                  "False",
			IsCorrect:             !answer,
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

func trueFalseRegionInCountry(dailyTriviaId int, answer bool) error {
	quizzes, err := getCountryRegionQuizzes()
	if err != nil {
		return err
	}

	max := len(quizzes)
	index := rand.Intn(max)
	quiz := quizzes[index]
	mappingOne := mapping.Mappings[quiz.APIPath]

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
		mappingTwo := mapping.Mappings[quizTwo.APIPath]

		max = len(mappingTwo)
		index = rand.Intn(max)
		region = mappingTwo[index].SVGName
	}

	quizNameSplit := strings.Split(quiz.Name, " ")
	question := DailyTriviaQuestion{
		DailyTriviaId: dailyTriviaId,
		Type:          "text",
		Question:      fmt.Sprintf("%s is a %s of %s", region, p.Singular(strings.ToLower(quizNameSplit[0])), quizNameSplit[len(quizNameSplit)-1]),
	}

	questionId, err := createQuestion(question)
	if err != nil {
		return err
	}

	answers := []DailyTriviaAnswer{
		{
			DailyTriviaQuestionID: questionId,
			Text:                  "True",
			IsCorrect:             answer,
		},
		{
			DailyTriviaQuestionID: questionId,
			Text:                  "False",
			IsCorrect:             !answer,
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

func trueFalseFlagForCountry(dailyTriviaId int, answer bool) error {
	quizzes, err := getFlagRegionQuizzes()
	if err != nil {
		return err
	}

	max := len(quizzes)
	index := rand.Intn(max)
	quiz := quizzes[index]
	mappingOne := mapping.Mappings[quiz.APIPath]

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
		mappingTwo := mapping.Mappings[quizTwo.APIPath]

		max = len(mappingTwo)
		index = rand.Intn(max)
		flagCode = mappingTwo[index].Code
	}

	quizNameSplit := strings.Split(quiz.Name, " ")
	question := DailyTriviaQuestion{
		DailyTriviaId: dailyTriviaId,
		Type:          "flag",
		Question:      fmt.Sprintf("This is a flag of %s", quizNameSplit[len(quizNameSplit)-1]),
		FlagCode:      flagCode,
	}

	questionId, err := createQuestion(question)
	if err != nil {
		return err
	}

	answers := []DailyTriviaAnswer{
		{
			DailyTriviaQuestionID: questionId,
			Text:                  "True",
			IsCorrect:             answer,
		},
		{
			DailyTriviaQuestionID: questionId,
			Text:                  "False",
			IsCorrect:             !answer,
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

func GetAllDailyTrivia() ([]DailyTrivia, error) {
	rows, err := Connection.Query("SELECT * FROM dailyTrivia ORDER BY date DESC LIMIT 20;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trivia = []DailyTrivia{}
	for rows.Next() {
		var quiz DailyTrivia
		if err = rows.Scan(&quiz.ID, &quiz.Name, &quiz.Date, &quiz.Plays); err != nil {
			return nil, err
		}
		trivia = append(trivia, quiz)
	}
	return trivia, rows.Err()
}

func GetDailyTrivia(date string) (*DailyTriviaDto, error) {
	var result DailyTriviaDto
	err := Connection.QueryRow("SELECT id, name from dailyTrivia WHERE date = $1;", date).Scan(&result.ID, &result.Name)
	if err != nil {
		return nil, err
	}

	rows, err := Connection.Query("SELECT id, type, question, map, highlighted, flagCode, imageUrl FROM dailyTriviaQuestions WHERE dailyTriviaId = $1;", result.ID)
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

		rows, err := Connection.Query("SELECT text, isCorrect, flagCode FROM dailyTriviaAnswers WHERE dailyTriviaQuestionId = $1;", question.ID)
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

func IncrementTriviaPlays(id int) error {
	statement := "UPDATE dailytrivia set plays = plays + 1 WHERE id = $1 RETURNING id;"
	var triviaId int
	return Connection.QueryRow(statement, id).Scan(&triviaId)
}
