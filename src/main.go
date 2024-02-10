package main

import (
	"crypto/rand"
	"crypto/rsa"
	"os"
	"github.com/PubStatic/PubStatic/activityPub"
	"github.com/PubStatic/PubStatic/repository"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()
var port = 8080
var version = "0.0.1"
var userName = "blog"
var summary = "Hello World!"
var publicKeyPem = ""

func main() {
	generateKeys()
	configureServer()
	startServer()
}

func init() {
	logger.Level = logrus.TraceLevel
}

func generateKeys() {
	fileContentPublicKey, _ := repository.ReadFile("publicKey.pem")
	fileContentPrivateKey, _ := repository.ReadFile("privateKey.pem")
	publicKeyPem = fileContentPublicKey

	if publicKeyPem == "" || fileContentPrivateKey == "" {
		key, rsaErr := rsa.GenerateKey(rand.Reader, 2048)
		if rsaErr != nil {
			logger.Fatal("RSA error")
			os.Exit(1)
		}

		publicKeyPemString, pemErr := activityPub.GetPublicKeyPem(*key)
		privateKeyPemString, pemPrivErr := activityPub.GetPrivateKeyPem(*key)

		if pemErr != nil || pemPrivErr != nil {
			logger.Fatal("PEM error")
			os.Exit(1)
		}

		publicKeyPem = *publicKeyPemString

		repository.WriteFile("publicKey.pem", publicKeyPem)
		repository.WriteFile("privateKey.pem", privateKeyPemString)
	}
}
