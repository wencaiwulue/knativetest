package main

import (
	"log"
	"net/http"
	"test/knative/route"
)

func main() {
	var d = route.Dispatch{}
	http.Handle("/{action}", &d)
	log.Fatal(http.ListenAndServe(":80", &d))
}
