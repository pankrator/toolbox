package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Peripli/itest-tools/deploy"

	"github.com/Peripli/service-manager-cli/pkg/auth"
	"github.com/Peripli/service-manager-cli/pkg/auth/oidc"
	"github.com/Peripli/service-manager-cli/pkg/smclient"

	dockerTypes "github.com/docker/docker/api/types"

	"github.com/Peripli/itest-tools/deploy/broker"
	"github.com/Peripli/itest-tools/deploy/config"
	"github.com/Peripli/itest-tools/deploy/pg"
	"github.com/Peripli/itest-tools/deploy/sm"
	"github.com/Peripli/itest-tools/docker"
)

func ping(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Could not ping %s", url))
	}

	return nil
}

// type Settings struct {
// 	SM sm.Settings

// 	TestID string

// 	// CFProxyDeploy Whether CF proxy should be deployed
// 	CFProxyDeploy       bool
// 	CFProxyAppName      string
// 	CFProxyManifestPath string
// 	CFProxyBuildpack    string
// 	CFProxyHost         string
// 	CFProxyPort         string

// 	BrokerImageTag string
// 	BrokerHost     string
// 	BrokerPort     string
// 	BrokerUser     string
// 	BrokerPassword string

// 	DockerMachineHost string
// }

// func DefaultSettings() Settings {
// 	dockerMachineIP := "192.168.99.102"

// 	return Settings{
// 		TestID: "pr-8",

// 		SM: sm.Settings{
// 			SMPostgresURI: "postgres://postgres:postgres@192.168.99.102:5432/postgres?sslmode=disable",
// 			SMRunPostgres: true,
// 			SMPath:        filepath.Join(os.Getenv("GOPATH"), filepath.Join("src", "github.com", "Peripli", "service-manager")),
// 			SMImageTag:    "service-manager:test",
// 			SMBuild:       true,
// 			SMRun:         true,
// 			SMHost:        "http://" + dockerMachineIP,
// 			SMPort:        "8080",
// 		},

// 		CFProxyDeploy:       true,
// 		CFProxyAppName:      "cfproxy2",
// 		CFProxyManifestPath: filepath.Join(os.Getenv("GOPATH"), filepath.Join("src", "github.com", "Peripli", "service-broker-proxy-cf", "manifest.yml")),
// 		CFProxyBuildpack:    "go_buildpack_new",

// 		BrokerImageTag: "sbf-broker:test",
// 		BrokerPort:     "5000",
// 		BrokerHost:     "",
// 		BrokerUser:     "admin",
// 		BrokerPassword: "admin",

// 		DockerMachineHost: dockerMachineIP,
// 	}
// }

