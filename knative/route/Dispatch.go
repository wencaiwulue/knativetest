package route

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"test/knative/controller"
)

type Dispatch struct {
	http.Handler
}

func (d *Dispatch) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var action = strings.ReplaceAll(req.RequestURI, "/", "")
	var value, _ = route.Load(action)
	var a, ok = value.(controller.Action)
	fmt.Printf("%v", ok)
	var buf, _ = ioutil.ReadAll(req.Body)
	_ = json.Unmarshal(buf, &a)
	var response = a.Process(req.Context())
	var bytes, _ = json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(bytes)
}

