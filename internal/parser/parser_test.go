package parser

import (
	"fmt"
	"testing"
)

func TestBuildAST(t *testing.T) {
	tests := []struct {
		expression string
		expected   [][]string
	}{
		{"2+3*4", [][]string{[]string{"2", "3", "4"}, []string{"+", "*"}}},
		{"(2+3)*4", [][]string{[]string{"(2+3)", "4"}, []string{"*"}}},
		{"10/(2+3)", [][]string{[]string{"10", "(2+3)"}, []string{"/"}}},
		{"(5+(6*7))-8", [][]string{[]string{"(5+(6*7))", "8"}, []string{"-"}}},
	}

	for _, test := range tests {
		v, ops, err := ParseExpression(test.expression)

		if err != nil {
			t.Fatalf("Failed to build AST for '%s': %v", test.expression, err)
		}
		if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", test.expected[0]) || fmt.Sprintf("%v", ops) != fmt.Sprintf("%v", test.expected[1]) {
			t.Errorf("For '%s', expected '%v', got '%v'", test.expression, test.expected, [][]string{v, ops})
		}

	}
}

// func astToString(nodes []interface{}) string {
// 	var result string
// 	for _, node := range nodes {
// 		switch n := node.(type) {
// 		case float64:
// 			result += fmt.Sprintf("%v ", n)
// 		case *Node:
// 			result += fmt.Sprintf("%s ", n.Operator)
// 		}
// 	}
// 	return result
// }
