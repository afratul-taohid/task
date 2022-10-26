package main

// #include <stdlib.h>
import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

type WorkspaceResponse struct {
	Status bool
	Data   []Workspace
}

type TaskResponse struct {
	Status bool
	Data   []Task
}

type TaskSaveResponse struct {
	Status bool
	Data   Task
}

type Workspace struct {
	MemberCount int8
	Team        string
	Message     string
}

type Task struct {
	Id        string
	Team      string
	Task      string
	Time      string
	StartDate string
	EndDate   string
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_, _ = fmt.Fprint(w, "Welcome!\n")
}

func GetWorkspaces(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var swdWorkspace = Workspace{
		MemberCount: 22,
		Team:        "SWD Team",
		Message:     "2 Project",
	}

	var supportWorkspace = Workspace{
		MemberCount: 10,
		Team:        "Support Team",
		Message:     "5 Support",
	}

	var resp = WorkspaceResponse{
		Status: true,
	}
	resp.Data = append(resp.Data, swdWorkspace)
	resp.Data = append(resp.Data, supportWorkspace)

	//w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.MarshalIndent(resp, "", " ")
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	_, _ = w.Write(jsonResp)
	return
}

var taskList []Task

func SaveTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	body, err := ioutil.ReadAll(r.Body)
	var task = Task{
		Id: *getUUID(),
	}
	_ = json.Unmarshal(body, &task)

	var resp = TaskSaveResponse{
		Status: true,
		Data:   task,
	}

	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.MarshalIndent(resp, "", " ")
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	_, _ = w.Write(jsonResp)
	return
}

func GetTasks(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var task1 = Task{
		Id:        *getUUID(),
		Team:      "System Support",
		Task:      "eKYC Support on NRB Bank",
		Time:      "19:24:05",
		StartDate: "12-09-2022",
		EndDate:   "12-09-2022",
	}
	var task2 = Task{
		Id:        *getUUID(),
		Team:      "eKYC Development",
		Task:      "eKYC Core Design Ready",
		Time:      "13:38:05",
		StartDate: "12-09-2022",
		EndDate:   "12-09-2022",
	}
	var task3 = Task{
		Id:        *getUUID(),
		Team:      "HSBC eKYC Deployment",
		Task:      "Ready to DevOps Team",
		Time:      "10:35:05",
		StartDate: "12-09-2022",
		EndDate:   "12-09-2022",
	}

	var resp = TaskResponse{
		Status: true,
	}
	taskList = append(taskList, task1, task2, task3)
	resp.Data = taskList

	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.MarshalIndent(resp, "", " ")
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	_, _ = w.Write(jsonResp)
	return
}

func getUUID() (response *string) {
	bytes, err := exec.Command("uuidgen").Output()
	if err == nil {
		uuid := strings.Trim(string(bytes), "\n")
		return &uuid
	}
	return nil
}

// export CreateServer
func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/api/workspace", GetWorkspaces)
	router.GET("/api/task", GetTasks)
	router.POST("/api/task", SaveTask)

	log.Fatal(http.ListenAndServe(":8080", router))
}
