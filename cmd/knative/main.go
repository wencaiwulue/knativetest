package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"knativetest/pkg/action"
	"knativetest/pkg/action/k8s"
	"knativetest/pkg/action/knative"
	"knativetest/pkg/action/tekton"
	"log"
	"net/http"
	"reflect"
	"strings"
	"sync"
)

func main() {
	mux := http.NewServeMux()
	//mux.HandleFunc("/webshell", webshell.Handle)
	mux.Handle("/CreateDockerImageAction", &tekton.CreateDockerImageAction{})
	mux.HandleFunc("/{action}", ServeHTTP)
	log.Fatal(http.ListenAndServe(":80", mux))
}

var route sync.Map

func ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var actionStr = strings.ReplaceAll(req.URL.Path, "/", "")
	var value, found = route.Load(actionStr)
	if !found {
		fmt.Printf("Can't found handler for action: %v\n", actionStr)
	}
	var a, ok = value.(action.Action)
	if !ok {
		fmt.Printf("Can't cast to a handler, value: %v\n", value)
	}
	var buf, _ = ioutil.ReadAll(req.Body)
	_ = json.Unmarshal(buf, &a)
	values := map[string]interface{}{
		"http.request":      req,
		"http.request.body": string(buf)}
	ctx := context.WithValue(req.Context(), "http", values)
	var response = a.Process(ctx)
	var bytes, _ = json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(bytes)
}

func RegisterAction(action action.Action) {
	name := reflect.TypeOf(action).Elem().Name()
	route.Store(name, action)
}

func init() {
	RegisterAction(&k8s.CreateNamespaceAction{})
	RegisterAction(&k8s.DeleteNamespaceAction{})
	RegisterAction(&knative.CreateServiceAction{})
	RegisterAction(&knative.CreateRevisionAction{})
	RegisterAction(&knative.CreateHPAAction{})
	RegisterAction(&knative.CreateKPAAction{})
	RegisterAction(&knative.ListServiceAction{})
	RegisterAction(&knative.ListRevisionAction{})
	RegisterAction(&knative.InvokeAction{})
	RegisterAction(&tekton.ListResourceAction{})
	RegisterAction(&tekton.CreateKanikoClusterTask{})
	RegisterAction(&tekton.CreateKanikoTaskRun{})
	RegisterAction(&k8s.KvAction{})
}
