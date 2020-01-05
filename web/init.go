package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/juju/loggo"
)

var (
	logger *loggo.Logger
	publicKeyPem string
)

func Init(pkp string) {
	newLogger := loggo.GetLogger("web")
	logger = &newLogger

	r := mux.NewRouter()

	r.HandleFunc("/.well-known/nodeinfo", HandleNodeInfoWellKnown).Methods("GET")
	r.HandleFunc("/.well-known/webfinger", HandleWebFinger).Methods("GET")

	r.HandleFunc("/actor", HandleActor).Methods("GET")
	r.HandleFunc("/nodeinfo/2.0.json", HandleNodeInfo).Methods("GET")

	go http.ListenAndServe(":8080", r)
}
