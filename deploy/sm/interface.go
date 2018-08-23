package sm

import (
	"github.com/Peripli/itest-tools/deploy"
)

type ServiceManager interface {
	deploy.Deployment
	deploy.Runnable
	deploy.DockerBuilder
}
