package server

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/byxorna/flow/config"
	"github.com/byxorna/flow/types/executor/shell"
	"github.com/byxorna/flow/types/storage"
	"github.com/byxorna/flow/version"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithFields(logrus.Fields{"module": "server"})
)

// for middleware logging
type responseLogger struct {
	http.ResponseWriter
	statusCode int
	timeStart  time.Time
	duration   time.Duration
}

func (rl responseLogger) WriteHeader(code int) {
	rl.statusCode = code
	rl.ResponseWriter.WriteHeader(code)
	rl.duration = time.Since(rl.timeStart)
}

type svr struct {
	config.Config
	mux *http.ServeMux
	// Store is the backend data access layer (etcd)
	store *storage.Store

	// executors the server needs to know about
	shellExecutor *shell.Executor
}

// Server ...
type Server interface {
	ListenAndServe() error
	RegisterShellExecutor(e *shell.Executor)
}

// RegisterShellExecutor ...
func (s *svr) RegisterShellExecutor(e *shell.Executor) {
	s.shellExecutor = e
}

// New returns a new server
func New(c config.Config, store *storage.Store) (Server, error) {
	mux := http.NewServeMux()

	s := svr{
		Config: c,
		mux:    mux,
		store:  store,
	}

	// register http handlers
	mux.HandleFunc("/", s.getVersion)
	mux.HandleFunc("/v1/jobs", s.getJobs)

	return &s, nil
}

// ListenAndServe calls http.ListenAndServe
func (s *svr) ListenAndServe() error {
	log.WithFields(
		logrus.Fields{"address": s.ServerListenAddr},
	).Infof("Listening for HTTP requests")
	return http.ListenAndServe(
		s.ServerListenAddr,
		logRequest(s.mux),
	)
}

func (s *svr) getVersion(w http.ResponseWriter, r *http.Request) {
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

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl := responseLogger{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			timeStart:      time.Now(),
		}
		log.WithFields(
			logrus.Fields{
				"remoteaddr": r.RemoteAddr,
				"method":     r.Method,
				"url":        r.URL,
			},
		).Infof("HTTP request received")
		handler.ServeHTTP(rl, r)
		log.WithFields(
			logrus.Fields{
				"status":     rl.statusCode,
				"statustext": http.StatusText(rl.statusCode),
				"duration":   rl.duration.String(),
			},
		).Infof("HTTP response")
	})
}
