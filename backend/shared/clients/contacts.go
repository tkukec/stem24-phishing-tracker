package clients

import (
	"fmt"
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/runtimebag"
	assecoContext "git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/context"
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/helpers"
	"net/http"
)

const (
	ContactsUriEnvKey = "CONTACTS_URI"

	numberMobileType   = "942d9fcd-9d78-48f2-b9f9-61ecc79fe108"
	numberLandLineType = "942d9fcd-9d72-412b-9aa9-9e1ce15bc50e"
)

type Contacts struct {
	uri string
}

func NewContacts() *Contacts {
	return &Contacts{uri: helpers.UrlParse(runtimebag.GetEnvString(
		ContactsUriEnvKey,
		"http://contacts:8080",
	))}
}

func (c *Contacts) SetUri(uri string) *Contacts {
	c.uri = helpers.UrlParse(uri)
	return c
}

func (c *Contacts) NewRequest(method string) *RequestBuilder {
	return NewRequestBuilder(c.uri, method)
}

func (c *Contacts) NewSearchRequest() *RequestBuilder {
	return c.NewRequest(http.MethodPost).AddPath("api/search/contact")
}

type ContactByNumberResponse struct {
	Response
	result ContactCollection
}

func (c ContactByNumberResponse) Result() ContactCollection {
	return c.result
}

func (c ContactByNumberResponse) FirstOfContacts() *Contact {
	if c.result == nil {
		return nil
	}
	return c.result.First()
}

func (c *Contacts) ContactByNumber(ctx *assecoContext.RequestContext, number string, headers map[string]string, extraQueries ...map[string]string) *ContactByNumberResponse {
	var result ContactCollection
	search := map[string]string{
		"media.media_type_id": fmt.Sprintf("=%s||=%s", numberMobileType, numberLandLineType),
		"media.value":         fmt.Sprintf("=%s", number),
	}

	for _, query := range extraQueries {
		for queryKey, queryValue := range query {
			search[queryKey] = queryValue
		}
	}

	req := c.NewSearchRequest().
		SetBody(map[string]interface{}{
			"search":    search,
			"relations": []string{"media", "addresses"},
			"order_by": map[string]string{
				"created_at": "asc",
			},
		}).
		Bind(&result).
		Print().
		Run(ctx, headers)

	return &ContactByNumberResponse{
		Response: req.Response(),
		result:   result,
	}
}

type ContactByIdResponse struct {
	Response
	result []*Contact
}

func (c ContactByIdResponse) Result() []*Contact {
	return c.result
}

func (c ContactByIdResponse) FirstOfContacts() *Contact {
	if c.result != nil && len(c.result) > 0 {
		return c.result[0]
	}
	return nil
}

func (c *Contacts) ContactById(ctx *assecoContext.RequestContext, id string, headers map[string]string) *ContactByIdResponse {
	var result ContactCollection
	req := c.NewSearchRequest().
		SetBody(map[string]interface{}{
			"search": map[string]string{
				"id": fmt.Sprintf("=%s", id),
			},
			"relations": []string{"media", "addresses"},
			"order_by": map[string]string{
				"created_at": "asc",
			},
		}).
		Bind(&result).
		Print().
		Run(ctx, headers)

	return &ContactByIdResponse{
		Response: req.Response(),
		result:   result,
	}
}

type ContactCollection []*Contact

func (c ContactCollection) First() *Contact {
	if c != nil && len(c) > 0 {
		return c[0]
	}
	return nil
}

func (c ContactCollection) IsEmpty() bool {
	if c == nil || len(c) == 0 {
		return true
	}
	return false
}

type Contact struct {
	Id               string            `json:"id"`
	TenantId         string            `json:"tenant_id"`
	Birthday         string            `json:"birthday"`
	Image            string            `json:"image"`
	Prefix           string            `json:"prefix"`
	FirstName        string            `json:"first_name"`
	MiddleName       string            `json:"middle_name"`
	LastName         string            `json:"last_name"`
	Suffix           string            `json:"suffix"`
	Description      string            `json:"description"`
	Emoji            string            `json:"emoji"`
	Timezone         string            `json:"timezone"`
	IsOrganization   bool              `json:"is_organization"`
	OrganizationName string            `json:"organization_name"`
	OrganizationUnit string            `json:"organization_unit"`
	PreferredAgentId string            `json:"preferred_agent_id"`
	AccountId        string            `json:"account_id"`
	TitleId          string            `json:"title_id"`
	Sequence         int               `json:"sequence"`
	Addresses        []*ContactAddress `json:"addresses"`
	Media            []*Media          `json:"media"`

	CreatedAt   string `json:"created_at"`
	CreatedBy   string `json:"created_by"`
	CreatorType string `json:"creator_type"`
	UpdatedAt   string `json:"updated_at"`
	UpdatedBy   string `json:"updated_by"`
	UpdaterType string `json:"updater_type"`
	DeletedAt   string `json:"deleted_at"`
	DeletedBy   string `json:"deleted_by"`
	DeleterType string `json:"deleter_type"`
	DisplayName string `json:"display_name"`
}

type ContactAddress struct {
	Id                string `json:"id"`
	Country           string `json:"country"`
	City              string `json:"city"`
	ZipCode           string `json:"zip_code"`
	Region            string `json:"region"`
	Street            string `json:"street"`
	StreetNumber      string `json:"street_number"`
	Longitude         string `json:"longitude"`
	Latitude          string `json:"latitude"`
	Primary           bool   `json:"primary"`
	ContactId         string `json:"contact_id"`
	AddressCategoryId string `json:"address_category_id"`
	CreatedAt         string `json:"created_at"`
	CreatedBy         string `json:"created_by"`
	CreatorType       string `json:"creator_type"`
	UpdatedAt         string `json:"updated_at"`
	UpdatedBy         string `json:"updated_by"`
	UpdaterType       string `json:"updater_type"`
	DeletedAt         string `json:"deleted_at"`
	DeletedBy         string `json:"deleted_by"`
	DeleterType       string `json:"deleter_type"`
}

type Media struct {
	Id          string   `json:"id"`
	Value       string   `json:"value"`
	Description string   `json:"description"`
	Preferred   bool     `json:"preferred"`
	Contact     *Contact `json:"contact"`
}
