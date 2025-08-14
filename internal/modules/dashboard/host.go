package dashboard

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"cayman"
	"cayman/internal/data/hardware"
	"cayman/internal/data/system"
	"cayman/internal/data/systemd"
	syssse "cayman/internal/sse"

	"github.com/elastic/go-sysinfo/types"
	"github.com/labstack/echo/v4"
	"github.com/tmaxmax/go-sse"
)

var (
	_          cayman.Module = (*DashboardModule)(nil)
	dashModule *DashboardModule
	topicHost  = "dashboard"
)

func init() {
	dashModule = &DashboardModule{}
	cayman.RegisterModule(dashModule)
}

type DashboardModule struct {
	ctx        context.Context
	sseHandler *sse.Server
	info       *cayman.HostState
}

func (h *DashboardModule) ShouldEnable() bool {
	return true
}

func (h *DashboardModule) Topics() []string {
	return []string{topicHost}
}

func (h *DashboardModule) Name() string {
	return "Dashboard"
}

func (h *DashboardModule) RegisterRoutes(ctx context.Context, parentRoute *echo.Group) {
	h.ctx = ctx
	h.sseHandler = syssse.NewSSE(topicHost)
	routeGroup := parentRoute.Group("/dashboard")
	go h.Poll()
	routeGroup.GET("/events", echo.WrapHandler(h.sseHandler))
	routeGroup.GET("/current", h.hostInfoHandler)
	cpustat, err := hardware.Info()
	if err != nil {
		slog.Error("failed to get cpu info", "error", err)
	}
	slog.Info("CPU Info", "count", len(cpustat))
	failed, active, err := systemd.UnitOverview(ctx)
	if err != nil {
		slog.Error("failed to get systemd unit status", "error", err)
	}
	sysinfo, err := system.HostInfo()
	if err != nil {
		slog.Error("failed to get host info", "error", err)
	}
	mem, err := sysinfo.Memory()
	if err != nil {
		slog.Error("failed to get memory info", "error", err)
	}
	var tmpLoad cayman.Load
	if loadaverage, ok := sysinfo.(types.LoadAverage); ok {
		loadavg, err := loadaverage.LoadAverage()
		if err != nil {
			slog.Error("failed to get load", "error", err)
		}
		tmpLoad = cayman.Load{
			Load1:  loadavg.One,
			Load5:  loadavg.Five,
			Load15: loadavg.Fifteen,
		}
	}
	domain, err := sysinfo.FQDNWithContext(ctx)
	if err != nil {
		slog.Error("failed to get FQDN", "error", err)
	}
	hi := &cayman.HostState{
		FQDN:     domain,
		CPUCount: len(cpustat),
		UnitStatus: cayman.UnitStatus{
			FailedCount: failed,
			ActiveCount: active,
		},

		Load:       tmpLoad,
		HostInfo:   sysinfo.Info(),
		MemoryInfo: *mem,
	}
	h.info = hi
}

func (h *DashboardModule) Poll() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.usage()
			h.stats()
		case <-h.ctx.Done():
			return
		}
	}
}

func (h *DashboardModule) usage() {
	usage, err := hardware.CPUUsage(h.ctx)
	if err != nil {
		slog.Error("failed to get cpu usage", "error", err)
	}

	e := &sse.Message{
		Type: sse.Type("cpu"),
	}
	bb, err := json.Marshal(int(usage))
	if err != nil {
		return
	}
	e.AppendData(string(bb))
	_ = h.sseHandler.Publish(e, topicHost)
}

func (h *DashboardModule) stats() {
	sysinfo, err := system.HostInfo()
	if err != nil {
		slog.Error("failed to get host info", "error", err)
	}
	h.info.HostInfo = sysinfo.Info()
	mem, err := sysinfo.Memory()
	if err != nil {
		slog.Error("failed to get memory info", "error", err)
	}
	h.info.MemoryInfo = *mem
	var tmpLoad cayman.Load
	if loadaverage, ok := sysinfo.(types.LoadAverage); ok {
		loadavg, err := loadaverage.LoadAverage()
		if err != nil {
			slog.Error("failed to get load", "error", err)
		}
		tmpLoad = cayman.Load{
			Load1:  loadavg.One,
			Load5:  loadavg.Five,
			Load15: loadavg.Fifteen,
		}
	}
	h.info.Load = tmpLoad
	e := &sse.Message{
		Type: sse.Type("mem"),
	}
	bb, err := json.Marshal(mem)
	if err != nil {
		return
	}
	e.AppendData(string(bb))

	_ = h.sseHandler.Publish(e, topicHost)

	e = &sse.Message{
		Type: sse.Type("load"),
	}
	bb, err = json.Marshal(tmpLoad)
	if err != nil {
		return
	}
	e.AppendData(string(bb))

	_ = h.sseHandler.Publish(e, topicHost)
}

func (h *DashboardModule) hostInfoHandler(c echo.Context) error {
	return c.JSON(200, h.info)
}
