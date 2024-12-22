package qstash

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const endTag string = "[DONE]"

var headerData []byte = []byte("data: ")

type StreamResponse struct {
	Response  *http.Response
	Reader    *bufio.Reader
	StreamEnd bool
}

func (c *StreamResponse) ReadStream(ctx context.Context, response interface{}) error {
	if c.Reader == nil {
		c.StreamEnd = true
		return fmt.Errorf("stream reader is nil")
	}

	for {
		select {
		case <-ctx.Done():
			c.StreamEnd = true
			return fmt.Errorf("context canceled: %s", ctx.Err().Error())
		default:
			line, err := c.Reader.ReadBytes('\n')
			if err != nil {
				c.StreamEnd = true
				return fmt.Errorf("failed to read from stream: %s", err.Error())
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
					return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
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
