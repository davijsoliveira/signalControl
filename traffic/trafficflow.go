package traffic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"signalControl/app"
	"time"
)

// tipo fluxo de trânsito representando o ambiente
type TrafficFlow struct {
	TrafficPerSemaphore []int
}

// instancia um novo fluxo de trânsito
func NewTrafficFlow() *TrafficFlow {
	return &TrafficFlow{}
}

// executa o fluxo de trânsito, gerando congestionamentos aleatórios
func (t *TrafficFlow) Exec(s *app.TrafficSignalSystem) {
	for {
		// gera um número aletório de congestionamento para cada sinal
		for i := range s.TrafficSignals {
			averageFlowRate, err := getAverageFlowRate(s.TrafficSignals[i].Id)
			if err != nil {
				log.Printf("Failed to get average flow rate for traffic signal %d: %v\n", s.TrafficSignals[i].Id, err)
			}
			//else {
			//	log.Printf("Average flow rate for Traffic Signal %d: %d\n", s.TrafficSignals[i].Id, averageFlowRate)
			//}
			s.TrafficSignals[i].Congestion = averageFlowRate
		}
		time.Sleep(5 * time.Second)
	}
}

func getAverageFlowRate(signalID int) (int, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8082/traffic/info?id=%d", signalID))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var response struct {
		AverageFlowRate int `json:"averageFlowRate"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}

	return response.AverageFlowRate, nil
}
