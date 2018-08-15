package infraserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mikloslorinczi/infra-exec/common"
)

// TaskDB is the interface of the JSON db containing the Tasks.
var TaskDB DBI

// AuthCheck check the HTTP request if "adminpassword" can be found
// in the HEader, and if it equas to the preset AdminPass.
func AuthCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Header.Get("adminpassword") != common.AdminPass {
			res.WriteHeader(http.StatusUnauthorized)
			encoder := json.NewEncoder(res)
			err := encoder.Encode(common.ResponseMsg{Msg: "Error: Wrong password"})
			if err != nil {
				fmt.Printf("Error encoding message %v", err)
			}
		} else {
			next.ServeHTTP(res, req)
		}
	})
}

// Custom404 handles all unhandlet request with a common Wrong way 404 message.
func Custom404() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusNotFound)
		encoder := json.NewEncoder(res)
		err := encoder.Encode(common.ResponseMsg{Msg: "Wrong way 404 🐸"})
		if err != nil {
			fmt.Printf("Error encoding message %v", err)
		}
	})
}

// ListTasks handles the api/task/list request
func ListTasks(res http.ResponseWriter, req *http.Request) {
	msg, _ := TaskDB.queryAll()
	encoder := json.NewEncoder(res)
	err := encoder.Encode(msg)
	if err != nil {
		fmt.Printf("Error encoding message %v", err)
	}
}

// QueryTask handles the api/task/query/{id} requests
func QueryTask(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	task, ok := TaskDB.query(params["id"])
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		encoder := json.NewEncoder(res)
		err := encoder.Encode("Task not found")
		if err != nil {
			fmt.Printf("Error encoding message %v", err)
		}
	} else {
		encoder := json.NewEncoder(res)
		err := encoder.Encode(task)
		if err != nil {
			fmt.Printf("Error encoding message %v", err)
		}
	}
}

// AddTask handles the api/task/add requests
func AddTask(res http.ResponseWriter, req *http.Request) {
	var (
		command common.CommandObj
		newTask common.Task
	)
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&command)
	if err != nil {
		fmt.Printf("Cannot decode JSON\n%v", err)
		res.WriteHeader(http.StatusInternalServerError)
		encoder := json.NewEncoder(res)
		err = encoder.Encode(common.ResponseMsg{Msg: "Invalid data"})
		if err != nil {
			fmt.Printf("Error encoding message %v", err)
		}
		return
	}
	newTask.Command = command.Command
	newTask.Tags = command.Tags
	id, err := TaskDB.add(newTask)
	if err != nil {
		fmt.Printf("Cannot add new task to DB\n%v", err)
		res.WriteHeader(http.StatusInternalServerError)
		encoder := json.NewEncoder(res)
		err = encoder.Encode(common.ResponseMsg{Msg: "Cannot save to DB"})
		if err != nil {
			fmt.Printf("Error encoding message %v", err)
		}
	}
	encoder := json.NewEncoder(res)
	err = encoder.Encode(common.ResponseMsg{Msg: id})
	if err != nil {
		fmt.Printf("Error encoding message %v", err)
	}
}
