package web

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/tyrm/go-activitypub-relay/activitypub"
	"github.com/tyrm/go-activitypub-relay/models"
)

func HandleInbox(w http.ResponseWriter, r *http.Request) {
	// Get Body
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Decode Activity
	var activity activitypub.Activity
	err = json.Unmarshal(b, &activity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ignore empty Actor. Not a properly formatted Activity
	if activity.Actor == "" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// TODO Validate Signature


	// Block non Follow Requests from unapproved Instances
	actor, err := url.Parse(activity.Actor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	exists, err := activitypub.ApprovedInstanceExists(actor.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Debugf("Host approval %s %v", actor.Host, exists)
	if activity.Type != "Follow" && !exists {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		logger.Warningf("Unauthorized actor [%s] tried to %s", actor, activity.Type)
		return
	}

	logger.Debugf(" > payload > %s", b)

	remoteActor, err := activitypub.FetchActor(activity.Actor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// process message
	switch activity.Type {
	case "Follow":
		go HandleInboxFollow(remoteActor, &activity, r.Host)
	}

	// return ok
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}


func HandleInboxFollow(actor *activitypub.Actor, activity *activitypub.Activity, reqHost string) {
	inbox := actor.GetActorInbox()
	inboxURL, err := url.Parse(inbox)
	if err != nil {
		logger.Errorf("Error parsing url %s: %s", inbox, err.Error())
		return
	}

	// Check Blacklist
	blacklisted, err := activitypub.OnBlacklist(inboxURL.Host)
	if err != nil {
		logger.Errorf("Could not check blacklist for %s: %s", inboxURL.Host, err.Error())
		return
	}
	if blacklisted {
		logger.Warningf("Blocked instance %s tried to Follow", inboxURL.Host)
		return
	}

	// if not in database add it
	exists, err := activitypub.InstanceExists(inboxURL.Host)
	if err != nil {
		logger.Errorf("Could not check instances for %s: %s", inboxURL.Host, err.Error())
		return
	}
	if !exists {
		logger.Debugf("Adding new instance %s", inboxURL.Host)

		_, err := models.CreateInstance(inboxURL.Host)
		if err != nil {
			logger.Errorf("Could not add instance %s: %s", inboxURL.Host, err.Error())
			return
		}

		// If Object ends with /actor then send follow reply
		logger.Tracef("HasSuffix(%s, /actor)", activity.Object.(string))
		if strings.HasSuffix(activity.Object.(string), "/actor") {
			logger.Debugf("Sending Follow to remote actor %s", actor.ID)
			activitypub.FollowRemoteActor(actor.ID, reqHost)
		}
	}


}