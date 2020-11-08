package Route

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"test/knative/controllers"
	"test/knative/controllers/k8s"
	"test/knative/controllers/knative"
	"test/knative/controllers/tekton"
)

type Dispatch struct {
	route sync.Map
}

func NewDispatch() *Dispatch {
	var d = &Dispatch{}
	d.init()
	return d
}

func (d *Dispatch) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var action = strings.ReplaceAll(req.RequestURI, "/", "")
	var value, _ = d.route.Load(action)
	var a = value.(controllers.Action)
	var buf, _ = ioutil.ReadAll(req.Body)
	_ = json.Unmarshal(buf, &a)
	var response = a.Process(req.Context())
	var bytes, _ = json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(bytes)
}

func (d *Dispatch) Registry(name string, action controllers.Action) {
	d.route.Store(name, action)
}

func (d *Dispatch) init() {
	d.Registry("createNamespace", &k8s.CreateNamespaceAction{})
	d.Registry("createService", &knative.CreateServiceAction{})
	d.Registry("list", &tekton.ListResourceAction{})
}
