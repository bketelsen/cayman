package cayman

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
)

type DockerInfo struct {
	Containers []container.Summary `json:"containers"`
	Images     []image.Summary     `json:"images"`
}
