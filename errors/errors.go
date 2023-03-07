package errors

import "fmt"

type ErrorResponse struct {
	StatusCode int    `json:"statusCode,omitempty"`
	Code       string `json:"code,omitempty"`
	Reason     string `json:"reason,omitempty"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("{\n code: %v \n reason: %v\n}", e.Code, e.Reason)

}
