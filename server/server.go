package server

import (
	"net/http"

	"github.com/byxorna/flow/config"
	etcd "github.com/coreos/etcd/client"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithFields(logrus.Fields{"module": "server"})
)

type svr struct {
	http.Server
	ServerConfig config.ServerConfig
	EtcdClient   client.Client
	KeysAPI      client.KeysAPI
}

// Server ...
type Server interface {
	ListenAndServe() error
}

// New returns a new server
func New(sc config.ServerConfig, ec config.EtcdConfig) (Server, error) {
	etcdClient, err := client.New(ec.ToEtcdConfig())
	if err != nil {
		return nil, err
	}
	s := svr{
		KeysAPI:      client.NewKeysAPIWithPrefix(etcdClient, ec.Prefix()),
		ServerConfig: sc,
		EtcdClient:   etcdClient,
	}
	return &s, nil
}

// ListenAndServe calls http.ListenAndServe
func (s *svr) ListenAndServe() error {
	log.Infof("Listening on %s", s.ServerConfig.ListenAddr())
	return s.ListenAndServe(s.ServerConfig.ListenAddr(), nil)
}

func (s *svr) handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Print("/")
}
