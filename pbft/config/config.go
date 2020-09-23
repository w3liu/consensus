package config

import "github.com/BurntSushi/toml"

type Config struct {
	Peer Peer
}

type Peer struct {
	Address string
	Seeds   []string
}

func New() *Config {
	return &Config{}
}

func Init(cfgFile string) (*Config, error) {
	cfg := &Config{}
	_, err := toml.DecodeFile(cfgFile, cfg)
	return cfg, err
}
