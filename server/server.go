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
	"github.com/mikloslorinczi/infra-exec/db"
)

const envFile = "server.env"

var port int

func main() {

	status, err1 := initConfig()
	if err1 != nil {
		fmt.Println(err1)
		printUsage()
		os.Exit(status)
	}

	taskDB = db.ConnectJSONDB("db.json")
	err2 := taskDB.Load()
	if err2 != nil {
		fmt.Println(err2)
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
	flag.IntVar(&port, "port", 8080, "Server PORT")
	flag.StringVar(&common.AdminPass, "pass", "", "The Admin Password")
}

func initConfig() (int, error) {

	err := common.LoadEnv(envFile)
	if err != nil {
		fmt.Printf("Cannot load %v\n%v\n", envFile, err)
	}

	common.AdminPass = os.Getenv("ADMIN_PASSWORD")
	// We don't need the error here, as if Atoi fails envPort will be 0.
	envPort, _ := strconv.Atoi(os.Getenv("SERVER_PORT")) // #nosec

	flag.Parse()

	// If envPort is set, and port is on default (not changed by flag), port will be set to envPort
	if envPort != 0 && port == 8080 {
		port = envPort
	}

	if common.AdminPass == "" {
		return 1, fmt.Errorf("No ADMIN_PASSWORD found, please export it or save it to %v or set it with the -pass flag", envFile)
	}

	return 0, nil
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
