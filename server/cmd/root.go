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
	Use:   "iserver [-flags]",
	Short: "Infra Server is responsible to distribute tasks to the Infra Clients",
	Long: `
Infra Server

The Infra Server is responsible to distribute tasks among the Infra Clients. Users can
issue new tasks with the help of the Infra Cli. Once the new task is saved to the db, it will
be retrived by clients when they poll the server. Clients can claim tasks and execute them,
updating its status and logs through the server.  
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
