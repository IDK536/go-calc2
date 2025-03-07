package agent

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/IDK536/go_calc2/internal/agent"
)

func Run() {
	cp := os.Getenv("COMPUTING_POWER")
	if cp == "" {
		cp = "1"
	}

	computingPower, err := strconv.Atoi(cp)
	if err != nil {
		log.Fatalf("Invalid COMPUTING_POWER value: %v", err)
	}
	wg := &sync.WaitGroup{}
	for i := 0; i < computingPower; i++ {
		go agent.AgentWorker(wg)
	}
	wg.Wait()
}
