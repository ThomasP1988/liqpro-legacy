package config

import "sort"

// AuthorisedInstrumentsAndLevels kind of instruments and level we offer for websocket
var AuthorisedInstrumentsAndLevels = map[string]map[float64]bool{
	"btceur": {
		0.001: true,
		0.005: true,
		0.01:  true,
		0.05:  true,
		0.1:   true,
		0.5:   true,
		1:     true,
		5:     true,
		10:    true,
		20:    true,
		30:    true,
		50:    true,
		100:   true,
		500:   true,
	},
	"btcusd": {
		1:   true,
		5:   true,
		10:  true,
		20:  true,
		30:  true,
		50:  true,
		100: true,
		500: true,
	},
}

// PairsInfo essentials information related to trading pairs
var PairsInfo map[string]PairInfo = map[string]PairInfo{
	"btceur": {
		Base:         BTC,
		Quote:        EUR,
		BaseDecimal:  10,
		QuoteDecinal: 2,
	},
	"btcusd": {
		Base:         BTC,
		Quote:        USD,
		BaseDecimal:  10,
		QuoteDecinal: 2,
	},
}

// AuthorisedInstrumentsAndLevelsArray for speed purporse we need array version
var AuthorisedInstrumentsAndLevelsArray = map[string][]float64{}

// SetConfig  generate some config obj or array
func SetConfig() {
	for instr, _ := range AuthorisedInstrumentsAndLevels {
		AuthorisedInstrumentsAndLevelsArray[instr] = []float64{}
		for k, _ := range AuthorisedInstrumentsAndLevels[instr] {
			AuthorisedInstrumentsAndLevelsArray[instr] = append(AuthorisedInstrumentsAndLevelsArray[instr], k)
		}
		sort.Float64s(AuthorisedInstrumentsAndLevelsArray[instr])
	}

}
