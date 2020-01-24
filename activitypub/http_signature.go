package activitypub

import (
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/yukimochi/httpsig"
)

type Signature struct {
	KeyID     string
	Algorithm string
	Headers   []string
	Signature string
}

func IsSignatureValid(r *http.Request) error {
	// Add missing host header
	r.Header.Set("Host", r.Host)
	verifier, err := httpsig.NewVerifier(r)
	if err != nil {
		return err
	}

	// Get Algorythm
	sig := parseSignature(r.Header["Signature"][0])
	logger.Debugf("sigdata: %s", sig)

	var algo httpsig.Algorithm

	switch sig.Algorithm {
	case "rsa-sha256":
		algo = httpsig.RSA_SHA256
	default:
		logger.Warningf("unkonwn signing algorythm: %s", sig.Algorithm)
	}

	// Get Actor's Key
	pubKeyId := verifier.KeyId()
	actor, err := FetchActor(pubKeyId)
	if err != nil {
		logger.Tracef("IsSignatureValid(%v) (nil)", &r)
		return err
	}

	pubKey, err := actor.GetPublicKey()
	if err != nil {
		logger.Tracef("IsSignatureValid(%v) (%s)", &r, err.Error())
		return err
	}

	err = verifier.Verify(pubKey, algo)
	if err != nil {
		logger.Tracef("IsSignatureValid(%v) (%s)", &r, err.Error())
		return err
	}
	return err
}

func AppendSignature(request *http.Request, body *[]byte, KeyID string, publicKey *rsa.PrivateKey) error {
	hash := sha256.New()
	hash.Write(*body)
	b := hash.Sum(nil)

	request.Header.Set("Digest", "SHA-256="+base64.StdEncoding.EncodeToString(b))
	request.Header.Set("Host", request.Host)

	signer, _, err := httpsig.NewSigner([]httpsig.Algorithm{httpsig.RSA_SHA256}, []string{httpsig.RequestTarget, "Host", "Date", "Digest", "Content-Type"}, httpsig.Signature)
	if err != nil {
		return err
	}
	err = signer.SignRequest(publicKey, KeyID, request)
	if err != nil {
		return err
	}
	return nil
}

func parseSignature(signature string) *Signature {
	newSig := Signature{}
	chunks := strings.SplitN(signature, ",", -1)

	for _, chunk := range chunks {
		field := strings.SplitN(chunk, "=", 2)

		switch field[0] {
		case "keyId":
			newSig.KeyID = field[1][1 : len(field[1])-1]
		case "algorithm":
			newSig.Algorithm = field[1][1 : len(field[1])-1]
		case "headers":
			headers := strings.SplitN(field[1][1:len(field[1])-1], " ", -1)
			newSig.Headers = headers
		case "signature":
			newSig.Signature = field[1][1 : len(field[1])-1]
		}
	}

	return &newSig
}
