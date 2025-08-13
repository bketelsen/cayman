package storage

import (
	"cayman"
	syssse "cayman/internal/sse"
	"context"

	"github.com/labstack/echo/v4"
	"github.com/tmaxmax/go-sse"
)

var (
	// compile time check for Module interface
	_       cayman.Module = (*StorageModule)(nil)
	lModule *StorageModule
)
var topicHost = "storage"

func init() {
	lModule = &StorageModule{}
	cayman.RegisterModule(lModule)
}

type StorageModule struct {
	ctx context.Context
	sse *sse.Server
}

func (p *StorageModule) ShouldEnable() bool {
	// Logic to determine if the Logs module should be enabled
	return true
}
func (p *StorageModule) RegisterRoutes(ctx context.Context, parentRoute *echo.Group) {
	p.ctx = ctx
	p.sse = syssse.NewSSE(topicHost)
	// Register Logs-specific routes here
	routeGroup := parentRoute.Group("/storage")
	go p.Poll()
	routeGroup.GET("/events", echo.WrapHandler(p.sse))
	routeGroup.GET("/current", p.storageInfoHandler)
}
func (p *StorageModule) Topics() []string {
	return []string{"storage"}
}

func (p *StorageModule) Name() string {
	return "Storage"
}

func (p *StorageModule) Poll() {
	// Logic to poll Storage for updates
}
func (p *StorageModule) storageInfoHandler(c echo.Context) error {
	// Logic to handle storage info requests
	return nil
}
