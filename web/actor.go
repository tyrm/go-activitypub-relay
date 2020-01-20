package web

import (
	"bytes"
	"encoding/asn1"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"time"

	"github.com/tyrm/go-activitypub-relay/activitypub"
)

func HandleActor(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// Create Public Key Block
	asn1Bytes, err := asn1.Marshal(serverRSA.PublicKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Errorf("Error Marshaling Public Key: %s", err.Error())
	}

	var pubPEM = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	var pubBuffer bytes.Buffer
	err = pem.Encode(&pubBuffer, pubPEM)

	webfinger := activitypub.Actor{
		Context: "https://www.w3.org/ns/activitystreams",
		Endpoints: &activitypub.Endpoints{
			SharedInbox: fmt.Sprintf("https://%s/inbox", r.Host),
		},
		Followers: fmt.Sprintf("https://%s/followers", r.Host),
		Following: fmt.Sprintf("https://%s/following", r.Host),
		Inbox:     fmt.Sprintf("https://%s/inbox", r.Host),
		Name:      "PettingZoo Relay",
		Type:      "Application",
		ID:        fmt.Sprintf("https://%s/actor", r.Host),
		PublicKey: activitypub.PublicKey{
			ID:           fmt.Sprintf("https://%s/actor#main-key", r.Host),
			Owner:        fmt.Sprintf("https://%s/actor", r.Host),
			PublicKeyPem: pubBuffer.String(),
		},
		Summary:           "ActivityRelay bot",
		PreferredUsername: "relay",
		URL:               fmt.Sprintf("https://%s/actor", r.Host),
	}

	js, err := json.Marshal(webfinger)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	elapsed := time.Since(start)
	logger.Infof("REQUEST HandleActor () %s", elapsed)
}
