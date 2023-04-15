package traffic

import (
	"math/rand"
	"signalControl/app"
	"signalControl/constants"
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
			time.Sleep(2 * time.Second)
			rand.Seed(time.Now().UnixNano())
			congestion := rand.Intn(constants.MaxTraffic)
			s.TrafficSignals[i].Congestion = congestion
		}
		time.Sleep(5 * time.Second)
	}
}
