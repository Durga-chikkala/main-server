package query

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/main-server/errors"
	"log"
	"net/http"
	"strings"

	"github.com/main-server/models"
	"github.com/main-server/services"
)

type Handler struct {
	Queries services.QueriesEndPoints
}

func New(Q services.QueriesEndPoints) Handler {
	return Handler{Queries: Q}
}

func (h Handler) Create(c *gin.Context) {
	var input models.QueryInfo

	err := c.Bind(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{Code: "BAD REQUEST", Reason: "Invalid Body"})
		return
	}

	err = validateInput(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	input.Id = uuid.NewString()
	resp, err := h.Queries.Create(c, input)
	if err != nil {
		err := err.(errors.ErrorResponse)
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h Handler) Get(c *gin.Context) {
	resp, err := h.Queries.Get(c)
	if err != nil {
		err := err.(errors.ErrorResponse)
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h Handler) GetByQuestion(c *gin.Context) {
	question := strings.TrimSpace(c.Param("question"))
	log.Print("question:", question)

	if question == "" {
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{Code: "BAD REQUEST", Reason: fmt.Sprintf("Missing parameter %v", "question")})
		return
	}

	resp, err := h.Queries.GetByQuestion(c, question)
	if err != nil {
		err := err.(errors.ErrorResponse)
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h Handler) PatchByQuestion(c *gin.Context) {
	question := c.Param("question")
	log.Print("question:", question)

	if question == "" {
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{Code: "BAD REQUEST", Reason: fmt.Sprintf("Missing parameter %v", "question")})
		return
	}

	question = refactorQuestion(question)

	var input models.QueryInfo
	err := c.Bind(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{Code: "BAD REQUEST", Reason: fmt.Sprintf("Missing parameter %v", "question")})
		return
	}

	resp, err := h.Queries.PatchByQuestion(c, input.Count, question)
	if err != nil {
		switch v := err.(type) {
		case errors.ErrorResponse:
			c.JSON(v.StatusCode, v)
			return
		default:
			c.JSON(http.StatusBadRequest, err)
			return
		}
	}

	c.JSON(http.StatusOK, resp)
}

func (h Handler) GetFrequentQuestions(c *gin.Context) {
	resp, err := h.Queries.GetFrequentQuestions(c)
	if err != nil {
		err := err.(errors.ErrorResponse)
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func validateInput(input *models.QueryInfo) error {
	if input.Question == "" {
		return errors.ErrorResponse{Code: "BAD REQUEST", Reason: "Missing Field Question"}
	}

	return nil
}

func refactorQuestion(question string) string {

	question = strings.TrimSpace(question)
	return strings.ReplaceAll(question, " ", "-")
}
