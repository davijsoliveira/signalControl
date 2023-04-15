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

// executa o sistema de sinais
func (s *TrafficSignalSystem) Exec(data chan []TrafficSignal, changes chan []TrafficSignal) {
	// imprime a configuração de tempo atual dos sinais
	fmt.Println("################### TRAFFIC SIGNAL CONTROL SYSTEM ####################################")
	for i := range s.TrafficSignals {
		fmt.Println("O semáforo de ID:", s.TrafficSignals[i].Id, "tem os seguintes tempos, Verde:", s.TrafficSignals[i].TimeGreen, "Amarelo:", s.TrafficSignals[i].TimeYellow, "Vermelho:", s.TrafficSignals[i].TimeRed)
	}
	fmt.Println("######################################################################################")
	fmt.Println("")

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
		// imprime a configuração de tempo atual dos sinais
		fmt.Println("################### TRAFFIC SIGNAL CONTROL SYSTEM ####################################")
		for i := range s.TrafficSignals {
			fmt.Println("O semáforo de ID:", s.TrafficSignals[i].Id, "tem os seguintes tempos, Verde:", s.TrafficSignals[i].TimeGreen, "Amarelo:", s.TrafficSignals[i].TimeYellow, "Vermelho:", s.TrafficSignals[i].TimeRed)
		}
		fmt.Println("######################################################################################")
		fmt.Println("")
	}

}
