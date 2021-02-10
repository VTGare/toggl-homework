package models

import (
	"time"

	"github.com/VTGare/toggl-homework/database"
)

type Option struct {
	ID         int64     `json:"id"`
	QuestionID int64     `json:"question_id"`
	Body       string    `json:"body"`
	Correct    bool      `json:"correct"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func InsertOptions(option *Option, questionID int64) error {
	query := "INSERT INTO options (question_id, body, correct) VALUES ($1, $2, $3)"

	_, err := database.DBConn.Exec(query, questionID, option.Body, option.Correct)
	if err != nil {
		return err
	}

	return nil
}

func ReplaceOption(option *Option, questionID int64) error {
	query := "REPLACE INTO options(question_id, body, correct) VALUES ($1, $2, $3)"

	_, err := database.DBConn.Exec(query, questionID, option.Body, option.Correct)
	if err != nil {
		return err
	}

	return nil
}
