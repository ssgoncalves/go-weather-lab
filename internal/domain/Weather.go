package domain

type Weather struct {
	City       string  `json:"city"`
	State      string  `json:"state"`
	Celsius    float64 `json:"temp_c"`
	Fahrenheit float64 `json:"temp_f"`
	Kelvin     float64 `json:"temp_k"`
}
