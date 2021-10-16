package filter

type Pagination struct {
	Limit  int64
	Offset int64
	Skip   int64
}

func NewPagination(limit int64, offset int64) *Pagination {

	if limit <= 0 {
		limit = 10
	}

	if offset <= 0 {
		offset = 1
	}

	skip := limit * (offset - 1)

	return &Pagination{Limit: limit, Offset: offset, Skip: skip}
}
