/*
***********************************************************************************************************************************************************
Author: Davi Oliveira
Description: This code implements a simple app for traffic signal timing control. The time of the signal traffic may change according to the traffic flow.
Date: 06/03/2023
***********************************************************************************************************************************************************
*/
package main

import (
	"fmt"
	"net/http"
	"signalControl/agent"
	"signalControl/app"
	"signalControl/constants"
	"signalControl/traffic"
	"sync"
)

func main() {
	app.Store = app.TrafficSignalStore{
		TrafficSignals: make(map[int]*app.TrafficSignalData),
	}

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
	wg.Add(4)
	go signalControl.Exec(appToAgent, agentToApp)
	go trafficFlw.Exec(signalControl)
	go agt.Exec(appToAgent, agentToApp)

	// Define a rota e o handler para expor os dados
	http.HandleFunc("/traffic-signals-current", signalControl.ExposeData)

	// Define a rota e o handler para atualizar um sinal de trânsito
	http.HandleFunc("/traffic-signals-update", signalControl.HandleTrafficSignal)

	// Inicia o servidor na porta 8081
	fmt.Println("Servidor iniciado na porta 8081")
	http.ListenAndServe(":8081", nil)

	wg.Wait()
}
