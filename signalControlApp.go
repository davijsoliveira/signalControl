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

	// instancia a aplicação e o agente
	signalControl := app.NewTrafficSignalSystem(constants.TrafficSignalNumber)
	trafficFlw := traffic.NewTrafficFlow()

	// coloca em execução a aplicação
	wg.Add(2)
	go trafficFlw.Exec(signalControl)

	// Define a rota e o handler para expor as informações dos sinais de trânsito
	http.HandleFunc("/traffic-signals-current", signalControl.ExposeData)

	// Define a rota e o handler para atualizar um sinal de trânsito
	http.HandleFunc("/traffic-signals-update", signalControl.HandleTrafficSignal)

	// Inicia o servidor na porta 8081
	fmt.Println("Configurator Microservice Started...")
	http.ListenAndServe(":8081", nil)

	wg.Wait()
}
