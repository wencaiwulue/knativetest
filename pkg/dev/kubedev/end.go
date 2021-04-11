package kubedev

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
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
			return errors.Errorf("%q requires at least 1 argument\n", cmd.CommandPath())
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("ending...")
		log.Println(endOption)
	},
}
