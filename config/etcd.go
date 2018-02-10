package config

import (
	"fmt"
	"time"

	"github.com/coreos/etcd/client"
)

// EtcdConfig ...
type EtcdConfig struct {
	EtcdEndpoints []string `yaml:"etcd-endpoints" arg:"--etcd-endpoints,required" help:"etcd endpoints for storage"`
	EtcdPrefix    string   `yaml:"etcd-prefix" arg:"--etcd-prefix" help:"etcd prefix for storage"`
}

// ValidateAndSetEtcdDefaults validates config, and sets defaults if possible
func (c *EtcdConfig) ValidateAndSetEtcdDefaults() error {
	if len(c.EtcdEndpoints) < 1 {
		return fmt.Errorf("Need to provide etcd-endpoints")
	}
	if c.EtcdPrefix == "" {
		c.EtcdPrefix = "/"
	}
	return nil
}

// ToEtcdConfig returns a etcd client Config structure
func (c *EtcdConfig) ToEtcdConfig() client.Config {
	x := client.Config{
		Endpoints:               c.EtcdEndpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	return x
}
