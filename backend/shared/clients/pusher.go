package clients

import (
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/runtimebag"
	assecoContext "git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/context"
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/helpers"
	"net/http"
)

const (
	PusherUriEnvKey = "PUSHER_URI"
)

func NewPusher() *Pusher {
	return &Pusher{
		uri: helpers.UrlParse(runtimebag.GetEnvString(
			PusherUriEnvKey,
			"http://pusher:8080",
		)),
	}
}

type Pusher struct {
	uri string
}

func (p *Pusher) SetUri(uri string) {
	p.uri = helpers.UrlParse(uri)
}

func (p *Pusher) NewRequest(method string) *RequestBuilder {
	return NewRequestBuilder(p.uri, method)
}

type PusherMessage struct {
	Namespace string      `json:"namespace"`
	Room      string      `json:"room"`
	Event     string      `json:"event"`
	Payload   interface{} `json:"payload"`
}

type SendNotificationToRoomResponse struct {
	Response
	result *PusherResponse
}

func (f SendNotificationToRoomResponse) Result() *PusherResponse {
	return f.result
}

func (p *Pusher) SendNotificationToRoom(ctx *assecoContext.RequestContext, msg *PusherMessage, headers map[string]string) *SendNotificationToRoomResponse {
	var result *PusherResponse
	req := p.NewRequest(http.MethodPost).
		AddPath("api/rooms").
		SetBody(msg).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)

	return &SendNotificationToRoomResponse{
		Response: req.Response(),
		result:   result,
	}
}

type PusherUserMessage struct {
	Namespace string      `json:"namespace"`
	UserID    string      `json:"user_id"`
	Event     string      `json:"event"`
	Payload   interface{} `json:"payload"`
}

type SendNotificationToUserResponse struct {
	Response
	result *PusherResponse
}

func (f SendNotificationToUserResponse) Result() *PusherResponse {
	return f.result
}

func (p *Pusher) SendNotificationToUser(ctx *assecoContext.RequestContext, msg *PusherUserMessage, headers map[string]string) *SendNotificationToUserResponse {
	var result *PusherResponse
	req := p.NewRequest(http.MethodPost).
		AddPath("api/users").
		SetBody(msg).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)

	return &SendNotificationToUserResponse{
		Response: req.Response(),
		result:   result,
	}
}

type PusherResponse struct {
	Message string `json:"message"`
}
