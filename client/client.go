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

	for running := true; running; {
		nextTick := time.Now().Add(time.Second * time.Duration(period))

		tasks, err := common.GetTasks()
		if err != nil {
			fmt.Println(err)
		} else {
			unassigned, found := getUnassigned(tasks)
			if found {
				task, ok := getFirstMatch(unassigned)
				if ok {
					task.Node = nodeName
					response, err := claimTask(task)
					if err != nil {
						fmt.Println(err)
					} else {
						fmt.Println(response.Msg)
						filePath, err := executeTask(task)
						if err != nil {
							fmt.Printf("Error during execution\n%v", err)
						} else {
							fmt.Printf("Task executed, logs can be found at %v\n", filePath)
						}
					}
				}
			}
		}

		time.Sleep(time.Until(nextTick))
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
