package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/juju/loggo"

	"github.com/tyrm/go-activitypub-relay/config"
	_ "github.com/tyrm/go-activitypub-relay/web"
)

var logger *loggo.Logger

func main() {
	newLogger := loggo.GetLogger("main")
	logger = &newLogger

	logger.Debugf("Gathering Config")
	_ = config.CollectConfig()

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	nch := make(chan os.Signal)
	signal.Notify(nch, syscall.SIGINT, syscall.SIGTERM)
	logger.Infof("%s", <-nch)

	logger.Infof("Done!")

}
