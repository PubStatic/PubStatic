package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/PubStatic/PubStatic/wellknown"
)

var port = 80

func configureFileServer() {
	fileserver := http.FileServer(http.Dir("./static"))

	http.Handle("/", fileserver)
}

func configureActivityPubServer() {
	http.HandleFunc("/.well-known/webfinger", func(w http.ResponseWriter, r *http.Request) {
		wellknown.GetWebfinger()
	})

	http.HandleFunc("/.well-known/nodeinfo", func(w http.ResponseWriter, r *http.Request) {
		nodeInfoLink := wellknown.GetLinkToNodeInfo(r.Host)

		jsonData, err := json.Marshal(nodeInfoLink)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		w.Write(jsonData)
	})
}

func startServer(){
	logger.Infof("Starting server at port %d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		logger.Fatal(err)
	}
}
