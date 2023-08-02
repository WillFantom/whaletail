package main

import (
	"net/netip"
	"os"
	"sync"
	"time"

	"github.com/willfantom/whaletail/pkg/config"
	"github.com/willfantom/whaletail/pkg/docker"
	"github.com/willfantom/whaletail/pkg/log"
	"github.com/willfantom/whaletail/pkg/tailscale"
)

type Whaletail struct {
	tailscaleClient *tailscale.Client
	dockerClient    *docker.Client

	configWatcher    <-chan any
	interruptWatcher <-chan os.Signal

	networks []netip.Prefix

	sync.Mutex
}

func (wt *Whaletail) DockerConnet() *docker.Client {
	for {
		log.Logger.Infoln("attempting to connect to a docker engine")
		client, err := docker.NewClient(config.AppConfig().GetString("docker.endpoint"))
		if err != nil {
			log.Logger.
				WithError(err).
				Errorln("failed to get docker endpoint")
			time.Sleep(30 * time.Second) //backoff
			continue
		}
		if !client.Online() {
			log.Logger.
				Errorln("no response from docker endpoint")
			time.Sleep(30 * time.Second) //backoff
			continue
		}
		log.Logger.Infoln("connection established to docker engine")
		return client
	}
}

func (wt *Whaletail) TailscaleConnect() *tailscale.Client {
	for {
		log.Logger.Infoln("attempting to connect to a tailscale server")
		client, err := tailscale.NewClient(config.AppConfig().GetString("tailscale.socket"))
		if err != nil {
			log.Logger.
				WithError(err).
				Errorln("failed to find tailscale socket")
			time.Sleep(30 * time.Second) //backoff
			continue
		}
		if !client.Online() {
			log.Logger.
				Errorln("tailscale machine is not online")
			time.Sleep(30 * time.Second) //backoff
			continue
		}
		log.Logger.Infoln("connection established to up tailscale server")
		return client
	}
}

func (wt *Whaletail) Run() {
	for {
		if !wt.dockerClient.Online() {
			wt.dockerClient = wt.DockerConnet()
		}
		if !wt.tailscaleClient.Online() {
			wt.tailscaleClient = wt.TailscaleConnect()
		}

		err := wt.updateTailscaleRoutes()
		if err != nil {
			log.Logger.
				WithError(err).
				Errorln("failed to advertise networks")
			continue
		}

		events, cancel, err := wt.dockerClient.WatchNetworkEvents()
		if err != nil {
			log.Logger.
				WithError(err).
				Errorln("failed to create docker network event watcher")
			continue
		}

		select {
		case _, ok := <-events:
			if !ok {
				log.Logger.
					Errorln("docker network event watcher has failed")
				cancel()
				continue
			}
			err := wt.updateTailscaleRoutes()
			if err != nil {
				log.Logger.
					WithError(err).
					Errorln("failed to advertise networks found through watch")
				cancel()
				continue
			}
		case <-wt.interruptWatcher:
			wt.Lock()
			defer wt.Unlock()
			if err := wt.tailscaleClient.SetRoutes([]netip.Prefix{}); err != nil {
				log.Logger.WithError(err).Fatalln("failed to remove routes on exit")
			}
			return
		case <-wt.configWatcher:
			log.Logger.Infoln("detected configuration change. restarting")
			cancel()
			continue
		}
	}
}

func (wt *Whaletail) updateTailscaleRoutes() error {
	networks, err := wt.dockerClient.GetNetworks(docker.NetworkFilterLabel("whaletail.enable", "true"))
	if err != nil {
		return err
	}
	log.Logger.WithField("count", len(networks)).Infoln("found docker networks")
	wt.Lock()
	// removedNetworks := make([]netip.Prefix, 0)
	// for _, n := range networks {
	// 	if !slices.Contains[netip.Prefix](wt.networks, n) {
	// 		removedNetworks = append(removedNetworks, n)
	// 	}
	// }
	wt.networks = networks
	wt.Unlock()
	log.Logger.WithField("count", len(networks)).Infoln("setting routes")
	if err = wt.tailscaleClient.SetRoutes(networks); err != nil {
		return err
	}
	// log.Logger.WithField("count", len(removedNetworks)).Infoln("removing networks")
	// if err = wt.tailscaleClient.DelRoutes(removedNetworks); err != nil {
	// 	return err
	// }
	return nil
}
