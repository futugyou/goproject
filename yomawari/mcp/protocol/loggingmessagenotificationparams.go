package protocol

import "encoding/json"

type LoggingMessageNotificationParams struct {
	Level  LoggingLevel    `json:"level"`
	Logger *string         `json:"logger"`
	Data   json.RawMessage `json:"data"`
}
