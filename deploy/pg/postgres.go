package pg

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/Peripli/itest-tools/deploy"
	"github.com/Peripli/itest-tools/deploy/config"
	"github.com/Peripli/itest-tools/docker"
)

type Postgres interface {
	deploy.Deployment
	deploy.Runnable
}

type PostgresImpl struct {
	settings     *config.Settings
	dockerClient *docker.Docker
}

var _ Postgres = &PostgresImpl{}

func New(settings *config.Settings,
	dockerClient *docker.Docker) Postgres {
	return &PostgresImpl{
		settings:     settings,
		dockerClient: dockerClient,
	}
}

func (p *PostgresImpl) Name() string {
	return "postgres-" + p.settings.Test.ID
}

func (p *PostgresImpl) Run() error {
	var err error

	portSet := nat.PortSet{}
	port := nat.Port(p.settings.PG.Port)
	portSet[port+"/tcp"] = struct{}{}

	if p.settings.PG.Run {
		ctx := context.Background()
		container, err := p.dockerClient.ContainerCreate(ctx,
			&container.Config{
				Image:        "postgres",
				ExposedPorts: portSet,
				User:         "postgres",
			},
			nil,
			nil,
			p.Name())

		if err != nil {
			return err
		}

		p.settings.PG.URI = fmt.Sprintf("postgres://postgres:postgres@%s:%s/postgres?sslmode=disable", p.Name(), p.settings.PG.Port)

		err = p.dockerClient.NetworkConnect(ctx, p.settings.Docker.NetworkID, container.ID, nil)
		if err != nil {
			return err
		}

		err = p.dockerClient.ContainerStart(ctx, container.ID, types.ContainerStartOptions{})
		if err != nil {
			return err
		}

		

		// TODO: Find a better way to check whether Postgres is running. May be try to connect to it
		fmt.Println("Wait 10 seconds for Postgres to start, otherwise SM will not be able to connect")
		time.Sleep(time.Second * 10)
	}

	return err
}

func (p *PostgresImpl) URI() string {
	return p.settings.PG.URI
}
