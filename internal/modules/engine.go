package modules

import (
	"cayman"
	"cayman/frontend"
	_ "cayman/internal/modules/dashboard"
	_ "cayman/internal/modules/docker"
	_ "cayman/internal/modules/podman"

	"cayman/internal/system"
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tmaxmax/go-sse"
)

type Engine struct {
	logger     *slog.Logger
	httpServer *http.Server
	listenAddr string
	port       string
}

func NewEngine(logger *slog.Logger, listenAddr string, port string) *Engine {
	return &Engine{
		logger:     logger,
		listenAddr: listenAddr,
		port:       port,
	}
}

func (e *Engine) Start(ctx context.Context) error {
	e.logger.Info("starting cayman engine")

	app := echo.New()
	app.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				slog.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				slog.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
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

	systemSSEHandler := system.GetSSEServer()
	api.GET("/systemevents", echo.WrapHandler(systemSSEHandler))
	// register api and sse routes for modules here
	//dashboard.RegisterRoutes(ctx, api)
	slog.Info("registering modules")

	slog.Info("available modules", "count", len(cayman.AvailableModules))
	for _, m := range cayman.AvailableModules {
		slog.Info("checking module", "name", m.Name())
		if m.ShouldEnable() {
			cayman.EnabledModules = append(cayman.EnabledModules, m)
		}
	}
	slog.Info("enabled modules", "count", len(cayman.EnabledModules))
	for _, m := range cayman.EnabledModules {
		slog.Info("enabling module", "name", m.Name())
		m.RegisterRoutes(ctx, api)
	}

	routes := app.Routes()
	slog.Info("registered routes", "count", len(routes))
	for _, route := range routes {
		slog.Info("route", "method", route.Method, "path", route.Path)
	}
	listenAddr := net.JoinHostPort(e.listenAddr, e.port)
	slog.Info("starting server", "address", listenAddr)

	e.httpServer = &http.Server{
		Addr:              listenAddr,
		ReadHeaderTimeout: time.Second * 10,
		Handler:           app,
	}
	e.httpServer.RegisterOnShutdown(func() {
		e := &sse.Message{Type: sse.Type("close")}
		// Adding data is necessary because spec-compliant clients
		// do not dispatch events without data.
		e.AppendData("bye")
		// Broadcast a close message so clients can gracefully disconnect.
		_ = systemSSEHandler.Publish(e)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		// We use a context with a timeout so the program doesn't wait indefinitely
		// for connections to terminate. There may be misbehaving connections
		// which may hang for an unknown timespan, so we just stop waiting on Shutdown
		// after a certain duration.
		_ = systemSSEHandler.Shutdown(ctx)
	})

	shutdownError := make(chan error)

	go func() {
		<-ctx.Done()

		sctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		slog.Info("starting graceful shutdown")
		shutdownError <- e.httpServer.Shutdown(sctx)
	}()

	if err := e.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return <-shutdownError
}
