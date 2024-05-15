package clients

import (
	"encoding/json"
	"fmt"
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/helpers"
)

// ApiError ...
type ApiError struct {
	Message string      `json:"message"`
	Data    interface{} `json:"errors"`
	Code    string      `json:"code"`
}

func mapError(body []byte) *ApiError {
	var apiError ApiError
	err := json.Unmarshal(body, &apiError)
	if err != nil {
		return nil
	}
	return &apiError
}

// Response ...
type Response struct {
	status        int
	error         error
	responseError *ApiError
}

func (r *Response) Failed() bool {
	if r.error != nil || r.responseError != nil || r.status >= 300 {
		return true
	}
	return false
}

func (r *Response) Status() int {
	return r.status
}

func (r *Response) Error() error {
	return r.error
}

func (r *Response) ResponseError() *ApiError {
	return r.responseError
}

// AsError returns either error (if exists) or ApiError as error
func (r *Response) AsError() error {
	if r.Error() != nil {
		return r.Error()
	}
	return fmt.Errorf(helpers.ToString(r.ResponseError()))
}
