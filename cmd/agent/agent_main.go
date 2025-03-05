package agent

import (
	"log"
	"os"
	"strconv"

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

	agent.AgentWorker(computingPower)

	// Чтобы агент не завершался
	select {}
}
