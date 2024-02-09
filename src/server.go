package main

import (
	"fmt"
	"net/http"
	"github.com/PubStatic/PubStatic/activityPub"
)

var port = 80

func configureFileServer() {
	fileserver := http.FileServer(http.Dir("./static"))

	http.Handle("/", fileserver)
}

func configureActivityPubServer() {
	http.HandleFunc("/.well-known/webfinger", func(w http.ResponseWriter, r *http.Request) {
		activityPub.GetWebfinger()
	})
}

func startServer(){
	logger.Infof("Starting server at port %d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		logger.Fatal(err)
	}
}
