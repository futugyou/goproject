package core

type Search struct {
	Sort   map[string]int
	Page   int
	Size   int
	Filter map[string]string
}
