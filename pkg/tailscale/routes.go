package tailscale

import (
	"context"
	"fmt"
	"net/netip"

	"golang.org/x/exp/slices"
	"tailscale.com/ipn"
)

// ListPrimaryRoutes returns the the routes in which the tailscale server associated
// with the client is the primary router for.
func (c *Client) ListPrimaryRoutes() ([]netip.Prefix, error) {
	status, err := c.local().StatusWithoutPeers(context.Background())
	if err != nil {
		return []netip.Prefix{}, err
	}
	if status.Self == nil {
		return []netip.Prefix{}, fmt.Errorf("failed to get route info")
	}
	if status.Self.PrimaryRoutes == nil {
		return []netip.Prefix{}, nil
	}
	return status.Self.PrimaryRoutes.AsSlice(), nil
}

// ListRoutes returns the the routes in which the tailscale server associated
// with the client is advertising. Note that this does not guarantee that the
// tailscale server is the primary router for the routes nor that the routes
// have been accepted by tailscale control.
func (c *Client) ListRoutes() ([]netip.Prefix, error) {
	prefs, err := c.local().GetPrefs(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get current server preferences: %w", err)
	}
	if prefs == nil {
		return nil, fmt.Errorf("received an empty preferences set from the tailscale server")
	}
	return prefs.AdvertiseRoutes, nil
}

// AddRoute tells the associated tailscale server to start advertising itself
// as a subnet router to the given route. If it is already doing so, no error is
// returned.
func (c *Client) AddRoute(subnet netip.Prefix) error {
	return c.AddRoutes([]netip.Prefix{subnet})
}

// AddRoutes tells the associated tailscale server to start advertising itself
// as a subnet router to the given routes. If it is already doing so, no error
// is returned.
func (c *Client) AddRoutes(subnets []netip.Prefix) error {
	routes, err := c.ListRoutes()
	if err != nil {
		return err
	}
	for _, route := range subnets {
		if !slices.Contains[netip.Prefix](routes, route) {
			routes = append(routes, route)
		}
	}
	return c.SetRoutes(routes)
}

// DelRoute updates the tailscale server configuration, stopping it from
// advertising the given route. If it already not advertsing the route, no error
// is returned.
func (c *Client) DelRoute(subnet netip.Prefix) error {
	return c.DelRoutes([]netip.Prefix{subnet})
}

// DelRoutes can be used to remove multiple route advertisements from the
// tailscale server and behaves the same as DelRoute.
func (c *Client) DelRoutes(subnets []netip.Prefix) error {
	routes, err := c.ListRoutes()
	if err != nil {
		return err
	}
	toDelete := make([]int, 0)
	for _, sb := range subnets {
		if slices.Contains[netip.Prefix](routes, sb) {
			idx := slices.Index[netip.Prefix](routes, sb)
			if idx >= 0 {
				toDelete = append(toDelete, idx)
			}
		}
	}
	for _, del := range toDelete {
		routes = slices.Delete[[]netip.Prefix](routes, del, del)
	}
	return c.SetRoutes(routes)
}

func (c *Client) SetRoutes(subnets []netip.Prefix) error {
	_, err := c.local().EditPrefs(context.Background(), &ipn.MaskedPrefs{
		Prefs: *&ipn.Prefs{
			AdvertiseRoutes: subnets,
		},
		AdvertiseRoutesSet: true,
	})
	if err != nil {
		return fmt.Errorf("failed to update server routes configuration: %w", err)
	}
	return nil
}
