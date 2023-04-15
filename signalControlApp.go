package main

import (
	"signalControl/agent"
	"signalControl/app"
	"signalControl/constants"
	"sync"
)

func main() {

	// variável wait group para controlar as go routines
	var wg sync.WaitGroup

	// cria os canais
	appToAgent := make(chan []app.TrafficSignal)
	agentToApp := make(chan []app.TrafficSignal)

	// instancia a aplicação e o agente
	signalControl := app.NewTrafficSignalSystem(constants.TrafficSignalNumber)
	agt := agent.NewAgent()

	// coloca em execução a aplicação e o agente
	wg.Add(2)
	go signalControl.Exec(appToAgent, agentToApp)
	go agt.Exec(appToAgent, agentToApp)
	wg.Wait()
}
