package sm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Peripli/itest-tools/deploy/config"
	"github.com/Peripli/itest-tools/docker"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"
)

type ServiceManagerImpl struct {
	settings     *config.Settings
	dockerClient *docker.Docker
}

var _ ServiceManager = &ServiceManagerImpl{}

func New(settings *config.Settings,
	dockerClient *docker.Docker) ServiceManager {
	return &ServiceManagerImpl{
		settings:     settings,
		dockerClient: dockerClient,
	}
}

func (sm *ServiceManagerImpl) Name() string {
	return "service-manager-" + sm.settings.Test.ID
}

func (sm *ServiceManagerImpl) Build() error {
	if sm.settings.SM.Build {
		dockerfilePath := filepath.Join(sm.settings.SM.Path, "Dockerfile")
		fmt.Printf("Building SM docker image %s from Dockerfile: %s\n", sm.settings.SM.ImageTag, dockerfilePath)
		err := sm.dockerClient.Build(dockerfilePath, sm.settings.SM.ImageTag, sm.settings.SM.Path+string(os.PathSeparator)+".")
		if err != nil {
			return err
		}
	}

	return nil
}

func (sm *ServiceManagerImpl) Run() error {
	if sm.settings.SM.Run {
		ctx := context.Background()
		portSet := nat.PortSet{}
		port := nat.Port(sm.settings.SM.Port)
		portSet[port+"/tcp"] = struct{}{}

		portBindings := nat.PortMap{}
		portBindings[port+"/tcp"] = []nat.PortBinding{
			{HostPort: ""},
		}

		smContainer, err := sm.dockerClient.ContainerCreate(ctx,
			&container.Config{
				Image:        sm.settings.SM.ImageTag,
				ExposedPorts: portSet,
				Cmd: strslice.StrSlice{
					"--server.port=" + sm.settings.SM.Port,
					"--api.skip_ssl_validation=t",
					"--api.security.encryption_key=ejHjRNHbS0NaqARSRvnweVV9zcmhQEa8",
					"--api.token_issuer_url=https://uaa.local.pcfdev.io",
					"--api.client_id=cf",
					"--storage.uri=" + sm.settings.PG.URI,
				},
			},
			&container.HostConfig{
				PortBindings: portBindings,
			},
			nil,
			sm.Name())

		if err != nil {
			return err
		}

		err = sm.dockerClient.NetworkConnect(ctx, sm.settings.Docker.NetworkID, smContainer.ID, nil)
		if err != nil {
			return err
		}

		err = sm.dockerClient.ContainerStart(ctx, smContainer.ID, types.ContainerStartOptions{})
		if err != nil {
			return err
		}

		containerJSON, err := sm.dockerClient.ContainerInspect(ctx, smContainer.ID)
		if err != nil {
			return err
		}

		// containerJSON.NetworkSettings.IPAddress

		sm.settings.SM.Port = containerJSON.NetworkSettings.Ports[port+"/tcp"][0].HostPort

		sm.settings.SM.URL = "http://" + sm.settings.Docker.Host + ":" + sm.settings.SM.Port

		return err
	}

	return nil
}

func (sm *ServiceManagerImpl) URI() string {
	return sm.settings.SM.URL
}
