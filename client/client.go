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
		fmt.Println()
		fmt.Println("No ADMIN_PASSWORD found, please set it with the -pass flag, or export it explicitly to the environment.")
		printUsage()
		os.Exit(1)
	}

	if nodeName == "" {
		printUsage()
		os.Exit(1)
	}

	for running := true; running; {
		nextTick := time.Now().Add(time.Second * time.Duration(period))

		tasks, err := fetchTasks()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("Tasks\n%v\n", tasks)

		time.Sleep(nextTick.Sub(time.Now()))
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
