package system

import (
	"github.com/elastic/go-sysinfo"
	"github.com/elastic/go-sysinfo/types"
)

func HostInfo() (types.Host, error) {
	return sysinfo.Host()
}
