package types

// PaginationQuery — embedded in list query structs
type PaginationQuery struct {
	Page    int `form:"page"`
	PerPage int `form:"per_page"`
}

// Normalize sets default values and clamps limits
func (q *PaginationQuery) Normalize() {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PerPage < 1 {
		q.PerPage = 20
	}
	if q.PerPage > 100 {
		q.PerPage = 100
	}
}

func (q *PaginationQuery) Offset() int {
	return (q.Page - 1) * q.PerPage
}

func (q *PaginationQuery) Limit() int {
	return q.PerPage
}

// PaginatedResponse — generic wrapper cho list endpoints
type PaginatedResponse[T any] struct {
	Data    []T `json:"data"`
	Count   int `json:"count"`
	Offset  int `json:"offset"`
	Limit   int `json:"limit"`
}

func NewPaginated[T any](data []T, total, offset, limit int) PaginatedResponse[T] {
	if data == nil {
		data = []T{}
	}
	return PaginatedResponse[T]{
		Data:   data,
		Count:  total,
		Offset: offset,
		Limit:  limit,
	}
}
