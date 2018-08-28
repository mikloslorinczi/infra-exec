package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-http-utils/logger"
	"github.com/gorilla/mux"
	"github.com/mikloslorinczi/infra-exec/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// taskDB is the interface of the JSON db containing the Tasks.
var taskDB db.Controller

// RootCmd is the basecommand, it is the "entrypoint" of iserver.
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
		initDB()
		runServer()
	},
}

func initDB() {
	// Connect Task Database
	taskDB = db.NewJSONDB("db.json")
	// Load Tasks
	if err := taskDB.Load(); err != nil {
		log.Fatalf("Error loading Database db.json %v\n", err)
	}
}

func runServer() {

	router := mux.NewRouter()

	router.NotFoundHandler = custom404()

	router.Use(authCheck)

	router.PathPrefix("/api/logs/").Handler(http.StripPrefix("/api/logs/", http.FileServer(http.Dir("./logs/"))))

	router.HandleFunc("/api/tasks", listTasks).Methods("GET")
	router.HandleFunc("/api/task/{id}", queryTask).Methods("GET")
	router.HandleFunc("/api/tasks", addTask).Methods("POST")
	router.HandleFunc("/api/tasks/claim", claimTask).Methods("POST")
	router.HandleFunc("/api/task/{id}/{status}", updateTaskStatus).Methods("PUT")
	router.HandleFunc("/api/log/{id}", uploadLog).Methods("POST")

	fmt.Printf("\nInfra server listening on PORT %v\n", viper.GetInt("PORT"))

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(viper.GetInt("PORT")), logger.Handler(router, os.Stdout, logger.CommonLoggerType)))

}
