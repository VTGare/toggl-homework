package controllers

import (
	"fmt"
	"strconv"

	"github.com/VTGare/toggl-homework/models"
	"github.com/VTGare/toggl-homework/views"
	"github.com/gofiber/fiber/v2"
)

type questionResp struct {
	Status     int   `json:"status"`
	QuestionID int64 `json:"question_id"`
}

func RegisterQuestion(app *fiber.App) {
	app.Get("/api/questions", GetQuestions)
	app.Post("/api/questions", AddQuestion)
	app.Put("/api/questions/:id", EditQuestion)
	app.Delete("/api/questions/:id", DeleteQuestion)
}

func GetQuestions(c *fiber.Ctx) error {
	mquestions, err := models.GetQuestions()
	if err != nil {
		return err
	}

	questions := views.BuildQuestions(mquestions)
	return c.JSON(questions)
}

func AddQuestion(c *fiber.Ctx) error {
	reqQuestion := &views.Question{}

	err := c.BodyParser(reqQuestion)
	if err != nil {
		return err
	}

	if len(reqQuestion.Options) == 0 {
		err := views.Error{Status: 400, Message: "Empty options"}
		return c.Status(err.Status).JSON(err)
	}

	question, err := models.InsertQuestion(&models.Question{Body: reqQuestion.Body})
	if err != nil {
		return err
	}

	for _, option := range reqQuestion.Options {
		err := models.InsertOptions(option.ToModel(), question.ID)
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	return c.JSON(questionResp{200, question.ID})
}

func EditQuestion(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 0, 64)
	if err != nil {
		return err
	}

	reqQuestion := &views.Question{}
	err = c.BodyParser(reqQuestion)
	if err != nil {
		return err
	}

	if len(reqQuestion.Options) == 0 {
		err := views.Error{Status: 400, Message: "Empty options"}
		return c.Status(err.Status).JSON(err)
	}

	_, err = models.EditQuestion(&models.Question{ID: id, Body: reqQuestion.Body})
	if err != nil {
		return err
	}

	for _, option := range reqQuestion.Options {
		err := models.ReplaceOption(option.ToModel(), id)
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	return c.JSON(questionResp{200, id})
}

func DeleteQuestion(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 0, 64)
	if err != nil {
		return err
	}

	aff, err := models.DeleteQuestion(id)
	if err != nil {
		return err
	}

	if aff == 0 {
		err := views.Error{Status: 400, Message: fmt.Sprintf("ID %v doesn't exist", id)}
		return c.Status(err.Status).JSON(err)
	}

	return c.JSON(questionResp{200, id})
}
