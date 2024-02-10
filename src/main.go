package main

import (
	"crypto/rand"
	"crypto/rsa"
	"os"
	"github.com/PubStatic/PubStatic/activityPub"
	"github.com/PubStatic/PubStatic/repository"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

var logger = logrus.New()
var port = 8080
var version = "0.0.1"
var userName = "blog"
var summary = "Hello World!"
var publicKeyPem = ""
var mongoConnectionString = ""

func main() {
	generateKeys()

	configureServer()
	startServer()
}

func init() {
	logger.Level = logrus.TraceLevel
}

func generateKeys() {
	fileContentPrivateKey, privErr := repository.ReadMongo[KeyValue]("Actor", "Key", bson.D{{"key", "privateKey"}}, mongoConnectionString)
	fileContentPublicKey, pubErr := repository.ReadMongo[KeyValue]("Actor", "Key", bson.D{{"key", "publicKey"}}, mongoConnectionString)

	if privErr != nil || pubErr != nil {
		logger.Error("Mongo error")
	}

	publicKeyPem = fileContentPublicKey.Value

	if publicKeyPem == "" || fileContentPrivateKey.Value == "" {
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

		repository.WriteMongo("Actor", "Key", KeyValue{Key: "publicKey", Value: publicKeyPem}, mongoConnectionString)
		repository.WriteMongo("Actor", "Key", KeyValue{Key: "privateKey", Value: privateKeyPemString}, mongoConnectionString)
	}
}
