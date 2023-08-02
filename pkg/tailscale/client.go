package tailscale

import (
	"fmt"
	"os"

	ts "tailscale.com/client/tailscale"
)

// Client can be used to interact and control a local tailscale server. This
// will work in most cases using the default configuration (without a socket
// path explicitly set), however for more fine grain control of what tailscale
// server is in use, this client allows that to be set and forced.
type Client struct {
	socketPath string
}

// NewClient returns a new tailscale client for interaction with a local
// tailscale server. Socket path can be empty to let tailscale try a range of
// defaults that should work on many platforms. Although, a socket path for a
// local tailscale server can be explicitly set.
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

// SetSocketPath lets the configuration of the client be reconfigured. This can
// be set to empty to go back to tailscale defaults.
func (c *Client) SetSocketPath(socketPath string) error {
	if socketPath != "" {
		if _, err := os.Stat(socketPath); err != nil {
			return fmt.Errorf("given tailscale socket path could not be used: %w", err)
		}
	}
	c.socketPath = socketPath
	return nil
}

func (c Client) local() *ts.LocalClient {
	lc := ts.LocalClient{}
	if c.socketPath != "" {
		lc.Socket = c.socketPath
		lc.UseSocketOnly = true
	}
	return &lc
}
