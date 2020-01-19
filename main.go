package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"

	"github.com/tyrm/go-activitypub-relay/models"
	"github.com/tyrm/go-activitypub-relay/web"
)

var logger *loggo.Logger

func main() {
	//
	newLogger := loggo.GetLogger("main")
	logger = &newLogger

	logger.Debugf("Gathering Config")
	cfg := CollectConfig()

	err := loggo.ConfigureLoggers(cfg.LoggerConfig)
	if err != nil {
		fmt.Printf("Error configurting Logger: %s", err.Error())
		return
	}
	_, err = loggo.ReplaceDefaultWriter(loggocolor.NewWriter(os.Stderr))
	if err != nil {
		fmt.Printf("Error configurting Color Logger: %s", err.Error())
		return
	}

	// Init Database
	models.Init(cfg.DBEngine)
	defer models.Close()

	// store server key
	var serverRSA *rsa.PrivateKey

	// Check for private key
	privateKey, err := models.ReadConfig("private_key")
	if err == sql.ErrNoRows {
		logger.Errorf("Private Key not in config", )

		// Generate RSA key
		reader := rand.Reader
		bitSize := 4096

		serverRSA, err = rsa.GenerateKey(reader, bitSize)
		if err != nil {
			logger.Errorf("Error Generating Private Key: %s", err.Error())
		}

		// Create Private Key Block
		var privPEM = &pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(serverRSA),
		}

		var privBuffer bytes.Buffer
		err = pem.Encode(&privBuffer, privPEM)

		logger.Errorf("%v", privBuffer.String())
		newPrivateKey, err := models.CreateConfig("private_key", privBuffer.String())
		if err != nil {
			logger.Errorf("Error Saving Public Key: %s", err.Error())
		}

		privateKey = newPrivateKey

	} else if err != nil {
		logger.Errorf("Error Reading Private Key: %s", err.Error())
	} else {
		// decode private key
		privBlock, _ := pem.Decode([]byte(privateKey.Value))
		if privBlock == nil || privBlock.Type != "PRIVATE KEY" {
			logger.Errorf("failed to decode PEM block containing private key")
		}

		priv, err := x509.ParsePKCS1PrivateKey(privBlock.Bytes)
		if err != nil {
			logger.Errorf("Error parsing Public Key: %v", err)
		}

		serverRSA = priv
	}

	// Init Web Server
	web.Init(serverRSA)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	nch := make(chan os.Signal)
	signal.Notify(nch, syscall.SIGINT, syscall.SIGTERM)
	logger.Infof("%s", <-nch)

	logger.Infof("Done!")

}
