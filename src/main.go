package main

import (
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()
var port = 80
var version = "0.0.1"
var userName = "blog"

func main() {
	port = 8080 // Override of default port 80

	configureFileServer()
	configureActivityPubServer()

	startServer()
}

func init() {
	logger.Level = logrus.TraceLevel
}
