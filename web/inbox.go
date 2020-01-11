package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tyrm/go-activitypub-relay/models"
)

func HandleInbox(w http.ResponseWriter, r *http.Request) {
	logger.Tracef("%v", r)

	// Get Body
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Decode Activity
	var activity models.Activity
	err = json.Unmarshal(b, &activity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ignore empty Actor. Not a properly formatted Activity
	if activity.Actor == "" {
		http.Error(w, "access denied", http.StatusUnauthorized)
		return
	}

	if activity.Type != "Follow" {
		http.Error(w, "access denied", http.StatusUnauthorized)
		
	}

	fmt.Printf("%s", activity)
}