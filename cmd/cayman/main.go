package main

import (
	"cayman/internal/modules"
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

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
	slog.SetDefault(logger)

	// app := echo.New()
	// app.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
	// 	LogStatus:   true,
	// 	LogURI:      true,
	// 	LogError:    true,
	// 	HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
	// 	LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
	// 		if v.Error == nil {
	// 			logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
	// 				slog.String("uri", v.URI),
	// 				slog.Int("status", v.Status),
	// 			)
	// 		} else {
	// 			logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
	// 				slog.String("uri", v.URI),
	// 				slog.Int("status", v.Status),
	// 				slog.String("err", v.Error.Error()),
	// 			)
	// 		}
	// 		return nil
	// 	},
	// }))
	// app.Use(middleware.Recover())
	// app.Use(middleware.StaticWithConfig(middleware.StaticConfig{
	// 	Filesystem: frontend.BuildHTTPFS(),
	// 	HTML5:      true,
	// }))

	// /// TODO: make sure this is correct/safe
	// app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins: []string{"*"},
	// 	AllowMethods: []string{"*"},
	// }))

	// api := app.Group("/api")

	// systemSSEHandler := system.GetSSEServer()
	// api.GET("/systemevents", echo.WrapHandler(systemSSEHandler))

	// // register api and sse routes for modules here
	// dashboard.RegisterRoutes(ctx, api)

	// routes := app.Routes()
	// logger.Info("registered routes", "count", len(routes))
	// for _, route := range routes {
	// 	logger.Info("route", "method", route.Method, "path", route.Path)
	// }

	// listenAddr := net.JoinHostPort(*addr, *port)
	// logger.Info("starting server", "address", listenAddr)

	// s := &http.Server{
	// 	Addr:              listenAddr,
	// 	ReadHeaderTimeout: time.Second * 10,
	// 	Handler:           app,
	// }
	// s.RegisterOnShutdown(func() {
	// 	e := &sse.Message{Type: sse.Type("close")}
	// 	// Adding data is necessary because spec-compliant clients
	// 	// do not dispatch events without data.
	// 	e.AppendData("bye")
	// 	// Broadcast a close message so clients can gracefully disconnect.
	// 	_ = systemSSEHandler.Publish(e)

	// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	// 	defer cancel()

	// 	// We use a context with a timeout so the program doesn't wait indefinitely
	// 	// for connections to terminate. There may be misbehaving connections
	// 	// which may hang for an unknown timespan, so we just stop waiting on Shutdown
	// 	// after a certain duration.
	// 	_ = systemSSEHandler.Shutdown(ctx)
	// })

	// if err := runServer(ctx, s); err != nil {
	// 	log.Println("server closed", err)
	// }

	engine := modules.NewEngine(logger, *addr, *port)
	if err := engine.Start(ctx); err != nil {
		logger.Error("failed to start engine", "error", err)
	}
}

// func runServer(ctx context.Context, s *http.Server) error {
// 	shutdownError := make(chan error)

// 	go func() {
// 		<-ctx.Done()

// 		sctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 		defer cancel()
// 		slog.Info("starting graceful shutdown")
// 		shutdownError <- s.Shutdown(sctx)
// 	}()

// 	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
// 		return err
// 	}

// 	return <-shutdownError
// }
