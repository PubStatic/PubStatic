package main

import (
	"fmt"
	"net/http"
	"github.com/sirupsen/logrus"
)

func main() {
    logger := logrus.New()

	port := 8080

	fileserver := http.FileServer(http.Dir("./static"))

	http.Handle("/", fileserver)

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello!")
	})

	logger.Infof("Starting server at port %d", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		logger.Fatal(err)
	}
}
