package activitypub

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/textproto"
	"sort"
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

func SignHeaders(headers map[string]string, key *rsa.PrivateKey, ourKeyId string) (string, error) {
	var headersLower map[string]string
	headersLower = make(map[string]string)
	for k, v := range headers {
		headersLower[strings.ToLower(k)] = v
	}

	var usedHeaders []string
	for k, _ := range headersLower {
		if k == "(request-target)" {
			continue
		}
		usedHeaders = append(usedHeaders, k)
	}
	sort.Strings(usedHeaders)
	usedHeaders = append([]string{"(request-target)"}, usedHeaders...)

	sig := map[string]string{
		"keyId":     ourKeyId,
		"algorithm": "rsa-sha256",
		"headers":   strings.Join(usedHeaders, " "),
	}

	sigstring := buildSigningString2(headers, usedHeaders)
	if sigstring == "" {
		logger.Tracef("SignHeaders(%v, %v, %v) (\"\", %s)", &headers, &key, &ourKeyId, "buildSigningString2 returned nothing")
		return "", errors.New("buildSigningString2 returned nothing")
	}

	signature, err := signSigningString(sigstring, key)
	if err != nil {
		logger.Tracef("SignHeaders(%v, %v, %v) (\"\", %s)", &headers, &key, &ourKeyId, err.Error())
		return "", err
	}
	sig["signature"] = string(signature)

	// Make String
	var chunks []string
	for k, v := range sig {
		chunks = append(chunks, fmt.Sprintf("%s=\"%s\"", k, v))
	}

	sigstr := strings.Join(chunks, ",")
	logger.Tracef("SignHeaders(%v, %v, %v) (\"%s\", nil)", &headers, &key, &ourKeyId, sigstr)
	return sigstr, nil
}

func buildSigningString2(headers map[string]string, usedHeaders []string) string {
	signingString := ""

	// put request-target first
	reqTgt, ok := headers["(request-target)"]
	if !ok {
		return ""
	}
	signingString = signingString + "(request-target): " + reqTgt + "\n"

	for i, h := range usedHeaders {
		if h == "(request-target)" {
			continue
		}

		key := strings.ToLower(h)
		logger.Tracef("trying to get process: %s", key)

		switch key {
		default:
			val := headers[textproto.CanonicalMIMEHeaderKey(h)]
			signingString = signingString + key + ": " + val
		}

		if i < len(usedHeaders)-1 {
			signingString = signingString + "\n"
		}
	}

	logger.Tracef("buildSigningString2(%v, %v) \"%s\"", &headers, &usedHeaders, signingString)
	return signingString
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

func signSigningString(sigstring string, key *rsa.PrivateKey) (string, error) {
	rng := rand.Reader
	hashed := sha256.Sum256([]byte(sigstring))

	signature, err := rsa.SignPKCS1v15(rng, key, crypto.SHA256, hashed[:])
	if err != nil {
		logger.Tracef("signSigningString(%v, %v) (\"\", %v)", len(sigstring), &key, err.Error())
		return "", err
	}

	sigstr := base64.StdEncoding.EncodeToString(signature)
	logger.Tracef("signSigningString(%v, %v) (\"%s\", nil)", len(sigstring), &key, sigstr)
	return sigstr, nil
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
