package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/mikloslorinczi/infra-exec/common"
)

// taskCmd represents the task command
// Without a subcommand it lists the Tasks from the Infra Server
var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Task is the base command to list, add and query tasks.",
	Long: `
	Simply calling <icli task> will list all the tasks
	<icli task add> will start the New Task dialog
	<icli task query ...ID> will query individual Task(s) by their ID (or part of ID)
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("List of Tasks\n\n")
		tasks, err := common.GetTasks()
		if err != nil {
			log.Fatalf("Error geting Task List %v\n", err)
		}
		for _, task := range tasks {
			fmt.Printf("\nID : %v\nNode : %v\nTags : %v\nStatus : %v\nCommand : %v\n", task.ID, task.Node, task.Tags, task.Status, task.Command)
		}
	},
}

func init() {
	RootCmd.AddCommand(taskCmd)
}
