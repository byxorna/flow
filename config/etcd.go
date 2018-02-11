package config

import (
	"fmt"
	"time"

	"github.com/docker/libkv/store"
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

// ToEtcdConfig returns a etcd client Config structure for libkv
func (c *EtcdConfig) ToLibKVConfig() store.Config {
	return store.Config{
		ClientTLS:         nil,
		TLS:               nil,
		ConnectionTimeout: 1 * time.Second,
		PersistConnection: true,
	}
}
