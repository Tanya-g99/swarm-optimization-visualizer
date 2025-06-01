package algos

func setDefault[T any](value *T, defaultValue T) T {
	if value == nil {
		return defaultValue
	}
	return *value
}

type Response struct {
	StepPositions [][]float64 `json:"stepPositions,omitempty"`
	BestPosition  []float64   `json:"bestPosition,omitempty"`
	BestValue     float64     `json:"bestValue"`
	Iteration     int         `json:"iteration"`
}
