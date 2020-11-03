package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/veupathdb/http2cli/v1/pkg/h2c"
)

var version = "development-build"

func main() {
	config := new(h2c.Config)
	config.Version = version
	h2c.InitCLI(config)



	r := mux.NewRouter()
	r.Handle("/{tool}/{job-id}", h2c.NewJobController(config))
	r.Handle("/", h2c.NewMetaController(config))

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", r))
}
