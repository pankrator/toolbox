package docker

import (
	"net/http"

	"github.com/Peripli/itest-tools/cmd"
	"github.com/docker/docker/client"
)

type Docker struct {
	*client.Client
}

type Settings struct {
	Host    string
	Version string
}

func New(settings *Settings) (*Docker, error) {
	var cli *client.Client
	var err error
	if settings == nil {
		cli, err = client.NewEnvClient()
	} else {
		cli, err = client.NewClient(settings.Host, settings.Version, http.DefaultClient, make(map[string]string))
	}

	if err != nil {
		return nil, err
	}

	return &Docker{
		Client: cli,
	}, nil
}

func (d *Docker) Build(dockerFilePath, imageTag, contextPath string) error {
	return dockerRun("build", "-f", dockerFilePath, "-t", imageTag, contextPath)
}

type RunOptions struct {
	ContainerName   string
	Network         string
	ExposePort      string
	Env             map[string]string
	DetachContainer bool
	// MapPort should be in the format hostPort:dockerPort i.e. 5000:5000
	MapPort string
	Cmd     []string
}

func (d *Docker) Run(image string, opts RunOptions) error {
	cArgs := []string{"run"}

	if opts.ExposePort != "" {
		cArgs = append(cArgs, "--expose", opts.ExposePort)
	}

	if opts.Env != nil {
		for k, v := range opts.Env {
			cArgs = append(cArgs, "-e", k+"="+v)
		}
	}

	if opts.DetachContainer {
		cArgs = append(cArgs, "-d")
	}

	if opts.ContainerName != "" {
		cArgs = append(cArgs, "--name", opts.ContainerName)
	}

	if opts.MapPort != "" {
		cArgs = append(cArgs, "-p", opts.MapPort)
	}

	if opts.Network != "" {
		cArgs = append(cArgs, "--network", opts.Network)
	}

	cArgs = append(cArgs, image)

	if opts.Cmd != nil {
		cArgs = append(cArgs, opts.Cmd...)
	}

	return dockerRun(cArgs...)
}

func dockerRun(args ...string) error {
	return cmd.DoOutput("docker", args...)
}
