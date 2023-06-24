package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"signalControl/constants"
	"sync"
)

type TrafficSignal struct {
	Id         int `json:"id"`
	TimeGreen  int `json:"timegreen"`
	TimeYellow int `json:"timeyellow"`
	TimeRed    int `json:"timered"`
	Congestion int `json:"congestion"`
}

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
