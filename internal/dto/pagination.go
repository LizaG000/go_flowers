package dto

type Pagination[T any] struct {
	Items       []T  `json:"items"`
	Limit       int  `json:"limit"`
	Offset      int  `json:"offset"`
	Total       int  `json:"total"`
	HasNext     bool `json:"has_next"`
	HasPrevious bool `json:"has_previous"`
}

type DictPagination struct {
	Items       []map[string]any `json:"items"`
	Limit       int              `json:"limit"`
	Offset      int              `json:"offset"`
	Total       int              `json:"total"`
	HasNext     bool             `json:"has_next"`
	HasPrevious bool             `json:"has_previous"`
}
