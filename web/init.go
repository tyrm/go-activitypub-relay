package web

import (
	"crypto/rsa"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/juju/loggo"
)

var (
	logger *loggo.Logger
	serverRSA *rsa.PrivateKey
)

func Init(pk *rsa.PrivateKey) {
	newLogger := loggo.GetLogger("web")
	logger = &newLogger

	// Save private key
	serverRSA = pk

	r := mux.NewRouter()

	r.HandleFunc("/.well-known/nodeinfo", HandleNodeInfoWellKnown).Methods("GET")
	r.HandleFunc("/.well-known/webfinger", HandleWebFinger).Methods("GET")

	r.HandleFunc("/actor", HandleActor).Methods("GET")
	r.HandleFunc("/inbox", HandleInbox).Methods("POST")
	r.HandleFunc("/nodeinfo/2.0.json", HandleNodeInfo).Methods("GET")

	go func() {
		err := http.ListenAndServe(":8080", r)
		if err != nil {
			logger.Errorf("Could not start web server %s", err.Error())
		}
	}()

	logger.Tracef("Init(%v)", &pk)
}
