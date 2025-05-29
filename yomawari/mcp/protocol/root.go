package protocol

type Root struct {
	Uri  string  `json:"uri"`
	Name *string `json:"name"`
	Meta any     `json:"meta"`
}
