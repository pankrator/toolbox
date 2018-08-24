package orchestrator

import (
	"github.com/Peripli/itest-tools/deploy"
	"github.com/docker/docker/api/types/container"
)

type DeployTarget string

var (
	CFTarget     DeployTarget = "cf"
	DockerTarget DeployTarget = "docker"
	K8STarget    DeployTarget = "k8s"
)

type Orchestrator struct {
	dockerDeployer *deploy.Deployer
}

type Config struct {
	Target     DeployTarget
	Name       string
	Image      string
	CFManifest string
}

func ()

func (o *Orchestrator) Add(c *Config) {

	// o.dockerDeployer.AddDockerRun(name string, dependencies []string, options deploy.DockerRunOptions, merge deploy.MergeConfigFunc)
}
