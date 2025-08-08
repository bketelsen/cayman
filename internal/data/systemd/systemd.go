package systemd

import (
	"context"
	"log/slog"

	"github.com/coreos/go-systemd/v22/dbus"
)

func UnitOverview(ctx context.Context) (failed, active int, err error) {
	dbusConn, err := dbus.NewWithContext(ctx)
	if err != nil {
		slog.Error("Failed to connect to systemd", "err", err)
		return
	}
	defer dbusConn.Close()

	uu, err := dbusConn.ListUnitsContext(ctx)
	if err != nil {
		slog.Error("Failed to list units", "err", err)
		return
	}
	failed = 0
	active = 0

	for _, u := range uu {

		if u.ActiveState == "failed" || u.SubState == "failed" {
			failed++
		}
		if u.ActiveState == "active" || u.SubState == "active" {
			active++
		}
	}
	return
}
