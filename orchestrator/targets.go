package orchestrator

func buildDockerConfig(c *Config) *deploy.DockerRunOptions {
	return &deploy.DockerRunOptions{
		Config: &container.Config{
			Image: c.Image,
		},
	}
}
