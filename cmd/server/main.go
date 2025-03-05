package main

import (
	// "log"
	// "net/http"

	"fmt"
	"log"
	"net/http"

	"github.com/IDK536/go_calc2/cmd/agent"
	"github.com/IDK536/go_calc2/internal/orchestrator"
	"github.com/IDK536/go_calc2/internal/parser"
	"github.com/go-chi/chi/v5"
)

func main() {
	go agent.Run()
	r := chi.NewRouter()

	r.Post("/api/v1/calculate", orchestrator.HandleCalculate)
	r.Get("/api/v1/expressions", orchestrator.HandleGetExpressions)
	r.Get("/api/v1/expression/{id}", orchestrator.HandleGetExpression)
	r.Get("/task", orchestrator.HandleGetTask)
	r.Post("/task/result", orchestrator.HandlePostTaskResult)

	log.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", r)

	fmt.Println(parser.ParseExpression("1+(1+1)"))
	// agent.Run()
	// expressions := []string{
	// 	"3 + (4 * 2) - 1",
	// 	"2+(3*4)",
	// 	"(2+3)*(4+5)",
	// 	"10-(2*3+4)",
	// }

	// for _, expr := range expressions {
	// 	fmt.Printf("Parsing expression: %s\n", expr)
	// 	values, ops, err := parser.ParseExpression(expr)
	// 	if err != nil {
	// 		fmt.Println("Error:", err)
	// 		continue
	// 	}
	// 	fmt.Println("Values:", values)
	// 	fmt.Println("Operators:", ops)
	// 	fmt.Println()
	// }
}
