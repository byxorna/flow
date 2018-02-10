package config

import (
	"github.com/cobratbq/flagtag"
)

type servercfg struct {
	listenAddr string `yaml:"listen-addr",flag:"listen-addr,:6969,What interface:port to listen on"`
	debug      bool   `yaml:"debug",flag:"debug,false,Enable debug logging"`
}

// ServerConfig ...
type ServerConfig interface {
	ListenAddr() string
	Debug() bool
}

// LoadServerConfigFromArgs ...
func LoadServerConfigFromArgs(args []string) ServerConfig {
	var c servercfg
	flagtag.MustConfigureAndParse(&c)
	return &c
}

// ListenAddr ...
func (c *servercfg) ListenAddr() string {
	return c.listenAddr
}

// Debug ...
func (c *servercfg) Debug() bool {
	return c.debug
}
