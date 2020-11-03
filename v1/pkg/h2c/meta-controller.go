package h2c

import (
	"encoding/json"
	"net/http"
	"time"
)

func NewMetaController(config *Config) http.Handler {
	return &meta{
		config: config,
		start:  time.Now(),
	}
}

type meta struct {
	config *Config
	start  time.Time
}

func (m *meta) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	out := new(metaInfo)
	out.Version = m.config.Version
	out.Uptime = time.Now().Sub(m.start).String()
	out.Tools = m.config.Tools

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

type metaInfo struct {
	Version string   `json:"version"`
	Uptime  string   `json:"uptime"`
	Tools   []string `json:"tools"`
}
