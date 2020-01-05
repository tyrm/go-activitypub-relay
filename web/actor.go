package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Actor struct {
	Context           string         `json:"@context"`
	Endpoints         ActorEndpoints `json:"endpoints"`
	Followers         string         `json:"followers"`
	Following         string         `json:"following"`
	Inbox             string         `json:"inbox"`
	Name              string         `json:"name"`
	Type              string         `json:"type"`
	ID                string         `json:"id"`
	PublicKey         ActorPublicKey `json:"publicKey"`
	Summary           string         `json:"summary"`
	PreferredUsername string         `json:"preferredUsername"`
	URL               string         `json:"url"`
}

type ActorEndpoints struct {
	SharedInbox string `json:"sharedInbox"`
}

type ActorPublicKey struct {
	ID           string `json:"id"`
	Owner        string `json:"owner"`
	PublicKeyPem string `json:"publicKeyPem"`
}

func HandleActor(w http.ResponseWriter, r *http.Request) {

	webfinger := Actor{
		Context: "https://www.w3.org/ns/activitystreams",
		Endpoints: ActorEndpoints{
			SharedInbox: fmt.Sprintf("https://%s/inbox", r.Host),
		},
		Followers: fmt.Sprintf("https://%s/followers", r.Host),
		Following: fmt.Sprintf("https://%s/following", r.Host),
		Inbox: fmt.Sprintf("https://%s/inbox", r.Host),
		Name: "PettingZoo Relay",
		Type: "Application",
		ID: fmt.Sprintf("https://%s/actor", r.Host),
		PublicKey: ActorPublicKey{
			ID: fmt.Sprintf("https://%s/actor#main-key", r.Host),
			Owner: fmt.Sprintf("https://%s/actor", r.Host),
			PublicKeyPem: publicKeyPem,
		},
		Summary: "ActivityRelay bot",
		PreferredUsername: "relay",
		URL: fmt.Sprintf("https://%s/actor", r.Host),
	}

	js, err := json.Marshal(webfinger)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
