package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mikloslorinczi/infra-exec/common"
)

// RootCmd is the Cobra "entrypoint" iclient
var RootCmd = &cobra.Command{
	Use:   "iclient",
	Short: "Infra Client is a node that executes Task polled from the Infra Server",
	Long: `
	Infra Client

	Periodically polls the Infra Server for new Tasks
	If new task is found, and its tags are matching with the Client's tags
	The Client will try to execute the Task and send the logfile back to the Server
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\nInfra Client - %v initialized.\nWith the tags: %v\nPolling Infra Server @ %v every %v seconds\n\n",
			viper.GetString("nodeName"),
			viper.GetString("nodeTags"),
			viper.GetString("apiURL"),
			viper.GetInt("period"))
		runClient()
	},
}

func runClient() {
	for running := true; running; {
		// Calculate next tick
		nextTick := time.Now().Add(time.Second * time.Duration(viper.GetInt("period")))
		// Polling Infra Server for Tasks
		fmt.Println("Polling Infra Server...")
		task, ok, err := getTaskToExec()
		if err != nil || !ok {
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("No matching Task found...")
			waitNext(nextTick)
			continue
		}
		// Try to claim matching Task
		fmt.Println("Matching Task found...")
		task.Node = viper.GetString("nodeName")
		response, err := claimTask(task)
		if err != nil {
			fmt.Println(err)
			waitNext(nextTick)
			continue
		}
		// Try to execute Task
		fmt.Println(response.Msg)
		filePath, err := executeTask(task)
		if err != nil {
			fmt.Printf("Error during execution %v\n", err)
			_, statusErr := common.UpdateTaskStatus(task.ID, "Execution error")
			if statusErr != nil {
				fmt.Printf("Error updating Task status %v\n", statusErr)
			}
			waitNext(nextTick)
			continue
		}
		// Try to update Task's status
		fmt.Printf("Task executed, logs can be found at %v\n", filePath)
		_, err = common.UpdateTaskStatus(task.ID, "Executed")
		if err != nil {
			fmt.Printf("Error updating Task status %v", err)
		}
		// Try to upload logfile to Infra Server
		if err = uploadLog(filePath, task.ID); err != nil {
			fmt.Printf("Error uploading logfile %v\n%v\n", filePath, err)
		} else {
			fmt.Println("Logfile successfully uploaded to the Infra Server")
		}
		// Wait next tick
		waitNext(nextTick)
	}
}

func waitNext(t time.Time) {
	time.Sleep(time.Until(t))
}
