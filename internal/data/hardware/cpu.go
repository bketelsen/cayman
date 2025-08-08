package hardware

import (
	"context"

	"github.com/shirou/gopsutil/v4/cpu"
	// "github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	// "github.com/shirou/gopsutil/v4/mem"
)

func Load(_ context.Context) (*load.AvgStat, error) {
	return load.Avg()
}

func CPUUsage(_ context.Context) (float64, error) {
	uu, err := cpu.Percent(0, false)
	if err != nil {
		return 0, err
	}
	return uu[0], nil
}

func Info() ([]cpu.InfoStat, error) {
	return cpu.Info()
}
