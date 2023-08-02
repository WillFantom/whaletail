package main

import (
	"net/netip"
	"os"
	"os/signal"

	"github.com/willfantom/whaletail/pkg/config"
	"github.com/willfantom/whaletail/pkg/docker"
	"github.com/willfantom/whaletail/pkg/log"
	"github.com/willfantom/whaletail/pkg/tailscale"
)

func main() {
	if err := config.Read(); err != nil {
		log.Logger.WithError(err).Fatalln("config read failure")
	}
	if err := log.SetLevel(config.AppConfig().GetString("log.level")); err != nil {
		log.Logger.WithError(err).Fatalln("log level could not be set")
	}

	configWatcher := config.Listen()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	wt := &Whaletail{
		tailscaleClient:  &tailscale.Client{},
		dockerClient:     &docker.Client{},
		configWatcher:    configWatcher,
		interruptWatcher: c,
		networks:         make([]netip.Prefix, 0),
	}

	wt.Run()

}
