package kubedev

import (
	"context"
	"encoding/json"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"knativetest/pkg/dev/util"
	"knativetest/pkg/dev/watch"
	"log"
)

var startOption = &StartOptions{}

type StartOptions struct {
	Namespace  string
	Deployment string
	Pod        string
	Container  string
	LocalDir   string
	RemoteDir  string
}

func init() {
	startCmd.PersistentFlags().StringVarP(&startOption.Namespace, "namespace", "n", "", "namespace")
	startCmd.Flags().StringVarP(&startOption.Deployment, "deployment", "d", "", "deployment")
	startCmd.Flags().StringVarP(&startOption.Container, "container", "c", "", "container")
	startCmd.Flags().StringVar(&startOption.LocalDir, "localdir", "", "local directory")
	startCmd.Flags().StringVar(&startOption.RemoteDir, "remotedir", "", "remote directory")
	Cmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start",
	Long:  `start`,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("starting...")
		log.Println(args)
		b, _ := json.Marshal(startOption)
		log.Println(string(b))
		client, err := util.InitClient(rootOption.Kubeconfig)
		if err != nil {
			log.Fatalf("init clientset error: %v\n", err)
		}

		deployment, err2 := client.ClientSet.AppsV1().Deployments(startOption.Namespace).Get(context.TODO(), startOption.Deployment, metav1.GetOptions{})
		if err2 != nil {
			log.Fatal(err2)
		}

		// todo rollback
		log.Println("prepare to update deployment")
		patch(client, deployment)
		log.Println("patch deployment ok, waiting for pod to be ready")

		util.WaitToBeStatus(client.ClientSet, startOption.Namespace, "deployments", "kubedev=debug", func(i interface{}) bool {
			return i.(*v1.Deployment).Status.ReadyReplicas == 1
		})
		log.Println("pod ready, finish patch deployment, try to synchronize file")
		watch.Watch(client, startOption)
	},
}

/**
1, update replica to 1
2, replace container to empty container
3, remove livenessProbe and readinessProbe
4, add deployment label and pod template label kubedev=debug
*/
func patch(client *util.ClientSet, r *v1.Deployment) {
	// todo why don't work
	_ = `{ "op": "remove", "path": "/spec/template/spec/containers/0/readinessProbe" }
            { "op": "remove", "path": "/spec/template/spec/containers/0/livenessProbe" },
            { "op": "add", "path": "/spec/template/metadata/labels/kubedev", "value":"debug" },
            { "op": "add", "path": "/metadata/labels/kubedev", "value":"debug" }`

	// already in debug mode
	alreadyInDebug := true
	one := int32(1)
	if r.Spec.Template.Spec.Containers[0].Image != "naison/empty-container:latest" || r.Spec.Replicas != &one {
		alreadyInDebug = false
		jsonPatch := `[
		  { "op": "replace", "path": "/spec/template/spec/containers/0/image", "value": "naison/empty-container:latest" },
          { "op": "replace", "path": "/spec/replicas", "value": 1 }
	    ]`
		res, err := client.ClientSet.AppsV1().Deployments(startOption.Namespace).
			Patch(context.TODO(), startOption.Deployment, types.JSONPatchType, []byte(jsonPatch), metav1.PatchOptions{})
		if err != nil {
			log.Fatalf("first patch deployment %v failed, error info: %v\n, response: %v", startOption.Deployment, err, res)
		}
	}
	if r.Labels["kubedev"] == "" || r.Spec.Template.Labels["kubedev"] == "" ||
		r.Spec.Template.Spec.Containers[0].LivenessProbe != nil ||
		r.Spec.Template.Spec.Containers[0].ReadinessProbe != nil {
		alreadyInDebug = false
		mergePatch := []string{
			`{"spec": {"template": {"metadata": {"labels":{"kubedev":"debug"}}}}}`,
			`{"spec": {"template": {"spec": {"containers": [{"name": "test","readinessProbe":null, "livenessProbe":null}]}}}}`,
			`{"metadata": {"labels": {"kubedev": "debug"}}}`,
		}
		for i, patch := range mergePatch {
			res, err := client.ClientSet.AppsV1().Deployments(startOption.Namespace).
				Patch(context.TODO(), startOption.Deployment, types.StrategicMergePatchType, []byte(patch), metav1.PatchOptions{})
			if err != nil {
				log.Fatalf("%v patch deployment %v failed, error info: %v\n, response: %v", startOption.Deployment, i, err, res)
			}
		}
	}
	/*s := `
	            metadata:
	              labels:
	                kubedev: debug
	            spec:
	              template:
	                metadata:
	                  labels:
	                    kubedev: debug
	                spec:
	                  containers:
	                  - name: test
	                    readinessProbe:
	                    livenessProbe:`
	  	res, err := util.Clients.ClientSet.AppsV1().Deployments(startOption.Namespace).
	  		Patch(context.TODO(), startOption.Deployment, types.StrategicMergePatchType, []byte(s), metav1.PatchOptions{})
	  	if err != nil {
	  		log.Fatalf("%v patch deployment %v third times failed, error info: %v\n, response: %v", startOption.Deployment, err, res)
	  	}*/
	if alreadyInDebug {
		log.Println("already in debug mode, don't needs to update")
	} else {
		log.Println("patch successfully")
	}
}

func (o *StartOptions) GetNamespace() string {
	return o.Namespace
}
func (o *StartOptions) GetDeployment() string {
	return o.Deployment
}
func (o *StartOptions) GetPod() string {
	return o.Pod
}
func (o *StartOptions) GetContainer() string {
	return o.Container
}
func (o *StartOptions) GetLocalDir() string {
	return o.LocalDir
}
func (o *StartOptions) GetRemoteDir() string {
	return o.RemoteDir
}
