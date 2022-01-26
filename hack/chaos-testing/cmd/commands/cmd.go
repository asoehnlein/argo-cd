package commands

import (
	"fmt"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"

	util2 "github.com/argoproj/argo-cd/v2/util"

	"github.com/spf13/cobra"

	"github.com/argoproj/argo-cd/v2/hack/chaos-testing/util"

	cmdutil "github.com/argoproj/argo-cd/v2/cmd/util"
	"github.com/argoproj/argo-cd/v2/util/cli"
)

const (
	cliName = "argocd-chaos-testing"
)

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	cli.SetLogFormat(cmdutil.LogFormat)
	cli.SetLogLevel(cmdutil.LogLevel)
}

func simulateReplicasChanges(token string, argoApi *util.ArgoApi, app v1alpha1.Application) {

	fmt.Printf("Pick application with name %s \n", app.Name)

	i := 0
	for {
		fmt.Printf("Execute replication emulation number: %v \n", i)

		fmt.Println("Scale to 2")

		_, err := argoApi.UpdateReplicas(token, 2, app.Name)
		if err != nil {
			fmt.Printf(err.Error())
			return
		}

		util2.Wait(30, func(bools chan<- bool) {

		})

		fmt.Println("Scale to 1")

		_, err = argoApi.UpdateReplicas(token, 1, app.Name)
		if err != nil {
			fmt.Printf(err.Error())
			return
		}

		i++

	}

}

// NewCommand returns a new instance of an argocd command
func NewCommand() *cobra.Command {

	var host, username, password string
	var threads int

	var command = &cobra.Command{
		Use:   cliName,
		Short: "Chaos testing",
		Run: func(c *cobra.Command, args []string) {
			argoApi := util.NewArgoApi(host)
			token, err := argoApi.GetToken(username, password)
			if err != nil {
				return
			}

			apps, err := argoApi.ListApplications(token)
			if err != nil {
				fmt.Printf(err.Error())
				return
			}

			for i := 0; i < threads; i++ {
				_, _ = argoApi.Sync(token, apps[i].Name)
				util2.Wait(5, func(bools chan<- bool) {
				})
				go simulateReplicasChanges(token, argoApi, apps[i])
			}

			select {}

		},
		DisableAutoGenTag: true,
	}

	command.PersistentFlags().StringVar(&host, "host", "http://localhost:8080", "Argo server host")
	command.PersistentFlags().StringVar(&username, "username", "admin", "Username")
	command.PersistentFlags().StringVar(&password, "password", "", "Argo server host")
	command.PersistentFlags().IntVar(&threads, "threads", 1, "Amount of threads")

	return command
}
