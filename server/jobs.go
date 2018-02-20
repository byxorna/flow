package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/byxorna/flow/types/job"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

func (s *svr) jobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	jobs, err := s.store.GetJobs()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorJSON(err))
		return
	}
	json.NewEncoder(w).Encode(jobs)
}

func (s *svr) job(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	namespace, nsok := vars["namespace"]
	name, idok := vars["name"]
	if !nsok && !idok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorJSON(fmt.Errorf("namespace and id required")))
		return
	}
	id := job.ID{Namespace: namespace, Name: name}

	switch r.Method {
	case http.MethodGet:
		j, err := s.store.GetJob(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write(errorJSON(err))
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(j)
	case http.MethodDelete:
		j, err := s.store.DeleteJob(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write(errorJSON(err))
			return
		}
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(j)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

}

func (s *svr) postJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var j job.Spec
	var err error
	switch r.Header.Get("Content-Type") {
	case "application/json":
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&j)
		defer r.Body.Close()
	case "application/yaml":
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write(errorJSON(err))
			return
		}
		err = yaml.Unmarshal(b, &j)
	default:
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(errorJSON(err))
		return
	}

	log.Debugf("storing a job %v", j)

	err = s.store.SetJob(&j)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorJSON(err))
		return
	}
	w.WriteHeader(http.StatusCreated)

}
