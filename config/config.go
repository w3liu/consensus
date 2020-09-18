package config

import "github.com/BurntSushi/toml"

type Config struct {
	Peer Peer
}

type Peer struct {
	Address string
	Seeds   []string
}

func (p *Config) Init(cfgFile string) error {
	_, err := toml.DecodeFile(cfgFile, p)
	return err
}
