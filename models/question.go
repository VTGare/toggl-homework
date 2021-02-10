package models

import (
	"time"

	"github.com/VTGare/toggl-homework/database"
)

type Question struct {
	ID        int64     `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type JoinedQuestion struct {
	ID        int64
	Body      string
	Option    string
	Correct   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func GetQuestions() ([]*JoinedQuestion, error) {
	query := `
	SELECT questions.id, questions.body, options.body, correct, questions.created_at, questions.updated_at
	FROM questions 
	INNER JOIN options ON options.question_id = questions.id
	ORDER BY options.created_at ASC
	`

	rows, err := database.DBConn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	questions := []*JoinedQuestion{}

	for rows.Next() {
		q := &JoinedQuestion{}
		err := rows.Scan(&q.ID, &q.Body, &q.Option, &q.Correct, &q.CreatedAt, &q.UpdatedAt)
		if err != nil {
			return nil, err
		}

		questions = append(questions, q)
	}

	return questions, nil
}

func InsertQuestion(quest *Question) (int64, error) {
	query := "INSERT INTO questions (body) VALUES ($1)"

	res, err := database.DBConn.Exec(query, quest.Body)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func DeleteQuestion(id int64) (int64, error) {
	query := "DELETE FROM questions WHERE id = $1"

	res, err := database.DBConn.Exec(query, id)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func EditQuestion(quest *Question) (int64, error) {
	query := "UPDATE questions SET body = $1 WHERE id = $2"

	res, err := database.DBConn.Exec(query, quest.Body, quest.ID)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
