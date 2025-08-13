package system

import (
	"cayman"
	syssse "cayman/internal/sse"
	"context"

	"github.com/labstack/echo/v4"
	"github.com/tmaxmax/go-sse"
)

var (
	// compile time check for Module interface
	_       cayman.Module = (*SystemModule)(nil)
	lModule *SystemModule
)
var topicHost = "system"

func init() {
	lModule = &SystemModule{}
	cayman.RegisterModule(lModule)
}

type SystemModule struct {
	ctx context.Context
	sse *sse.Server
}

func (p *SystemModule) ShouldEnable() bool {
	// Logic to determine if the Logs module should be enabled
	return true
}
func (p *SystemModule) RegisterRoutes(ctx context.Context, parentRoute *echo.Group) {
	p.ctx = ctx
	p.sse = syssse.NewSSE(topicHost)
	// Register Logs-specific routes here
	routeGroup := parentRoute.Group("/system")
	go p.Poll()
	routeGroup.GET("/events", echo.WrapHandler(p.sse))
	routeGroup.GET("/current", p.systemInfoHandler)
}
func (p *SystemModule) Topics() []string {
	return []string{"system"}
}

func (p *SystemModule) Name() string {
	return "System"
}

func (p *SystemModule) Poll() {
	// Logic to poll System for updates
}
func (p *SystemModule) systemInfoHandler(c echo.Context) error {
	// Logic to handle system info requests
	return nil
}
