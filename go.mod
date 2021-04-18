module knativetest

go 1.15

require (
	github.com/docker/docker v1.13.1
	github.com/fsnotify/fsnotify v1.4.9
	github.com/google/uuid v1.1.2
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.0.0
	github.com/tektoncd/pipeline v0.18.1
	github.com/wencaiwulue/stream4go v0.0.4
	go.etcd.io/etcd v0.5.0-alpha.5.0.20200819165624-17cef6e3e9d5
	k8s.io/api v0.20.5
	k8s.io/apimachinery v0.20.5
	k8s.io/cli-runtime v0.19.0
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	k8s.io/klog v1.0.0
	k8s.io/klog/v2 v2.4.0 // indirect
	k8s.io/kubectl v0.19.0
	k8s.io/utils v0.0.0-20200729134348-d5654de09c73
	knative.dev/networking v0.0.0-20201103163404-b9f80f4537af
	knative.dev/serving v0.19.0
	sigs.k8s.io/yaml v1.2.0
)

// Pin k8s deps to v0.19.0
replace (
	k8s.io/api => k8s.io/api v0.19.0
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.19.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.0
	k8s.io/apiserver => k8s.io/apiserver v0.19.0
	k8s.io/client-go => k8s.io/client-go v0.19.0
	k8s.io/code-generator => k8s.io/code-generator v0.19.0
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20210323165736-1a6458611d18
)

replace (
	github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.8.1
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.5.4
	golang.org/x/sys => golang.org/x/sys v0.0.0-20210415045647-66c3f260301c
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
