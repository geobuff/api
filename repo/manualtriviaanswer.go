package repo

type ManualTriviaAnswer struct {
	ID                     int    `json:"id"`
	ManualTriviaQuestionID int    `json:"manualTriviaQuestionId"`
	Text                   string `json:"text"`
	IsCorrect              bool   `json:"isCorrect"`
	FlagCode               string `json:"flagCode"`
}
