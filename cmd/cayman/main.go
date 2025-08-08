package main

import (
	"cayman/frontend"
	"cayman/internal/modules/dashboard"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tmaxmax/go-sse"
)

const (
	topicRandomNumbers = "numbers"
	topicMetrics       = "metrics"
	topicSystem        = "system"
)

func newSSE() *sse.Server {
	rp, _ := sse.NewValidReplayer(time.Minute*5, true)
	rp.GCInterval = time.Minute

	return &sse.Server{
		Provider: &sse.Joe{Replayer: rp},
		// If you are using a 3rd party library to generate a per-request logger, this
		// can just be a simple wrapper over it.
		Logger: func(r *http.Request) *slog.Logger {
			return getLogger(r.Context())
		},
		OnSession: func(w http.ResponseWriter, r *http.Request) (topics []string, permitted bool) {
			topics = r.URL.Query()["topic"]
			slog.Info("new session", "topics", topics)
			for _, topic := range topics {
				if topic != topicRandomNumbers && topic != topicMetrics {
					fmt.Fprintf(w, "invalid topic %q; supported are %q, %q", topic, topicRandomNumbers, topicMetrics)

					// NOTE: if you are returning false to reject the subscription, we strongly recommend writing
					// your own response code. Clients will receive a 200 code otherwise, which may be confusing.
					w.WriteHeader(http.StatusBadRequest)
					return nil, false
				}
			}
			if len(topics) == 0 {
				// Provide default topics, if none are given.
				topics = []string{topicRandomNumbers}
			}
			// add system topic
			topics = append(topics, topicSystem)

			// the shutdown message is sent on the default topic
			return append(topics, sse.DefaultTopic), true
		},
	}
}

func main() {
	var (
		addr = flag.String("addr", "0.0.0.0", "listen address")
		port = flag.String("port", "8080", "listen port")
	)
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(handler)

	sseHandler := newSSE()

	app := echo.New()
	app.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))
	app.Use(middleware.Recover())
	app.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Filesystem: frontend.BuildHTTPFS(),
		HTML5:      true,
	}))

	/// TODO: make sure this is correct/safe
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
	}))

	api := app.Group("/api")
	api.GET("/stop", func(c echo.Context) error {
		cancel()
		return c.NoContent(http.StatusOK)
	})

	// register api and sse routes for modules here
	dashboard.RegisterRoutes(ctx, api)

	routes := app.Routes()
	logger.Info("registered routes", "count", len(routes))
	for _, route := range routes {
		logger.Info("route", "method", route.Method, "path", route.Path)
	}

	listenAddr := net.JoinHostPort(*addr, *port)
	logger.Info("starting server", "address", listenAddr)

	s := &http.Server{
		Addr:              listenAddr,
		ReadHeaderTimeout: time.Second * 10,
		Handler:           app,
	}
	s.RegisterOnShutdown(func() {
		e := &sse.Message{Type: sse.Type("close")}
		// Adding data is necessary because spec-compliant clients
		// do not dispatch events without data.
		e.AppendData("bye")
		// Broadcast a close message so clients can gracefully disconnect.
		_ = sseHandler.Publish(e)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		// We use a context with a timeout so the program doesn't wait indefinitely
		// for connections to terminate. There may be misbehaving connections
		// which may hang for an unknown timespan, so we just stop waiting on Shutdown
		// after a certain duration.
		_ = sseHandler.Shutdown(ctx)
	})

	if err := runServer(ctx, s); err != nil {
		log.Println("server closed", err)
	}
}

func runServer(ctx context.Context, s *http.Server) error {
	shutdownError := make(chan error)

	go func() {
		<-ctx.Done()

		sctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		slog.Info("starting graceful shutdown")
		shutdownError <- s.Shutdown(sctx)
	}()

	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return <-shutdownError
}

type loggerCtxKey struct{}

// getLogger retrieves the request-specific logger from a request's context. This is
// similar to how existing per-request http logging libraries work, just very simplified.
func getLogger(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(loggerCtxKey{}).(*slog.Logger)
	if !ok {
		// We are accepting an arbitrary context object, so it's better to explicitly return
		// nil here since the exact behavior of getting the value of an undefined key is undefined
		return nil
	}
	return logger
}
