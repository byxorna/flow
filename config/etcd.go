package config

import (
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
}

// LoadEtcdConfigFromArgs ...
func LoadEtcdConfigFromArgs(args []string) EtcdConfig {
	var c etcdcfg
	flagtag.MustConfigureAndParse(&c)
	return &c
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
