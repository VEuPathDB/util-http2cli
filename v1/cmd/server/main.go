package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/veupathdb/http2cli/v1/pkg/h2c"
	"github.com/x-cray/logrus-prefixed-formatter"
)

var version = "development-build"

func init() {
	format := new(prefixed.TextFormatter)
	format.FullTimestamp = true
	format.TimestampFormat = time.RFC3339Nano

	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(format)
}

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
	logrus.Info("Starting server on port " + sPort)
	logrus.Fatal(http.ListenAndServe("0.0.0.0:"+sPort, r))
}
