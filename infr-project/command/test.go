package command

import (
	"context"
	"fmt"
	"time"
)

type BookRoom struct {
	RoomId    string    `json:"room_id,omitempty"`
	GuestName string    `json:"guest_name,omitempty"`
	StartDate time.Time `json:"start_date,omitempty"`
	EndDate   time.Time `json:"end_date,omitempty"`
}

type BookRoomHandler struct {
}

func NewBookRoomHandler() BookRoomHandler {
	fmt.Println("NewBookRoomHandler called by system")
	return BookRoomHandler{}
}

func (b BookRoomHandler) HandlerName() string {
	return "BookRoomHandler"
}

// NewCommand returns type of command which this handle should handle. It must be a pointer.
func (b BookRoomHandler) NewCommand() interface{} {
	return &BookRoom{}
}

func (b BookRoomHandler) Handle(ctx context.Context, c interface{}) error {
	// c is always the type returned by `NewCommand`, so casting is always safe
	cmd := c.(*BookRoom)

	fmt.Println("handle message: ", cmd)

	return nil
}
