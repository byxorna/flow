package server

import (
	"encoding/json"
	"net/http"
)

func (s *svr) handleJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jobs, err := s.store.GetJobs()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorJSON(err))
		return
	}
	json.NewEncoder(w).Encode(jobs)
}
