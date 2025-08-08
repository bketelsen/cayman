package dashboard

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/tmaxmax/go-sse"
)

var topicHost = "host"

func newSSE() *sse.Server {
	rp, _ := sse.NewValidReplayer(time.Minute*5, true)
	rp.GCInterval = time.Minute

	return &sse.Server{
		Provider: &sse.Joe{Replayer: rp},
		// If you are using a 3rd party library to generate a per-request logger, this
		// can just be a simple wrapper over it.
		Logger: func(r *http.Request) *slog.Logger {
			handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
			logger := slog.New(handler)
			return logger
		},
		OnSession: func(w http.ResponseWriter, r *http.Request) (topics []string, permitted bool) {
			topics = []string{}

			// add system topic
			topics = append(topics, topicHost)

			// the shutdown message is sent on the default topic
			return append(topics, sse.DefaultTopic), true
		},
	}
}
