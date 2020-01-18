package activitypub

import (
	"time"

	"github.com/juju/loggo"
	"github.com/patrickmn/go-cache"
)

var (
	logger *loggo.Logger
	cRemoteActors *cache.Cache
)

func init() {
	newLogger := loggo.GetLogger("activitypub")
	logger = &newLogger

	// init cache
	cRemoteActors = cache.New(60*time.Minute, 60*time.Minute)
}