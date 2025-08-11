package system

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/tmaxmax/go-sse"
)

var (
	sseServer *sse.Server
)

const (
	topicSystem = "system"
)

func init() {
	rp, _ := sse.NewValidReplayer(time.Minute*5, true)
	rp.GCInterval = time.Minute

	sseServer = &sse.Server{
		Provider: &sse.Joe{Replayer: rp},
		// If you are using a 3rd party library to generate a per-request logger, this
		// can just be a simple wrapper over it.
		Logger: func(r *http.Request) *slog.Logger {
			handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
			logger := slog.New(handler)
			return logger
		},
		OnSession: func(w http.ResponseWriter, r *http.Request) (topics []string, permitted bool) {
			slog.Info("new root sse session")

			// add system topic
			topics = append(topics, topicSystem)

			// the shutdown message is sent on the default topic
			return append(topics, sse.DefaultTopic), true
		},
	}
}

func GetSSEServer() *sse.Server {
	if sseServer == nil {
		return nil
	}
	return sseServer
}

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
	// SystemEventTypeMessage is the event type for general system messages.
	SystemEventTypeMessage SystemEventType = "systemmessage"
)

// PublishSystemEvent handles a system event and sends it to the SSE server.
func PublishSystemEvent(eventType SystemEventType, data string) error {
	slog.Info("sending system event", "type", eventType, "data", data)

	if sseServer == nil {
		return errors.New("system sse server is not initialized")
	}
	e := &sse.Message{
		Type: sse.Type(string(eventType)),
	}
	e.AppendData(data)

	return sseServer.Publish(e, TopicSystem)
}
