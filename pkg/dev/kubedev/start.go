package kubedev

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"knativetest/pkg/dev/util"
	"knativetest/pkg/dev/watch"
	"log"
)

var startOption = &startOptions{}

type startOptions struct {
	NameSpace  string
	Deployment string
	Pod        string
	Container  string
	LocalDir   string
	RemoteDir  string
}

func init() {
	startCmd.PersistentFlags().StringVarP(&startOption.NameSpace, "namespace", "n", "", "namespace")
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
		if len(args) < 1 {
			return errors.Errorf("%s requires at least 1 argument\n", cmd.CommandPath())
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("starting...")
		b, _ := json.Marshal(startOption)
		log.Println(string(b))
		client, err := util.InitClient(rootOption.Kubeconfig)
		if err != nil {
			log.Fatalf("init clientset error: %v\n", err)
		}

		_ = `{ "op": "remove", "path": "/spec/template/spec/containers/0/readinessProbe" }
            { "op": "remove", "path": "/spec/template/spec/containers/0/livenessProbe" },
            { "op": "add", "path": "/spec/template/metadata/labels/kubedev", "value":"debug" },
            { "op": "add", "path": "/metadata/labels/kubedev", "value":"debug" }`
		mergePatch := []string{
			`{"spec": {"template": {"metadata": {"labels":{"kubedev":"debug"}}}}}`,
			`{"spec": {"template": {"spec": {"containers": [{"name": "test","readinessProbe":null, "livenessProbe":null}]}}}}`,
			`{"metadata": {"labels": {"kubedev": "debug"}}}`,
		}
		jsonPatch := `[
		  { "op": "replace", "path": "/spec/template/spec/containers/0/image", "value": "naison/empty-container:latest" },
          { "op": "replace", "path": "/spec/replicas", "value": 1 }
	    ]`
		r, e := client.ClientSet.AppsV1().Deployments(startOption.NameSpace).Get(context.TODO(), startOption.Deployment, metav1.GetOptions{})
		if e != nil {
			log.Fatal(err)
		} else {
			log.Println(r.Name)
		}
		res, err := client.ClientSet.AppsV1().Deployments(startOption.NameSpace).
			Patch(context.TODO(), startOption.Deployment, types.JSONPatchType, []byte(jsonPatch), metav1.PatchOptions{})
		if err != nil {
			log.Fatalf("first patch deployment %v failed, error info: %v\n, response: %v", startOption.Deployment, err, res)
		}
		for i, patch := range mergePatch {
			res, err = client.ClientSet.AppsV1().Deployments(startOption.NameSpace).
				Patch(context.TODO(), startOption.Deployment, types.StrategicMergePatchType, []byte(patch), metav1.PatchOptions{})
			if err != nil {
				log.Fatalf("%v patch deployment %v failed, error info: %v\n, response: %v", startOption.Deployment, i, err, res)
			}
		}

		log.Println("patch deployment ok, waiting for pod to be ready")

		util.WaitToBeStatus(client.ClientSet, startOption.NameSpace, "deployments", "kubedev=debug", func(i interface{}) bool {
			return i.(*v1.Deployment).Status.ReadyReplicas == 1
		})
		log.Println("pod ready, finish patch deployment, try to synchronize file")
		watch.Watch(startOption.LocalDir)
	},
}
