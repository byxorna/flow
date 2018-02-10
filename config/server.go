package config

import (
	"fmt"

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
	Validate() error
}

// LoadServerConfigFromArgs ...
func LoadServerConfigFromArgs(args []string) (ServerConfig, error) {
	var c servercfg
	flagtag.MustConfigureAndParseArgs(&c, args)
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return &c, nil
}

// Validate validates config
func (c *servercfg) Validate() error {
	if c.listenAddr == "" {
		return fmt.Errorf("Need to provide listen-addr")
	}
	return nil
}

// ListenAddr ...
func (c *servercfg) ListenAddr() string {
	return c.listenAddr
}

// Debug ...
func (c *servercfg) Debug() bool {
	return c.debug
}
