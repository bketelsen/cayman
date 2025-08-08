package dashboard

import (
	"cayman/internal/data/hardware"
	"cayman/internal/data/system"
	"cayman/internal/data/systemd"
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/elastic/go-sysinfo/types"
	"github.com/labstack/echo/v4"
	"github.com/tmaxmax/go-sse"
)

type HostHandler struct {
	ctx        context.Context
	sseHandler *sse.Server
	info       *HostState
}

func (h *HostHandler) Poll() {
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

func (h *HostHandler) usage() {
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

func (h *HostHandler) stats() {
	sysinfo, err := system.HostInfo()
	if err != nil {
		slog.Error("failed to get host info", "error", err)
	}
	mem, err := sysinfo.Memory()
	if err != nil {
		slog.Error("failed to get memory info", "error", err)
	}
	var tmpLoad Load
	if loadaverage, ok := sysinfo.(types.LoadAverage); ok {
		loadavg, err := loadaverage.LoadAverage()
		if err != nil {
			slog.Error("failed to get load", "error", err)

		}
		tmpLoad = Load{
			Load1:  loadavg.One,
			Load5:  loadavg.Five,
			Load15: loadavg.Fifteen,
		}
	}
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

func NewHostHandler(ctx context.Context) *HostHandler {

	cpustat, err := hardware.Info()
	if err != nil {
		slog.Error("failed to get cpu info", "error", err)
	}
	slog.Info("CPU Info", "info", cpustat, "count", len(cpustat))
	failed, active, _ := systemd.UnitOverview(ctx)

	sysinfo, err := system.HostInfo()
	if err != nil {
		slog.Error("failed to get host info", "error", err)
	}
	mem, err := sysinfo.Memory()
	if err != nil {
		slog.Error("failed to get memory info", "error", err)
	}
	var tmpLoad Load
	if loadaverage, ok := sysinfo.(types.LoadAverage); ok {
		loadavg, err := loadaverage.LoadAverage()
		if err != nil {
			slog.Error("failed to get load", "error", err)

		}
		tmpLoad = Load{
			Load1:  loadavg.One,
			Load5:  loadavg.Five,
			Load15: loadavg.Fifteen,
		}
	}
	domain, err := sysinfo.FQDNWithContext(ctx)
	if err != nil {
		slog.Error("failed to get FQDN", "error", err)
	}
	hi := &HostState{
		FQDN: domain,
		UnitStatus: UnitStatus{
			FailedCount: failed,
			ActiveCount: active,
		},

		Load:       tmpLoad,
		HostInfo:   sysinfo.Info(),
		MemoryInfo: *mem,
	}
	return &HostHandler{
		ctx:        ctx,
		sseHandler: newSSE(),
		info:       hi,
	}
}

func RegisterRoutes(ctx context.Context, parentRoute *echo.Group) {
	routeGroup := parentRoute.Group("/host")
	hostHandler := NewHostHandler(ctx)
	go hostHandler.Poll()
	routeGroup.GET("/events", echo.WrapHandler(hostHandler.sseHandler))
	routeGroup.GET("/current", hostHandler.hostInfoHandler)
}

func (h *HostHandler) hostInfoHandler(c echo.Context) error {

	return c.JSON(200, h.info)
}
