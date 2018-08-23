package broker

import (
	"context"
	"fmt"

	"github.com/Peripli/itest-tools/deploy"
	"github.com/Peripli/itest-tools/deploy/config"
	"github.com/Peripli/itest-tools/docker"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

type Broker interface {
	deploy.Deployment
	deploy.Runnable
}

type BrokerImpl struct {
	settings     *config.Settings
	dockerClient *docker.Docker
}

var _ Broker = &BrokerImpl{}

func New(settings *config.Settings,
	dockerClient *docker.Docker) Broker {
	return &BrokerImpl{
		settings:     settings,
		dockerClient: dockerClient,
	}
}

func (b *BrokerImpl) Run() error {
	var err error
	if b.settings.Broker.Run {
		ctx := context.Background()
		portSet := nat.PortSet{}
		portSet["8080/tcp"] = struct{}{}

		brContainer, err := b.dockerClient.ContainerCreate(ctx,
			&container.Config{
				Image:        "sbf-broker:test",
				ExposedPorts: portSet,
			},
			nil,
			nil,
			b.Name())

		if err != nil {
			return err
		}

		err = b.dockerClient.NetworkConnect(ctx, b.settings.Docker.NetworkID, brContainer.ID, nil)
		if err != nil {
			return err
		}

		err = b.dockerClient.ContainerStart(ctx, brContainer.ID, types.ContainerStartOptions{})
		if err != nil {
			return err
		}

		b.settings.Broker.URL = fmt.Sprintf("http://%s:8080", b.Name())
	}
	return err
}

func (b *BrokerImpl) Name() string {
	return "broker-" + b.settings.Test.ID
}

func (b *BrokerImpl) URI() string {
	return b.settings.Broker.URL
}
