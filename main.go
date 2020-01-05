package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/juju/loggo"

	"github.com/tyrm/go-activitypub-relay/config"
	_ "github.com/tyrm/go-activitypub-relay/web"
)

var logger *loggo.Logger

func main() {
	//
	newLogger := loggo.GetLogger("main")
	logger = &newLogger

	logger.Debugf("Gathering Config")
	cfg := config.CollectConfig()

	err := loggo.ConfigureLoggers(cfg.LoggerConfig)
	if err != nil {
		fmt.Printf("Error configurting Logger: %s", err.Error())
		return
	}

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	nch := make(chan os.Signal)
	signal.Notify(nch, syscall.SIGINT, syscall.SIGTERM)
	logger.Infof("%s", <-nch)

	logger.Infof("Done!")

}
