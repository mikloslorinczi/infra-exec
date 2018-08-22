package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/mikloslorinczi/infra-exec/common"
	"github.com/pkg/errors"
)

const envFile = "client.env"

var (
	nodeName string
	nodeTags string
	period   int
)

func main() {

	status, err := initConfig()
	if err != nil {
		fmt.Println(err)
		printUsage()
		os.Exit(status)
	}

	fmt.Printf("\nInfra Client - %v initialized.\nWith the tags: %v\nPolling Infra Server @ %v every %v seconds\n\n", nodeName, nodeTags, common.APIURL, period)

	for running := true; running; {

		nextTick := time.Now().Add(time.Second * time.Duration(period))

		fmt.Println("Polling Infra Server...")
		task, ok, err1 := getTaskToExec()
		if err1 != nil || !ok {
			fmt.Println(err1)
			waitNext(nextTick)
			continue
		}

		fmt.Println("Matching Task found...")
		task.Node = nodeName
		response, err2 := claimTask(task)
		if err2 != nil {
			fmt.Println(err2)
			waitNext(nextTick)
			continue
		}

		fmt.Println(response.Msg)
		filePath, err3 := executeTask(task)
		if err3 != nil {
			fmt.Printf("Error during execution %v\n", err3)
			_, err4 := common.UpdateTaskStatus(task.ID, "Execution error")
			if err4 != nil {
				fmt.Printf("Error updating Task status %v\n", err4)
			}
			waitNext(nextTick)
			continue
		}

		fmt.Printf("Task executed, logs can be found at %v\n", filePath)
		_, err5 := common.UpdateTaskStatus(task.ID, "Executed")
		if err5 != nil {
			fmt.Printf("Error updating Task status %v", err5)
		}

		err6 := uploadLog(filePath, task.ID)
		if err6 != nil {
			fmt.Printf("Error uploading logfile %v\n%v\n", filePath, err6)
		} else {
			fmt.Println("Logfile successfully uploaded to the Infra Server")
		}

		waitNext(nextTick)
	}
}

func init() {
	flag.StringVar(&nodeName, "name", "", "Name of this node")
	flag.StringVar(&nodeTags, "tags", "", "Tags of this node")
	flag.IntVar(&period, "period", 0, "Poll period in seconds")
	flag.StringVar(&common.AdminPass, "pass", "", "Admin password")
	flag.StringVar(&common.APIURL, "u", "", "URL addres of the api")
}

func waitNext(t time.Time) {
	time.Sleep(time.Until(t))
}

func initConfig() (int, error) {

	err := common.LoadEnv(envFile)
	if err != nil {
		fmt.Printf("Cannot load %v %v\n", envFile, err)
	}
	common.AdminPass = os.Getenv("ADMIN_PASSWORD")
	common.APIURL = os.Getenv("API_URL")
	nodeName = os.Getenv("NODE_NAME")
	nodeTags = os.Getenv("NODE_TAGS")
	envPeriod, err := strconv.Atoi(os.Getenv("NODE_PERIOD"))
	if err != nil {
		return 1, errors.Wrap(err, "Error loading poll-period from environment")
	}

	flag.Parse()

	if period == 0 {
		period = envPeriod
	}

	if common.AdminPass == "" {
		return 1, fmt.Errorf("No ADMIN_PASSWORD found, please export it or save it to %v or set it with the -pass flag", envFile)
	}
	if common.APIURL == "" {
		return 1, fmt.Errorf("No API_URL found, please export it or save it to %v or set it with the -u flag", envFile)
	}
	if nodeName == "" {
		return 1, fmt.Errorf("No NODE_NAME found, please export it or save it to %v or set it with the -name flag", envFile)
	}

	return 0, nil
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
