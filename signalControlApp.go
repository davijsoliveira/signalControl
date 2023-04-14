package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"signalControl/constants"
)

// tipo semáforo
type TrafficSignal struct {
	Id         int
	TimeGreen  int
	TimeYellow int
	TimeRed    int
}

// tipo sistema de semáforos
type TrafficSignalSystem struct {
	TrafficSignals []TrafficSignal
}

// instancia um semáforo
func NewTrafficSignal(id int) TrafficSignal {
	s := TrafficSignal{Id: id, TimeGreen: constants.DefaultGreen, TimeYellow: constants.DefaultYellow, TimeRed: constants.DefaultRed}

	return s
}

// instancia um sistema de semáforos
func NewTrafficSignalSystem(num int) *TrafficSignalSystem {
	s := make([]TrafficSignal, num)
	system := TrafficSignalSystem{TrafficSignals: s}
	for i := 0; i < num; i++ {
		system.TrafficSignals[i] = NewTrafficSignal(i)
	}
	return &system
}

// executa o sistema de semáforos
func (s *TrafficSignalSystem) Exec() {
	for {
		//toMonitor <- s.TrafficSignals
		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		// Envia a mensagem ao servidor
		//msg := Message{Text: "Olá, servidor!"}
		encoder := json.NewEncoder(conn)
		if err := encoder.Encode(s.TrafficSignals); err != nil {
			log.Fatal(err)
		}

		// Lê a resposta do servidor
		decoder := json.NewDecoder(conn)
		var ts []TrafficSignal
		if err := decoder.Decode(&ts); err != nil {
			log.Fatal(err)
		}

		//ts := <-fromExecutor

		// itera sobre os semáforos alterados e os pertencentes ao sistema para aplicar as alterações
		for _, signalsChange := range ts {
			for j, signals := range s.TrafficSignals {
				if signalsChange.Id == signals.Id {
					s.TrafficSignals[j].TimeGreen = signalsChange.TimeGreen
					s.TrafficSignals[j].TimeYellow = signalsChange.TimeYellow
					s.TrafficSignals[j].TimeRed = signalsChange.TimeRed
				}

			}
		}
		fmt.Println("################### APP TRAFFIC SIGNAL CONTROL #######################################")
		for i := range s.TrafficSignals {
			fmt.Println("O semáforo de ID:", s.TrafficSignals[i].Id, "tem agora os seguintes tempos, Verde:", s.TrafficSignals[i].TimeGreen, "Amarelo:", s.TrafficSignals[i].TimeYellow, "Vermelho:", s.TrafficSignals[i].TimeRed)
		}
		fmt.Println("######################################################################################")
	}

}

func main() {
	trafSystem := NewTrafficSignalSystem(constants.TrafficSignalNumber)
	trafSystem.Exec()
}
