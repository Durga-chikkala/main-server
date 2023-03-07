package services

import (
	"github.com/gin-gonic/gin"
	"github.com/main-server/models"
)

type User interface {
	Create(c *gin.Context, data models.UserInfo) (models.UserInfo, error)
	GetByID(c *gin.Context, ID string) (models.UserInfo, error)
	Get(c *gin.Context, username string, password string) (models.UserInfo, error)
	PatchByID(c *gin.Context, ID string, user models.UserInfo) (models.UserInfo, error)
}

type QueriesEndPoints interface {
	Create(c *gin.Context, data models.QueryInfo) (models.QueryInfo, error)
	GetByQuestion(c *gin.Context, question string) (models.QueryInfo, error)
	Get(c *gin.Context) ([]models.QueryInfo, error)
	PatchByQuestion(c *gin.Context, count int64, question string) (models.QueryInfo, error)
	GetFrequentQuestions(c *gin.Context) ([]models.QueryInfo, error)
}
