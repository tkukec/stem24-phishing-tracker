package clients

import (
	"fmt"
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/runtimebag"
	assecoContext "git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/context"
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/helpers"
	"net/http"
	"strings"
)

const (
	TicketingUriEnvKey = "TICKETING_URI"
)

type Tickets struct {
	uri string
}

func NewTickets() *Tickets {
	return &Tickets{uri: helpers.UrlParse(runtimebag.GetEnvString(
		TicketingUriEnvKey,
		"http://ticketing:8080",
	))}
}

func (t *Tickets) SetUri(uri string) *Tickets {
	t.uri = helpers.UrlParse(uri)
	return t
}

func (t *Tickets) NewRequest(method string) *RequestBuilder {
	return NewRequestBuilder(t.uri, method)
}

type FindTicketResponse struct {
	Response
	result []*Ticket
}

func (f FindTicketResponse) Result() []*Ticket {
	return f.result
}

func (f FindTicketResponse) FirstOfTickets() *Ticket {
	if f.result != nil && len(f.result) > 0 {
		return f.result[0]
	}
	return nil
}

func (t *Tickets) FindById(ctx *assecoContext.RequestContext, id string, headers map[string]string) *FindTicketResponse {
	var result TicketCollection
	searchPayload := map[string]interface{}{
		"search": map[string]string{
			"id": fmt.Sprintf("=%s", id),
		},
		"relations": []string{
			"type",
			"stage",
		},
	}

	req := t.NewRequest(http.MethodPost).
		AddPath("api/search/ticket").
		SetBody(searchPayload).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)

	return &FindTicketResponse{
		Response: req.Response(),
		result:   result,
	}
}

func (t *Tickets) FindByShortId(ctx *assecoContext.RequestContext, shortId string, headers map[string]string) *FindTicketResponse {
	var result TicketCollection
	searchPayload := map[string]interface{}{
		"search": map[string]string{
			"short_id": fmt.Sprintf("=%s", shortId),
		},
		"relations": []string{
			"type",
			"stage",
		},
	}
	req := t.NewRequest(http.MethodPost).
		AddPath("api/search/ticket").
		SetBody(searchPayload).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)

	return &FindTicketResponse{
		Response: req.Response(),
		result:   result,
	}
}
func (t *Tickets) FindByContactId(ctx *assecoContext.RequestContext, contactId string, headers map[string]string) *FindTicketResponse {
	var result TicketCollection
	searchPayload := map[string]interface{}{
		"search": map[string]string{
			"contact_id": fmt.Sprintf("=%s", contactId),
		},
		"relations": []string{
			"type",
			"stage",
		},
	}
	req := t.NewRequest(http.MethodPost).
		AddPath("api/search/ticket").
		SetBody(searchPayload).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)

	return &FindTicketResponse{
		Response: req.Response(),
		result:   result,
	}
}

type CreateTicketResponse struct {
	Response
	result *Ticket
}

func (c CreateTicketResponse) Result() *Ticket {
	return c.result
}

func (t *Tickets) CreateTicket(ctx *assecoContext.RequestContext, payload map[string]interface{}, headers map[string]string) *CreateTicketResponse {
	var result *Ticket

	req := t.NewRequest(http.MethodPost).
		AddPath("api/tickets").
		SetBody(payload).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)

	return &CreateTicketResponse{
		Response: req.Response(),
		result:   result,
	}
}

type TicketCollection []*Ticket

func (c TicketCollection) First() *Ticket {
	if c != nil && len(c) > 0 {
		return c[0]
	}
	return nil
}

func (c TicketCollection) IsEmpty() bool {
	if c == nil || len(c) == 0 {
		return true
	}
	return false
}

