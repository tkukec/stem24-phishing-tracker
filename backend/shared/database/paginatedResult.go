package database

import (
	"github.com/andrezz-b/stem24-phishing-tracker/shared/constants"
	"math"
)

type PaginatedResult struct {
	Items      []interface{} `json:"items"`
	Page       *int          `json:"page"`
	PerPage    *int          `json:"per_page"`
	TotalCount *int64        `json:"total_count"`
	TotalPages *int          `json:"total_pages"`
}

func NewPaginatedResult(items Items, page *int, perPage *int, totalCount *int64) *PaginatedResult {
	return &PaginatedResult{
		Items:      items.GetItems(),
		Page:       page,
		PerPage:    perPage,
		TotalCount: totalCount,
		TotalPages: calculateNrOfPages(perPage, totalCount),
	}
}

func calculateNrOfPages(perPage *int, totalCount *int64) *int {
	if perPage == nil {
		defaultValue := constants.DefaultPerPage
		perPage = &defaultValue
	}
	totalPages := int(math.Ceil(float64(*totalCount) / float64(*perPage)))
	return &totalPages
}

type Items interface {
	GetItems() []interface{}
}
