package dto

type BaseListResponse[T any] struct {
	Total   int64 `json:"total"`
	Results T     `json:"results"`
}
