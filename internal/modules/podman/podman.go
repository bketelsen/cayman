package podman

import (
	"cayman"
	syssse "cayman/internal/sse"
	"context"
	"log/slog"

	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/labstack/echo/v4"
	"github.com/tmaxmax/go-sse"
)

var (
	// compile time check for Module interface
	_       cayman.Module = (*PodmanModule)(nil)
	pModule *PodmanModule
)

func init() {
	pModule = &PodmanModule{}
	cayman.RegisterModule(pModule)
}

type PodmanModule struct {
	ctx context.Context
	sse *sse.Server
}

func (p *PodmanModule) ShouldEnable() bool {
	// Logic to determine if the Podman module should be enabled
	// TODO: Implement logic to determine if the Podman module should be enabled
	_, err := bindings.NewConnection(context.Background(), "unix:///run/user/1000/podman/podman.sock")
	if err != nil {
		slog.Error("failed to connect to podman socket", "error", err)
		return false
	}
	slog.Info("podman socket found, enabling podman module")
	return true
}
func (p *PodmanModule) RegisterRoutes(ctx context.Context, parentRoute *echo.Group) {
	p.ctx = ctx
	p.sse = syssse.NewSSE("podman")
	// Register Podman-specific routes here
	routeGroup := parentRoute.Group("/virt/podman")
	go p.Poll()
	routeGroup.GET("/events", echo.WrapHandler(p.sse))
	routeGroup.GET("/current", p.podmanInfoHandler)
}
func (p *PodmanModule) Topics() []string {
	return []string{"podman"}
}

func (p *PodmanModule) Name() string {
	return "Podman"
}

func (p *PodmanModule) Poll() {
	// Logic to poll Podman for updates
}
func (p *PodmanModule) podmanInfoHandler(c echo.Context) error {
	// Logic to handle podman info requests
	return nil
}
