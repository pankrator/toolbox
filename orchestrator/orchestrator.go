package orchestrator

import (
	"github.com/Peripli/itest-tools/deploy"
)

type DeployTarget string

var (
	CFTarget     DeployTarget = "cf"
	DockerTarget DeployTarget = "docker"
	K8STarget    DeployTarget = "k8s"
)

type Orchestrator struct {
	dockerDeployer *deploy.Deployer
	deployables    map[string]PreparationDeployment
	indexed        []string
}

type PreparationDeployment struct {
	name         string
	dependencies []string
	config       *Config
}

type Config struct {
	Target DeployTarget

	Name  string
	Image string

	NetworkID string

	Port       string
	Env        map[string]string
	CmdArgs    map[string]string
	CFManifest string

	dockerConfigCreate func(*Config, map[string]deploy.DeploymentResult) *deploy.DockerRunOptions
	// cfConfigCreate func(*deploy.DockerRunOptions, dependencies map[string]deploy.DeploymentResult) *deploy.DockerRunOptions
}

func (o *Orchestrator) Add(name string, dependencies []string, c *Config) {
	deployment := PreparationDeployment{
		name:         name,
		dependencies: dependencies,
		config:       c,
	}
	o.deployables[name] = deployment
}

func (o *Orchestrator) Run() {
	results := make(map[string]deploy.DeploymentResult)

	for _, deployment := range o.deployables {
		o.runOne(deployment, results)
	}
}

func (o *Orchestrator) runOne(deployment PreparationDeployment, results map[string]deploy.DeploymentResult) {
	for _, depName := range deployment.dependencies {
		_, exists := results[depName]
		if !exists {
			o.runOne(o.deployables[depName], results)
		}
	}

	switch deployment.config.Target {
	case DockerTarget:
		dockerRunOptions := deployment.config.dockerConfigCreate(deployment.config, results)
		// buildDockerConfig(deployment.config)

		o.dockerDeployer.Run
	case CFTarget:

	}
}
