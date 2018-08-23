package config

import (
	"os"
	"path/filepath"
)

type Settings struct {
	Docker *DockerSettings
	Test   *TestSettings
	SM     *SMSettings
	PG     *PGSettings
	Broker *BrokerSettings
}

func DefaultSettings() *Settings {
	dockerMachineIP := "192.168.99.102"

	return &Settings{
		Docker: &DockerSettings{
			Host: dockerMachineIP,
		},
		Test: &TestSettings{
			ID: "pr-3",
		},
		Broker: &BrokerSettings{
			Run: true,
		},
		SM: &SMSettings{
			Path:     filepath.Join(os.Getenv("GOPATH"), filepath.Join("src", "github.com", "Peripli", "service-manager")),
			ImageTag: "service-manager:test",
			Build:    true,
			Run:      true,
			Port:     "8080",
		},
		PG: &PGSettings{
			Run:  true,
			Port: "5432",
		},
	}
}

type BrokerSettings struct {
	// ImageTag string
	URL string
	Run bool
}

type TestSettings struct {
	ID string
}

type DockerSettings struct {
	NetworkID string
	Host      string
}

type PGSettings struct {
	URI  string
	Port string
	Run  bool
}

type SMSettings struct {
	Path string
	// Build SM image
	Build bool
	// Whether SM should be started
	Run bool
	// SM image name + tag
	ImageTag string
	// SMHost is the URL to the service manager
	Host string
	// SMPort
	Port string
	URL  string
}
