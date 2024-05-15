package clients

import (
	"fmt"
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/runtimebag"
	assecoContext "git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/context"
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/helpers"
	"net/http"
)

const (
	AssecoChatUriEnvKey = "ASSECO_CHAT_URI"
)

type AssecoChat struct {
	uri string
}

func NewAssecoChat() *AssecoChat {
	return &AssecoChat{uri: helpers.UrlParse(runtimebag.GetEnvString(
		AssecoChatUriEnvKey,
		"http://social-networks:8080",
	))}
}

func (t *AssecoChat) SetUri(uri string) *AssecoChat {
	t.uri = helpers.UrlParse(uri)
	return t
}

func (t *AssecoChat) NewRequest(method string) *RequestBuilder {
	return NewRequestBuilder(t.uri, method)
}

func (t *AssecoChat) NewSearchRequest() *RequestBuilder {
	return t.NewRequest(http.MethodPost).AddPath("api/search/conversation")
}

type FindConversationByIdResponse struct {
	Response
	result SocialNetworksConversationCollection
}

func (r *FindConversationByIdResponse) Result() SocialNetworksConversationCollection {
	return r.result
}

func (t *AssecoChat) FindConversationById(ctx *assecoContext.RequestContext, id string, headers map[string]string) *FindConversationByIdResponse {
	var result SocialNetworksConversationCollection
	req := t.NewSearchRequest().
		SetBody(map[string]interface{}{
			"search": map[string]string{
				"id": fmt.Sprintf("=%s", id),
			},
			"relations": []string{
				"remote_relations",
				"network",
				"messages",
			},
		}).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)

	return &FindConversationByIdResponse{
		Response: req.Response(),
		result:   result,
	}
}

func (t *AssecoChat) FindConversationByProviderConversationId(ctx *assecoContext.RequestContext, providerConversationId string, headers map[string]string) *FindConversationByIdResponse {
	var result SocialNetworksConversationCollection
	req := t.NewRequest(http.MethodPost).
		AddPath("api/search/conversation").
		SetBody(map[string]interface{}{
			"search": map[string]string{
				"provider_conversation_id": fmt.Sprintf("=%s", providerConversationId),
			},
			"relations": []string{
				"remote_relations",
				"network",
				"messages",
			},
		}).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)

	return &FindConversationByIdResponse{
		Response: req.Response(),
		result:   result,
	}
}

type SocialNetworksConversationCollection []*SocialNetworksConversation

func (c SocialNetworksConversationCollection) First() *SocialNetworksConversation {
	if c != nil && len(c) > 0 {
		return c[0]
	}
	return nil
}

func (c SocialNetworksConversationCollection) IsEmpty() bool {
	if c == nil || len(c) == 0 {
		return true
	}
	return false
}

type SocialNetworksConversation struct {
	Id                     string                   `json:"id"`
	NetworkId              string                   `json:"network_id"`
	Client                 string                   `json:"client"`
	ProviderConversationId string                   `json:"provider_conversation_id"`
	SkillGroupId           string                   `json:"skill_group_id"`
	ResponsibleAgentId     string                   `json:"responsible_agent_id"`
	Closed                 bool                     `json:"closed"`
	CreatedAt              string                   `json:"created_at"`
	CreatedBy              string                   `json:"created_by"`
	CreatorType            string                   `json:"creator_type"`
	UpdatedAt              string                   `json:"updated_at"`
	UpdatedBy              string                   `json:"updated_by"`
	UpdaterType            string                   `json:"updater_type"`
	DeletedAt              string                   `json:"deleted_at"`
	DeletedBy              string                   `json:"deleted_by"`
	DeleterType            string                   `json:"deleter_type"`
	Messages               []*SocialNetworksMessage `json:"messages"`
	Network                *SocialNetworkNetwork    `json:"network"`
}

type SocialNetworksMessage struct {
	Id                string `json:"id"`
	From              string `json:"from"`
	To                string `json:"to"`
	Message           string `json:"message"`
	Read              bool   `json:"read"`
	ProviderMessageId string `json:"provider_message_id"`
	ConversationId    string `json:"conversation_id"`
	StatusId          string `json:"status_id"`
	TypeId            string `json:"type_id"`
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

type SocialNetworkNetwork struct {
	Id                 string      `json:"id"`
	Name               string      `json:"name"`
	Driver             string      `json:"driver"`
	Sender             interface{} `json:"sender"`
	InboundBackground  string      `json:"inbound_background"`
	InboundFontColor   string      `json:"inbound_font_color"`
	OutboundBackground string      `json:"outbound_background"`
	OutboundFontColor  string      `json:"outbound_font_color"`
	CreatedAt          string      `json:"created_at"`
	CreatedBy          string      `json:"created_by"`
	CreatorType        string      `json:"creator_type"`
	UpdatedAt          string      `json:"updated_at"`
	UpdatedBy          string      `json:"updated_by"`
	UpdaterType        string      `json:"updater_type"`
	DeletedAt          string      `json:"deleted_at"`
	DeletedBy          string      `json:"deleted_by"`
	DeleterType        string      `json:"deleter_type"`
}
