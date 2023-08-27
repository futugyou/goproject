package core

type Paging struct {
	Page      int64
	Limit     int64
	SortField string
	Direct    SortDirect
}

const ASC SortDirect = "ASC"
const DESC SortDirect = "DESC"

type SortDirect string
