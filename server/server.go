package server

import (
	"fmt"
	"io"
	"net/http"

	"github.com/byxorna/flow/config"
	"github.com/byxorna/flow/version"
	etcd "github.com/coreos/etcd/client"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithFields(logrus.Fields{"module": "server"})
)

type svr struct {
	mux          *http.ServeMux
	ServerConfig config.ServerConfig
	EtcdClient   etcd.Client
	KeysAPI      etcd.KeysAPI
}

// Server ...
type Server interface {
	ListenAndServe() error
}

// New returns a new server
func New(sc config.ServerConfig, ec config.EtcdConfig) (Server, error) {
	etcdClient, err := etcd.New(ec.ToEtcdConfig())
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	s := svr{
		mux:          mux,
		KeysAPI:      etcd.NewKeysAPIWithPrefix(etcdClient, ec.Prefix()),
		ServerConfig: sc,
		EtcdClient:   etcdClient,
	}

	// register http handlers
	mux.HandleFunc("/", s.handleVersion)

	return &s, nil
}

// ListenAndServe calls http.ListenAndServe
func (s *svr) ListenAndServe() error {
	log.Infof("Listening on %s", s.ServerConfig.ListenAddr())
	return http.ListenAndServe(s.ServerConfig.ListenAddr(), s.mux)
}

func (s *svr) handleVersion(w http.ResponseWriter, r *http.Request) {
	io.WriteString(
		w,
		fmt.Sprintf(
			"Version: %s\nDate: %s\nBranch: %s\nCommit: %s\n",
			version.Version,
			version.BuildDate,
			version.Branch,
			version.Commit,
		),
	)
}
