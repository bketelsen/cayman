package docker

import (
	"cayman"
	"context"
	"encoding/json"
	"log/slog"
	"sort"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
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
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker.C:
			info, err := getDockerInfo()
			if err != nil {
				slog.Error("docker poll error", "error", err)
				continue
			}
			cc, err := json.Marshal(info.Containers)
			if err != nil {
				slog.Error("docker marshal error", "error", err)
				continue
			}
			event := &sse.Message{
				Type: sse.Type("containers"),
			}
			event.AppendData(string(cc))
			p.sse.Publish(event, topicHost)
			ii, err := json.Marshal(info.Images)
			if err != nil {
				slog.Error("docker marshal error", "error", err)
				continue
			}
			event = &sse.Message{
				Type: sse.Type("images"),
			}
			event.AppendData(string(ii))
			p.sse.Publish(event, topicHost)
		}
	}
}
func (p *DockerModule) dockerInfoHandler(c echo.Context) error {
	info, err := getDockerInfo()
	if err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}
	return c.JSON(200, info)
}

func getDockerInfo() (*cayman.DockerInfo, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	// Sort containers by Created field in descending order (newest first)
	sort.Slice(containers, func(i, j int) bool {
		return containers[i].Created > containers[j].Created
	})

	images, err := cli.ImageList(context.Background(), image.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	// Sort images by Created field in descending order (newest first)
	sort.Slice(images, func(i, j int) bool {
		return images[i].Created > images[j].Created
	})

	return &cayman.DockerInfo{
		Containers: containers,
		Images:     images,
	}, nil
}
