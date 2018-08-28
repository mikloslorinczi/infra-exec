package cmd

import (
	"fmt"

	"github.com/mikloslorinczi/infra-exec/common"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query <icli task query ID...>",
	Short: "Queryes the Infra Server for Task by ID",
	Run: func(cmd *cobra.Command, args []string) {
		for _, query := range args {
			task, err := queryTask(query)
			if err != nil {
				fmt.Printf("Error querying Task %v %v\n", query, err)
			} else {
				fmt.Printf("\nQuery : %v\nTask Found!\nID      : %v\nNode    : %v\nTags    : %v\nStatus  : %v\nCommand : %v\n", query, task.ID, task.Node, task.Tags, task.Status, task.Command)
			}
		}
	},
}

func init() {
	taskCmd.AddCommand(queryCmd)
}

func queryTask(query string) (common.Task, error) {
	var task common.Task
	taskJSON, err := common.SendRequest("GET", viper.GetString("apiUrl")+"/task/"+query, nil)
	if err != nil {
		return task, errors.Wrap(err, "Cannot query task")
	}
	if err := common.FromJSON(&task, taskJSON); err != nil {
		return task, errors.Wrap(err, "Cannot query task")
	}
	return task, nil
}
