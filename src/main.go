package main

import (
	"crypto/rand"
	"crypto/rsa"
	"os"
	"github.com/PubStatic/PubStatic/activityPub"
	"github.com/PubStatic/PubStatic/models"
	"github.com/PubStatic/PubStatic/repository"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/yaml.v2"
)

var logger = logrus.New()
var settings = models.Settings{}

var version = "0.0.1"
var publicKeyPem = ""
var mongoConnectionString = ""

func main() {
	loadConfig()
	generateKeys()

	configureServer()
	startServer()
}

func init() {
	logger.Level = logrus.TraceLevel
}

func loadConfig(){
	settingsString, _ := repository.ReadFile("settings.dev.yaml")

	if settingsString == "" {
		settingsString, _ = repository.ReadFile("settings.yaml")
	}

	err := yaml.Unmarshal([]byte(settingsString), &settings)
	if err != nil {
		panic(err)
	}

	mongoConnectionString = os.Getenv("MONGODB")
}

func generateKeys() {
	fileContentPrivateKey, privErr := repository.ReadMongo[models.KeyValue]("Actor", "Key", bson.D{{"key", "privateKey"}}, mongoConnectionString)
	fileContentPublicKey, pubErr := repository.ReadMongo[models.KeyValue]("Actor", "Key", bson.D{{"key", "publicKey"}}, mongoConnectionString)

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

		repository.WriteMongo("Actor", "Key", models.KeyValue{Key: "publicKey", Value: publicKeyPem}, mongoConnectionString)
		repository.WriteMongo("Actor", "Key", models.KeyValue{Key: "privateKey", Value: privateKeyPemString}, mongoConnectionString)
	}
}
