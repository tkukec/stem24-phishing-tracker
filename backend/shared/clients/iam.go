package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/runtimebag"
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/constants"
	assecoContext "git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/context"
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/helpers"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	GrantTypeCredentials    = "client_credentials"
	GrantTypeRefresherToken = "refresh_token"
)

func NewIam(url, clientId, clientSecret, realm string) *Iam {
	return &Iam{
		url:          helpers.UrlParse(url),
		clientId:     clientId,
		clientSecret: clientSecret,
		realm:        realm,
	}
}

func NewEnvIam() *Iam {
	return NewIam(
		helpers.UrlParse(runtimebag.GetEnvString(constants.IamUri, "http://iam:8080")),
		runtimebag.GetEnvString(constants.ClientId, ""),
		runtimebag.GetEnvString(constants.ClientSecret, ""),
		runtimebag.GetEnvString(constants.IamRealm, "live"),
	)
}

type Iam struct {
	url          string
	clientId     string
	clientSecret string
	realm        string
}

type ServiceLoginResponse struct {
	Response
	result *IamResponse
}

func (s ServiceLoginResponse) Result() *IamResponse {
	return s.result
}

type IamResponse struct {
	AccessToken      string `json:"access_token" binding:"required"`
	ExpiresIn        int32  `json:"expires_in"`
	RefreshExpiresIn int32  `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int32  `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

func (p *Iam) ServiceLogin(ctx *assecoContext.RequestContext, timeout time.Duration, headers map[string]string, params map[string]string) *ServiceLoginResponse {
	data := url.Values{}
	data.Set("client_id", p.clientId)
	data.Set("client_secret", p.clientSecret)
	data.Set("grant_type", GrantTypeCredentials)
	if params[constants.RefreshToken] != "" {
		data.Set("refresh_token", params[constants.RefreshToken])
	}

	urlLogin := fmt.Sprintf("/auth/realms/%s/protocol/openid-connect/token", p.realm)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", p.url, urlLogin), strings.NewReader(data.Encode()))
	if err != nil {
		return &ServiceLoginResponse{
			Response: Response{
				error: fmt.Errorf("error creating request, error %s", err.Error()),
			},
		}
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add(constants.XCorrelationID, ctx.XCorrelationID())
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	rCtx, cancel := context.WithTimeout(ctx.Context(), timeout)
	defer cancel()
	resp, err := client.Do(req.WithContext(rCtx))

	if err != nil {
		return &ServiceLoginResponse{
			Response: Response{
				error: fmt.Errorf("error executing request, error %s", err.Error()),
			},
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &ServiceLoginResponse{
			Response: Response{
				status: resp.StatusCode,
				error:  fmt.Errorf("error reading body, error %s", err.Error()),
			},
		}
	}
	if resp.StatusCode >= 300 {
		return &ServiceLoginResponse{
			Response: Response{
				status:        resp.StatusCode,
				responseError: mapError(body),
			},
		}
	}

	var result *IamResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return &ServiceLoginResponse{
			Response: Response{
				status: resp.StatusCode,
				error:  fmt.Errorf("error un-marshaling body, error %s body %s", err.Error(), string(body)),
			},
		}
	}
	_ = resp.Body.Close()

	return &ServiceLoginResponse{
		Response: Response{
			status: resp.StatusCode,
		},
		result: result,
	}
}

func (p *Iam) ServiceToken(ctx *assecoContext.RequestContext, timeout time.Duration, headers map[string]string, params map[string]string) (string, error) {
	resp := p.ServiceLogin(ctx, timeout, headers, params)
	if resp.error != nil {
		return "", resp.error
	}

	return resp.Result().AccessToken, nil
}
