package app

import (
	"fmt"
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
func (s *TrafficSignalSystem) Exec(data chan []TrafficSignal, changes chan []TrafficSignal) {
	for {
		data <- s.TrafficSignals
		ts := <-changes

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
		fmt.Println("################### TRAFFIC SIGNAL CONTROL SYSTEM ####################################")
		for i := range s.TrafficSignals {
			fmt.Println("O semáforo de ID:", s.TrafficSignals[i].Id, "tem os seguintes tempos, Verde:", s.TrafficSignals[i].TimeGreen, "Amarelo:", s.TrafficSignals[i].TimeYellow, "Vermelho:", s.TrafficSignals[i].TimeRed)
		}
		fmt.Println("######################################################################################")
		fmt.Println("")
	}

}

/*func (s *TrafficSignalSystem) Exec() {
	for {
		//toMonitor <- s.TrafficSignals
		time.Sleep(10 * time.Second)
		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		// Envia a mensagem ao servidor
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
		fmt.Println("################### TRAFFIC SIGNAL CONTROL SYSTEM ####################################")
		for i := range s.TrafficSignals {
			fmt.Println("O semáforo de ID:", s.TrafficSignals[i].Id, "tem os seguintes tempos, Verde:", s.TrafficSignals[i].TimeGreen, "Amarelo:", s.TrafficSignals[i].TimeYellow, "Vermelho:", s.TrafficSignals[i].TimeRed)
		}
		fmt.Println("######################################################################################")
	}

}*/
