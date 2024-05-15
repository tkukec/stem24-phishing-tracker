package clients

import (
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/runtimebag"
	assecoContext "git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/context"
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/helpers"
	"net/http"
	"time"
)

const (
	AgentManagementUriEnvKey = "AGENT_MANAGEMENT_URI"

	agentsEndpoint      = "api/agents"
	skillGroupsEndpoint = "api/skill-groups"
	queuesEndpoint      = "api/queues"
	skillGroups         = "/skill-groups"
)

type AgentManagement struct {
	uri string
}

func NewAgentManagement() *AgentManagement {
	return &AgentManagement{uri: helpers.UrlParse(runtimebag.GetEnvString(
		AgentManagementUriEnvKey,
		"http://agent-management:8080",
	))}
}

func (a *AgentManagement) SetUri(uri string) *AgentManagement {
	a.uri = helpers.UrlParse(uri)
	return a
}

func (a *AgentManagement) NewRequest(method string) *RequestBuilder {
	return NewRequestBuilder(a.uri, method)
}

type AgentsOfSkillGroups struct {
	Response
	result Agents
}

func (q AgentsOfSkillGroups) Result() Agents {
	return q.result
}

func (a *AgentManagement) AgentsOfSkillGroups(ctx *assecoContext.RequestContext, headers map[string]string, includeBlocked bool, skillGroupIds ...string) *AgentsOfSkillGroups {
	var result Agents
	req := a.NewRequest(http.MethodGet).
		AddPath(agentsEndpoint).
		AddBoolQuery("include_blocked", includeBlocked).
		AddMultipleQueryValues("skill_group_id", skillGroupIds).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)

	return &AgentsOfSkillGroups{
		Response: req.Response(),
		result:   result,
	}
}

type SkillGroupsResponse struct {
	Response
	result SkillGroups
}

func (q SkillGroupsResponse) Result() SkillGroups {
	return q.result
}

func (a *AgentManagement) SkillGroups(ctx *assecoContext.RequestContext, headers map[string]string, skillGroupIds ...string) *SkillGroupsResponse {
	var result SkillGroups
	req := a.NewRequest(http.MethodGet).
		AddPath(skillGroupsEndpoint).
		AddMultipleQueryValues("id", skillGroupIds).
		Bind(&result).
		Print().
		Run(ctx, headers)

	return &SkillGroupsResponse{
		Response: req.Response(),
		result:   result,
	}
}

type QueueCountOfSkillGroupIdResponse struct {
	Response
	result int
}

func (q QueueCountOfSkillGroupIdResponse) Result() int {
	return q.result
}

func (a *AgentManagement) QueueCountOfSkillGroupId(ctx *assecoContext.RequestContext, skillGroupId string, headers map[string]string) *QueueCountOfSkillGroupIdResponse {
	type QueueCountOfSkillGroup struct {
		Count int `json:"count"`
	}
	var result QueueCountOfSkillGroup
	req := a.NewRequest(http.MethodGet).
		AddPath(queuesEndpoint).
		AddPath(skillGroupId).
		AddPath("items/count").
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)

	if req.Error() != nil {
		return &QueueCountOfSkillGroupIdResponse{
			Response: req.Response(),
		}
	}

	return &QueueCountOfSkillGroupIdResponse{
		Response: req.Response(),
		result:   result.Count,
	}
}

type ContactPositionInQueueResponse struct {
	Response
	result *ContactPositionInQueue
}

func (q ContactPositionInQueueResponse) Result() *ContactPositionInQueue {
	return q.result
}

type ContactPositionInQueue struct {
	Position  int `json:"position"`
	QueueTime int `json:"queue_time"`
}

func (a *AgentManagement) ContactPositionInQueue(ctx *assecoContext.RequestContext, skillGroupId, contact string, headers map[string]string) *ContactPositionInQueueResponse {
	var result *ContactPositionInQueue
	req := a.NewRequest(http.MethodGet).
		AddPath(queuesEndpoint).
		AddPath(skillGroupId).
		AddPath("position").
		AddPath(contact).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)

	return &ContactPositionInQueueResponse{
		Response: req.Response(),
		result:   result,
	}
}

func (s SkillGroups) Count() int {
	return len(s)
}

func (a Agents) Count() int {
	return len(a)
}

func (a Agents) First() *Agent {
	if a.Count() == 0 {
		return nil
	}
	return a[0]
}

func (a Agents) Unblocked() Agents {
	unblocked := make(Agents, 0)
	for _, agent := range a {
		if agent.GlobalStatus != nil && !agent.GlobalStatus.Blocked {
			unblocked = append(unblocked, agent)
		}
	}
	return unblocked
}

func (a Agents) UnblockedOnChannel(channelId string) Agents {
	unblocked := make(Agents, 0)
	for _, agent := range a {
		if agent.GlobalStatus != nil && !agent.GlobalStatus.Blocked && !agent.IsBlockedOnChannel(channelId) {
			unblocked = append(unblocked, agent)
		}
	}
	return unblocked
}

func (a Agents) Blocked() Agents {
	blocked := make(Agents, 0)
	for _, agent := range a {
		if agent.GlobalStatus != nil && agent.GlobalStatus.Blocked {
			blocked = append(blocked, agent)
		}
	}
	return blocked
}

func (a *Agent) StatusOfChannel(channelId string) *Status {
	for _, status := range a.Statuses {
		if status.ChannelID == channelId {
			return status
		}
	}
	return nil
}

