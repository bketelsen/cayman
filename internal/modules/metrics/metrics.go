package metrics

import (
	"cayman"
	syssse "cayman/internal/sse"
	"context"

	"github.com/labstack/echo/v4"
	"github.com/tmaxmax/go-sse"
)

var (
	// compile time check for Module interface
	_       cayman.Module = (*MetricsModule)(nil)
	lModule *MetricsModule
)
var topicHost = "metrics"

func init() {
	lModule = &MetricsModule{}
	cayman.RegisterModule(lModule)
}

type MetricsModule struct {
	ctx context.Context
	sse *sse.Server
}

func (p *MetricsModule) ShouldEnable() bool {
	// Logic to determine if the Logs module should be enabled
	return true
}
func (p *MetricsModule) RegisterRoutes(ctx context.Context, parentRoute *echo.Group) {
	p.ctx = ctx
	p.sse = syssse.NewSSE(topicHost)
	// Register Logs-specific routes here
	routeGroup := parentRoute.Group("/metrics")
	go p.Poll()
	routeGroup.GET("/events", echo.WrapHandler(p.sse))
	routeGroup.GET("/current", p.metricsInfoHandler)
}
func (p *MetricsModule) Topics() []string {
	return []string{"metrics"}
}

func (p *MetricsModule) Name() string {
	return "Metrics"
}

func (p *MetricsModule) Poll() {
	// Logic to poll Metrics for updates
}
func (p *MetricsModule) metricsInfoHandler(c echo.Context) error {
	// Logic to handle metrics info requests
	return nil
}
