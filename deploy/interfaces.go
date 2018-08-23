package deploy

type Deployment interface {
	Name() string
	URIService
}

type URIService interface {
	URI() string
}

type HostService interface {
	Host() string
	Port() string
}

type DockerBuilder interface {
	Build() error
}

type Runnable interface {
	Run() error
}

type CFDeployer interface {
	DeployCF() error
}
