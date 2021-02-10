package views

import (
	"sort"
	"time"

	"github.com/VTGare/toggl-homework/models"
)

type Question struct {
	ID        int64     `json:"id"`
	Body      string    `json:"body"`
	Options   []*Option `json:"options"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Option struct {
	Body    string `json:"body"`
	Correct bool   `json:"correct"`
}

func BuildQuestions(qModels []*models.JoinedQuestion) []*Question {
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

	sorted := make([]*Question, 0, len(views))
	for _, v := range views {
		sorted = append(sorted, v)
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].ID < sorted[j].ID
	})

	return sorted
}

func (o *Option) ToModel() *models.Option {
	return &models.Option{Body: o.Body, Correct: o.Correct}
}
