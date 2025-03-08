/**
Copyright 2025 Google LLC.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// simple request counter metric
var pingCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "ping_request_count",
		Help: "No of request handled by Ping handler",
	},
)

// check if this is a COR preflight request
func isPreflight(r *http.Request) bool {
	return r.Method == "OPTIONS" &&
	  r.Header.Get("Origin") != "" &&
	  r.Header.Get("Access-Control-Request-Method") != ""
  }

func ping(w http.ResponseWriter, req *http.Request) {
	// increment req counter
	pingCounter.Inc()
	log.Printf("Handling request from %s", req.RemoteAddr)
	// lazy CORS handling
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// lazy way to return CORS preflight headers
	if isPreflight(req){
		 w.Header().Set("Access-Control-Allow-Methods","GET, OPTIONS")
		 if reqHeaders, found := req.Header["Access-Control-Request-Headers"]; found {
			w.Header().Set("Access-Control-Allow-Headers",strings.Join(reqHeaders,", "))
		 }
	}
	fmt.Fprintf(w, "pong: "+runtime.GOARCH+"\n")
}

func main() {
	log.SetOutput(os.Stdout)

	prometheus.MustRegister(pingCounter)

	http.HandleFunc("/ping", ping)
	http.Handle("/metrics", promhttp.Handler())
	log.Print("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
	log.Print("Server running on port 8080")
}