package deploy

import (
	"context"
	"fmt"

	"github.com/Peripli/itest-tools/docker"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

type Deployer struct {
	builders    []DockerBuilder
	runnables   []Runnable
	deployments map[string]DeploymentResult

	dockerClient *docker.Docker
}

func NewDeployer(dockerClient *docker.Docker) *Deployer {
	return &Deployer{
		dockerClient: dockerClient,
		deployments:  make(map[string]DeploymentResult),
	}
}

type DeploymentResult struct {
	URL string
}

type DockerRunOptions struct {
	Run bool
	*container.Config
	*container.HostConfig
	*network.NetworkingConfig
	ContainerName string
}

type MergeConfigFunc func(currOptions DockerRunOptions, dependencies map[string]DeploymentResult) DockerRunOptions

func (d *Deployer) AddDockerRun(name string, dependencies []string, options DockerRunOptions, merge MergeConfigFunc) {
	finalOptions := merge(options, d.deployments)

	ctx := context.Background()
	if finalOptions.Run {
		d.dockerClient.ContainerCreate(ctx,
			finalOptions.Config,
			finalOptions.HostConfig,
			finalOptions.NetworkingConfig,
			finalOptions.ContainerName)
	}

	d.deployments[name] = DeploymentResult{
		URL: "asd",
	}
}

func (d *Deployer) AddBuilder(b DockerBuilder) {
	d.builders = append(d.builders, b)
}

func (d *Deployer) AddRunnable(r Runnable) {
	d.runnables = append(d.runnables, r)
}

func (d *Deployer) Run() error {
	for _, r := range d.runnables {
		deployment, ok := r.(Deployment)
		if ok {
			fmt.Printf("Starting %s...\n", deployment.Name())
		}
		err := r.Run()
		if err != nil {
			return err
		}
		if ok {
			fmt.Printf("%s is running and accessible at %s\n", deployment.Name(), deployment.URI())
		}
	}
	return nil
}

func (d *Deployer) Build() error {
	for _, b := range d.builders {
		err := b.Build()
		if err != nil {
			return err
		}
	}
	return nil
}
