package config

import ()

// ServerConfig ...
type ServerConfig struct {
	ServerListenAddr string `yaml:"listen-addr" arg:"--server-listen-addr" help:"What interface:port to listen on"`
	Debug            bool   `yaml:"debug" arg:"-d" help:"Enable debug logging"`
}

// ValidateAndSetServerDefaults validates config and sets defaults if possible
func (c *ServerConfig) ValidateAndSetServerDefaults() error {
	if c.ServerListenAddr == "" {
		c.ServerListenAddr = ":6969"
	}
	return nil
}
