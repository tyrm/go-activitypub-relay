package activitypub

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/textproto"
	"strings"
)

type Signature struct {
	KeyID     string
	Algorithm string
	Headers   []string
	Signature string
}

func IsSignatureValid(r *http.Request, actorURI string) error {
	actor, err := FetchActor(actorURI)
	if err != nil {
		logger.Tracef("IsSignatureValid(%v) (nil)", &r)
		return err
	}

	sig := parseSignature(r.Header["Signature"][0])
	logger.Debugf("sigdata: %s", sig)

	sigstring := buildSigningString(r, sig.Headers)
	logger.Debugf("sigstring : \"%s\"", sigstring)

	sigalg := strings.Split(sig.Algorithm, "-")
	logger.Debugf("sign alg: %s, hash alg: %s", sigalg[0], sigalg[1])

	sigdatadata, err := base64.StdEncoding.DecodeString(sig.Signature)
	if err != nil {
		logger.Tracef("IsSignatureValid(%v) (%s)", &r, err.Error())
		return err
	}

	var alg crypto.Hash
	var hashed [32]byte

	switch sigalg[1] {
	case "sha256":
		alg = crypto.SHA256
		hashed = sha256.Sum256([]byte(sigstring))
	default:
		logger.Warningf("unkonwn signing algorythm: %s", sigalg[0])
	}

	pubKey, err := actor.GetPublicKey()
	if err != nil {
		logger.Tracef("IsSignatureValid(%v) (%s)", &r, err.Error())
		return err
	}

	err = rsa.VerifyPKCS1v15(pubKey, alg, hashed[:], sigdatadata)
	if err != nil {
		logger.Tracef("IsSignatureValid(%v) (%s)", &r, err.Error())
		return err
	}

	logger.Tracef("IsSignatureValid(%v) (nil)", &r)
	return nil
}

func buildSigningString(r *http.Request, usedHeaders []string) string {
	signingString := ""

	for i, h := range usedHeaders {
		key := strings.ToLower(h)
		logger.Tracef("trying to get process: %s", key)

		switch key {
		case "(request-target)":
			signingString = signingString + key + ": " + strings.ToLower(r.Method) + " " + r.URL.Path
		case "host":
			signingString = signingString + key + ": " + r.Host
		default:
			val := r.Header[textproto.CanonicalMIMEHeaderKey(h)]
			signingString = signingString + key + ": " + val[0]
		}

		if i < len(usedHeaders)-1 {
			signingString = signingString + "\n"
		}
	}

	logger.Tracef("buildSigningString(%v, %v) \"%s\"", &r, &usedHeaders, signingString)
	return signingString
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
