package logs

import (
	"cayman"
	syssse "cayman/internal/sse"
	"context"

	"github.com/labstack/echo/v4"
	"github.com/tmaxmax/go-sse"
)

var (
	// compile time check for Module interface
	_         cayman.Module = (*LogsModule)(nil)
	lModule   *LogsModule
	topicHost = "logs"
)

func init() {
	lModule = &LogsModule{}
	cayman.RegisterModule(lModule)
}

type LogsModule struct {
	ctx context.Context
	sse *sse.Server
}

func (p *LogsModule) ShouldEnable() bool {
	// Logic to determine if the Logs module should be enabled
	return true
}
func (p *LogsModule) RegisterRoutes(ctx context.Context, parentRoute *echo.Group) {
	p.ctx = ctx
	p.sse = syssse.NewSSE(topicHost)
	// Register Logs-specific routes here
	routeGroup := parentRoute.Group("/logs")
	go p.Poll()
	routeGroup.GET("/events", echo.WrapHandler(p.sse))
	routeGroup.GET("/current", p.logsInfoHandler)
}
func (p *LogsModule) Topics() []string {
	return []string{"logs"}
}

func (p *LogsModule) Name() string {
	return "Logs"
}

func (p *LogsModule) Poll() {
	// Logic to poll Logs for updates
}
func (p *LogsModule) logsInfoHandler(c echo.Context) error {
	// Logic to handle logs info requests
	return nil
}
