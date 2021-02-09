package route

import (
	"reflect"
	"sync"
	"test/knative/controller"
	"test/knative/controller/k8s"
	"test/knative/controller/knative"
	"test/knative/controller/tekton"
)

var route sync.Map

func Register(name string, action controller.Action) {
	route.Store(name, action)
}

func RegisterAction(action controller.Action) {
	name := reflect.TypeOf(action).Elem().Name()
	Register(name, action)
}

func init() {
	RegisterAction(&k8s.CreateNamespaceAction{})
	RegisterAction(&k8s.CreateNamespaceAction{})
	RegisterAction(&knative.CreateServiceAction{})
	RegisterAction(&knative.ListServiceAction{})
	RegisterAction(&tekton.ListResourceAction{})
	RegisterAction(&tekton.CreateKanikoClusterTask{})
	RegisterAction(&tekton.CreateKanikoTaskRun{})
}
