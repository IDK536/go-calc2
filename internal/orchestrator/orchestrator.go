package orchestrator

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/IDK536/go_calc2/internal/parser"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Node struct {
	Valuse    []string
	Operators []string
}

type Task struct {
	ID           int
	Arg1         interface{} // Может быть int (идентификатор задачи) или float64 (число)
	Arg2         interface{} // Может быть int (идентификатор задачи) или float64 (число)
	Operation    string
	Duration     time.Duration
	Result       float64
	Status       string
	Dependencies []int
	Index        int
}

type Expression struct {
	ID      string
	Status  string
	Result  float64
	TaskIDs []int
	DepMap  map[int][]int
}

var (
	expressions = make(map[string]*Expression)
	tasks       = make(map[int]Task)
	mu          sync.RWMutex
	taskID      int
	thisnode    Node
	IsResult    = false
	glavIndex   []int
	glavResult  []string
	// done        = make(chan struct{})
	inx = 0
)

func getOperationDuration(op string) time.Duration {
	switch op {
	case "+", "-":
		return time.Duration(getEnvInt("TIME_ADDITION_MS", 3000)) * time.Millisecond
	case "*":
		return time.Duration(getEnvInt("TIME_MULTIPLICATIONS_MS", 4000)) * time.Millisecond
	case "/":
		return time.Duration(getEnvInt("TIME_DIVISIONS_MS", 5000)) * time.Millisecond
	default:
		return 0
	}
}

