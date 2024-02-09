package main

import (
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()
var port = 8080
var version = "0.0.1"
var userName = "blog"
var summary = "Hello World!"
var publicKeyPem = ""

func main() {
	configureServer()
	startServer()
}

func init() {
	logger.Level = logrus.TraceLevel
}
