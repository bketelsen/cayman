package incus

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"cayman"
	syssse "cayman/internal/sse"

	"github.com/lxc/incus/v6/shared/api"
	config "github.com/lxc/incus/v6/shared/cliconfig"

	"github.com/bketelsen/inclient"
	"github.com/labstack/echo/v4"
	"github.com/tmaxmax/go-sse"
)

var (
	// compile time check for Module interface
	_         cayman.Module = (*IncusModule)(nil)
	iModule   *IncusModule
	topicHost = "incus"
)

func init() {
	iModule = &IncusModule{}
	cayman.RegisterModule(iModule)
}

type IncusModule struct {
	ctx context.Context
	sse *sse.Server
}

func (p *IncusModule) ShouldEnable() bool {
	cfg, err := config.LoadConfig("")
	if err != nil {
		return false
	}
	client, err := inclient.NewClient(cfg)
	if err != nil {
		return false
	}
	_, err = client.Instances(context.Background())
	return err == nil
}

func (p *IncusModule) RegisterRoutes(ctx context.Context, parentRoute *echo.Group) {
	p.ctx = ctx
	p.sse = syssse.NewSSE(topicHost)
	// Register Incus-specific routes here
	routeGroup := parentRoute.Group("/virt/incus")
	go p.Poll()
	routeGroup.GET("/events", echo.WrapHandler(p.sse))
	routeGroup.GET("/current", p.incusInfoHandler)
}

func (p *IncusModule) Topics() []string {
	return []string{"incus"}
}

func (p *IncusModule) Name() string {
	return "Incus"
}

func (p *IncusModule) Poll() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker.C:
			info, err := getIncusInfo()
			if err != nil {
				slog.Error("incus poll error", "error", err)
				continue
			}
			cc, err := json.Marshal(info.Instances)
			if err != nil {
				slog.Error("incus marshal error", "error", err)
				continue
			}
			event := &sse.Message{
				Type: sse.Type("instances"),
			}
			event.AppendData(string(cc))
			p.sse.Publish(event, topicHost)
			// ii, err := json.Marshal(info.Images)
			// if err != nil {
			// 	slog.Error("docker marshal error", "error", err)
			// 	continue
			// }
			// event = &sse.Message{
			// 	Type: sse.Type("images"),
			// }
			// event.AppendData(string(ii))
			// p.sse.Publish(event, topicHost)

		}
	}
}

func (p *IncusModule) incusInfoHandler(c echo.Context) error {
	info, err := getIncusInfo()
	if err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}
	return c.JSON(200, info)
}

func getIncusInfo() (*cayman.IncusInfo, error) {
	cfg, err := config.LoadConfig("")
	if err != nil {
		return nil, err
	}
	client, err := inclient.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	ii, err := client.Instances(context.Background())
	if err != nil {
		return nil, err
	}

	return &cayman.IncusInfo{
		Instances: ii,
		Images:    []api.Image{},
	}, nil
}
