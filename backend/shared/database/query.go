package database

import (
	"github.com/andrezz-b/stem24-phishing-tracker/shared/constants"
	"strings"
	"time"
)

type Pagination struct {
	Page    *int `json:"page" form:"page"`
	PerPage *int `json:"per_page" form:"per_page"`
}

type WhereQuery interface {
	Operator() string
	Key() string
	Value() interface{}
}

type Query interface {
	Build() []WhereQuery
	OrderBy() string
	Offset() int
	Limit() int
}

func (q Pagination) Offset() int {
	if q.Page == nil {
		return constants.DefaultPage
	}

	return *q.Page * q.Limit()
}

func (q Pagination) Limit() int {
	if q.PerPage == nil {
		return constants.DefaultPerPage
	}
	return *q.PerPage
}

type Sort struct {
	Sort string `json:"sort" form:"sort"`
}

func (q Sort) OrderBy() string {
	if q.Sort == "" {
		return "created_at DESC"
	}
	return q.Sort
}

type GetAllAgentsRequest struct {
	Pagination
	Sort
	Id              []string `json:"id" form:"id"`
	GlobalStatusId  []string `json:"global_status_id" form:"global_status_id"`
	Extension       []string `json:"extension" form:"extension"`
	Number          []string `json:"number" form:"number"`
	DisplayName     string   `json:"display_name" form:"display_name"`
	FirstName       string   `json:"first_name" form:"first_name"`
	LastName        string   `json:"last_name" form:"last_name"`
	SkillGroupId    []string `json:"skill_group_id" form:"skill_group_id"`
	IsBlocked       *bool    `json:"is_blocked" form:"is_blocked"`
	ChannelId       string   `json:"channel_id" form:"channel_id"`
	IncludeBlocked  bool     `json:"include_blocked" form:"include_blocked"`
	ChannelStatusId []string `json:"channel_status_id" form:"channel_status_id"`
}

type QueryItem struct {
	value    interface{}
	key      string
	operator string
}

func (q QueryItem) Value() interface{} {
	return q.value
}

func (q QueryItem) Key() string {
	return q.key
}

func (q QueryItem) Operator() string {
	return q.operator
}

func (g GetAllAgentsRequest) Build() []WhereQuery {
	qMap := make([]WhereQuery, 0)
	if g.Id != nil {
		qMap = append(qMap, QueryItem{
			value:    g.Id,
			key:      "agents.id",
			operator: "IN",
		})
	}
	if g.GlobalStatusId != nil {
		qMap = append(qMap, QueryItem{
			value:    g.GlobalStatusId,
			key:      "agents.global_status_id",
			operator: "IN",
		})
	}
	if g.Extension != nil {
		qMap = append(qMap, QueryItem{
			value:    g.Extension,
			key:      "agents.extension",
			operator: "IN",
		})
	}
	if g.Number != nil {
		qMap = append(qMap, QueryItem{
			value:    g.Number,
			key:      "agents.number",
			operator: "IN",
		})
	}
	if g.DisplayName != "" {
		qMap = append(qMap, QueryItem{
			value:    strings.ToLower("%" + string(g.DisplayName) + "%"),
			key:      "LOWER(agents.display_name)",
			operator: "LIKE",
		})
	}
	if g.FirstName != "" {
		qMap = append(qMap, QueryItem{
			value:    strings.ToLower("%" + string(g.FirstName) + "%"),
			key:      "LOWER(agents.first_name)",
			operator: "LIKE",
		})
	}
	if g.LastName != "" {
		qMap = append(qMap, QueryItem{
			value:    strings.ToLower("%" + string(g.LastName) + "%"),
			key:      "LOWER(agents.last_name)",
			operator: "LIKE",
		})
	}
	if g.SkillGroupId != nil {
		qMap = append(qMap, QueryItem{
			value:    g.SkillGroupId,
			key:      "sga.skill_group_id",
			operator: "IN",
		})
	}
	if !g.IncludeBlocked {
		qMap = append(qMap,
			QueryItem{
				value:    false,
				key:      "st.blocked",
				operator: "=",
			},
			QueryItem{
				value:    false,
				key:      "gst.blocked",
				operator: "=",
			},
		)
	}

	if g.ChannelId != "" {
		qMap = append(qMap, QueryItem{
			value:    g.ChannelId,
			key:      "st.channel_id",
			operator: "=",
		})
	}

	if g.IsBlocked != nil {
		qMap = append(qMap, QueryItem{
			value:    *g.IsBlocked,
			key:      "st.blocked",
			operator: "=",
		})
	}

	if g.ChannelStatusId != nil {
		qMap = append(qMap, QueryItem{
			value:    g.ChannelStatusId,
			key:      "ast.status_id",
			operator: "IN",
		})
	}

	return qMap
}

