package user

import (
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/main-server/errors"
	"github.com/main-server/models"
	"github.com/main-server/services/auth"
	"golang.org/x/exp/slices"
	"log"
	"net/http"
	"regexp"
	"time"
)

const (
	BASEURL = "http://localhost:8080/user"
	EMAIL   = `^[a-zA-Z0-9+_.-]+@[a-zA-Z0-9.-]+$`
	PHONE   = `^[6-9]\d{9}$`
)

type service struct {
	*resty.Client
}

func New(client *resty.Client) service {
	return service{client}
}

func (s service) Create(c *gin.Context, user models.UserInfo) (models.UserInfo, error) {
	var clientErr errors.ErrorResponse
	err := validateUser(user)
	if err != nil {
		return models.UserInfo{}, err
	}

	user.Password = encrypt(user.Password)
	user.ID = uuid.NewString()

	_, err = s.Client.R().SetBody(user).SetHeader("Accept", "application/json").
		SetResult(&user).SetError(&clientErr).Post(BASEURL + "/signup")
	if err != nil || (clientErr.StatusCode >= 400 && clientErr.StatusCode < 600) {
		if err != nil {
			return models.UserInfo{}, err
		}
		return models.UserInfo{}, clientErr
	}

	return user, nil
}

func (s service) GetByID(c *gin.Context, ID string) (models.UserInfo, error) {
	var userInfo models.UserInfo
	var clientErr errors.ErrorResponse

	_, err := s.Client.R().SetHeader("Accept", "application/json").SetResult(&userInfo).SetError(&clientErr).Get(BASEURL + "/" + ID)
	if err != nil || (clientErr.StatusCode >= 400 && clientErr.StatusCode < 600) {
		if err != nil {
			return models.UserInfo{}, err
		}

		return models.UserInfo{}, clientErr
	}

	return userInfo, nil
}

func (s service) Get(c *gin.Context, email string, password string) (models.UserInfo, error) {
	password = encrypt(password)
	var clientErr errors.ErrorResponse
	var userInfo models.UserInfo

	_, err := s.Client.R().
		SetQueryParams(map[string]string{
			"email":    email,
			"password": password,
		}).SetHeader("Accept", "application/json").SetResult(&userInfo).SetError(&clientErr).Get(BASEURL + "/login")
	if err != nil || (clientErr.StatusCode >= 400 && clientErr.StatusCode < 600) {
		if err != nil {
			return models.UserInfo{}, err
		}
		return models.UserInfo{}, clientErr
	}

	tkn, err := auth.GenerateJWT(userInfo.Email, userInfo.Password)
	if err != nil {
		return models.UserInfo{}, errors.ErrorResponse{StatusCode: http.StatusInternalServerError, Code: "Internal Server Error", Reason: "Problem in Login"}
	}

	http.SetCookie(c.Writer, &http.Cookie{Name: "token", Value: tkn, Expires: time.Now().Add(1 * time.Hour)})

	return userInfo, nil
}

func (s service) PatchByID(c *gin.Context, ID string, user models.UserInfo) (models.UserInfo, error) {
	user.Password = encrypt(user.Password)
	var clientErr errors.ErrorResponse
	var userInfo models.UserInfo

	_, err := s.Client.R().SetHeader("Accept", "application/json").SetBody(user).SetError(&clientErr).SetResult(&userInfo).Patch(BASEURL + "/" + ID)
	if err != nil || (clientErr.StatusCode >= 400 && clientErr.StatusCode < 600) {
		if err != nil {
			log.Println(err)
			return models.UserInfo{}, err
		}
		return models.UserInfo{}, clientErr
	}

	return userInfo, nil
}

func encrypt(password string) string {

	return password
}

func validateUser(user models.UserInfo) error {
	switch {
	case user.FirstName == "":
		return errors.ErrorResponse{StatusCode: http.StatusBadRequest, Code: "BAD REQUEST", Reason: "Missing field firstName"}
	case user.LastName == "":
		return errors.ErrorResponse{StatusCode: http.StatusBadRequest, Code: "BAD REQUEST", Reason: "Missing field lastName"}
	case !isEmail(user.Email):
		return errors.ErrorResponse{StatusCode: http.StatusBadRequest, Code: "BAD REQUEST", Reason: "Invalid Email"}
	case !isPhone(user.Phone):
		return errors.ErrorResponse{StatusCode: http.StatusBadRequest, Code: "BAD REQUEST", Reason: "Invalid Phone"}
	case !validGender(user.Gender):
		return errors.ErrorResponse{StatusCode: http.StatusBadRequest, Code: "BAD REQUEST", Reason: "Invalid Gender"}
	}
	return nil
}

func isEmail(email string) bool {
	regex := regexp.MustCompile(EMAIL)
	if regex.MatchString(email) {
		return true
	}

	return false
}

func isPhone(phone string) bool {
	regex := regexp.MustCompile(PHONE)
	if regex.MatchString(phone) {
		return true
	}

	return false
}

func validGender(gender string) bool {
	genders := []string{"MALE", "FEMALE"}
	if slices.Contains(genders, gender) {
		return true
	}
	return false
}
