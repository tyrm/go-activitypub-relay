package web

import (
	"crypto/rsa"
	"net/http"
	"strings"

	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
	"github.com/juju/loggo"
	"github.com/ryanuber/go-glob"
)

var (
	logger *loggo.Logger
	serverRSA *rsa.PrivateKey
	templates *packr.Box

	descriptionText = ""
)

func Init(pk *rsa.PrivateKey) {
	newLogger := loggo.GetLogger("web")
	logger = &newLogger

	// load templates
	templates = packr.New("htmlTemplates", "./templates")

	// Save private key
	serverRSA = pk

	r := mux.NewRouter()
	r.HandleFunc("/", HandleIndex).Methods("GET")

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

func isAccepteType(r *http.Request, mimeType string) bool {

	for _, accepts := range r.Header["Accept"] {
		accept := strings.SplitN(accepts, ",", -1)
		for _, a := range accept {
			data := strings.SplitN(a, ";", -1)

			if glob.Glob(data[0], mimeType) {
				return true
			}
		}
	}
	return false
}
