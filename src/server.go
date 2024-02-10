package main

import (
	"encoding/json"
	"fmt"
	"github.com/PubStatic/PubStatic/activityPub"
	"github.com/PubStatic/PubStatic/wellknown"
	"io"
	"net/http"
	"time"
)

var server = http.Server{}

func configureServer() {
	mux := http.NewServeMux()
	fileserver := http.FileServer(http.Dir("./static"))

	loggedHandler := loggingMiddleware(mux)

	server = http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: loggedHandler,
	}

	mux.HandleFunc("/.well-known/webfinger", func(w http.ResponseWriter, r *http.Request) {
		webfinger := wellknown.GetWebfinger(r.Host, userName)

		jsonData, err := json.Marshal(webfinger)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		w.Write(jsonData)
	})

	mux.HandleFunc("/.well-known/nodeinfo", func(w http.ResponseWriter, r *http.Request) {
		nodeInfoLink := wellknown.GetLinkToNodeInfo(r.Host)

		jsonData, err := json.Marshal(nodeInfoLink)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		w.Write(jsonData)
	})

	mux.HandleFunc("/nodeinfo/2.1", func(w http.ResponseWriter, r *http.Request) {
		nodeInfo := wellknown.GetNodeInfo2_1(version)

		jsonData, err := json.Marshal(nodeInfo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		w.Write(jsonData)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		acceptHeader := r.Header["Accept"]

		if len(acceptHeader) > 0 && r.Header["Accept"][0] == "application/json" {
			actor := activityPub.GetActor(r.Host, userName, userName, summary, publicKeyPem)

			jsonData, err := json.Marshal(actor)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")

			w.Write(jsonData)
		} else {
			fileserver.ServeHTTP(w, r)
		}
	})

	mux.HandleFunc("/inbox", func(w http.ResponseWriter, r *http.Request) {

		body, err := io.ReadAll(io.Reader(r.Body))
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		var activity activityPub.Activity

		err = json.Unmarshal(body, &activity)
		if err != nil {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}

		if activityPub.ReceiveActivity(activity, r.Header) != nil {
			http.Error(w, "Invalid signature", http.StatusForbidden)
			return
		}
	})
}

func startServer() {
	logger.Infof("Starting server at port %d", port)
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}

// Logging middleware logs the incoming request method and URL
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call the next handler in the chain
		next.ServeHTTP(w, r)

		logger.Debugf("%s %s %v",
			r.Method,
			r.URL,
			time.Since(start))
	})
}
