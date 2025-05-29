package shared

import (
	"io"
	"net/http"
)

type DuplexPipe struct {
	Input  io.Reader
	Output io.Writer
}

func NewInMemoryDuplexPipe() (*DuplexPipe, io.Closer) {
	pr, pw := io.Pipe()
	return &DuplexPipe{
		Input:  pr,
		Output: pw,
	}, pw
}

func NewHttpDuplexPipe(r *http.Request, w http.ResponseWriter) *DuplexPipe {
	return &DuplexPipe{
		Input:  r.Body,
		Output: w,
	}
}