type GetAllQueuedItemsRequest struct {
	Pagination
	Sort
	ID           []string `json:"id" form:"id"`
	SessionUuid  []string `json:"session_uuid" form:"session_uuid"`
	QueueCount   []*int   `json:"queue_count" form:"queue_count"`
	SkillGroupID []string `json:"skill_group_id" form:"skill_group_id"`
	Queued       *bool    `json:"queued" form:"queued"`
	Answered     *bool    `json:"answered" form:"answered"`
	AgentID      []string `json:"agent_id" form:"agent_id"`
	PonderValue  []*int   `json:"ponder_value" form:"ponder_value"`
}

func (g GetAllQueuedItemsRequest) Build() []WhereQuery {
	qMap := make([]WhereQuery, 0)
	if g.ID != nil {
		qMap = append(qMap, QueryItem{
			value:    g.ID,
			key:      "id",
			operator: "IN",
		})
	}
	if g.SessionUuid != nil {
		qMap = append(qMap, QueryItem{
			value:    g.SessionUuid,
			key:      "session_uuid",
			operator: "IN",
		})
	}
	if g.QueueCount != nil {
		qMap = append(qMap, QueryItem{
			value:    g.QueueCount,
			key:      "queue_count",
			operator: "IN",
		})
	}
	if g.SkillGroupID != nil {
		qMap = append(qMap, QueryItem{
			value:    g.SkillGroupID,
			key:      "skill_group_id",
			operator: "IN",
		})
	}
	if g.Queued != nil {
		qMap = append(qMap, QueryItem{
			value:    g.Queued,
			key:      "queued",
			operator: "=",
		})
	}
	if g.Answered != nil {
		qMap = append(qMap, QueryItem{
			value:    g.Answered,
			key:      "answered",
			operator: "=",
		})
	}
	if g.AgentID != nil {
		qMap = append(qMap, QueryItem{
			value:    g.AgentID,
			key:      "agent_id",
			operator: "IN",
		})
	}
	if g.PonderValue != nil {
		qMap = append(qMap, QueryItem{
			value:    g.PonderValue,
			key:      "ponder_value",
			operator: "IN",
		})
	}
	return qMap
}

type GetAllSkillGroupsRequest struct {
	Pagination
	Sort
	ID            []string `json:"id" form:"id"`
	Name          []string `json:"name" form:"name"`
	Description   []string `json:"description" form:"description"`
	ChannelID     []string `json:"channel_id" form:"channel_id"`
	PromptTimeout []*int64 `json:"prompt_timeout" form:"prompt_timeout"`
	QueueTimeout  []*int64 `json:"queue_timeout" form:"queue_timeout"`
	From          []string `json:"from" form:"from"`
	Disabled      *bool    `json:"disabled" form:"disabled"`
}

func (q GetAllSkillGroupsRequest) Build() []WhereQuery {
	qMap := make([]WhereQuery, 0)
	if q.ID != nil {
		qMap = append(qMap, QueryItem{
			value:    q.ID,
			key:      "id",
			operator: "IN",
		})
	}
	if q.Name != nil {
		qMap = append(qMap, QueryItem{
			value:    q.Name,
			key:      "name",
			operator: "IN",
		})
	}
	if q.Description != nil {
		qMap = append(qMap, QueryItem{
			value:    q.Description,
			key:      "description",
			operator: "IN",
		})
	}
	if q.ChannelID != nil {
		qMap = append(qMap, QueryItem{
			value:    q.ChannelID,
			key:      "channel_id",
			operator: "IN",
		})
	}
	if q.From != nil {
		qMap = append(qMap, QueryItem{
			value:    q.From,
			key:      "from",
			operator: "IN",
		})
	}
	if q.PromptTimeout != nil {
		qMap = append(qMap, QueryItem{
			value:    q.PromptTimeout,
			key:      "prompt_timeout",
			operator: "IN",
		})
	}
	if q.QueueTimeout != nil {
		qMap = append(qMap, QueryItem{
			value:    q.QueueTimeout,
			key:      "queue_timeout",
			operator: "IN",
		})
	}

	if q.Disabled != nil {
		value := *q.Disabled
		qMap = append(qMap, QueryItem{
			value:    value,
			key:      "disabled",
			operator: "=",
		})
	}
	return qMap
}

type GetAllChannelsRequest struct {
	Pagination
	Sort
	ID   []string `json:"id" form:"id"`
	Name []string `json:"name" form:"name"`
}

