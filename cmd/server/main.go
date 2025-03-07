package main

import (
	"net/http"

	"github.com/IDK536/go_calc2/cmd/agent"
	"github.com/IDK536/go_calc2/internal/orchestrator"
)

func main() {
	go agent.Run()

	http.HandleFunc("/api/v1/calculate", orchestrator.HandleCalculate)
	http.HandleFunc("/api/v1/expressions", orchestrator.HandleGetExpressions)
	http.HandleFunc("/api/v1/expressions/{id}", orchestrator.HandleGetExpression)
	http.HandleFunc("/task", orchestrator.HandleGetTask)
	http.HandleFunc("/task/result", orchestrator.HandlePostTaskResult)

	http.ListenAndServe(":8080", nil)
}
