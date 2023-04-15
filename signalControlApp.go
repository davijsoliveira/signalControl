/*
***********************************************************************************************************************************************************
Author: Davi Oliveira
Description: This code implements a simple app for traffic signal timing control. The time of the signal traffic may change according to the traffic flow.
Date: 06/03/2023
***********************************************************************************************************************************************************
*/
package main

import (
	"signalControl/agent"
	"signalControl/app"
	"signalControl/constants"
	"signalControl/traffic"
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
	trafficFlw := traffic.NewTrafficFlow()
	agt := agent.NewAgent()

	// coloca em execução a aplicação e o agente
	wg.Add(2)
	go signalControl.Exec(appToAgent, agentToApp)
	go trafficFlw.Exec(signalControl)
	go agt.Exec(appToAgent, agentToApp)
	wg.Wait()
}
