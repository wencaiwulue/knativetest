package main

import (
	"log"
	"net/http"
	"test/knative/controller/tekton"
	"test/knative/route"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/CreateDockerImageAction", &tekton.CreateDockerImageAction{})
	mux.Handle("/{action}", &route.Dispatch{})
	log.Fatal(http.ListenAndServe(":80", mux))
}
