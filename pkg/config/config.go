package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	appConfig            *viper.Viper = viper.New()
	listeners            []chan any   = make([]chan any, 0)
	defaultConfiguration Config       = Config{
		Log: LogConfig{
			Level: "info",
			File:  "",
		},
		Tailscale: TsConfig{
			SocketPath: "",
		},
		Docker: DockerConfig{
			Endpoint: "",
		},
	}
)

// AppConfig is the viper instance in use for the main application
// configuration. This should be used to access all configuration information.
func AppConfig() *viper.Viper {
	return appConfig
}

// Read attempts to read in a config file. If no config file is found, the
// default config is used. If a config file is present, but can not be read
// in or parsed, and error is returned.
func Read() error {
	if err := appConfig.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}
	return nil
}

// Write attempts to write the current configuration state to file. It does so
// even if a current config file exists. An error is returned if the write
// fails.
func Write(c *Config) error {
	return appConfig.WriteConfig()
}

// Listen adds a new listener to the config. When a config change occurs some
// data (pointless) will be pushed to the returned channel to notify the
// listener.
func Listen() <-chan any {
	listener := make(chan any)
	listeners = append(listeners, listener)
	return listener
}

func init() {
	// general config
	appConfig.SetConfigName("config")
	appConfig.AddConfigPath("/etc/whaletail")
	appConfig.AddConfigPath("$HOME/.config/whaletail")
	appConfig.AddConfigPath(".")
	appConfig.WatchConfig()
	appConfig.OnConfigChange(func(e fsnotify.Event) {
		for _, listener := range listeners {
			listener <- 69
		}
	})

	// log defaults
	appConfig.SetDefault("log.level", defaultConfiguration.Log.Level)
	appConfig.SetDefault("log.file", defaultConfiguration.Log.File)

	// tailscale defaults
	appConfig.SetDefault("tailscale.socket", defaultConfiguration.Tailscale.SocketPath)

	// docker defaults
	appConfig.SetDefault("docker.endpoint", defaultConfiguration.Docker.Endpoint)
}
