package main

import (
	"log"
	"net/http"
	"test/knative/Route"
)

func main() {
	var d = Route.NewDispatch()
	http.Handle("/{action}", d)
	log.Fatal(http.ListenAndServe(":80", d))
}
