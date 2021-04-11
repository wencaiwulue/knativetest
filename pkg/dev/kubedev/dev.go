package kubedev

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"log"
)

var rootOption = &RootOptions{}

type RootOptions struct {
	Kubeconfig string
}

func init() {
	Cmd.PersistentFlags().StringVar(&rootOption.Kubeconfig, "kubeconfig", "", "kubeconfig file")
}

var Cmd = &cobra.Command{
	Use:   "kubedev",
	Short: "kubedev",
	Long:  `kubedev`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.Errorf("%q requires at least 1 argument\n", cmd.CommandPath())
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("kubedev starting...")
		log.Println(rootOption)
	},
}