func (q GetAllChannelsRequest) Build() []WhereQuery {
	qMap := make([]WhereQuery, 0)
	if q.ID != nil {
		qMap = append(qMap, QueryItem{
			value:    q.ID,
			key:      "id",
			operator: "IN",
		})
	}
	if q.Name != nil {
		qMap = append(qMap, QueryItem{
			value:    q.Name,
			key:      "name",
			operator: "IN",
		})
	}

	return qMap
}

type GetAllAgentPonderValueSnapshotsRequest struct {
	Pagination
	Sort
	ID                 []string `json:"id" form:"id"`
	AgentID            []string `json:"agent_id" form:"agent_id"`
	QueuedItemUuid     []string `json:"queued_item_uuid" form:"queued_item_uuid"`
	SessionUuid        []string `json:"session_uuid" form:"session_uuid"`
	Value              []int64  `json:"value" form:"value"`
	Blocked            []bool   `json:"blocked" form:"blocked"`
	GlobalStatusID     []string `json:"global_status_id" form:"global_status_id"`
	ChannelStatusID    []string `json:"channel_status_id" form:"channel_status_id"`
	Matched            []bool   `json:"matched" form:"matched"`
	SatisfiesThreshold []bool   `json:"satisfies_threshold" form:"satisfies_threshold"`
	Error              []string `json:"error" form:"error"`
	Message            []string `json:"message" form:"message"`
	DateTimeFrom       string   `json:"date_time_from" form:"date_time_from"`
	DateTimeTo         string   `json:"date_time_to" form:"date_time_to"`
	SkillGroupId       []string `json:"skill_group_id" form:"skill_group_id"`
}

func (g GetAllAgentPonderValueSnapshotsRequest) Build() []WhereQuery {
	qMap := make([]WhereQuery, 0)
	if g.ID != nil {
		qMap = append(qMap, QueryItem{
			value:    g.ID,
			key:      "ap.id",
			operator: "IN",
		})
	}
	if g.AgentID != nil {
		qMap = append(qMap, QueryItem{
			value:    g.AgentID,
			key:      "agent_id",
			operator: "IN",
		})
	}
	if g.QueuedItemUuid != nil {
		qMap = append(qMap, QueryItem{
			value:    g.QueuedItemUuid,
			key:      "queued_item_uuid",
			operator: "IN",
		})
	}
	if g.SessionUuid != nil {
		qMap = append(qMap, QueryItem{
			value:    g.SessionUuid,
			key:      "session_uuid",
			operator: "IN",
		})
	}
	if g.Value != nil {
		qMap = append(qMap, QueryItem{
			value:    g.Value,
			key:      "value",
			operator: "IN",
		})
	}
	if g.Blocked != nil {
		qMap = append(qMap, QueryItem{
			value:    g.Blocked,
			key:      "blocked",
			operator: "IN",
		})
	}
	if g.GlobalStatusID != nil {
		qMap = append(qMap, QueryItem{
			value:    g.GlobalStatusID,
			key:      "global_status",
			operator: "IN",
		})
	}
	if g.ChannelStatusID != nil {
		qMap = append(qMap, QueryItem{
			value:    g.ChannelStatusID,
			key:      "channel_status",
			operator: "IN",
		})
	}
	if g.Matched != nil {
		qMap = append(qMap, QueryItem{
			value:    g.Matched,
			key:      "matched",
			operator: "=",
		})
	}
	if g.SatisfiesThreshold != nil {
		qMap = append(qMap, QueryItem{
			value:    g.SatisfiesThreshold,
			key:      "satisfies_threshold",
			operator: "=",
		})
	}
	if g.Error != nil {
		qMap = append(qMap, QueryItem{
			value:    g.Error,
			key:      "error",
			operator: "IN",
		})
	}
	if g.Message != nil {
		qMap = append(qMap, QueryItem{
			value:    g.Message,
			key:      "message",
			operator: "IN",
		})
	}
	if g.DateTimeFrom != "" && g.DateTimeTo != "" {
		layout := "2006-01-02T15:04:05.000Z"
		timeFrom, _ := time.Parse(layout, g.DateTimeFrom)
		timeTo, _ := time.Parse(layout, g.DateTimeTo)

		qMap = append(qMap, QueryItem{
			value:    timeFrom,
			key:      "ap.created_at",
			operator: ">=",
		})
		qMap = append(qMap, QueryItem{
			value:    timeTo,
			key:      "ap.created_at",
			operator: "<=",
		})
	}
	if g.SkillGroupId != nil {
		qMap = append(qMap, QueryItem{
			value:    g.SkillGroupId,
			key:      "skill_group_id",
			operator: "IN",
		})
	}
	return qMap
}

