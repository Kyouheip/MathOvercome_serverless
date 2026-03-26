package dto

// --- Requests ---

type AnswerRequest struct {
	SelectedChoiceID *int64 `json:"selectedChoiceId"`
}

// --- Responses ---

type Choice struct {
	ID         int64  `json:"id"`
	ChoiceText string `json:"choiceText"`
}

type SessionProblem struct {
	ID         int64    `json:"id"`
	Question   string   `json:"question"`
	Choices    []Choice `json:"choices"`
	Hint       string   `json:"hint"`
	SelectedID *int64   `json:"selectedId"`
	Total      int      `json:"total"`
}

type Category struct {
	CategoryName string `json:"categoryName"`
	Total        int    `json:"total"`
	CorrectCount int    `json:"correctCount"`
}

type ProblemCategory struct {
	IsCorrect    *bool  `json:"isCorrect"`
	CategoryName string `json:"categoryName"`
}

type TestSession struct {
	SessionID        int64             `json:"sessionId"`
	StartTime        string            `json:"startTime"`
	Total            int               `json:"total"`
	CorrectCount     int               `json:"correctCount"`
	ProbCategoryDtos []ProblemCategory `json:"-"`
	CategoryDtos     []Category        `json:"categoryDtos"`
	WeakCategories   []string          `json:"weakCategories"`
}

type User struct {
	UserName     string        `json:"userName"`
	TestSessDtos []TestSession `json:"testSessDtos"`
}
