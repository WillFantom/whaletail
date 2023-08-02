package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

func (c Client) WatchNetworkEvents() (<-chan any, context.CancelFunc, error) {
	docker, err := c.local()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect with a docker engine: %w", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	eventFilters := filters.NewArgs(
		filters.KeyValuePair{"type", "network"},
		filters.KeyValuePair{"event", "create"},
		filters.KeyValuePair{"event", "destroy"},
	)
	reportingChan := make(chan any)
	messageC, errC := docker.Events(ctx, types.EventsOptions{Filters: eventFilters})
	go func() {
		defer close(reportingChan)
		for {
			select {
			case _, ok := <-messageC:
				if !ok {
					cancel()
					return
				}
				reportingChan <- 69
			case <-errC:
				cancel()
				return
			}
		}
	}()
	return reportingChan, cancel, nil
}
