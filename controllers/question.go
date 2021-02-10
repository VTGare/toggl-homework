package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/VTGare/toggl-homework/models"
	"github.com/VTGare/toggl-homework/views"
	"github.com/gofiber/fiber/v2"
)

type questionResp struct {
	*views.Response
	QuestionID int64 `json:"question_id"`
}

//RegisterQuestion registers question controllers
func RegisterQuestion(app *fiber.App) {
	app.Get("/api/questions", GetQuestions)
	app.Post("/api/questions", AddQuestion)
	app.Put("/api/questions/:id", EditQuestion)
	app.Delete("/api/questions/:id", DeleteQuestion)
}

func GetQuestions(c *fiber.Ctx) error {
	mquestions, err := models.GetQuestions()
	if err != nil {
		return views.Error(500, err).Send(c)
	}

	questions := views.BuildQuestions(mquestions)
	return c.JSON(questions)
}

func AddQuestion(c *fiber.Ctx) error {
	reqQuestion := &views.Question{}

	err := c.BodyParser(reqQuestion)
	if err != nil {
		return views.Error(400, err).Send(c)
	}

	if len(reqQuestion.Options) == 0 {
		return views.ErrNoOptions.Send(c)
	}

	id, err := models.InsertQuestion(&models.Question{Body: reqQuestion.Body})
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "UNIQUE constraint failed"):
			return views.Message(400, fmt.Sprintf("Question %v already exists", reqQuestion.Body)).Send(c)
		default:
			return views.Error(500, err).Send(c)
		}
	}

	//I don't think SQLite supports bulk insert so we'll have to do it this way.
	for _, option := range reqQuestion.Options {
		err := models.InsertOptions(option.ToModel(), id)
		if err != nil {
			return views.Error(500, err).Send(c)
		}
	}

	return c.JSON(questionResp{views.Message(200, "success"), id})
}

func EditQuestion(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 0, 64)
	if err != nil {
		return views.ErrParseID.Send(c)
	}

	reqQuestion := &views.Question{}
	err = c.BodyParser(reqQuestion)
	if err != nil {
		return views.Error(400, err).Send(c)
	}

	_, err = models.EditQuestion(&models.Question{ID: id, Body: reqQuestion.Body})
	if err != nil {
		return views.Error(500, err).Send(c)
	}

	for _, option := range reqQuestion.Options {
		err := models.ReplaceOption(option.ToModel(), id)
		if err != nil {
			return views.Error(500, err).Send(c)
		}
	}

	if err != nil {
		return views.Error(500, err).Send(c)
	}

	return c.JSON(questionResp{views.Message(200, "success"), id})
}

func DeleteQuestion(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 0, 64)
	if err != nil {
		return views.ErrParseID.Send(c)
	}

	aff, err := models.DeleteQuestion(id)
	if err != nil {
		return views.Error(500, err).Send(c)
	}

	if aff == 0 {
		return views.Message(400, fmt.Sprintf("ID %v doesn't exist", id)).Send(c)
	}

	return c.JSON(questionResp{views.Message(200, "success"), id})
}
