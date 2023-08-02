package tailscale

import (
	"context"
	"fmt"

	"tailscale.com/client/tailscale"
	"tailscale.com/ipn"
	"tailscale.com/ipn/ipnstate"
)

// Status requests the status of the tailscale server and returns it. If the
// request fails, an error is returned.
func (c *Client) Status() (*ipnstate.Status, error) {
	status, err := c.local().Status(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to obtain status from tailscale server: %w", err)
	}
	return status, nil
}

// Online returns true if it can be determined that both the tailscale server
// backend is running and the machine is flagged as online by the control
// server.
func (c *Client) Online() bool {
	status, err := c.local().StatusWithoutPeers(context.Background())
	if err != nil || status == nil {
		return false
	}
	return (status.Self.Online && status.BackendState == ipn.Running.String())
}

// ListenServerState creates a listener on the tailscale server that watches for
// state changes. At each state change, the given state is passed on the
// returned channel. The watcher can be exited by calling the returned cancel
// function. If the watcher can not be setup an error is returned.
func (c *Client) ListenServerState() (<-chan ipn.State, context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(context.Background())
	watcher, err := c.local().WatchIPNBus(ctx, 0)
	if err != nil {
		cancel()
		return nil, nil, fmt.Errorf("failed to create a tailscale server state watcher: %w", err)
	}
	reportingChan := make(chan ipn.State)
	go func(watcher *tailscale.IPNBusWatcher, reporter chan<- ipn.State) {
		defer close(reporter)
		for {
			n, err := watcher.Next()
			if err != nil {
				cancel()
				watcher.Close()
				return
			}
			if n.State != nil {
				reporter <- *n.State
			}
		}
	}(watcher, reportingChan)
	return reportingChan, cancel, nil
}