type Ticket struct {
	Id                string              `json:"id"`
	ShortId           string              `json:"short_id"`
	ProcessInstanceId string              `json:"process_instance_id"`
	TenantId          string              `json:"tenant_id"`
	ContactId         string              `json:"contact_id"`
	Name              string              `json:"name"`
	TypeId            string              `json:"type_id"`
	StateId           string              `json:"state_id"`
	UrgencyId         string              `json:"urgency_id"`
	StageId           string              `json:"stage_id"`
	Sequence          int                 `json:"sequence"`
	CreatedAt         string              `json:"created_at"`
	CreatedBy         string              `json:"created_by"`
	CreatorType       string              `json:"creator_type"`
	UpdatedAt         string              `json:"updated_at"`
	UpdatedBy         string              `json:"updated_by"`
	UpdaterType       string              `json:"updater_type"`
	DeletedAt         string              `json:"deleted_at"`
	DeletedBy         string              `json:"deleted_by"`
	DeleterType       string              `json:"deleter_type"`
	SkillGroupId      string              `json:"skill_group_id"`
	OwnerId           string              `json:"owner_id"`
	SourceId          string              `json:"source_id"`
	Type              *TicketType         `json:"type"`
	Stage             *TicketStage        `json:"stage"`
	CustomFieldValues []*CustomFieldValue `json:"custom_field_values"`
	Urgency           *Urgency            `json:"urgency"`
	Source            *Source             `json:"source"`
}

