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
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	timeStart  time.Time
	duration   time.Duration
}

func (lrw loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
	lrw.duration = time.Since(lrw.timeStart)
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
	mux.HandleFunc("/", s.handleVersion)
	mux.HandleFunc("/v1/jobs", s.handleJobs)

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

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := loggingResponseWriter{
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
		handler.ServeHTTP(lrw, r)
		log.WithFields(
			logrus.Fields{
				"status":     lrw.statusCode,
				"statustext": http.StatusText(lrw.statusCode),
				"duration":   lrw.duration.String(),
			},
		).Infof("HTTP response")
	})
}
