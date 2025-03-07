package agent

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type Task struct {
	ID           int
	Arg1         float64
	Arg2         float64
	Operation    string
	Duration     time.Duration
	Result       float64
	Status       string
	Dependencies []int
}

func AgentWorker(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		task := fetchTask()
		if task == nil {
			time.Sleep(1 * time.Second)
			continue
		}
		result := simulateWork(task)
		sendResult(task.ID, result)
	}
}

func fetchTask() *Task {
	resp, err := http.Get("http://localhost:8080/task")
	if err != nil {
		log.Println("Failed to fetch task:", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil
	}

	var taskResp map[string]Task
	if err := json.NewDecoder(resp.Body).Decode(&taskResp); err != nil {
		log.Println("Failed to decode task:", err)
		return nil
	}

	taskCopy := taskResp["task"]
	return &taskCopy
}

func simulateWork(task *Task) float64 {
	time.Sleep(task.Duration)
	switch task.Operation {
	case "+":
		return task.Arg1 + task.Arg2
	case "-":
		return task.Arg1 - task.Arg2
	case "*":
		return task.Arg1 * task.Arg2
	case "/":
		if task.Arg2 == 0 {
			panic("division by zero")
		}
		return task.Arg1 / task.Arg2
	default:
		panic("unknown operator")
	}
}

func sendResult(taskID int, result float64) {
	data := map[string]interface{}{
		"id":     taskID,
		"result": result,
	}

	jsonData, _ := json.Marshal(data)
	resp, err := http.Post("http://localhost:8080/task/result", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Failed to send result:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Failed to send result, status:", resp.Status)
	}
}
