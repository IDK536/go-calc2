package orchestrator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandleCalculate(t *testing.T) {
	reqBody := `{"expression": "2 + (3 * 4)"}` // Пример с скобками
	req, _ := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer([]byte(reqBody)))
	w := httptest.NewRecorder()

	HandleCalculate(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	t.Errorf(resp["id"])
	if resp["id"] == "" {
		t.Errorf("Expected non-empty ID, got empty")
	}
}

func TestHandleGetExpressions(t *testing.T) {
	reqBody := `{"expression": "2 + (3 * 4)"}` // Пример с скобками
	req, _ := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer([]byte(reqBody)))
	w := httptest.NewRecorder()

	HandleCalculate(w, req)

	req, _ = http.NewRequest("GET", "/api/v1/expressions", nil)
	w = httptest.NewRecorder()

	HandleGetExpressions(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp map[string][]*Expression
	json.Unmarshal(w.Body.Bytes(), &resp)

	if len(resp["expressions"]) != 1 {
		t.Errorf("Expected 1 expression, got %d", len(resp["expressions"]))
	}
}

func TestHandleGetExpression(t *testing.T) {
	reqBody := `{"expression": "2 + (3 * 4)"}` // Пример с скобками
	req, _ := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer([]byte(reqBody)))
	w := httptest.NewRecorder()

	HandleCalculate(w, req)

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)

	exprID := resp["id"]

	req, _ = http.NewRequest("GET", "/api/v1/expression/"+exprID, nil)
	w = httptest.NewRecorder()

	HandleGetExpression(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var exprResp map[string]*Expression
	json.Unmarshal(w.Body.Bytes(), &exprResp)

	expr := exprResp["expression"]
	if expr.ID != exprID {
		t.Errorf("Expected ID %s, got %s", exprID, expr.ID)
	}
	if expr.Status != "pending" {
		t.Errorf("Expected status pending, got %s", expr.Status)
	}
	if expr.Result != 0 {
		t.Errorf("Expected result 0, got %v", expr.Result)
	}
}

func TestHandleGetTask(t *testing.T) {
	reqBody := `{"expression": "2 + (3 * 4)"}` // Пример с скобками
	req, _ := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer([]byte(reqBody)))
	w := httptest.NewRecorder()

	HandleCalculate(w, req)

	req, _ = http.NewRequest("GET", "/task", nil)
	w = httptest.NewRecorder()

	HandleGetTask(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var taskResp map[string]Task
	json.Unmarshal(w.Body.Bytes(), &taskResp)

	task := taskResp["task"]
	if task.ID == 0 {
		t.Errorf("Expected non-zero task ID, got %d", task.ID)
	}
	if task.Operation != "*" {
		t.Errorf("Expected operation *, got %s", task.Operation)
	}
	if arg1, ok := task.Arg1.(float64); !ok || arg1 != 3 {
		t.Errorf("Expected Arg1 3, got %v", task.Arg1)
	}
	if arg2, ok := task.Arg2.(float64); !ok || arg2 != 4 {
		t.Errorf("Expected Arg2 4, got %v", task.Arg2)
	}
}

func TestHandlePostTaskResult(t *testing.T) {
	reqBody := `{"expression": "2 + (3 * 4)"}` // Пример с скобками
	req, _ := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer([]byte(reqBody)))
	w := httptest.NewRecorder()

	HandleCalculate(w, req)

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)

	exprID := resp["id"]

	req, _ = http.NewRequest("GET", "/task", nil)
	w = httptest.NewRecorder()

	HandleGetTask(w, req)

	var taskResp map[string]Task
	json.Unmarshal(w.Body.Bytes(), &taskResp)

	task := taskResp["task"]
	taskID := task.ID

	// Отправляем результат умножения 3 * 4 = 12
	resultReqBody := fmt.Sprintf(`{"id": %d, "result": 12}`, taskID)
	resultReq, _ := http.NewRequest("POST", "/task/result", bytes.NewBuffer([]byte(resultReqBody)))
	resultW := httptest.NewRecorder()

	HandlePostTaskResult(resultW, resultReq)

	if resultW.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resultW.Code)
	}

	// Получаем следующую задачу (2 + 12)
	req, _ = http.NewRequest("GET", "/task/result", nil)
	w = httptest.NewRecorder()

	HandleGetTask(w, req)

	json.Unmarshal(w.Body.Bytes(), &taskResp)
	task = taskResp["task"]
	taskID = task.ID

	// Отправляем результат сложения 2 + 12 = 14
	resultReqBody = fmt.Sprintf(`{"id": %d, "result": 14}`, taskID)
	resultReq, _ = http.NewRequest("POST", "/task", bytes.NewBuffer([]byte(resultReqBody)))
	resultW = httptest.NewRecorder()

	HandlePostTaskResult(resultW, resultReq)

	if resultW.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resultW.Code)
	}

	time.Sleep(2 * time.Second) // Ожидание завершения всех задач

	req, _ = http.NewRequest("GET", "/api/v1/expression/"+exprID, nil)
	w = httptest.NewRecorder()

	HandleGetExpression(w, req)

	var exprResp map[string]*Expression
	json.Unmarshal(w.Body.Bytes(), &exprResp)

	expr := exprResp["expression"]
	if expr.Status != "completed" {
		t.Errorf("Expected status completed, got %s", expr.Status)
	}
	if expr.Result != 14 {
		t.Errorf("Expected result 14, got %v", expr.Result)
	}
}
