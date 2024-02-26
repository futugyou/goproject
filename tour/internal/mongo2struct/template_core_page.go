package mongo2struct

const core_page_TplString = `
package core

type Paging struct {
	Page      int64
	Limit     int64
	SortField string
	Direct    SortDirect
}

const ASC sortDirect = "ASC"
const DESC sortDirect = "DESC"

type SortDirect interface {
	privateSortDirect()
	String() string
}

type sortDirect string

func (c sortDirect) privateSortDirect() {}
func (c sortDirect) String() string {
	return string(c)
}
`
