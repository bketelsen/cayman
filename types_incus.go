package cayman

import (
	"github.com/lxc/incus/v6/shared/api"
)

type IncusInfo struct {
	Instances []api.InstanceFull `json:"instances"`
	Images    []api.Image        `json:"images"`
}
