package host

import (
	"cayman"
	"context"

	"github.com/labstack/echo/v4"
	"github.com/tmaxmax/go-sse"
)

var (
	// compile time check for Module interface
	_       cayman.Module = (*HostModule)(nil)
	hModule *HostModule
)

func init() {
	hModule = &HostModule{}
	cayman.RegisterModule(hModule)
}

type HostModule struct {
	ctx context.Context
	sse *sse.Server
}

func (p *HostModule) ShouldEnable() bool {
	// Logic to determine if the Podman module should be enabled
	return true
}
func (p *HostModule) RegisterRoutes(ctx context.Context, parentRoute *echo.Group) {
	p.ctx = ctx
	p.sse = newSSE()
	// Register Podman-specific routes here
	routeGroup := parentRoute.Group("/host")
	go p.Poll()
	routeGroup.GET("/events", echo.WrapHandler(p.sse))
	routeGroup.GET("/current", p.hostInfoHandler)
}
func (p *HostModule) Topics() []string {
	return []string{"host"}
}

func (p *HostModule) Name() string {
	return "Host"
}

func (p *HostModule) Poll() {
	// Logic to poll Host for updates
}
func (p *HostModule) hostInfoHandler(c echo.Context) error {
	// Logic to handle host info requests
	return nil
}
