package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-http-utils/logger"
	"github.com/gorilla/mux"
	"github.com/mikloslorinczi/infra-exec/common"
	"github.com/mikloslorinczi/infra-exec/infraserver"
)

var port int

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

	infraserver.TaskDB = infraserver.ConnectJSONDB("db.json")
	err := infraserver.TaskDB.Load()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	router := mux.NewRouter()
	router.NotFoundHandler = infraserver.Custom404()
	router.Use(infraserver.AuthCheck)
	router.HandleFunc("/api/task/list", infraserver.ListTasks).Methods("GET")
	router.HandleFunc("/api/task/query/{id}", infraserver.QueryTask).Methods("GET")
	router.HandleFunc("/api/task/add", infraserver.AddTask).Methods("POST")

	fmt.Printf("\nInfra server listening on PORT %v\n", port)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), logger.Handler(router, os.Stdout, logger.CommonLoggerType)))

}

func init() {
	flag.IntVar(&port, "port", 7474, "Server PORT")
	flag.StringVar(&common.AdminPass, "pass", "", "The Admin Password")
}

func printUsage() {
	fmt.Println()
	fmt.Println("Infra Server")
	fmt.Println()
	fmt.Println("Hosts the Task Database")
	fmt.Println("New tasks can be added with the Infra CLI")
	fmt.Println("Infra Client(s) periodically polls the Task Database for commands to execute")
	fmt.Println()
	fmt.Printf("Usage: %s [Option] [PORT or AdminPassword]", os.Args[0])
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println()
	fmt.Println()
	flag.PrintDefaults()
}
