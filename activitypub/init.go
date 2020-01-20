package activitypub

import (
	"crypto/rsa"
	"time"

	"github.com/juju/loggo"
	"github.com/patrickmn/go-cache"
)

var (
	logger *loggo.Logger
	cRemoteActors *cache.Cache
	serverRSA *rsa.PrivateKey
)

func Init(sr *rsa.PrivateKey) {
	newLogger := loggo.GetLogger("activitypub")
	logger = &newLogger

	// Store server RSA
	serverRSA = sr

	// init cache
	cRemoteActors = cache.New(60*time.Minute, 60*time.Minute)
}