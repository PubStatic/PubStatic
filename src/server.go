package main

import (
	"encoding/json"
	"fmt"
	"github.com/PubStatic/PubStatic/activityPub"
	"github.com/PubStatic/PubStatic/wellknown"
	"io"
	"net/http"
	"strings"
	"time"
)

var server = http.Server{}

func configureServer() {
	mux := http.NewServeMux()
	fileserver := http.FileServer(http.Dir("./static"))

	mux.HandleFunc("/.well-known/webfinger", func(w http.ResponseWriter, r *http.Request) {
		webfinger := wellknown.GetWebfinger(r.Host, settings.ActivityPubSettings.UserName)

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
		acceptHeader := r.Header.Get("Accept")

		if len(acceptHeader) > 0 && (strings.Contains(acceptHeader, "application/json") || strings.Contains(acceptHeader, "application/activity+json")) {
			actor := activityPub.GetActor(r.Host, settings, publicKeyPem)

			jsonData, err := json.Marshal(actor)
			if err != nil {
				logger.Error("InternalServer Error. Could not retrieve Actor.", err)
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
			logger.Error(err)
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		var activity activityPub.Activity

		err = json.Unmarshal(body, &activity)
		if err != nil {
			logger.Error(err)
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}

		err = activityPub.ReceiveActivity(activity, r.Header, r.Host)

		if err != nil {
			logger.Error(err)
			http.Error(w, "Invalid signature", http.StatusForbidden)
			return
		}
	})

	server = http.Server{
		Addr:    fmt.Sprintf(":%d", settings.ServerSettings.Port),
		Handler: gzipHandler(loggingMiddleware(mux)),
	}
}

func startServer() {
	logger.Infof("Starting server at port %d", settings.ServerSettings.Port)
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

		logger.Debugf("Headers: %s", r.Header)

		logger.Infof("%s %s %v",
			r.Method,
			r.URL,
			time.Since(start))
	})
}
