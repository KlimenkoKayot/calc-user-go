package main

import (
	"log"

	agent "github.com/klimenkokayot/calc-net-go/internal/agent"
)

func main() {
	// Создаем нового агента
	agent, err := agent.NewAgent()
	if err != nil {
		log.Fatal(err)
	}
	// Запускаем агента
	if err := agent.Run(); err != nil {
		log.Fatal(err)
	}
}