func getEnvInt(key string, defaultValue int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	v, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return v
}
func HandleCalculate(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received expression: %s", r.Body)
	// fmt.Println("go run")
	thisnode = Node{}
	IsResult = false
	glavIndex = []int{}
	glavResult = []string{}
	var req struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusUnprocessableEntity)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	// fmt.Println("go run1")
	id := uuid.New().String()
	expr := &Expression{
		ID:      id,
		Status:  "pending",
		TaskIDs: []int{},
		DepMap:  make(map[int][]int),
	}
	expressions[id] = expr
	go func() {
		values, ops, err := parser.ParseExpression(req.Expression)

		if err != nil {
			log.Printf("Error building AST: %v", err)
			http.Error(w, "Invalid expression", http.StatusUnprocessableEntity)
			return
		}
		// fmt.Println("go run3", values, ops)
		// var fulltasks []Task
		fulldepmap := make(map[int][]int)
		for i, n := range values {
			// fmt.Println("go run4")
			if n[0] == '(' {
				for len(thisnode.Valuse) != 1 {
					values1, ops1, err := parser.ParseExpression(n[1 : len(n)-1])
					if err != nil {
						log.Printf("Error building AST: %v", err)
						http.Error(w, "Invalid expression", http.StatusUnprocessableEntity)
						return
					}
					thisnode = Node{
						Valuse:    values1,
						Operators: ops1,
					}

					tasksList1, depMap1, err := createTasksFromAST(fulldepmap)
					fulldepmap = depMap1
					// tasks = append(tasks, tasksList1...)
					for _, n := range tasksList1 {
						tasks[len(tasks)] = n
					}
					// for _, t := range fulltasks {
					// 	tasks[len(tasks)] = t
					// }
					for !IsResult {
						time.Sleep(time.Second)
					}
					// for _, n := range tasks {
					// 	n.Status = "completed"
					// }
					IsResult = false
					// fmt.Println("gosssssss run5")
					inx = 0
				}
				glavIndex = append(glavIndex, i)
				// fmt.Println("fdhbrthrtherhhrregerg", thisnode.Valuse, tasks, thisnode)
				glavResult = append(glavResult, thisnode.Valuse[0])

			}
		}
		// fmt.Println("ababbababab", glavIndex)
		// fmt.Println("rtrtrtrtrtr", glavIndex)

		fmt.Println("grgerg", tasks)
		fmt.Println(thisnode)
		thisnode = Node{}
		proshI := 0
		if len(glavIndex) != 0 {
			for i, n := range glavIndex {
				// fmt.Println(values, glavResult, n)
				don := append(values[proshI:n], glavResult[i])
				thisnode.Valuse = append(thisnode.Valuse, don...)
				proshI = n + 1
			}
		} else {
			thisnode.Valuse = values
		}
		// fmt.Println("dghrdhrdhhrherreherherherh", thisnode)
		thisnode = Node{
			Valuse:    thisnode.Valuse,
			Operators: ops,
		}
		inx = 0
		fmt.Println("aaarseniy", thisnode)
		fmt.Println("gergergergergergergergejnswtertgewrger", thisnode)
		for len(thisnode.Valuse) != 1 {
			tasksList, depMap, err := createTasksFromAST(fulldepmap)
			inx = 0
			if err != nil {
				log.Printf("Error building AST: %v", err)
				http.Error(w, "Invalid expression", http.StatusUnprocessableEntity)
				return
			}
			fulldepmap = depMap
			for _, n := range tasksList {
				tasks[len(tasks)] = n
			}
			for !IsResult {
				time.Sleep(time.Second)
			}
			// for _, n := range tasks {
			// 	n.Status = "completed"
			// }
			IsResult = false
			fmt.Println("arseniy", tasks)
			// fmt.Println("gosssssss run5")
		}
		// fmt.Println("res", thisnode.Valuse)
		res, err := strconv.ParseFloat(thisnode.Valuse[0], 64)
		if err != nil {
			log.Printf("Error parsing result: %v", err)
		}
		expr.Result = res
		if err != nil {
			http.Error(w, "Invalid expression", http.StatusUnprocessableEntity)
			return
		}
		taskID := 0
		for range tasks {
			taskID++
			// tasks[taskID] = task
			expr.TaskIDs = append(expr.TaskIDs, taskID)
		}
		fmt.Println("huihuihuihui", tasks, expr.TaskIDs)
		expr.DepMap = fulldepmap
		expr.Status = "completed"
		expressions[id] = expr
		fmt.Println("huihuihuihui1", expr)
	}()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func createTasksFromAST(depmap1 map[int][]int) ([]Task, map[int][]int, error) {
	var tasks []Task
	values := thisnode.Valuse
	ops := thisnode.Operators
	fmt.Println("ars", values, ops)
	depMap := make(map[int][]int)
	dm := make(map[int][]int)
	inx := -1
	// fmt.Println("11111111go run 5", values, ops)
	for i, op := range ops {
		// fmt.Println("go run 8")
		// fmt.Println("go run 9")
		inx++
		if i != 0 && i != len(ops)-1 && len(ops) != 1 {
			if op == "+" || op == "-" {
				if ops[i-1] != "*" && ops[i+1] != "/" && ops[i+1] != "*" && ops[i-1] != "/" {
					fmt.Println("go run 10", values, ops, op, values[i])
					value1, err := strconv.ParseFloat(values[i], 64)
					value2, err := strconv.ParseFloat(values[i+1], 64)
					if err != nil {
						return nil, nil, err
					}
					task := Task{
						Arg1:         value1,
						Arg2:         value2,
						Operation:    op,
						Duration:     getOperationDuration(op),
						Status:       "pending",
						Dependencies: []int{},
						ID:           taskID,
						Index:        i,
					}
					taskID++
					dm[len(depMap)] = append(dm[len(depMap)], taskID)
					tasks = append(tasks, task)
					i++
				}
			} else if op == "*" || op == "/" {
				if ops[i-1] != "/" && ops[i-1] != "*" {
					value1, err := strconv.ParseFloat(values[i], 64)
					value2, err := strconv.ParseFloat(values[i+1], 64)
					if err != nil {
						return nil, nil, err
					}
					task := Task{
						Arg1:         value1,
						Arg2:         value2,
						Operation:    op,
						Duration:     getOperationDuration(op),
						Status:       "pending",
						Dependencies: []int{},
						ID:           taskID,
						Index:        i,
					}
					taskID++
					dm[len(depMap)] = append(dm[len(depMap)], taskID)
					tasks = append(tasks, task)
					i++
				}
			}
		} else if len(ops) == 1 {
			// fmt.Println("go run 7")
			value1, err := strconv.ParseFloat(values[i], 64)
			// fmt.Println("go run 11", err != nil)
			if err != nil {
				return nil, nil, err
			}
			value2, err := strconv.ParseFloat(values[i+1], 64)
			// fmt.Println("go run 10", err != nil, values[i+1])
			if err != nil {
				return nil, nil, err
			}
			task := Task{
				Arg1:         value1,
				Arg2:         value2,
				Operation:    op,
				Duration:     getOperationDuration(op),
				Status:       "pending",
				Dependencies: []int{},
				ID:           taskID,
				Index:        i,
			}
			taskID++
			dm[len(depMap)] = append(dm[len(depMap)], taskID)
			tasks = append(tasks, task)
			i++
			// fmt.Println("go run 10")
		} else if i == 0 {
			if op == "+" || op == "-" {
				if ops[i+1] != "/" && ops[i+1] != "*" {
					fmt.Println("eshkere", values[i])
					value1, err := strconv.ParseFloat(values[i], 64)
					value2, err := strconv.ParseFloat(values[i+1], 64)
					if err != nil {
						return nil, nil, err
					}
					task := Task{
						Arg1:         value1,
						Arg2:         value2,
						Operation:    op,
						Duration:     getOperationDuration(op),
						Status:       "pending",
						Dependencies: []int{},
						ID:           taskID,
						Index:        i,
					}
					taskID++
					dm[len(depMap)] = append(dm[len(depMap)], taskID)
					tasks = append(tasks, task)
					i++
				}
			} else if op == "*" || op == "/" {
				value1, err := strconv.ParseFloat(values[i], 64)
				value2, err := strconv.ParseFloat(values[i+1], 64)
				if err != nil {
					return nil, nil, err
				}
				task := Task{
					Arg1:         value1,
					Arg2:         value2,
					Operation:    op,
					Duration:     getOperationDuration(op),
					Status:       "pending",
					Dependencies: []int{},
					ID:           taskID,
					Index:        i,
				}
				taskID++
				dm[len(depMap)] = append(dm[len(depMap)], taskID)
				tasks = append(tasks, task)
				i++
			}
		} else if i == len(ops) {
			if op == "+" || op == "-" {
				if ops[i-1] != "*" && ops[i-1] != "/" {
					fmt.Println("eshkere1", values)
					value1, err := strconv.ParseFloat(values[i], 64)
					value2, err := strconv.ParseFloat(values[i+1], 64)
					if err != nil {
						return nil, nil, err
					}
					task := Task{
						Arg1:         value1,
						Arg2:         value2,
						Operation:    op,
						Duration:     getOperationDuration(op),
						Status:       "pending",
						Dependencies: []int{},
						ID:           taskID,
						Index:        i,
					}
					taskID++
					dm[len(depMap)] = append(dm[len(depMap)], taskID)
					tasks = append(tasks, task)
					i++
				}
			} else if op == "*" || op == "/" {
				if ops[i-1] != "/" && ops[i-1] != "*" {
					value1, err := strconv.ParseFloat(values[i], 64)
					value2, err := strconv.ParseFloat(values[i+1], 64)
					if err != nil {
						return nil, nil, err
					}
					task := Task{
						Arg1:         value1,
						Arg2:         value2,
						Operation:    op,
						Duration:     getOperationDuration(op),
						Status:       "pending",
						Dependencies: []int{},
						ID:           taskID,
						Index:        i,
					}
					taskID++
					dm[len(depMap)] = append(dm[len(depMap)], taskID)
					tasks = append(tasks, task)
					i++
				}
			}
		}
	}
	fmt.Println("go run 6ftghtfghfghbfgbijkfijoijkfgbkoldfgvkiodfgviklofbgkoifbgklgbfklgblkgbfklgbfklgklfbklgbklgbh", tasks)
	// var createTask func(node *parser.Node) (interface{}, error)
	// createTask = func(node *parser.Node) (interface{}, error) {
	// 	if node == nil {
	// 		return 0, errors.New("invalid node")
	// 	}
	// 	if node.Operator == "" {
	// 		return node.Value, nil
	// 	}

	// 	left, err := createTask(node.Left)
	// 	if err != nil {
	// 		return 0, err
	// 	}

	// 	right, err := createTask(node.Right)
	// 	if err != nil {
	// 		return 0, err
	// 	}

	// 	task := Task{
	// 		Arg1:         left,
	// 		Arg2:         right,
	// 		Operation:    node.Operator,
	// 		Duration:     getOperationDuration(node.Operator),
	// 		Status:       "pending",
	// 		Dependencies: []int{},
	// 	}
	// 	taskID++
	// 	tasks = append(tasks, task)

	// 	if leftID, ok := left.(int); ok {
	// 		depMap[leftID] = append(depMap[leftID], taskID)
	// 		tasks[leftID-1].Dependencies = append(tasks[leftID-1].Dependencies, taskID)
	// 	}
	// 	if rightID, ok := right.(int); ok {
	// 		depMap[rightID] = append(depMap[rightID], taskID)
	// 		tasks[rightID-1].Dependencies = append(tasks[rightID-1].Dependencies, taskID)
	// 	}
	// 	return taskID, nil
	// }

	// // Создаем последнюю задачу (конечный результат)
	// rootID, err := createTask(node)
	// if err != nil {
	// 	return nil, nil, err
	// }

	// // Устанавливаем Arg1 последней задачи на корневую задачу
	// if rootID != 0 {
	// 	tasks[len(tasks)-1].Arg1 = rootID.(int)
	// }

	depMap[len(dm)] = append(depMap[len(dm)], dm[len(dm)-1]...)
	for _, t := range tasks {
		for _, values := range depMap {
			for _, value := range values {
				t.Dependencies = append(t.Dependencies, value)
			}
		}
	}
	// if len(tasks1) != 0 {
	// 	tasks = append(tasks1, tasks...)
	// }
	// fmt.Println("tasks", tasks)
	return tasks, depMap, nil
}

