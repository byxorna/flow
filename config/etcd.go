package config

import (
	"fmt"
	"time"

	"github.com/cobratbq/flagtag"
	"github.com/coreos/etcd/client"
)

type etcdcfg struct {
	endpoints []string `yaml:"etcd-endpoints",flag:"endpoints,,etcd endpoints for storage"`
	prefix    string   `yaml:"etcd-prefix",flag:"etcd-prefix,/,etcd prefix for storage"`
}

// EtcdConfig ...
type EtcdConfig interface {
	Endpoints() []string
	Prefix() string
	ToEtcdConfig() client.Config
	Validate() error
}

// LoadEtcdConfigFromArgs ...
func LoadEtcdConfigFromArgs(args []string) (EtcdConfig, error) {
	var c etcdcfg
	flagtag.MustConfigureAndParseArgs(&c, args)

	err := c.Validate()
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// Validate validates config
func (c *etcdcfg) Validate() error {
	if len(c.endpoints) < 1 {
		return fmt.Errorf("Need to provide etcd-endpoints")
	}
	if c.prefix == "" {
		return fmt.Errorf("Need to provide valid etcd-prefix")
	}
	return nil
}

// Endpoints ...
func (c *etcdcfg) Endpoints() []string { return c.endpoints }

// Prefix ...
func (c *etcdcfg) Prefix() string { return c.prefix }

// ToEtcdConfig ...
func (c *etcdcfg) ToEtcdConfig() client.Config {
	x := client.Config{
		Endpoints:               c.endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	return x
}
