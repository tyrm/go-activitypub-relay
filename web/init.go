package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/juju/loggo"
)

var logger *loggo.Logger

func init() {
	newLogger := loggo.GetLogger("web")
	logger = &newLogger

	r := mux.NewRouter()

	rWellKnown := r.PathPrefix("/.well-known").Subrouter()
	rWellKnown.HandleFunc("/nodeinfo", HandleNodeInfoWellKnown).Methods("GET")

	go http.ListenAndServe(":8080", r)
}
