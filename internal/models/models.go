package models

import "time"

type Question struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Text      string    `json:"text" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	Answers   []Answer  `json:"answers,omitempty" gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE"`
}

type Answer struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	QuestionID uint      `json:"question_id" gorm:"not null;index"`
	UserID     string    `json:"user_id" gorm:"not null;index"`
	Text       string    `json:"text" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`
}
