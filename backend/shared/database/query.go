package database

import (
	"github.com/andrezz-b/stem24-phishing-tracker/shared/constants"
	"strings"
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

type GetAllEventsRequest struct {
	Pagination
	Sort
	ID       []string `json:"id" form:"id"`
	Name     []string `json:"name" form:"name"`
	Date     string   `json:"description" form:"description"`
	Brand    string   `json:"brand" form:"brand"`
	Keywords []string `json:"keywords" form:"keywords"`
}

func (q GetAllEventsRequest) Build() []WhereQuery {
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
	if q.Brand != "" {
		qMap = append(qMap, QueryItem{
			value:    strings.ToLower("%" + string(q.Brand) + "%"),
			key:      "LOWER(brand)",
			operator: "LIKE",
		})
	}
	if q.Keywords != nil {
		qMap = append(qMap, QueryItem{
			value:    q.Keywords,
			key:      "keywords",
			operator: "IN",
		})
	}
	return qMap
}

type GetAllStatusesRequest struct {
	Pagination
	Sort
	ID   []string `json:"id" form:"id"`
	Name []string `json:"name" form:"name"`
}

func (q GetAllStatusesRequest) Build() []WhereQuery {
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

type GetAllCommentsRequest struct {
	Pagination
	Sort
	ID          []string `json:"id" form:"id"`
	Description []string `json:"description" form:"description"`
	EventID     *string  `json:"event_id" form:"event_id"`
}

func (q GetAllCommentsRequest) Build() []WhereQuery {
	qMap := make([]WhereQuery, 0)
	if q.ID != nil {
		qMap = append(qMap, QueryItem{
			value:    q.ID,
			key:      "id",
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
	if q.EventID != nil {
		qMap = append(qMap, QueryItem{
			value:    q.EventID,
			key:      "EventID",
			operator: "IN",
		})
	}
	return qMap
}
