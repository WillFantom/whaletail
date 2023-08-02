package docker

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/docker/docker/api/types"
)

// NetworkFilter can be used to filter a set of docker network resources.
type NetworkFilter func([]types.NetworkResource) []types.NetworkResource

// GetNetworks returns the subnets managed y any docker networks that can be
// found from the associated engine. The networks searched can be filtered by
// any network filters.
func (c Client) GetNetworks(filters ...NetworkFilter) ([]netip.Prefix, error) {
	prefixes := make([]netip.Prefix, 0)
	docker, err := c.local()
	if err != nil {
		return prefixes, fmt.Errorf("failed to connect with a docker engine: %w", err)
	}
	networks, err := docker.NetworkList(context.Background(), types.NetworkListOptions{})
	if err != nil {
		return prefixes, fmt.Errorf("failed to network set from docker: %w", err)
	}
	for _, f := range filters {
		networks = f(networks)
	}
	for _, network := range networks {
		if network.IPAM.Config == nil {
			continue
		}
		for _, config := range network.IPAM.Config {
			prefix, err := netip.ParsePrefix(config.Subnet)
			if err != nil {
				continue
			}
			prefixes = append(prefixes, prefix)
		}
	}
	return prefixes, nil
}

// NetworkFilterLabel will reduce the set of networks based on their requirement
// to have the given label.
func NetworkFilterLabel(key, value string) NetworkFilter {
	return func(networks []types.NetworkResource) []types.NetworkResource {
		validNetworks := make([]types.NetworkResource, 0)
		for _, network := range networks {
			if v, ok := network.Labels[key]; ok {
				if value == v {
					validNetworks = append(validNetworks, network)
				}
			}
		}
		return validNetworks
	}
}
