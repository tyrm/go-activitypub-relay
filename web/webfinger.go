package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type WebFinger struct {
	Aliases []string `json:"aliases"`
	Links   []WebFingerLinks `json:"links"`
	Subject string   `json:"subject"`
}

type WebFingerLinks struct {
	Rel  string `json:"rel"`
	HREF string `json:"href"`
	Type string `json:"type"`
}

func HandleWebFinger(w http.ResponseWriter, r *http.Request) {
	subject := r.URL.Query().Get("resource")

	if subject != fmt.Sprintf("acct:relay@%s", r.Host) {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "{'error': 'user not found'}", http.StatusNotFound)
		return
	}

	actorURL := fmt.Sprintf("https://%s/actor", r.Host)

	webfinger := WebFinger{
		Aliases: []string{actorURL},
		Links: []WebFingerLinks{
			{
				HREF: actorURL,
				Rel: "self",
				Type: "application/activity+json",
			},
			{
				HREF: actorURL,
				Rel: "self",
				Type: "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"",
			},
		},
		Subject: subject,
	}

	js, err := json.Marshal(webfinger)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
