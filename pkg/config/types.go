package config

type Config struct {
	Log       LogConfig    `toml:"log,omitempty" json:"log,omitempty" yaml:"log,omitempty"`
	Tailscale TsConfig     `toml:"tailscale,omitempty" json:"tailscale,omitempty" yaml:"tailscale,omitempty"`
	Docker    DockerConfig `toml:"docker,omitempty" json:"docker,omitempty" yaml:"docker,omitempty"`
}

type LogConfig struct {
	Level string `toml:"level,omitempty" json:"level,omitempty" yaml:"level,omitempty"`
	File  string `toml:"file,omitempty" json:"file,omitempty" yaml:"file,omitempty"`
}

type TsConfig struct {
	SocketPath string `toml:"socket,omitempty" json:"socket,omitempty" yaml:"socket,omitempty"`
}

type DockerConfig struct {
	Endpoint string `toml:"endpoint,omitempty" json:"endpoint,omitempty" yaml:"endpoint,omitempty"`
}