func HandleGetExpressions(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("go run 10rtghrtijrtdghiugrdtiriutghduihgrtduihtgruhirgthuirgtuihrgtuhirgtuihtrghiugrthuiorgtuihtrguhirgthui", expressions)

	mu.RLock()
	defer mu.RUnlock()
	resp := map[string][]*Expression{"expressions": {}}
	for _, expr := range expressions {
		resp["expressions"] = append(resp["expressions"], expr)
	}
	json.NewEncoder(w).Encode(resp)
}

func HandleGetExpression(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	id := chi.URLParam(r, "id")
	expr, exists := expressions[id]
	fmt.Println("expr", expressions)
	if !exists {
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]*Expression{"expression": expr})
}

func HandleGetTask(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	// fmt.Println("go run 11", tasks)
	for _, task := range tasks {
		if task.Status == "pending" {
			// fmt.Println("go run 12", task)
			taskCopy := task
			taskCopy.Status = "in_progress"
			json.NewEncoder(w).Encode(map[string]Task{"task": taskCopy})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "No tasks available"})
}

func HandlePostTaskResult(w http.ResponseWriter, r *http.Request) {
	var resultData map[string]interface{}
	// fmt.Println("gooo run 2")
	if err := json.NewDecoder(r.Body).Decode(&resultData); err != nil {
		http.Error(w, "Invalid data", http.StatusUnprocessableEntity)
		return
	}
	// fmt.Println("gooo run 1")

	id, ok := resultData["id"].(float64)
	// fmt.Println("gooo run 3hrehherheregrhergegrgergerrgegregreregrg", id)
	if !ok {
		http.Error(w, "Invalid task ID", http.StatusUnprocessableEntity)
		return
	}

	taskID := int(id)
	task, exists := tasks[taskID]
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	result, ok := resultData["result"].(float64)
	if !ok {
		http.Error(w, "Invalid result", http.StatusUnprocessableEntity)
		return
	}
	// fmt.Println("eshkereeeee", thisnode.Valuse, fmt.Sprintf("%f", result), task.Index)
	var newSlice []string
	// newSlice = append(thisnode.Valuse[:task.Index], fmt.Sprintf("%f", result))
	fmt.Println("gpggoggoggooggogoo", newSlice, thisnode, task.Index, taskID)
	// if task.Index+2 < len(thisnode.Valuse) {
	// 	newSlice = append(newSlice, thisnode.Valuse[task.Index+2:]...)
	// 	thisnode = Node{
	// 		Valuse:    newSlice,
	// 		Operators: append(thisnode.Operators[:task.Index], thisnode.Operators[task.Index+1:]...),
	// 	}
	// } else {
	newSlice = append(thisnode.Valuse[:task.Index], fmt.Sprintf("%f", result))
	thisnode = Node{
		Valuse:    append(thisnode.Valuse[:task.Index+1], thisnode.Valuse[task.Index+2:]...),
		Operators: append(thisnode.Operators[:task.Index], thisnode.Operators[task.Index+1:]...),
	}
	fmt.Println("gpggoggoggooggogoo", newSlice, thisnode, task.Index, taskID)
	// }

	fmt.Println("feefwfeef", thisnode)
	// fmt.Println("ghghgghg", newSlice, thisnode)
	task.Status = "completed"
	task.Result = result
	// fmt.Println("go run 11rtghrghtrtgrgttgrtgrrgtgrtrerew", taskID)
	tasks[taskID] = task
	// fmt.Println("go run 11rtghrghtrtgrgttgrtgrrgtgrtreregerrgtergew", tasks)
	// checkTaskCompletion(taskID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Result received"})
	// go func() {
	// 	// Отправка сигнала о завершении
	// 	done <- struct{}{}
	// }()
	// a := false
	// fmt.Println("tafgdhhfgtrdhrthtfghtrthrhrtdthrfhthtsks", tasks)
	// for _, n := range tasks {
	// 	if n.Status == "completed" {
	// 		a = true
	// 	} else {
	// 		a = false
	// 		break
	// 	}
	// }
	// fmt.Println("tafgdhhfgtrdhrthtfghtrthrhrtdthrfhthtsks", a)
	// if a == true {
	IsResult = true
	// }
	// for i, n := range tasks {
	// 	if i >= taskID {
	// 		n.Index -= 1
	// 	}

	// }
	inx += 1
}

