package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"signalControl/constants"
	"sync"
)

// tipo semáforo
//type TrafficSignal struct {
//	Id         int
//	TimeGreen  int
//	TimeYellow int
//	TimeRed    int
//	Congestion int
//}

type TrafficSignal struct {
	Id         int `json:"id"`
	TimeGreen  int `json:"timegreen"`
	TimeYellow int `json:"timeyellow"`
	TimeRed    int `json:"timered"`
	Congestion int `json:"congestion"`
}

// tipo sistema de semáforos
//type TrafficSignalSystem struct {
//	TrafficSignals []TrafficSignal
//}

type TrafficSignalSystem struct {
	TrafficSignals []TrafficSignal `json:"trafficsignals"`
}

type TrafficSignalStore struct {
	mu             sync.Mutex
	TrafficSignals map[int]*TrafficSignalData
	ActiveRequests int64
}

type TrafficSignalData struct {
	TrafficSignal   TrafficSignal
	FlowData        []int
	AverageFlowRate int
}

var Store TrafficSignalStore

// instancia um semáforo
func NewTrafficSignal(id int) TrafficSignal {
	s := TrafficSignal{Id: id, TimeGreen: constants.DefaultGreen, TimeYellow: constants.DefaultYellow, TimeRed: constants.DefaultRed, Congestion: constants.DefaultTraffic}

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
func (s *TrafficSignalSystem) ExposeData(w http.ResponseWriter, r *http.Request) {
	// Converte os dados da struct para JSON
	jsonData, err := json.Marshal(s.TrafficSignals)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Define o cabeçalho da resposta como JSON
	w.Header().Set("Content-Type", "application/json")

	// Escreve o JSON na resposta
	w.Write(jsonData)
}

func (s *TrafficSignalSystem) UpdateSignal(w http.ResponseWriter, r *http.Request) {
	// Lê o corpo da solicitação POST
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Decodifica o JSON recebido no corpo da solicitação
	var updatedSignal TrafficSignal
	err = json.Unmarshal(body, &updatedSignal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(updatedSignal.Id)
	fmt.Println(updatedSignal.Congestion)
	fmt.Println(updatedSignal.TimeRed)
	fmt.Println(updatedSignal.TimeYellow)
	fmt.Println(updatedSignal.TimeGreen)
	// Atualiza o sinal de trânsito correspondente no slice TrafficSignals
	//for i, signal := range s.TrafficSignals {
	//	if signal.Id == updatedSignal.Id {
	//		s.TrafficSignals[i] = updatedSignal
	//		break
	//	}
	//}

	// Retorna uma resposta de sucesso
	fmt.Fprint(w, "Sinal de trânsito atualizado com sucesso!")
}

func (s *TrafficSignalSystem) HandleTrafficSignal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Method not allowed")
		return
	}

	var trafficSignal TrafficSignal
	err := json.NewDecoder(r.Body).Decode(&trafficSignal)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request payload")
		return
	}

	Store.mu.Lock()
	defer Store.mu.Unlock()

	trafficSignalData, ok := Store.TrafficSignals[trafficSignal.Id]
	if !ok {
		trafficSignalData = &TrafficSignalData{
			TrafficSignal: trafficSignal,
			FlowData:      make([]int, 0),
		}
		Store.TrafficSignals[trafficSignal.Id] = trafficSignalData
	}

	trafficSignalData.FlowData = append(trafficSignalData.FlowData, trafficSignal.Congestion)

	if len(trafficSignalData.FlowData) > 10 {
		trafficSignalData.FlowData = trafficSignalData.FlowData[len(trafficSignalData.FlowData)-10:]
	}

	total := 0
	for _, flow := range trafficSignalData.FlowData {
		total += flow
	}

	trafficSignalData.AverageFlowRate = total / len(trafficSignalData.FlowData)

	for i, signal := range s.TrafficSignals {
		if signal.Id == trafficSignal.Id {
			s.TrafficSignals[i] = trafficSignal
			break
		}
	}

	// Exibir informações da requisição POST
	fmt.Printf("Traffic Signal ID: %d, tem os seguintes tempos, Verde: %d, Amarelo: %d, Vermelho: %d\n",
		trafficSignal.Id, trafficSignal.TimeGreen, trafficSignal.TimeYellow, trafficSignal.TimeRed)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Traffic signal data stored successfully")
}
