package config

// Config is a union of all configuration structs
type Config struct {
	EtcdConfig
	ServerConfig
}

// ValidateAndSetDefaults validates all embedded structs and sets defaults where applicable
func (c *Config) ValidateAndSetDefaults() error {
	if err := c.ValidateAndSetEtcdDefaults(); err != nil {
		return err
	}
	err := c.ValidateAndSetServerDefaults()
	return err
}
