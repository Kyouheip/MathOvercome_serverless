package model

import (
	"encoding/gob"
	"time"
)

func init() {
	// gin-contrib/sessions の cookie store は gob でシリアライズするため
	// セッションに保存する型を事前登録する必要がある
	gob.Register(&User{})
	gob.Register(uint64(0))
}

type User struct {
	ID           uint64 `gorm:"primaryKey"`
	UserName     string
	UserID       string        `gorm:"unique"`
	Password     string        `json:"-"`
	TestSessions []TestSession `json:",omitempty"`
}

type Category struct {
	ID   uint64 `gorm:"primaryKey"`
	Name string
}

type Problem struct {
	ID         uint64 `gorm:"primaryKey"`
	CategoryID int
	Question   string   `gorm:"type:text"`
	Hint       string
	Choices    []Choice `json:",omitempty"`
}

type Choice struct {
	ID         uint64 `gorm:"primaryKey"`
	ProblemID  uint64
	Problem    Problem `json:"-"`
	ChoiceText string
	IsCorrect  bool
}

type TestSession struct {
	ID              uint64    `gorm:"primaryKey"`
	UserID          uint64
	User            User      `json:"-"`
	IncludeIntegers bool
	StartTime       time.Time        `gorm:"default:CURRENT_TIMESTAMP"`
	SessionProblems []SessionProblem `json:",omitempty"`
}

type SessionProblem struct {
	ID               uint64      `gorm:"primaryKey"`
	TestSessionID    uint64      `gorm:"column:session_id"`
	TestSession      TestSession `json:"-"`
	ProblemID        uint64
	Problem          Problem
	SelectedChoiceID *uint64 `gorm:"column:selected_choice_id"`
	SelectedChoice   *Choice
	IsCorrect        *bool
}

// Java 側の DB テーブル名に合わせる (GORM デフォルトだと session_problems になる)
func (SessionProblem) TableName() string {
	return "sessionproblems"
}
