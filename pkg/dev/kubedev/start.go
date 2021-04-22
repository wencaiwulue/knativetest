package kubedev

import (
	"context"
	"encoding/json"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knativetest/pkg/dev/network"
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
		if startOption.Deployment == "" && startOption.Pod == "" {
			b, _ := json.Marshal(startOption)
			log.Fatalf("Deployment and Pod can not be null at the same time, options: %s", string(b))
		}
		client, err := util.InitClient(rootOption.Kubeconfig)
		if err != nil {
			log.Fatalf("init clientset error: %v", err)
		}

		deployment, err2 := client.ClientSet.AppsV1().Deployments(startOption.Namespace).Get(context.TODO(), startOption.Deployment, metav1.GetOptions{})
		if err2 != nil {
			log.Fatal(err2)
		}

		log.Println("prepare to update deployment")
		patch(client, deployment)
		log.Println("patch deployment ok, waiting for pod to be ready")

		util.WaitToBeStatus(client.ClientSet, startOption.Namespace, "deployments", "kubedev=debug", func(i interface{}) bool {
			return i.(*v1.Deployment).Status.ReadyReplicas == 1
		})
		log.Println("pod ready, finish patch deployment, try to synchronize file")
		go watch.Watch(client, startOption)
		network.PortForward(client, startOption)
		//if err = extra.Shell(client, startOption); err != nil {
		//    log.Printf("open shell error, info: %v", err)
		//}

	},
}

/**
1, update replica to 1
2, replace container to empty container
3, remove livenessProbe and readinessProbe
4, add deployment label and pod template label kubedev=debug
*/
func patch(client *util.ClientSet, r *v1.Deployment) {
	alreadyInDebug := true
	if r.Labels == nil || r.Labels["kubedev"] == "" {
		alreadyInDebug = false
		backup, _ := json.Marshal(r)

		deployment := r.DeepCopy()
		one := int32(1)
		deployment.Spec.Template.Spec.Containers[0].ReadinessProbe = nil
		deployment.Spec.Template.Spec.Containers[0].LivenessProbe = nil
		deployment.Spec.Template.Spec.Containers[0].Image = "naison/empty-container:latest"
		deployment.Spec.Replicas = &one
		if deployment.Labels == nil {
			deployment.Labels = make(map[string]string)
		}
		deployment.Labels["kubedev"] = "debug"
		if deployment.Annotations == nil {
			deployment.Annotations = make(map[string]string)
		}
		deployment.Annotations["revision/backup"] = string(backup)
		if deployment.Spec.Template.Labels == nil {
			deployment.Spec.Template.Labels = make(map[string]string)
		}
		deployment.Spec.Template.Labels["kubedev"] = "debug"
		// already in debug mode
		if _, err := client.ClientSet.AppsV1().Deployments(r.Namespace).Update(context.Background(), deployment, metav1.UpdateOptions{}); err != nil {
			log.Fatalln(err)
		}
	}
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