type Source struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Label     string `json:"label"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Urgency struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Level       int    `json:"level"`
	CreatedAt   string `json:"created_at"`
	CreatedBy   string `json:"created_by"`
	CreatorType string `json:"creator_type"`
	UpdatedAt   string `json:"updated_at"`
	UpdatedBy   string `json:"updated_by"`
	UpdaterType string `json:"updater_type"`
	DeletedAt   string `json:"deleted_at"`
	DeletedBy   string `json:"deleted_by"`
	DeleterType string `json:"deleter_type"`
	Default     bool   `json:"default"`
}

type TicketType struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	TenantId      string `json:"tenant_id"`
	FormId        string `json:"form_id"`
	Enabled       bool   `json:"enabled"`
	ProcessTypeId string `json:"process_type_id"`
	CreatedAt     string `json:"created_at"`
	CreatedBy     string `json:"created_by"`
	CreatorType   string `json:"creator_type"`
	UpdatedAt     string `json:"updated_at"`
	UpdatedBy     string `json:"updated_by"`
	UpdaterType   string `json:"updater_type"`
	DeletedAt     string `json:"deleted_at"`
	DeletedBy     string `json:"deleted_by"`
	DeleterType   string `json:"deleter_type"`
}

type TicketStage struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	CreatedAt   string `json:"created_at"`
	CreatedBy   string `json:"created_by"`
	CreatorType string `json:"creator_type"`
	UpdatedAt   string `json:"updated_at"`
	UpdatedBy   string `json:"updated_by"`
	UpdaterType string `json:"updater_type"`
	DeletedAt   string `json:"deleted_at"`
	DeletedBy   string `json:"deleted_by"`
	DeleterType string `json:"deleter_type"`
	IsFinal     bool   `json:"is_final"`
}

type CustomFieldValue struct {
	Id            string       `json:"id"`
	CustomFieldId string       `json:"custom_field_id"`
	ModelType     string       `json:"model_type"`
	ModelId       string       `json:"model_id"`
	String        string       `json:"string"`
	Integer       int          `json:"integer"`
	Float         float64      `json:"float"`
	Text          string       `json:"text"`
	Boolean       bool         `json:"boolean"`
	Datetime      string       `json:"datetime"`
	Date          string       `json:"date"`
	Time          string       `json:"time"`
	CreatedAt     string       `json:"created_at"`
	CreatedBy     string       `json:"created_by"`
	CreatorType   string       `json:"creator_type"`
	UpdatedAt     string       `json:"updated_at"`
	UpdatedBy     string       `json:"updated_by"`
	UpdaterType   string       `json:"updater_type"`
	DeletedAt     string       `json:"deleted_at"`
	DeletedBy     string       `json:"deleted_by"`
	DeleterType   string       `json:"deleter_type"`
	Value         string       `json:"value"`
	CustomField   *CustomField `json:"custom_field"`
}

func (v *CustomFieldValue) FindValue() string {
	if v.CustomField == nil {
		return v.Value
	}
	if v.CustomField.Selectable == nil {
		return v.Value
	}

	if v.CustomField.Selectable.Multiselect {
		var values []string
		for _, i := range strings.Split(v.Value, ";") {
			values = append(values, v.CustomField.Selectable.ValueToLabel(i))
		}
		return strings.Join(values, ", ")
	}

	return v.CustomField.Selectable.ValueToLabel(v.Value)
}

func (v *CustomFieldSelectable) ValueToLabel(value string) string {
	for _, i := range v.Values {
		if value == i.Value {
			return i.Label
		}
	}
	return value
}

type CustomField struct {
	Id             string                 `json:"id"`
	SelectableType string                 `json:"selectable_type"`
	SelectableId   string                 `json:"selectable_id"`
	ValidationId   string                 `json:"validation_id"`
	Name           string                 `json:"name"`
	Label          string                 `json:"label"`
	Placeholder    string                 `json:"placeholder"`
	Model          string                 `json:"model"`
	Required       bool                   `json:"required"`
	Group          string                 `json:"group"`
	Order          int                    `json:"order"`
	CreatedAt      string                 `json:"created_at"`
	CreatedBy      string                 `json:"created_by"`
	CreatorType    string                 `json:"creator_type"`
	UpdatedAt      string                 `json:"updated_at"`
	UpdatedBy      string                 `json:"updated_by"`
	UpdaterType    string                 `json:"updater_type"`
	DeletedAt      string                 `json:"deleted_at"`
	DeletedBy      string                 `json:"deleted_by"`
	DeleterType    string                 `json:"deleter_type"`
	Hidden         bool                   `json:"hidden"`
	Selectable     *CustomFieldSelectable `json:"selectable"`
}

type CustomFieldSelectable struct {
	Id          string                         `json:"id"`
	PlainTypeId string                         `json:"plain_type_id"`
	Multiselect bool                           `json:"multiselect"`
	CreatedAt   string                         `json:"created_at"`
	CreatedBy   string                         `json:"created_by"`
	CreatorType string                         `json:"creator_type"`
	UpdatedAt   string                         `json:"updated_at"`
	UpdatedBy   string                         `json:"updated_by"`
	UpdaterType string                         `json:"updater_type"`
	DeletedAt   string                         `json:"deleted_at"`
	DeletedBy   string                         `json:"deleted_by"`
	DeleterType string                         `json:"deleter_type"`
	Name        string                         `json:"name"`
	Values      []*CustomFieldSelectableValues `json:"values"`
}

type CustomFieldSelectableValues struct {
	Id              string                          `json:"id"`
	SelectionTypeId string                          `json:"selection_type_id"`
	Label           string                          `json:"label"`
	Value           string                          `json:"value"`
	Preselect       bool                            `json:"preselect"`
	CreatedAt       string                          `json:"created_at"`
	CreatedBy       string                          `json:"created_by"`
	CreatorType     string                          `json:"creator_type"`
	UpdatedAt       string                          `json:"updated_at"`
	UpdatedBy       string                          `json:"updated_by"`
	UpdaterType     string                          `json:"updater_type"`
	DeletedAt       string                          `json:"deleted_at"`
	DeletedBy       string                          `json:"deleted_by"`
	DeleterType     string                          `json:"deleter_type"`
	Type            *CustomFieldSelectableValueType `json:"type"`
}

type CustomFieldSelectableValueType struct {
	Id                int    `json:"id"`
	Name              string `json:"name"`
	CreatedAt         string `json:"created_at"`
	CreatedBy         string `json:"created_by"`
	CreatorType       string `json:"creator_type"`
	UpdatedAt         string `json:"updated_at"`
	UpdatedBy         string `json:"updated_by"`
	UpdaterType       string `json:"updater_type"`
	DeletedAt         string `json:"deleted_at"`
	DeletedBy         string `json:"deleted_by"`
	DeleterType       string `json:"deleter_type"`
	LaravelThroughKey string `json:"laravel_through_key"`
}
