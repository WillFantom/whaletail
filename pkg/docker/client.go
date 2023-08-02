package docker

import (
	"fmt"
	"os"

	"github.com/docker/docker/client"
)

// Client defines a connection to a docker engine. A specific socket path can be
// used to bypass the given env values used by default.
type Client struct {
	socketPath string
}

// NewClient creates and returns a new client for interaction with the docker
// socket. If a specific socket path is given yet does not exist, an error is
// returned. Socket path can be empty to use the docker defaults.
func NewClient(socketPath string) (*Client, error) {
	client := Client{}
	if socketPath != "" {
		if err := client.SetSocketPath(socketPath); err != nil {
			return nil, err
		}
		client.socketPath = socketPath
	}
	return &client, nil
}

func (c *Client) SetSocketPath(socketPath string) error {
	if socketPath != "" {
		if _, err := os.Stat(socketPath); err != nil {
			return fmt.Errorf("given docker socket path could not be used: %w", err)
		}
	}
	c.socketPath = socketPath
	return nil
}

func (c Client) local() (*client.Client, error) {
	if c.socketPath == "" {
		return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	} else {
		return client.NewClientWithOpts(client.WithHost(c.socketPath), client.WithAPIVersionNegotiation())
	}
}
