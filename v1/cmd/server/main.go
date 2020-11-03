package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/veupathdb/http2cli/v1/pkg/h2c"
)

var version = "development-build"

func main() {
	config := new(h2c.Config)
	config.Version = version
	h2c.InitCLI(config)

	r := mux.NewRouter()
	r.Handle("/{tool}/{job-id}", h2c.NewJobController(config)).
		Methods(http.MethodPost)
	r.Handle("/", h2c.NewMetaController(config)).
		Methods(http.MethodGet)

	sPort := strconv.FormatUint(uint64(config.ServerPort), 10)
	log.Println("Starting server on port " + sPort)
	log.Fatal(http.ListenAndServe("0.0.0.0:" + sPort, r))
}
