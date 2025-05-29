package protocol

type ListRootsResult struct {
	Meta  any    `json:"meta"`
	Roots []Root `json:"roots"`
}
