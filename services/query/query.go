package query

import (
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/main-server/errors"
	"github.com/main-server/models"
	"log"
)

const (
	BASEURL = "http://localhost:8080/chatbot"
)

type services struct {
	*resty.Client
}

func New(client *resty.Client) services {
	return services{client}
}

func (s services) Create(c *gin.Context, input models.QueryInfo) (models.QueryInfo, error) {
	var clientErr errors.ErrorResponse

	_, err := s.Client.R().SetBody(input).SetHeader("Accept", "application/json").
		SetResult(&input).SetError(&clientErr).Post(BASEURL)
	if err != nil || (clientErr.StatusCode >= 400 && clientErr.StatusCode < 600) {
		if err != nil {
			return models.QueryInfo{}, err
		}
		return models.QueryInfo{}, clientErr
	}

	return input, nil
}

func (s services) Get(c *gin.Context) ([]models.QueryInfo, error) {
	var clientErr errors.ErrorResponse
	var queryInfo []models.QueryInfo

	_, err := s.Client.R().
		SetHeader("Accept", "application/json").SetResult(&queryInfo).SetError(&clientErr).Get(BASEURL)
	if err != nil || (clientErr.StatusCode >= 400 && clientErr.StatusCode < 600) {
		if err != nil {
			return []models.QueryInfo{}, err
		}
		return []models.QueryInfo{}, clientErr
	}

	return queryInfo, nil
}

func (s services) GetByQuestion(c *gin.Context, question string) (models.QueryInfo, error) {
	var queryInfo models.QueryInfo
	var clientErr errors.ErrorResponse

	_, err := s.Client.R().SetHeader("Accept", "application/json").SetResult(&queryInfo).SetError(&clientErr).Get(BASEURL + "/" + question)
	if err != nil || (clientErr.StatusCode >= 400 && clientErr.StatusCode < 600) {
		if err != nil {
			return models.QueryInfo{}, err
		}

		return models.QueryInfo{}, clientErr
	}

	return queryInfo, nil
}

func (s services) PatchByQuestion(c *gin.Context, count int64, question string) (models.QueryInfo, error) {
	var clientErr errors.ErrorResponse
	var queryInfo models.QueryInfo

	_, err := s.Client.R().SetHeader("Accept", "application/json").SetBody(models.QueryInfo{Count: count}).SetError(&clientErr).SetResult(&queryInfo).Patch(BASEURL + "/" + question)
	if err != nil || (clientErr.StatusCode >= 400 && clientErr.StatusCode < 600) {
		if err != nil {
			log.Println(err)
			return models.QueryInfo{}, err
		}
		return models.QueryInfo{}, clientErr
	}

	return queryInfo, nil
}

func (s services) GetFrequentQuestions(c *gin.Context) ([]models.QueryInfo, error) {
	var clientErr errors.ErrorResponse
	var queryInfo []models.QueryInfo

	_, err := s.Client.R().
		SetHeader("Accept", "application/json").SetResult(&queryInfo).SetError(&clientErr).Get(BASEURL + "/frequentQuestions")
	if err != nil || (clientErr.StatusCode >= 400 && clientErr.StatusCode < 600) {
		if err != nil {
			return []models.QueryInfo{}, err
		}
		return []models.QueryInfo{}, clientErr
	}

	return queryInfo, nil
}
