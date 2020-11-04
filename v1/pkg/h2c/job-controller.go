package h2c

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
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
		logrus.Infof("Call to disallowed tool %s by %s", tool, r.RemoteAddr)
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

	if err := job.Run(); err != nil {
		logrus.Error(err)
		logrus.Warnf("Job %s with tool %s failed.", job.ID, tool)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Job failed: " + err.Error() + ".\n Error output will be in job directory."))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Job \"%s\" completed successfully.", job.ID)))
}

func (e *endpoint) toolAllowed(tool string) bool {
	for i := range e.config.Tools {
		if tool == e.config.Tools[i] {
			return true
		}
	}

	return false
}
