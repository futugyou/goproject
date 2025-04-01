package messages

type CancelledNotification struct {
	RequestId RequestId `json:"requestId"`
	Reason    *string   `json:"reason"`
}
