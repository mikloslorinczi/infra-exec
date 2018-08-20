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

	taskDB = connectJSONDB("db.json")
	err := taskDB.Load()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	router := mux.NewRouter()
	router.NotFoundHandler = custom404()
	router.Use(authCheck)
	router.HandleFunc("/api/task/list", listTasks).Methods("GET")
	router.HandleFunc("/api/task/query/{id}", queryTask).Methods("GET")
	router.HandleFunc("/api/task/add", addTask).Methods("POST")
	router.HandleFunc("/api/task/claim", claimTask).Methods("POST")
	router.HandleFunc("/api/task/status/{id}/{status}", updateTaskStatus).Methods("POST")
	router.HandleFunc("/api/log/upload/{id}", uploadLog).Methods("POST")
	router.HandleFunc("/api/log/download/{id}", downloadLog).Methods("GET")

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
