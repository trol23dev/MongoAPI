package domain

type QueryFilter struct {
	Limit     int64
	Skip      int64
	SortField string
	SortOrder int // 1 for asc, -1 for desc
	RawFilter string
}
