package docker

import (
	"cayman"
	"context"
	"log/slog"

	"github.com/docker/docker/client"
	"github.com/labstack/echo/v4"
	"github.com/tmaxmax/go-sse"
)

var (
	// compile time check for Module interface
	_       cayman.Module = (*DockerModule)(nil)
	dModule *DockerModule
)

func init() {
	dModule = &DockerModule{}
	cayman.RegisterModule(dModule)
}

type DockerModule struct {
	ctx context.Context
	sse *sse.Server
}

func (p *DockerModule) ShouldEnable() bool {
	// Logic to determine if the Docker module should be enabled
	// TODO: Implement logic to determine if the Podman module should be enabled
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		slog.Error("failed to connect to docker socket", "error", err)
		return false
	}
	defer apiClient.Close()
	slog.Info("docker socket found, enabling docker module")
	return true
}
func (p *DockerModule) RegisterRoutes(ctx context.Context, parentRoute *echo.Group) {
	p.ctx = ctx
	p.sse = newSSE()
	// Register Docker-specific routes here
	routeGroup := parentRoute.Group("/virt/docker")
	go p.Poll()
	routeGroup.GET("/events", echo.WrapHandler(p.sse))
	routeGroup.GET("/current", p.dockerInfoHandler)
}
func (p *DockerModule) Topics() []string {
	return []string{"docker"}
}

func (p *DockerModule) Name() string {
	return "Docker"
}

func (p *DockerModule) Poll() {
	// Logic to poll Docker for updates
}
func (p *DockerModule) dockerInfoHandler(c echo.Context) error {
	// Logic to handle docker info requests
	return nil
}
