package openai

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"

	er "github.com/futugyousuzu/go-openai/internal"
)

const endTag string = "[DONE]"

var headerData []byte = []byte("data: ")

type StreamResponse struct {
	Response  *http.Response
	Reader    *bufio.Reader
	StreamEnd bool
}

func (c *StreamResponse) ReadStream(response interface{}) (e error) {
	reader := c.Reader
	if reader == nil {
		c.StreamEnd = true
		return
	}

	line, err := reader.ReadBytes('\n')
	responseStr := ""

	// for loop is for skip the row which is not start with 'data:'
	for {
		if err != nil {
			c.StreamEnd = true
			return er.SystemError(err.Error())
		}

		line = bytes.TrimSpace(line)
		if bytes.HasPrefix(line, headerData) {
			line = bytes.TrimPrefix(line, headerData)
			responseStr = string(line)
			break
		} else {
			line, err = reader.ReadBytes('\n')
		}
	}

	if responseStr == endTag {
		c.StreamEnd = true
		return
	}

	if err = json.Unmarshal(line, response); err != nil {
		return er.SystemError(err.Error())
	}

	return
}

func (c *StreamResponse) Close() {
	if c.Response != nil {
		c.Response.Body.Close()
	}
}

func (c *StreamResponse) CanReadStream() bool {
	return !c.StreamEnd
}
