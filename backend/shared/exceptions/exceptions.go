package exceptions

import (
	"fmt"
	"github.com/andrezz-b/stem24-phishing-tracker/domain/models"
	"net/http"
)

type HttpErrors map[string][]string

type ApplicationException interface {
	Status() int
	Message() string
	Data() interface{}
	Code() string
	ToDto() *ApiError
	Error() error
}

type Exception struct {
	status  int
	message string
	data    interface{}
	code    string
}

func (a *Exception) Status() int {
	return a.status
}

func (a *Exception) Message() string {
	return a.message
}

func (a *Exception) Data() interface{} {
	return a.data
}

func (a *Exception) Code() string {
	return a.code
}

func (a *Exception) Error() error {
	return fmt.Errorf("message:%s status:%d code:%s", a.Message(), a.Status(), a.Code())
}

func (a *Exception) ToDto() *ApiError {
	return NewApiError(a.Message(), a.Code(), a.Data())
}

func New(status int, message, code string, data interface{}) ApplicationException {
	if data == nil {
		data = make(map[string]interface{})
	}
	return &Exception{status: status, message: message, data: data, code: code}
}

func NewApiError(msg, code string, data interface{}) *ApiError {
	return &ApiError{
		Message: msg,
		Data:    data,
		Code:    code,
	}
}

type ApiError struct {
	Message string      `json:"message"`
	Data    interface{} `json:"errors"`
	Code    string      `json:"code"`
}

const (
	ResourceNotFoundCode    = "agent_management_404"
	InternalCode            = "agent_management_500"
	BadRequestCode          = "agent_management_400"
	UnprocessableEntityCode = "agent_management_422"
	ConflictCode            = "agent_management_409"
	ForbiddenCode           = "agent_management_403"

	InvalidInputFormat       = "agent_management_111"
	FilenameMissingExtension = "agent_management_112"
	MissingConfiguration     = "agent_management_113"
	MediaServerIssue         = "agent_management_114"
	AgentBusyCode            = "agent_management_115"

	FailedPersistingCode = "agent_management_116"
	FailedDeletingCode   = "agent_management_117"
	FailedUpdatingCode   = "agent_management_118"
	FailedQueryingCode   = "agent_management_119"

	MissingExtensionCode          = "agent_management_120"
	FailedGeneratingExtensionCode = "agent_management_121"
	MissingAgentSkillGroupCode    = "agent_management_122"
	FailedAddingTimerCode         = "agent_management_123"
	InvalidActivityTypeCode       = "agent_management_124"
	CanNotTransferActivityCode    = "agent_management_125"
	CanNotTransitionStatusCode    = "agent_management_126"
	InvalidExtensionCode          = "agent_management_127"
	AgentNotLoggedIn              = "agent_management_128"
	WrongRecordingState           = "agent_management_129"
	SessionNotFound               = "agent_management_130"

	recordNotFound   = "record not found"
	itemNotAvailable = "item not available"
	agentBusy        = "agent busy"
)

func NotFound(data interface{}, code string) ApplicationException {
	if code == "" {
		code = ResourceNotFoundCode
	}
	return New(http.StatusNotFound, recordNotFound, code, data)
}

func Internal(message string, data interface{}, code string) ApplicationException {
	if code == "" {
		code = InternalCode
	}
	return New(http.StatusInternalServerError, message, code, data)
}

func BadRequest(data interface{}, code string) ApplicationException {
	if code == "" {
		code = BadRequestCode
	}
	return New(http.StatusBadRequest, "bad request", code, data)
}

func Forbidden(data interface{}, code string) ApplicationException {
	if code == "" {
		code = ForbiddenCode
	}
	return New(http.StatusForbidden, "forbidden", code, data)
}

func UnprocessableEntity(data interface{}, code string) ApplicationException {
	if code == "" {
		code = UnprocessableEntityCode
	}
	return New(http.StatusUnprocessableEntity, "422 unprocessable entity", code, data)
}

func Conflict(message string, data interface{}, code string) ApplicationException {
	if code == "" {
		code = ConflictCode
	}
	return New(http.StatusConflict, message, code, data)
}

func TenantNotFound(errors ...error) ApplicationException {
	if len(errors) == 0 {
		errors = []error{fmt.Errorf(recordNotFound)}
	}

	data := map[string][]string{
		models.TenantModelName: make([]string, len(errors)),
	}

	for i, err := range errors {
		data[models.TenantModelName][i] = err.Error()
	}

	return NotFound(data, ResourceNotFoundCode)
}

func FailedPersisting(model string, err error) ApplicationException {
	return Internal("failed persisting record", map[string]string{model: err.Error()}, FailedPersistingCode)
}

func FailedDeleting(model string, err error) ApplicationException {
	return Internal("failed deleting record", map[string]string{model: err.Error()}, FailedDeletingCode)
}

func FailedUpdating(model string, err error) ApplicationException {
	return Internal("failed updating record", map[string]string{model: err.Error()}, FailedUpdatingCode)
}

func FailedQuerying(model string, err error) ApplicationException {
	return Internal("failed querying database", map[string]string{model: err.Error()}, FailedQueryingCode)
}

func FailedFetchingServiceToken(errors ...error) ApplicationException {
	if len(errors) == 0 {
		errors = []error{fmt.Errorf("missing activity configuration")}
	}

	data := map[string][]string{
		"iam": make([]string, len(errors)),
	}

	for i, err := range errors {
		data["iam"][i] = err.Error()
	}

	return Internal("failed fetching service token", data, MissingConfiguration)
}
