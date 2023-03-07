package user

import (
	"fmt"
	"github.com/main-server/services"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/main-server/errors"
	"github.com/main-server/models"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	User services.User
}

func New(s services.User) Handler {
	return Handler{User: s}
}

func (h Handler) Create(c *gin.Context) {
	var input models.UserInfo

	err := c.Bind(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{Code: "BAD REQUEST", Reason: "Invalid Body"})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	resp, err := h.User.Create(c, input)
	if err != nil {
		err := err.(errors.ErrorResponse)
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h Handler) Get(c *gin.Context) {
	userName := c.Query("email")
	password := c.Query("password")

	if userName == "" || password == "" {
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{Code: "BAD REQUEST", Reason: "Invalid Parameters"})
		return
	}

	resp, err := h.User.Get(c, userName, password)
	if err != nil {
		err := err.(errors.ErrorResponse)
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h Handler) GetByID(c *gin.Context) {
	ID := strings.TrimSpace(c.Param("id"))
	log.Print("ID:", ID)

	if ID == "" {
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{Code: "BAD REQUEST", Reason: fmt.Sprintf("Missing parameter %v", "id")})
		return
	}

	resp, err := h.User.GetByID(c, ID)
	if err != nil {
		err := err.(errors.ErrorResponse)
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h Handler) Logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{Name: "token", Value: "expired", Expires: time.Now()})
}

func (h Handler) PatchByID(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	log.Print("ID:", id)

	if id == "" {
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{Code: "BAD REQUEST", Reason: fmt.Sprintf("Missing parameter %v", "id")})
		return
	}

	var input models.UserInfo
	err := c.Bind(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{Code: "BAD REQUEST", Reason: fmt.Sprintf("Missing parameter %v", "question")})
		return
	}

	resp, err := h.User.PatchByID(c, id, input)
	if err != nil {
		err := err.(errors.ErrorResponse)
		c.JSON(err.StatusCode, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}
