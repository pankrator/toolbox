package cf

import "github.com/Peripli/itest-tools/cmd"

type Client struct {
}

func New() *Client {
	return &Client{}
}

func (c *Client) Push(name, manifestPath, buildpack string) error {
	return cfRun("push", name, "-f", manifestPath, "-b", buildpack, "--no-start")
}

func (c *Client) Start(name string) error {
	return cfRun("start", name)
}

func (c *Client) Env(appName, envName, envValue string) error {
	return cfRun("set-env", appName, envName, envValue)
}

func cfRun(args ...string) error {
	return cmd.DoOutput("cf", args...)
}
