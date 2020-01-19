package activitypub

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/patrickmn/go-cache"
	"github.com/satori/go.uuid"
)

type Actor struct {
	Context           interface{} `json:"@context,omitempty"`
	ID                string      `json:"id,omitempty"`
	Type              string      `json:"type,omitempty"`
	Name              string      `json:"name,omitempty"`
	PreferredUsername string      `json:"preferredUsername,omitempty"`
	Summary           string      `json:"summary,omitempty"`
	Inbox             string      `json:"inbox,omitempty"`
	Endpoints         *Endpoints  `json:"endpoints,omitempty"`
	Followers         string      `json:"followers,omitempty"`
	Following         string      `json:"following,omitempty"`
	PublicKey         PublicKey   `json:"publicKey,omitempty"`
	Icon              Image       `json:"icon,omitempty"`
	Image             Image       `json:"image,omitempty"`
	URL               string      `json:"url,omitempty"`
}

type PublicKey struct {
	ID           string `json:"id,omitempty"`
	Owner        string `json:"owner,omitempty"`
	PublicKeyPem string `json:"publicKeyPem,omitempty"`
}

type Endpoints struct {
	SharedInbox string `json:"sharedInbox,omitempty"`
}

type Image struct {
	URL string `json:"url,omitempty"`
}

func (a *Actor) GetActorInbox() string {
	return a.Endpoints.SharedInbox
}

func (a *Actor) GetPublicKey() (*rsa.PublicKey, error) {
	var parsedKey interface{}
	var err error

	logger.Tracef("Trying to decode PEM: \"%v\"", a.PublicKey.PublicKeyPem)

	block, _ := pem.Decode([]byte(a.PublicKey.PublicKeyPem))
	if block == nil {
		logger.Errorf("failed to parse PEM block containing the public key %s", err)
		return nil, err
	}

	if parsedKey, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		logger.Errorf("Unable to parse RSA public key %s", err)
		return nil, err
	}

	var pubKey *rsa.PublicKey
	var ok bool
	if pubKey, ok = parsedKey.(*rsa.PublicKey); !ok {
		logger.Errorf("Unable to parse RSA public key")
		return nil, err
	}

	return pubKey, nil
}

func FetchActor(url string) (*Actor, error) {
	if a, found := cRemoteActors.Get(url); found {
		actor := a.(*Actor)
		logger.Tracef("FetchActor(%s) (%v, nil) [HIT]", url, actor.ID)
		return actor, nil
	}

	// Fetch Actor
	resp, err := http.Get(url)
	if err != nil {
		logger.Tracef("FetchActor(%s) (nil, %s) [MISS]", url, err)
		return nil, err
	}

	// Ready Body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Tracef("FetchActor(%s) (nil, %s) [MISS]", url, err)
		return nil, err
	}

	logger.Debugf("FetchActor > Remote Actor Payload > %s", body)
	var actor Actor
	err = json.Unmarshal([]byte(body), &actor)
	if err != nil {
		logger.Tracef("FetchActor(%s) (nil, %s) [MISS]", url, err)
		return nil, err
	}

	cRemoteActors.Set(url, &actor, cache.DefaultExpiration)
	logger.Tracef("FetchActor(%s) (%v, nil) [MISS]", url, actor.ID)
	return &actor, nil
}

func FollowRemoteActor(url string, reqHost string) {
	remoteActor, err := FetchActor(url)
	if err != nil {
		logger.Errorf("failed to fetch actor at %s: %s", url, err.Error())
		return
	}

	logger.Infof("Following %s", url)

	message := &Activity{
		Context: "https://www.w3.org/ns/activitystreams",
		Type:    "Follow",
		To:      []string{remoteActor.ID},
		Object:  remoteActor.ID,
		ID:      fmt.Sprintf("https://%s/activities/%s", reqHost, uuid.NewV4()),
		Actor:   fmt.Sprintf("https://%s/actor", reqHost),
	}

	PushMessageToActor(remoteActor, message, fmt.Sprint("https://%s/actor#main-key", reqHost))
}

func PushMessageToActor(actor *Actor, message *Activity, outKeyId string) {
	inbox := actor.GetActorInbox()

	// Encode message to JSON
	data, err := json.Marshal(message)
	if err != nil {
		logger.Errorf("Error encoding message to JSON: %s", err.Error())
		return
	}

	// Send POST
	logger.Debugf("POST %s << %s", inbox, data)
	buf := bytes.NewBuffer(data)
	req, err := http.NewRequest("POST", inbox, buf)
	if err != nil {
		logger.Errorf("Could not create request: %s", err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/activity+json")
	req.Header.Add("User-Agent", "ActivityRelay")

	// TODO sign message

	client := &http.Client{}
	resp, err := client.Do(req)

	// Read Body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Error reading response body: %s", err.Error())
		return
	}

	logger.Debugf("response %s >> (%d) %s", inbox, resp.StatusCode, body)
}
