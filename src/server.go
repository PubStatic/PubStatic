package main

import (
	"encoding/json"
	"fmt"
	"github.com/PubStatic/PubStatic/wellknown"
	"net/http"
)

func configureFileServer() {
	fileserver := http.FileServer(http.Dir("./static"))

	http.Handle("/", fileserver)
}

func configureActivityPubServer() {
	http.HandleFunc("/.well-known/webfinger", func(w http.ResponseWriter, r *http.Request) {
		webfinger := wellknown.GetWebfinger(r.Host, userName)

		jsonData, err := json.Marshal(webfinger)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		w.Write(jsonData)
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

	http.HandleFunc("/nodeinfo/2.1", func(w http.ResponseWriter, r *http.Request) {
		nodeInfo := wellknown.GetNodeInfo2_1(version)

		jsonData, err := json.Marshal(nodeInfo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		w.Write(jsonData)
	})
}

func startServer() {
	logger.Infof("Starting server at port %d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		logger.Fatal(err)
	}
}
