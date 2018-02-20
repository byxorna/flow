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
	"github.com/gorilla/mux"
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

func (rl *responseLogger) WriteHeader(code int) {
	rl.statusCode = code
	rl.ResponseWriter.WriteHeader(code)
	rl.duration = time.Since(rl.timeStart)
}

type svr struct {
	config.Config
	router *mux.Router
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
	router := mux.NewRouter()

	s := svr{
		Config: c,
		router: router,
		store:  store,
	}

	// register http handlers
	s.router.HandleFunc("/", s.getVersion)
	v1api := s.router.PathPrefix("/v1").Subrouter()
	v1api.Path("/jobs/{namespace}").Methods("GET").
		HandlerFunc(s.jobs)
	v1api.Path("/jobs").Methods("GET").
		HandlerFunc(s.jobs)
	v1api.Path("/job").Methods("POST").
		HandlerFunc(s.postJob)
	v1api.Path("/job/{namespace}/{name}").Methods("GET", "DELETE").
		HandlerFunc(s.job)

	return &s, nil
}

// ListenAndServe calls http.ListenAndServe
func (s *svr) ListenAndServe() error {
	log.WithFields(
		logrus.Fields{"address": s.ServerListenAddr},
	).Infof("Listening for HTTP requests")
	return http.ListenAndServe(
		s.ServerListenAddr,
		logRequest(s.router),
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

func logRequest(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(
			logrus.Fields{
				"remoteaddr": r.RemoteAddr,
				"method":     r.Method,
				"url":        r.URL,
			},
		).Infof("HTTP request received")
		rl := responseLogger{
			ResponseWriter: w,
			timeStart:      time.Now(),
		}
		next.ServeHTTP(&rl, r)
		log.WithFields(
			logrus.Fields{
				"status":     rl.statusCode,
				"statustext": http.StatusText(rl.statusCode),
				"duration":   rl.duration.String(),
			},
		).Infof("HTTP response")
	}
}
