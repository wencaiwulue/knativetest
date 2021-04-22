package kubedev

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knativetest/pkg/dev/util"
	"log"
)

var endOption = &endOptions{}

type endOptions struct {
	NameSpace  string
	Deployment string
	Pod        string
	Container  string
	LocalDir   string
	RemoteDir  string
}

func init() {
	endCmd.PersistentFlags().StringVarP(&endOption.NameSpace, "namespace", "n", "", "namespace")
	endCmd.Flags().StringVarP(&endOption.Deployment, "deployment", "d", "", "deployment")
	endCmd.Flags().StringVarP(&endOption.Container, "container", "c", "", "container to develop")
	Cmd.AddCommand(endCmd)
}

var endCmd = &cobra.Command{
	Use:   "end",
	Short: "end",
	Long:  `end`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.Errorf("%s requires at least 1 argument\n", cmd.CommandPath())
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("ending...")
		log.Println(endOption)
		client, err := util.InitClient(rootOption.Kubeconfig)
		if err != nil {
			log.Fatalf("init clientset error: %v", err)
		}
		rollback(client, endOption)
	},
}

func rollback(client *util.ClientSet, o *endOptions) {
	deployment, err := client.ClientSet.AppsV1().Deployments(startOption.Namespace).Get(context.TODO(), startOption.Deployment, metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}
	backup := deployment.Annotations["revision/backup"]
	d := v1.Deployment{}
	_ = json.Unmarshal([]byte(backup), &d)
	_ = client.ClientSet.AppsV1().Deployments(o.NameSpace).Delete(context.Background(), d.Name, metav1.DeleteOptions{})
	if _, err := client.ClientSet.AppsV1().Deployments(o.NameSpace).Create(context.Background(), &d, metav1.CreateOptions{}); err != nil {
		log.Fatalln(err)
	}
}
