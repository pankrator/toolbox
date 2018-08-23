package deploy

import "fmt"

type Deployer struct {
	builders  []DockerBuilder
	runnables []Runnable
}

func NewDeployer() *Deployer {
	return &Deployer{}
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
