package route

import (
	"reflect"
	"sync"
	"test/pkg/action"
	"test/pkg/action/k8s"
	"test/pkg/action/knative"
	"test/pkg/action/tekton"
)

var route sync.Map

func Register(name string, action action.Action) {
	route.Store(name, action)
}

func RegisterAction(action action.Action) {
	name := reflect.TypeOf(action).Elem().Name()
	Register(name, action)
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
}
