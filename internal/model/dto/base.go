package dto

type BaseListResponse[T any] struct {
	Total   int64 `json:"total,omitempty"`
	Results T     `json:"results"`
}