type GetAllAgentGlobalHistoryItemsRequest struct {
	Pagination
	Sort
	ID           []string `json:"id" form:"id"`
	AgentID      []string `json:"agent_id" form:"agent_id"`
	StatusID     []string `json:"status_id" form:"status_id"`
	ExecutedBy   []string `json:"executed_by" form:"executed_by"`
	Reason       []string `json:"reason" form:"reason"`
	DateTimeFrom string   `json:"date_time_from" form:"date_time_from"`
	DateTimeTo   string   `json:"date_time_to" form:"date_time_to"`
}

func (h GetAllAgentGlobalHistoryItemsRequest) Build() []WhereQuery {
	qMap := make([]WhereQuery, 0)
	if h.ID != nil {
		qMap = append(qMap, QueryItem{
			value:    h.ID,
			key:      "ah.id",
			operator: "IN",
		})
	}
	if h.AgentID != nil {
		qMap = append(qMap, QueryItem{
			value:    h.AgentID,
			key:      "ah.agent_id",
			operator: "IN",
		})
	}
	if h.StatusID != nil {
		qMap = append(qMap, QueryItem{
			value:    h.StatusID,
			key:      "ah.status_id",
			operator: "IN",
		})
	}
	if h.ExecutedBy != nil {
		qMap = append(qMap, QueryItem{
			value:    h.ExecutedBy,
			key:      "ah.executed_by",
			operator: "IN",
		})
	}
	if h.Reason != nil {
		qMap = append(qMap, QueryItem{
			value:    h.Reason,
			key:      "ah.reason",
			operator: "IN",
		})
	}
	if h.DateTimeFrom != "" && h.DateTimeTo != "" {
		layout := "2006-01-02T15:04:05.000Z"
		timeFrom, _ := time.Parse(layout, h.DateTimeFrom)
		timeTo, _ := time.Parse(layout, h.DateTimeTo)

		qMap = append(qMap, QueryItem{
			value:    timeFrom,
			key:      "ah.created_at",
			operator: ">=",
		})
		qMap = append(qMap, QueryItem{
			value:    timeTo,
			key:      "ah.created_at",
			operator: "<=",
		})
	}
	return qMap
}

type GetAllAgentChannelHistoryItemsRequest struct {
	Pagination
	Sort
	ID           []string `json:"id" form:"id"`
	AgentID      []string `json:"agent_id" form:"agent_id"`
	StatusID     []string `json:"status_id" form:"status_id"`
	ChannelID    []string `json:"channel_id" form:"channel_id"`
	ExecutedBy   []string `json:"executed_by" form:"executed_by"`
	Reason       []string `json:"reason" form:"reason"`
	DateTimeFrom string   `json:"date_time_from" form:"date_time_from"`
	DateTimeTo   string   `json:"date_time_to" form:"date_time_to"`
}

func (h GetAllAgentChannelHistoryItemsRequest) Build() []WhereQuery {
	qMap := make([]WhereQuery, 0)
	if h.ID != nil {
		qMap = append(qMap, QueryItem{
			value:    h.ID,
			key:      "ah.id",
			operator: "IN",
		})
	}
	if h.AgentID != nil {
		qMap = append(qMap, QueryItem{
			value:    h.AgentID,
			key:      "ah.agent_id",
			operator: "IN",
		})
	}
	if h.StatusID != nil {
		qMap = append(qMap, QueryItem{
			value:    h.StatusID,
			key:      "ah.status_id",
			operator: "IN",
		})
	}
	if h.ChannelID != nil {
		qMap = append(qMap, QueryItem{
			value:    h.ChannelID,
			key:      "ah.channel_id",
			operator: "IN",
		})
	}
	if h.ExecutedBy != nil {
		qMap = append(qMap, QueryItem{
			value:    h.ExecutedBy,
			key:      "ah.executed_by",
			operator: "IN",
		})
	}
	if h.Reason != nil {
		qMap = append(qMap, QueryItem{
			value:    h.Reason,
			key:      "ah.reason",
			operator: "IN",
		})
	}
	if h.DateTimeFrom != "" && h.DateTimeTo != "" {
		layout := "2006-01-02T15:04:05.000Z"
		timeFrom, _ := time.Parse(layout, h.DateTimeFrom)
		timeTo, _ := time.Parse(layout, h.DateTimeTo)

		qMap = append(qMap, QueryItem{
			value:    timeFrom,
			key:      "ah.created_at",
			operator: ">=",
		})
		qMap = append(qMap, QueryItem{
			value:    timeTo,
			key:      "ah.created_at",
			operator: "<=",
		})
	}
	return qMap
}
