package views

import "github.com/gofiber/fiber/v2"

//Response is a standard response struct.
type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

//Collection of common errors
var (
	ErrParseID   = &Response{Status: 400, Message: "Failed to parse an ID"}
	ErrNoOptions = &Response{Status: 400, Message: "Request body doesn't contain options. Please add at least one option."}
)

//Message returns a Response struct with a custom string.
func Message(status int, message string) *Response {
	return &Response{status, message}
}

//Error returns a Response struct with an error
func Error(status int, err error) *Response {
	return &Response{status, err.Error()}
}

//Send responds to Fiber.
func (r *Response) Send(c *fiber.Ctx) error {
	return c.Status(r.Status).JSON(r)
}
