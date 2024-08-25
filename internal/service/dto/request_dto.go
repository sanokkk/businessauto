package dto

import "autoshop/internal/storage/filters"

type Request struct {
	Body filters.FilterBody `json:"body"`
}
