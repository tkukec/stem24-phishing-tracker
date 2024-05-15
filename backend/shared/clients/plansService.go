package clients

import (
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/eventMessages/telephone"
	assecoContext "git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/context"
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/helpers"
	"net/http"
)

type PlansService struct {
}

func NewPlansService() *PlansService {
	return &PlansService{}
}

func (s *PlansService) NewRequest(uri string, method string) *RequestBuilder {
	return NewRequestBuilder(uri, method)
}

type FetchPlansResponse struct {
	Response
	//Plan struct is same for video chat and telephone service
	result []*telephone.Plan
}

func (s FetchPlansResponse) Result() []*telephone.Plan {
	return s.result
}

func (s *PlansService) FetchPlansFromServiceUri(ctx *assecoContext.RequestContext, serviceUri string, headers map[string]string) *FetchPlansResponse {
	var result []*telephone.Plan

	searchPayload := map[string]interface{}{
		"relations": []string{"rules", "values"},
	}

	req := s.NewRequest(helpers.UrlParse(serviceUri), http.MethodPost).
		AddPath(searchEndpoint).
		AddPath(plans).
		SetBody(searchPayload).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)

	return &FetchPlansResponse{
		Response: req.Response(),
		result:   result,
	}
}
