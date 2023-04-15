package agent

import (
	"encoding/json"
	"log"
	"net"
	"signalControl/app"
	"time"
)

type Agent struct{}

func NewAgent() *Agent {
	return &Agent{}
}

func (Agent) Exec(s chan []app.TrafficSignal, changes chan []app.TrafficSignal) {
	for {
		// Intervalo entre as coletas
		time.Sleep(10 * time.Second)

		// recebe da aplicação as informações dos sinais
		data := <-s

		// Estabelece a conexão com o controller
		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			log.Fatal(err)
		}

		defer conn.Close()

		// Envia as informações dos sinais ao servidor
		encoder := json.NewEncoder(conn)
		if err := encoder.Encode(data); err != nil {
			log.Fatal(err)
		}

		// recebe do controller as alterações a serem realizadas
		decoder := json.NewDecoder(conn)
		var ts []app.TrafficSignal
		if err := decoder.Decode(&ts); err != nil {
			log.Fatal(err)
		}

		// envia para a aplicação as alterações
		changes <- ts
	}
}
