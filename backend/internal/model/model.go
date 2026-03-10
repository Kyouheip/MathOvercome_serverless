package model

import (
	"encoding/gob"
	"time"
)

func init() {
	gob.Register(&User{})
	gob.Register(uint64(0))
}

type User struct {
	ID       uint64
	UserName string
	UserID   string
	Password string `json:"-"`
}

type Category struct {
	ID   uint64
	Name string
}

type Problem struct {
	ID         uint64
	CategoryID int
	Question   string
	Hint       string
	Choices    []Choice `json:",omitempty"`
}

type Choice struct {
	ID         uint64
	ProblemID  uint64
	ChoiceText string
	IsCorrect  bool
}

type TestSession struct {
	ID              uint64
	UserID          uint64
	IncludeIntegers bool
	StartTime       time.Time
	SessionProblems []SessionProblem `json:",omitempty"`
}

type SessionProblem struct {
	ID               uint64
	TestSessionID    uint64
	ProblemID        uint64
	Problem          Problem
	SelectedChoiceID *uint64
	IsCorrect        *bool
	CategoryName     string
	CategoryID       int
}
