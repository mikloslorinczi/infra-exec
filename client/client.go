package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/mikloslorinczi/infra-exec/common"
)

var (
	nodeName string
	nodeTags string
	period   int
)

func main() {

	flag.Parse()

	initClient()

	fmt.Printf("\nInfra Client - %v initialized.\nWith the tags: %v\nPolling Infra Server @ %v every %v seconds\n\n", nodeName, nodeTags, common.APIURL, period)

	for running := true; running; {
		nextTick := time.Now().Add(time.Second * time.Duration(period))
		fmt.Println("Polling server ...")
		task, ok, err := getTaskToExec()
		if err != nil {
			fmt.Println(err)
		}
		if ok {
			fmt.Println("Matching Task found ...")
			task.Node = nodeName
			response, err := claimTask(task)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(response.Msg)
				filePath, err := executeTask(task)
				if err != nil {
					fmt.Printf("Error during execution %v\n", err)
					_, err := common.UpdateTaskStatus(task.ID, "Execution error")
					if err != nil {
						fmt.Printf("Error updating Task status %v\n", err)
					}
				} else {
					fmt.Printf("Task executed, logs can be found at %v\n", filePath)
					_, err := common.UpdateTaskStatus(task.ID, "Executed")
					if err != nil {
						fmt.Printf("Error updating Task status %v", err)
					}
					err = uploadLog(filePath, task.ID)
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}

		time.Sleep(time.Until(nextTick))
	}
}

func initClient() {
	if common.AdminPass == "" {
		common.AdminPass = os.Getenv("ADMIN_PASSWORD")
	}
	if common.AdminPass == "" {
		fmt.Printf("\nNo ADMIN_PASSWORD found, please set it with the -pass flag, or export it explicitly to the environment.")
		printUsage()
		os.Exit(1)
	}
	if nodeName == "" {
		fmt.Printf("\nThe node must have a name, please set it with the -name flag\n")
		printUsage()
		os.Exit(1)
	}
}

func init() {
	flag.StringVar(&nodeName, "name", "", "Name of this node")
	flag.StringVar(&nodeTags, "tags", "", "Tags of this node")
	flag.IntVar(&period, "period", 15, "Poll period in seconds")
	flag.StringVar(&common.AdminPass, "pass", "", "Admin password")
	flag.StringVar(&common.APIURL, "u", "http://localhost:7474/api", "URL addres of the api")
}

func printUsage() {
	fmt.Println()
	fmt.Println("Infra Client")
	fmt.Println()
	fmt.Println("Periodically polls the Infra Server for new Tasks")
	fmt.Println("If new task is found, and its tags are matching with the Client's tags")
	fmt.Println("The Client will try to execute the Task and send the logfile to the Server")
	fmt.Println()
	fmt.Printf("Usage: %s [Option] [string]", os.Args[0])
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println()
	fmt.Println()
	flag.PrintDefaults()
}
