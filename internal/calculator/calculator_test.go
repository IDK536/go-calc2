package calculator

import (
	"testing"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		node     *Node
		expected float64
	}{
		{&Node{Operator: "+", Left: &Node{Value: 2}, Right: &Node{Value: 3}}, 5},
		{&Node{Operator: "*", Left: &Node{Value: 2}, Right: &Node{Value: 3}}, 6},
		{&Node{Operator: "/", Left: &Node{Value: 10}, Right: &Node{Value: 2}}, 5},
		{&Node{Operator: "+", Left: &Node{Operator: "*", Left: &Node{Value: 2}, Right: &Node{Value: 3}}, Right: &Node{Value: 4}}, 10},
		{&Node{Operator: "*", Left: &Node{Operator: "+", Left: &Node{Value: 2}, Right: &Node{Value: 3}}, Right: &Node{Value: 4}}, 20},
		{&Node{Operator: "-", Left: &Node{Operator: "/", Left: &Node{Value: 10}, Right: &Node{Value: 2}}, Right: &Node{Value: 3}}, 2},
		{&Node{Operator: "-", Left: &Node{Operator: "+", Left: &Node{Value: 5}, Right: &Node{Operator: "*", Left: &Node{Value: 6}, Right: &Node{Value: 7}}}, Right: &Node{Value: 8}}, 39},
	}

	for _, tt := range tests {
		result, err := Calculate(tt.node)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if result != tt.expected {
			t.Errorf("For %+v: expected %v, got %v", tt.node, tt.expected, result)
		}
	}
}