func main() {
	// TODO: make it better
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	settings := config.DefaultSettings()

	var err error

	dockerClient, _ := docker.New(nil)

	network, err := dockerClient.NetworkCreate(context.Background(), "test-net-"+settings.Test.ID, dockerTypes.NetworkCreate{
		Attachable:     true,
		CheckDuplicate: true,
		Driver:         "bridge",
	})
	if err != nil {
		panic(err)
	}

	settings.Docker.NetworkID = network.ID
	pgOperator := pg.New(settings, dockerClient)
	smOperator := sm.New(settings, dockerClient)
	brokerOperator := broker.New(settings, dockerClient)

	deployer := deploy.NewDeployer()
	deployer.AddBuilder(smOperator)
	deployer.AddRunnable(brokerOperator)
	deployer.AddRunnable(pgOperator)
	deployer.AddRunnable(smOperator)

	err = deployer.Build()
	if err != nil {
		panic(err)
	}

	err = deployer.Run()
	if err != nil {
		panic(err)
	}

	/*

		fmt.Printf("Starting broker image %s in container\n", settings.BrokerImageTag)
		err = dockerClient.Run(settings.BrokerImageTag, docker.RunOptions{
			ContainerName:   "broker-" + settings.TestID,
			DetachContainer: true,
			Env: map[string]string{
				"PORT":               settings.BrokerPort,
				"SBF_CATALOG_SUFFIX": settings.TestID,
			},
			ExposePort: settings.BrokerPort,
			Network:    network.ID,
			// MapPort:         settings.BrokerPort + ":" + settings.BrokerPort,
		})
		if err != nil {
			panic(err)
		}

		if settings.SMBuild {
			dockerfilePath := filepath.Join(settings.SMPath, "Dockerfile")
			fmt.Printf("Building SM docker image %s from Dockerfile: %s\n", settings.SMImageTag, dockerfilePath)
			err = dockerClient.Build(dockerfilePath, settings.SMImageTag, settings.SMPath+string(os.PathSeparator)+".")
			if err != nil {
				panic(err)
			}
		}

		if settings.SMRun {
			fmt.Printf("Starting SM image %s in container\n", settings.SMImageTag)

			postgresURL := settings.SMPostgresURI
			if postgresURL == "" {
				postgresURL = "postgres://postgres:postgres@postgres-" + settings.TestID + ":5432/postgres?sslmode=disable"
			}

			err = dockerClient.Run(settings.SMImageTag, docker.RunOptions{
				ContainerName:   "service-manager-" + settings.TestID,
				DetachContainer: true,
				MapPort:         settings.SMPort + ":8080",
				Network:         network.ID,
				// TODO: SM command line args should be configurable
				Cmd: []string{
					"--api.skip_ssl_validation=t",
					"--api.security.encryption_key=ejHjRNHbS0NaqARSRvnweVV9zcmhQEa8",
					"--api.token_issuer_url=https://uaa.local.pcfdev.io",
					"--api.client_id=cf",
					"--storage.uri=" + postgresURL,
				},
			})
			if err != nil {
				panic(err)
			}
			// TODO: Show where SM is accessible
		}

		fmt.Print("SM-Username:")
		user, _ := util.ReadInput(os.Stdin)
		fmt.Print("SM-Password:")
		password, _ := util.ReadPassword()
		fmt.Println()

		smClient := getAuthenticatedSmClient(settings.SMHost+":"+settings.SMPort, user, password, http.DefaultClient)
		cfClient := cf.New()

		// TODO: fix http://
		brokerURL := "http://" + settings.BrokerHost + ":" + settings.BrokerPort
		if settings.BrokerHost == "" {
			brokerURL = "http://broker-" + settings.TestID + ":" + settings.BrokerPort
		}
		fmt.Printf("Registering broker with URL: %s in SM\n", brokerURL)
		_, err = smClient.RegisterBroker(&types.Broker{
			Name: "sbf-test",
			URL:  brokerURL,
			Credentials: &types.Credentials{
				Basic: types.Basic{
					User:     settings.BrokerUser,
					Password: settings.BrokerPassword,
				},
			},
		})
		if err != nil {
			fmt.Println("Could not register broker. Reason: ", err)
		}

		if settings.CFProxyDeploy {
			fmt.Printf("Registering platform in SM\n")
			platform, err := smClient.RegisterPlatform(&types.Platform{
				// TODO
				Name: "asd",
				Type: "asd",
			})
			if err != nil {
				panic(err)
			}

			fmt.Printf("Push CF proxy with manifest from: %s\n", settings.CFProxyManifestPath)
			err = cfClient.Push(settings.CFProxyAppName+settings.TestID, settings.CFProxyManifestPath, settings.CFProxyBuildpack)
			if err != nil {
				panic(err)
			}

			cfClient.Env(settings.CFProxyAppName, "SM_HOST", settings.SMHost+":"+settings.SMPort)
			cfClient.Env(settings.CFProxyAppName, "SM_USER", platform.Credentials.Basic.User)
			cfClient.Env(settings.CFProxyAppName, "SM_PASSWORD", platform.Credentials.Basic.Password)

			cfClient.Start(settings.CFProxyAppName)
		}
	*/
}

func getAuthenticatedSmClient(smURL, username, password string, httpClient *http.Client) smclient.Client {
	client := smclient.NewClient(httpClient, &smclient.ClientConfig{
		URL: smURL,
	})
	info, err := client.GetInfo()
	if err != nil {
		panic(err)
	}
	// TODO: Use clientid from configuration
	authOptions := &auth.Options{
		ClientID:     "cf",
		ClientSecret: "",
		IssuerURL:    info.TokenIssuerURL,
		SSLDisabled:  true,
		Timeout:      time.Second * 10,
	}
	authenticator, _, err := oidc.NewOpenIDStrategy(authOptions)
	if err != nil {
		panic(err)
	}
	token, err := authenticator.Authenticate(username, password)
	if err != nil {
		panic(fmt.Sprintf("Could not authenticate. Reason %s", err))
	}

	return smclient.NewClient(oidc.NewClient(authOptions, token), &smclient.ClientConfig{
		URL: smURL,
	})
}
