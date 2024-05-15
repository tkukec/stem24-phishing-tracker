package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/constants"
	assecoContext "git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/context"
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/helpers"
	"io"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func NewRequestBuilder(base, method string) *RequestBuilder {
	return &RequestBuilder{
		client:  &http.Client{},
		base:    base,
		method:  strings.ToUpper(method),
		path:    make([]string, 0),
		queries: make(map[string][]string),
		headers: make(map[string]string),
	}
}

type RequestBuilder struct {
	client *http.Client
	ctx    *assecoContext.RequestContext
	error  error

	request *http.Request

	response       *http.Response
	responseBody   []byte
	bindResponseTo interface{}

	base    string
	method  string
	path    []string
	queries map[string][]string
	headers map[string]string
	body    any
	build   string
}

func (u *RequestBuilder) SetDefaultHeaders() {
	if u.request.Header.Get(constants.ContentType) == "" {
		u.request.Header.Set(constants.ContentType, constants.JsonContentType)
	}
	if u.request.Header.Get(constants.XCorrelationID) == "" && u.ctx != nil {
		u.request.Header.Set(constants.XCorrelationID, u.ctx.XCorrelationID())
	}
	for key, header := range u.headers {
		u.request.Header.Add(key, header)
	}
}

func (u *RequestBuilder) SetHeaders(headers map[string]string) *RequestBuilder {
	u.headers = headers
	return u
}

func (u *RequestBuilder) SetBody(body any) *RequestBuilder {
	u.body = body
	return u
}

func (u *RequestBuilder) AddPath(path string) *RequestBuilder {
	u.path = append(u.path, path)
	return u
}

func (u *RequestBuilder) AddMultipleQueryValues(key string, value []string) *RequestBuilder {
	if len(value) == 0 {
		return u
	}
	_, ok := u.queries[key]
	if !ok {
		u.queries[key] = value
		return u
	}
	u.queries[key] = append(u.queries[key], value...)
	return u
}

func (u *RequestBuilder) AddQuery(key, value string) *RequestBuilder {
	_, ok := u.queries[key]
	if !ok {
		u.queries[key] = []string{value}
		return u
	}
	u.queries[key] = append(u.queries[key], value)
	return u
}

func (u *RequestBuilder) AddBoolQuery(key string, value bool) *RequestBuilder {
	return u.AddQuery(key, strconv.FormatBool(value))
}

func (u *RequestBuilder) Url() string {
	url := u.base
	for _, path := range u.path {
		url += "/" + strings.TrimRight(strings.TrimLeft(path, "/"), "/")
	}
	return u.buildQueries(url)
}

func (u *RequestBuilder) StatusCode() int {
	if u.response == nil {
		return 0
	}
	return u.response.StatusCode
}

func (u *RequestBuilder) Status() string {
	if u.response == nil {
		return ""
	}
	return u.response.Status
}

func (u *RequestBuilder) Bind(value interface{}) *RequestBuilder {
	u.bindResponseTo = value
	return u
}

func (u *RequestBuilder) Run(ctx *assecoContext.RequestContext, headers ...map[string]string) *RequestBuilder {
	u.ctx = ctx
	if u.error != nil {
		return u
	}
	if u.request == nil {
		_, err := u.Build()
		if err != nil {
			u.error = err
			return u
		}
	}

	if u.request == nil {
		u.error = errors.New("missing request object")
		return u
	}
	for _, header := range headers {
		for i, v := range header {
			u.request.Header.Add(i, v)
		}
	}
	u.SetDefaultHeaders()
	response, err := u.client.Do(u.clientRequest())
	if err != nil {
		u.error = err
		return u
	}
	u.responseBody, err = io.ReadAll(response.Body)
	if err != nil {
		u.error = err
		return u
	}
	defer response.Body.Close()
	u.response = response

	if u.response.StatusCode >= 300 {
		u.error = fmt.Errorf("status code not ok, status: %s (code: %d)", u.response.Status, u.response.StatusCode)
		return u
	}
	u.error = json.Unmarshal(u.responseBody, &u.bindResponseTo)
	return u
}

func (u *RequestBuilder) clientRequest() *http.Request {
	if u.ctx != nil {
		return u.request.WithContext(u.ctx.Context())
	}
	return u.request
}

func (u *RequestBuilder) Error() error {
	return u.error
}

func (u *RequestBuilder) Print() *RequestBuilder {
	u.PrintUrl()
	u.PrintBody()
	u.PrintHeaders()
	u.PrintBind()
	return u
}

func (u *RequestBuilder) PrintUrl() *RequestBuilder {
	log.Println("Request URL: ", u.Url())
	return u
}

func (u *RequestBuilder) PrintBody() *RequestBuilder {
	log.Println("Request body: ", helpers.ToPrettyPrint(u.body))
	return u
}

func (u *RequestBuilder) PrintBind() *RequestBuilder {
	log.Println("Bind: ", reflect.TypeOf(u.bindResponseTo))
	return u
}

func (u *RequestBuilder) PrintHeaders() *RequestBuilder {
	log.Println("Headers: ", helpers.ToPrettyPrint(u.headers))
	return u
}

func (u *RequestBuilder) buildQueries(url string) string {
	if u.queries == nil || len(u.queries) == 0 {
		return url
	}
	url += "?"
	queries := make([]string, 0)
	for key, values := range u.queries {
		for _, value := range values {
			queries = append(queries, fmt.Sprintf("%s=%s", key, value))
		}
	}
	return url + strings.Join(queries, "&")
}

func (u *RequestBuilder) Build() (*http.Request, error) {
	request, err := http.NewRequest(u.method, u.Url(), bytes.NewBuffer(helpers.ToByte(u.body)))
	u.request = request
	u.error = err
	return request, err
}

func (u *RequestBuilder) Response() Response {
	if u.StatusCode() < 200 || u.StatusCode() >= 400 {
		return Response{
			status: u.StatusCode(),
			error:  u.Error(),
			responseError: &ApiError{
				Message: string(u.responseBody),
				Code:    u.Status(),
			},
		}
	}
	return Response{
		status: u.StatusCode(),
		error:  u.Error(),
	}
}
