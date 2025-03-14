package httputils

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/futugyou/yomawari/common/errorutils"
)

const endTag string = "[DONE]"

var headerData []byte = []byte("data: ")

type StreamResponse struct {
	Response  *http.Response
	Reader    *bufio.Reader
	StreamEnd bool
}

func (c *StreamResponse) ReadStream(ctx context.Context, response interface{}) *errorutils.OpenaiError {
	if c.Reader == nil {
		c.StreamEnd = true
		return errorutils.SystemError("stream reader is nil")
	}

	for {
		select {
		case <-ctx.Done():
			c.StreamEnd = true
			return errorutils.SystemError("context canceled: " + ctx.Err().Error())
		default:
			line, err := c.Reader.ReadBytes('\n')
			if err != nil {
				c.StreamEnd = true
				return errorutils.SystemError("failed to read from stream: " + err.Error())
			}

			line = bytes.TrimSpace(line)
			if bytes.HasPrefix(line, headerData) {
				line = bytes.TrimPrefix(line, headerData)
				responseStr := string(line)

				if responseStr == endTag {
					c.StreamEnd = true
					return nil
				}

				if err := json.Unmarshal(line, response); err != nil {
					c.StreamEnd = true
					return errorutils.SystemError("failed to unmarshal JSON: " + err.Error())
				}

				return nil
			}
		}
	}
}

func (c *StreamResponse) Close() {
	if c.Response != nil {
		c.Response.Body.Close()
	}
}

func (c *StreamResponse) CanReadStream() bool {
	return !c.StreamEnd
}
