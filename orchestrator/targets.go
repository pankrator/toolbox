package orchestrator

import (
	"github.com/Peripli/itest-tools/deploy"
	"github.com/docker/docker/api/types/container"
)

func buildDockerConfig(c *Config) *deploy.DockerRunOptions {
	return &deploy.DockerRunOptions{
		Config: &container.Config{
			Image: c.Image,
			
		},
	}
}
