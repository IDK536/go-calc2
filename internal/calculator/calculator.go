package calculator

import (
	"errors"
)

type Node struct {
	Value    float64
	Operator string
	Left     *Node
	Right    *Node
}

// Вычисляет значение AST
func Calculate(node *Node) (float64, error) {
	if node == nil {
		return 0, errors.New("invalid node")
	}
	if node.Operator == "" {
		return node.Value, nil
	}

	left, err := Calculate(node.Left)
	if err != nil {
		return 0, err
	}

	right, err := Calculate(node.Right)
	if err != nil {
		return 0, err
	}

	switch node.Operator {
	case "+":
		return left + right, nil
	case "-":
		return left - right, nil
	case "*":
		return left * right, nil
	case "/":
		if right == 0 {
			return 0, errors.New("division by zero")
		}
		return left / right, nil
	default:
		return 0, errors.New("unknown operator")
	}
}
