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

func Register0(action controller.Action) {
	name := reflect.TypeOf(action).Elem().Name()
	Register(name, action)
}

func init() {
	Register0(&k8s.CreateNamespaceAction{})
	Register0(&k8s.CreateNamespaceAction{})
	Register0(&knative.CreateServiceAction{})
	Register0(&knative.ListServiceAction{})
	Register0(&tekton.ListResourceAction{})
	Register0(&tekton.CreateKanikoClusterTask{})
	Register0(&tekton.CreateKanikoTaskRun{})
}

//func read(path string, ch chan string) {
//	list, err := ioutil.ReadDir(path)
//	if err != nil {
//		fmt.Printf("read dir: %v, error: %v\n", path, err)
//		return
//	}
//	for _, j := range list {
//		if j.IsDir() {
//			read(filepath.Join(path, j.Name()), ch)
//		} else {
//			if strings.Contains(j.Name(), "Action") {
//				ch <- strings.TrimSuffix(j.Name(), ".go")
//			}
//		}
//	}
//}
