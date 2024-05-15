package clients

import (
	assecoContext "git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/context"
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/helpers"
	"net/http"
)

const (
	searchEndpoint = "/api/search"
)

type SkillGroupService struct {
}

func NewSkillGroupService() *SkillGroupService {
	return &SkillGroupService{}
}

func (s *SkillGroupService) NewRequest(uri string, method string) *RequestBuilder {
	return NewRequestBuilder(uri, method)
}

type FetchSkillGroupsResponse struct {
	Response
	result []*SkillGroup
}

func (s FetchSkillGroupsResponse) Result() []*SkillGroup {
	return s.result
}

func (s *SkillGroupService) FetchSkillGroupsFromServiceUri(ctx *assecoContext.RequestContext, serviceUri string, headers map[string]string) *FetchSkillGroupsResponse {
	var result []*SkillGroup

	searchPayload := map[string]interface{}{
		"relations": []string{"members"},
	}

	req := s.NewRequest(helpers.UrlParse(serviceUri), http.MethodPost).
		AddPath(searchEndpoint).
		AddPath(skillGroups).
		SetBody(searchPayload).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)

	return &FetchSkillGroupsResponse{
		Response: req.Response(),
		result:   result,
	}
}

type SkillGroup struct {
	ID              string              `json:"id"`
	Name            string              `json:"name"`
	Description     string              `json:"description"`
	PromptTimeout   int64               `json:"prompt_timeout"`
	TelephoneNumber string              `json:"telephone_number"`
	From            string              `json:"from"`
	Disabled        bool                `json:"disabled"`
	Members         []*SkillGroupMember `json:"members"`
	UseAmd          bool                `json:"use_amd"`
}

type SkillGroupMember struct {
	SkillGroupId string `json:"skill_group_id"`
	SkillLevel   int    `json:"level"`
	Primary      bool   `json:"primary"`
	UserUuid     string `json:"user_uuid"`
}
