package clients

import (
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/runtimebag"
	assecoContext "git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/context"
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/helpers"
	"net/http"
)

const (
	DirectoryUriEnvKey = "DIRECTORY_URI"
)

type DirectoryService struct {
	uri string
}

func NewDirectoryService() *DirectoryService {
	return &DirectoryService{
		uri: helpers.UrlParse(runtimebag.GetEnvString(
			DirectoryUriEnvKey,
			"http://directory:8080",
		)),
	}
}

func (d *DirectoryService) SetUri(uri string) *DirectoryService {
	d.uri = helpers.UrlParse(uri)
	return d
}

func (d *DirectoryService) NewRequest(method string) *RequestBuilder {
	return NewRequestBuilder(d.uri, method)
}

type FindByUsernameResponse struct {
	Response
	result *DirectoryServiceResponse
}

func (r *FindByUsernameResponse) Result() *DirectoryServiceResponse {
	return r.result
}

type DirectoryServiceResponse struct {
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
	TotalCount int               `json:"total_count"`
	Items      []*DirectoryAgent `json:"items"`
}

type FindByIdResponse struct {
	Response
	result *DirectoryAgent
}

func (r *FindByIdResponse) Result() *DirectoryAgent {
	return r.result
}

func (d *DirectoryServiceResponse) First() *DirectoryAgent {
	if len(d.Items) == 0 {
		return nil
	}
	return d.Items[0]
}

func (d *DirectoryService) FindByUsername(ctx *assecoContext.RequestContext, username string, headers map[string]string) *FindByUsernameResponse {
	var result *DirectoryServiceResponse
	req := d.NewRequest(http.MethodGet).AddPath("v1/directory/default/users").
		AddQuery("casing", "snake").
		AddQuery("username", username).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)

	return &FindByUsernameResponse{
		Response: req.Response(),
		result:   result,
	}
}

func (d *DirectoryService) FindById(ctx *assecoContext.RequestContext, id string, headers map[string]string) *FindByIdResponse {
	var result *DirectoryAgent
	req := d.NewRequest(http.MethodGet).AddPath("v1/directory/default/users/"+id).
		AddQuery("casing", "snake").
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)

	return &FindByIdResponse{
		Response: req.Response(),
		result:   result,
	}
}

type DirectoryAgent struct {
	Id                  string        `json:"id"`
	FirstName           string        `json:"first_name"`
	LastName            string        `json:"last_name"`
	DisplayName         string        `json:"display_name"`
	Username            string        `json:"username"`
	UserType            string        `json:"user_type"`
	AccountEnabled      bool          `json:"account_enabled"`
	Email               string        `json:"email"`
	EmailVerified       bool          `json:"email_verified"`
	PhoneNumberVerified bool          `json:"phone_number_verified"`
	HintPictureUrl      string        `json:"hint_picture_url"`
	ProfilePictureUrl   string        `json:"profile_picture_url"`
	LinkedLogins        []interface{} `json:"linked_logins"`
	AccountLocked       bool          `json:"account_locked"`
	MultifactorEnabled  bool          `json:"multifactor_enabled"`
	RequiredActions     interface{}   `json:"required_actions"`
	Attributes          interface{}   `json:"attributes"`
	Groups              interface{}   `json:"groups"`
	Roles               interface{}   `json:"roles"`
	OuCode              string        `json:"ou_code"`
}