func checkTaskCompletion(taskID int) {
	mu.Lock()
	defer mu.Unlock()

	exprID := ""
	for exprID, expr := range expressions {
		for _, tid := range expr.TaskIDs {
			if tid == taskID {
				break
			}
		}
		if exprID != "" {
			break
		}
	}

	if exprID == "" {
		return
	}

	expr := expressions[exprID]
	allCompleted := true
	for _, tid := range expr.TaskIDs {
		if tasks[tid].Status != "completed" {
			allCompleted = false
			break
		}
	}

	if allCompleted {
		finalResult := calculateFinalResult(expr.TaskIDs, expr.DepMap)
		expr.Status = "completed"
		expr.Result = finalResult
		expressions[exprID] = expr
	}
}

func calculateFinalResult(taskIDs []int, depMap map[int][]int) float64 {
	results := make(map[int]float64)
	for _, tid := range taskIDs {
		task := tasks[tid]
		if task.Operation == "" {
			continue
		}
		// Получаем результаты зависимых задач
		var arg1, arg2 float64
		switch v := task.Arg1.(type) {
		case int:
			arg1 = results[v]
		case float64:
			arg1 = v
		default:
			panic("unknown type for Arg1")
		}
		switch v := task.Arg2.(type) {
		case int:
			arg2 = results[v]
		case float64:
			arg2 = v
		default:
			panic("unknown type for Arg2")
		}
		results[tid] = performOperation(arg1, arg2, task.Operation)
	}

	// Финальный результат
	return results[taskIDs[len(taskIDs)-1]]
}

func topologicalSort(taskIDs []int, depMap map[int][]int) []int {
	inDegree := make(map[int]int)
	for _, deps := range depMap {
		for _, dep := range deps {
			inDegree[dep]++
		}
	}

	queue := []int{}
	for _, tid := range taskIDs {
		if inDegree[tid] == 0 {
			queue = append(queue, tid)
		}
	}

	sorted := []int{}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		sorted = append(sorted, current)

		for _, tid := range taskIDs {
			if contains(depMap[current], tid) {
				inDegree[tid]--
				if inDegree[tid] == 0 {
					queue = append(queue, tid)
				}
			}
		}
	}

	return sorted
}

func contains(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func performOperation(arg1, arg2 float64, op string) float64 {
	switch op {
	case "+":
		return arg1 + arg2
	case "-":
		return arg1 - arg2
	case "*":
		return arg1 * arg2
	case "/":
		if arg2 == 0 {
			panic("division by zero")
		}
		return arg1 / arg2
	default:
		panic("unknown operator")
	}
}