func (a *Agent) IsBlockedOnChannel(channelId string) bool {
	status := a.StatusOfChannel(channelId)
	if status == nil {
		return true
	}
	return false
}

type Agents []*Agent

type SkillGroups []*AgentSkillGroup

type Agent struct {
	ID                    string             `json:"id"`
	GlobalStatus          *GlobalStatus      `json:"global_status"`
	GlobalStatusChangedAT time.Time          `json:"global_status_changed_at"`
	Statuses              []*Status          `json:"statuses"`
	Extension             string             `json:"extension"`
	Number                string             `json:"number"`
	DisplayName           string             `json:"display_name"`
	Activities            []*Activity        `json:"activities"`
	SkillGroups           []*AgentSkillGroup `json:"skill_group,omitempty"`

	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
	Utilization int       `json:"utilization"`
}

type GlobalStatus struct {
	ID              string          `json:"id"`
	Blocked         bool            `json:"blocked"`
	OnReject        bool            `json:"on_reject"`
	OnTimeout       bool            `json:"on_timeout"`
	StartingStatus  bool            `json:"starting_status"`
	Name            string          `json:"name"`
	Reason          string          `json:"reason"`
	Label           string          `json:"label"`
	Transitions     []*GlobalStatus `json:"transitions"`
	TimerTransition *GlobalStatus   `json:"timer_transition"`
	Timer           int64           `json:"timer"`
	SystemStatus    bool            `json:"system_status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type Status struct {
	ID              string    `json:"id"`
	Blocked         bool      `json:"blocked"`
	System          bool      `json:"system_status"`
	OnReject        bool      `json:"on_reject"`
	OnTimeout       bool      `json:"on_timeout"`
	StartingStatus  bool      `json:"starting_status"`
	Reason          string    `json:"reason"`
	Label           string    `json:"label"`
	Name            string    `json:"name"`
	ChannelID       string    `json:"channel_id"`
	Transitions     []*Status `json:"transitions"`
	TimerTransition *Status   `json:"timer_transition"`
	Timer           int64     `json:"timer"`
	ChangedAt       time.Time `json:"changed_at,omitempty"`

	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	DeletedAt        time.Time `json:"deleted_at"`
	DefaultBlocked   bool      `json:"default_blocked"`
	DefaultUnblocked bool      `json:"default_unblocked"`
}

type Activity struct {
	ID           string          `json:"id"`
	ChannelId    string          `json:"channel_id"`
	AgentId      string          `json:"agent_id"`
	SessionUuid  string          `json:"session_uuid"`
	Status       *ActivityStatus `json:"status,omitempty"`
	UacdQueued   bool            `json:"uacd_queued"`
	Closed       bool            `json:"closed"`
	QueuedItemId string          `json:"queued_item_id"`
	QueuedItem   *QueuedItem     `json:"queued_item"`

	ConferenceUuid string `json:"conference_uuid,omitempty"`
	RecordingState string `json:"recording_state"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`

	Parent      *Activity   `json:"parent"`
	Children    []*Activity `json:"children"`
	Type        string      `json:"type"`
	Contact     string      `json:"contact"`
	DisplayName string      `json:"display_name"`

	Muted        bool                   `json:"muted"`
	SkillGroupId string                 `json:"skill_group_id"`
	ContextData  map[string]interface{} `json:"context_data"`
	IvrData      map[string]interface{} `json:"ivr_data"`
}

type QueuedItem struct {
	ID             string                 `json:"id"`
	SessionUuid    string                 `json:"session_uuid"`
	ContextData    map[string]interface{} `json:"context_data"`
	IvrData        map[string]interface{} `json:"ivr_data"`
	QueueID        string                 `json:"queue_id"`
	Queued         bool                   `json:"queued"`
	Answered       bool                   `json:"answered"`
	QueueCount     int                    `json:"queue_count"`
	AgentID        string                 `json:"agent_id,omitempty"`
	Agent          *Agent                 `json:"agent,omitempty"`
	AdditionalData interface{}            `json:"additional_data"`
	Session        interface{}            `json:"session"`
	PonderValue    int64                  `json:"ponder_value"`
	AutoCallBack   bool                   `json:"auto_call_back"`
	Contact        string                 `json:"contact"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	DeletedAt      time.Time              `json:"deleted_at"`
	DisplayName    string                 `json:"display_name"`
	ExternalPonder int64                  `json:"external_ponder"`
	PreferredAgent string                 `json:"preferred_agent"`
	Matched        bool                   `json:"matched"`
}

type ActivityStatus struct {
	ID     string `json:"id"`
	Label  string `json:"label"`
	Name   string `json:"name"`
	Reason string `json:"reason"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type AgentSkillGroup struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	TelephoneNumber string    `json:"telephone_number"`
	Agents          []*Agent  `json:"agents,omitempty"`
	QueueTimeout    int64     `json:"queue_timeout"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       time.Time `json:"deleted_At"`
	ChannelID       string    `json:"channel_id"`
	PromptTimeout   int64     `json:"prompt_timeout"`
	UseAmd          bool      `json:"use_amd"`

	//SkillLevel only exists and is visible when in context of agent relationship
	SkillLevel int64 `json:"skill_level,omitempty"`
	//Primary only exists and is visible when in context of agent relationship
	Primary bool `json:"primary,omitempty"`
}
