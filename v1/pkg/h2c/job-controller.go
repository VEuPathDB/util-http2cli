package h2c

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const (
	hb = ".\n"
)

func NewJobController(config *Config) http.Handler {
	return &endpoint{config: config}
}

type endpoint struct {
	config *Config
}

func (e *endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		defer r.Body.Close()
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Path Variables
	vars := mux.Vars(r)
	tool := vars["tool"]
	if !e.toolAllowed(tool) {
		logrus.Info("Call to disallowed tool %s by %s", tool, r.RemoteAddr)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	args := make([]string, 0, 64)
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&args); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	job := NewJob()
	job.ID = vars["job-id"]
	job.Tool = tool
	job.Args = args
	job.Config = e.config

	done := make(chan error, 1)
	defer close(done)

	if err := job.Run(done); err != nil {
		logrus.Error(err)
		logrus.Warnf("Failed to run %s.", tool)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	for {
		select {
		case err := <-done:
			if err == nil {
				w.WriteHeader(200)
			} else {
				logrus.Infof("Job %s failed with %s", vars["job-id"], err.Error())
				w.Write([]byte(err.Error()))
				w.WriteHeader(500)
			}
			return
		default:
			time.Sleep(10 * time.Second)
			w.Write([]byte(hb))
		}
	}
}

func (e *endpoint) toolAllowed(tool string) bool {
	for i := range e.config.Tools {
		if tool == e.config.Tools[i] {
			return true
		}
	}

	return false
}
