package event

import (
	"errors"

	"clc.hmu/app/frp/frpc/models/msg"
)

type EventType int

const (
	EvStartProxy EventType = iota
	EvCloseProxy
)

var (
	ErrPayloadType = errors.New("error payload type")
)

type EventHandler func(evType EventType, payload interface{}) error

type StartProxyPayload struct {
	NewProxyMsg *msg.NewProxy
}

type CloseProxyPayload struct {
	CloseProxyMsg *msg.CloseProxy
}
