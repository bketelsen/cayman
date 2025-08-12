package cayman

import (
	"context"
	"log/slog"

	"github.com/labstack/echo/v4"
)

var (
	AvailableModules = make([]Module, 0)
	EnabledModules   = make([]Module, 0)
)

type Module interface {
	ShouldEnable() bool
	RegisterRoutes(ctx context.Context, parentRoute *echo.Group)
	Topics() []string
	Name() string
}

func RegisterModule(m Module) {
	slog.Info("registering module", "name", m.Name())
	AvailableModules = append(AvailableModules, m)
}

// func RegisterRoutes(ctx context.Context, parentRoute *echo.Group) {
// 	for _, module := range AvailableModules {
// 		module.RegisterRoutes(ctx, parentRoute)
// 	}
// }
