package agent

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"
)

func TestAgentWorker(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/internal/task" {
			task := Task{
				ID:        1,
				Arg1:      2,
				Arg2:      3,
				Operation: "*",
				Duration:  1 * time.Second,
			}
			json.NewEncoder(w).Encode(map[string]Task{"task": task})
		} else if r.Method == http.MethodPost && r.URL.Path == "/internal/task" {
			var resultData map[string]interface{}
			json.NewDecoder(r.Body).Decode(&resultData)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Result received"})
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	os.Setenv("AGENT_TASK_URL", server.URL+"/internal/task")

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		AgentWorker(wg)
	}()

	time.Sleep(5 * time.Second)

	wg.Wait()
}

func TestSimulateWork(t *testing.T) {
	task := Task{
		ID:        1,
		Arg1:      2,
		Arg2:      3,
		Operation: "*",
		Duration:  1 * time.Second,
	}

	result := simulateWork(&task)
	if result != 6 {
		t.Errorf("Expected result 6, got %v", result)
	}
}

func TestSendResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/internal/task" {
			var resultData map[string]interface{}
			json.NewDecoder(r.Body).Decode(&resultData)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Result received"})
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	os.Setenv("AGENT_TASK_URL", server.URL+"/internal/task")

	task := Task{
		ID:        1,
		Arg1:      2,
		Arg2:      3,
		Operation: "*",
		Duration:  1 * time.Second,
	}

	sendResult(task.ID, 6)
}
