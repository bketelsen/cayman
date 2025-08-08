package system

import (
	"log/slog"

	"github.com/tmaxmax/go-sse"
)

const (
	// TopicSystem is the topic for system-related events.
	TopicSystem = "system"
)

// enum SystemEventType represents the type of system events.
type SystemEventType string

const (
	// SystemEventTypeInfo is the event type for informational system events.
	SystemEventTypeInfo SystemEventType = "systeminfo"
	// SystemEventTypeWarning is the event type for warning system events.
	SystemEventTypeWarning SystemEventType = "systemwarning"
	// SystemEventTypeError is the event type for system errors.
	SystemEventTypeError SystemEventType = "systemerror"
)

// SystemHandler handles system-related events and publishes them to the SSE server.
type SystemHandler struct {
	sse *sse.Server
}

// NewSystemHandler creates a new SystemHandler with the given SSE server.
func NewSystemHandler(sse *sse.Server) *SystemHandler {
	return &SystemHandler{
		sse: sse,
	}
}

// PublishSystemEvent handles a system event and sends it to the SSE server.
func (h *SystemHandler) PublishSystemEvent(eventType SystemEventType, data string) {
	slog.Info("sending system event", "type", eventType, "data", data)

	if h.sse == nil {
		return
	}
	e := &sse.Message{
		Type: sse.Type(string(eventType)),
	}
	e.AppendData(data)

	_ = h.sse.Publish(e, TopicSystem)
}
