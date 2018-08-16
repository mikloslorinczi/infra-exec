package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mikloslorinczi/infra-exec/common"
)

// taskDB is the interface of the JSON db containing the Tasks.
var taskDB DBI

// authCheck check the HTTP request if "adminpassword" can be found
// in the HEader, and if it equas to the preset AdminPass.
func authCheck(next http.Handler) http.Handler {
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

// custom404 handles all unhandlet request with a common Wrong way 404 message.
func custom404() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusNotFound)
		encoder := json.NewEncoder(res)
		err := encoder.Encode(common.ResponseMsg{Msg: "Wrong way 404 🐸"})
		if err != nil {
			fmt.Printf("Error encoding message %v", err)
		}
	})
}

// listTasks handles the api/task/list request
func listTasks(res http.ResponseWriter, req *http.Request) {
	msg, _ := taskDB.queryAll()
	encoder := json.NewEncoder(res)
	err := encoder.Encode(msg)
	if err != nil {
		fmt.Printf("Error encoding message %v", err)
	}
}

// queryTask handles the api/task/query/{id} requests
func queryTask(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	task, ok := taskDB.query(params["id"])
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

// addTask handles the api/task/add requests
func addTask(res http.ResponseWriter, req *http.Request) {
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
	id, err := taskDB.add(newTask)
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

// claimTask handles the api/task/claim request
func claimTask(res http.ResponseWriter, req *http.Request) {
	var taskToClaim common.Task
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&taskToClaim)
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
	taskToClaim.Status = "Assigned"
	_, err = taskDB.update(taskToClaim.ID, taskToClaim)
	if err != nil {
		fmt.Printf("Cannot claim Task\n%v\n", err)
		res.WriteHeader(http.StatusInternalServerError)
		encoder := json.NewEncoder(res)
		err = encoder.Encode(common.ResponseMsg{Msg: "Database error"})
		if err != nil {
			fmt.Printf("Error encoding message %v", err)
		}
		return
	}
	encoder := json.NewEncoder(res)
	err = encoder.Encode(common.ResponseMsg{Msg: "Successfully claimed task " + taskToClaim.ID})
	if err != nil {
		fmt.Printf("Error encoding message %v", err)
	}
}
