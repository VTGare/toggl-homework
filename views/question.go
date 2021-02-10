package views

import (
	"sort"
	"time"

	"github.com/VTGare/toggl-homework/models"
)

//Question is a question model
type Question struct {
	ID        int64     `json:"id"`
	Body      string    `json:"body"`
	Options   []*Option `json:"options"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//Option is an option model
type Option struct {
	Body    string `json:"body"`
	Correct bool   `json:"correct"`
}

//BuildQuestions turns joined questions into Question models.
func BuildQuestions(qModels []*models.JoinedQuestion) []*Question {
	//Use map to store slices to prevent duplicates and faster access.
	views := make(map[int64]*Question)

	for _, m := range qModels {
		if v, ok := views[m.ID]; ok {
			v.Options = append(v.Options, &Option{m.Option, m.Correct})
		} else {
			q := &Question{ID: m.ID, Body: m.Body, Options: make([]*Option, 1), CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}

			q.Options[0] = &Option{m.Option, m.Correct}
			views[m.ID] = q
		}
	}

	//Prepare a slice to get sorted.
	sorted := make([]*Question, 0, len(views))
	for _, v := range views {
		sorted = append(sorted, v)
	}

	//Sort questions by ID to preserve insert order.
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].ID < sorted[j].ID
	})

	return sorted
}

//ToModel turns a Question into a Question model.
func (o *Option) ToModel() *models.Option {
	return &models.Option{Body: o.Body, Correct: o.Correct}
}
