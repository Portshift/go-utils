package healthz

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	listenAddress string
	isReady       bool
}

func (s *Server) SetIsReady(isReady bool) {
	s.isReady = isReady
}

// Start starts the server run
func (s *Server) Start() {
	log.Infof("Starting healthz server. listenAddr=%v", s.listenAddress)

	http.HandleFunc("/healthz/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/healthz/ready", func(w http.ResponseWriter, r *http.Request) {
		if s.isReady {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	})

	go func() {
		if err := http.ListenAndServe(s.listenAddress, nil); err != nil {
			log.WithError(err).Error("Failed to serve.")
		}
	}()
}

func NewHealthServer(listenAddress string) *Server {
	return &Server{
		listenAddress: listenAddress,
	}
}
