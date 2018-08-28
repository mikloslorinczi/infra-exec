package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/viper"

	"github.com/gorilla/mux"
	"github.com/mikloslorinczi/infra-exec/common"
)

// authCheck check the HTTP request if "APITOKEN" can be found
// in the Header, and if its value is equals to the one set by Viper.
func authCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Header.Get("APITOKEN") != viper.GetString("apiToken") {
			res.WriteHeader(http.StatusUnauthorized)
			encoder := json.NewEncoder(res)
			if err := encoder.Encode(common.ResponseMsg{Msg: "Error: Wrong password"}); err != nil {
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
		if err := encoder.Encode(common.ResponseMsg{Msg: "Wrong way 404 🐸"}); err != nil {
			fmt.Printf("Error encoding message %v", err)
		}
	})
}

// listTasks handles the api/tasks request.
func listTasks(res http.ResponseWriter, req *http.Request) {
	msg, _ := taskDB.QueryAll()
	encoder := json.NewEncoder(res)
	if err := encoder.Encode(msg); err != nil {
		fmt.Printf("Error encoding message %v", err)
	}
}

// queryTask handles the api/task/{id} requests.
func queryTask(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	task, ok := taskDB.Query(params["id"])
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		encoder := json.NewEncoder(res)
		if err := encoder.Encode("Task not found"); err != nil {
			fmt.Printf("Error encoding message %v", err)
		}
	} else {
		encoder := json.NewEncoder(res)
		if err := encoder.Encode(task); err != nil {
			fmt.Printf("Error encoding message %v", err)
		}
	}
}

// addTask handles the api/tasks POST requests.
func addTask(res http.ResponseWriter, req *http.Request) {
	var (
		command common.CommandObj
		newTask common.Task
	)
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&command); err != nil {
		fmt.Printf("Cannot decode JSON\n%v", err)
		res.WriteHeader(http.StatusInternalServerError)
		encoder := json.NewEncoder(res)
		if err := encoder.Encode(common.ResponseMsg{Msg: "Invalid data"}); err != nil {
			fmt.Printf("Error encoding message %v", err)
		}
		return
	}
	newTask.Command = command.Command
	newTask.Tags = command.Tags
	id, err := taskDB.Add(newTask)
	if err != nil {
		fmt.Printf("Cannot add new task to DB\n%v", err)
		res.WriteHeader(http.StatusInternalServerError)
		encoder := json.NewEncoder(res)
		if err := encoder.Encode(common.ResponseMsg{Msg: "Cannot save to DB"}); err != nil {
			fmt.Printf("Error encoding message %v", err)
		}
	}
	encoder := json.NewEncoder(res)
	if err := encoder.Encode(common.ResponseMsg{Msg: id}); err != nil {
		fmt.Printf("Error encoding message %v", err)
	}
}

// claimTask handles the api/tasks/claim request.
func claimTask(res http.ResponseWriter, req *http.Request) {
	var taskToClaim common.Task
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&taskToClaim); err != nil {
		fmt.Printf("Cannot decode JSON\n%v", err)
		res.WriteHeader(http.StatusInternalServerError)
		encoder := json.NewEncoder(res)
		if err := encoder.Encode(common.ResponseMsg{Msg: "Invalid data"}); err != nil {
			fmt.Printf("Error encoding message %v", err)
		}
		return
	}
	taskToClaim.Status = "Assigned"
	_, err := taskDB.Update(taskToClaim.ID, taskToClaim)
	if err != nil {
		fmt.Printf("Cannot claim Task\n%v\n", err)
		res.WriteHeader(http.StatusInternalServerError)
		encoder := json.NewEncoder(res)
		if err = encoder.Encode(common.ResponseMsg{Msg: "Database error"}); err != nil {
			fmt.Printf("Error encoding message %v", err)
		}
		return
	}
	encoder := json.NewEncoder(res)
	if err = encoder.Encode(common.ResponseMsg{Msg: "Successfully claimed task " + taskToClaim.ID}); err != nil {
		fmt.Printf("Error encoding message %v", err)
	}
}

// updateTaskStatus handles the api/task/update/{id}/{status} requests.
func updateTaskStatus(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	task, ok := taskDB.Query(params["id"])
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		encoder := json.NewEncoder(res)
		if err := encoder.Encode(common.ResponseMsg{Msg: "Task not found"}); err != nil {
			fmt.Printf("Error encoding message %v", err)
		}
	} else {
		task.Status = params["status"]
		_, dbErr := taskDB.Update(params["id"], task)
		if dbErr != nil {
			res.WriteHeader(http.StatusInternalServerError)
			encoder := json.NewEncoder(res)
			if err := encoder.Encode(common.ResponseMsg{Msg: "Cannot save to Task Database"}); err != nil {
				fmt.Printf("Error encoding message %v", err)
			}
			return
		}
		encoder := json.NewEncoder(res)
		if err := encoder.Encode(common.ResponseMsg{Msg: "Task status updated"}); err != nil {
			fmt.Printf("Error encoding message %v", err)
			return
		}
	}
}

// uploadLog handes the api/log/{id} requests.
func uploadLog(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	path := "logs/" + params["id"] + ".log"
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Printf("Cannot create logfile %v\n%v\n", path, err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Fatalf("Error closing logfile %v\n%v\n", path, closeErr)
		}
	}()
	n, err := io.Copy(file, req.Body)
	if err != nil {
		fmt.Printf("Cannot write data to file \n%v", err)
	}
	task, _ := taskDB.Query(params["id"])
	task.Status += ", log available"
	_, err = taskDB.Update(task.ID, task)
	if err != nil {
		fmt.Printf("Error updating task status %v", err)
	}
	_, err = res.Write([]byte(fmt.Sprintf("%d bytes are recieved.\n", n)))
	if err != nil {
		fmt.Printf("Error sending response\n%v", err)
	}
}
