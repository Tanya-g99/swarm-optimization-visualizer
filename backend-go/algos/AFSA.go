package algos

import (
	"fmt"
	"math"
)

type AFSARequest struct {
	AlgoRequest

	Eta      *float64  `json:"eta,omitempty"`
	MaxTries *int      `json:"maxTryNum,omitempty"`
	Visual   []float64 `json:"visual,omitempty"`
	Teta     *float64  `json:"teta,omitempty"`
}

type AFSA struct {
	Algo

	Eta           float64
	MaxTries      int
	MinVisual     float64
	InitialVisual float64
	Teta          float64

	HistoryBest []float64
	Visual      float64
}

func NewAFSA(request AFSARequest) (Algorithm, error) {

	algo, err := NewAlgo(request.AlgoRequest)
	if err != nil {
		return nil, err
	}

	vis := setDefault(&request.Visual, []float64{1, 8})
	return &AFSA{
		Algo:          *algo,
		HistoryBest:   []float64{algo.GlobalBestValue},
		MaxTries:      setDefault(request.MaxTries, 5),
		Eta:           setDefault(request.Eta, 1e-4),
		Teta:          setDefault(request.Teta, 1),
		MinVisual:     vis[0],
		InitialVisual: vis[1],
	}, nil
}

func (afsa *AFSA) Run(send func(Response) error) ([]float64, float64) {
	err := send(Response{
		StepPositions: afsa.Population,
		BestPosition:  afsa.GlobalBestPosition,
		BestValue:     afsa.GlobalBestValue,
		Iteration:     0,
	})
	if err != nil {
		fmt.Println("CONN ERR", err.Error())
		return afsa.GlobalBestPosition, afsa.GlobalBestValue
	}
	stagnationCount := 0

	distanceMatrix := make([][]float64, afsa.PopulationSize)
	for i := 0; i < afsa.PopulationSize; i++ {
		distanceMatrix[i] = make([]float64, afsa.PopulationSize)
		for j := 0; j < afsa.PopulationSize; j++ {
			distanceMatrix[i][j] = math.Sqrt(afsa.calculateDistance(afsa.Population[i], afsa.Population[j]))
		}
	}

	for t := range afsa.Iterations {
		stepPositions := make([][]float64, 0)

		afsa.Visual = math.Max(afsa.MinVisual, afsa.InitialVisual*(1-float64(t)/float64(afsa.Iterations)))

		for i := range afsa.PopulationSize {
			var newPosition []float64
			neighbors := afsa.findNeighbors(i, distanceMatrix)

			if len(neighbors) == 0 {
				newPosition = afsa.randomMove(afsa.Population[i])
			} else {
				if float64(len(neighbors))/float64(afsa.PopulationSize) > afsa.Teta {
					newPosition = afsa.searchBehavior(i, neighbors)
				} else {
					c_i := afsa.meanPosition(neighbors)
					if afsa.Func(c_i) < afsa.Func(afsa.Population[i]) {
						newPosition = afsa.swarmBehavior(c_i, afsa.Population[i])
					} else {
						newPosition = afsa.searchBehavior(i, neighbors)
					}

					jStar := afsa.bestNeighbor(neighbors)
					if jStar != -1 {
						if afsa.Func(afsa.Population[jStar]) < afsa.Func(afsa.Population[i]) {
							newPosition = afsa.chaseBehavior(i, jStar)
						} else {
							newPosition = afsa.searchBehavior(i, neighbors)
						}
					}

				}
			}

			fNewPosition := afsa.Func(newPosition)
			fCurrentPosition := afsa.Func(afsa.Population[i])

			if fNewPosition < fCurrentPosition {
				afsa.Population[i] = newPosition
			}

			if fNewPosition < afsa.GlobalBestValue {
				afsa.GlobalBestValue = fNewPosition
				afsa.GlobalBestPosition = newPosition
			}

			stepPositions = append(stepPositions, afsa.Population[i])
		}

		err = send(Response{
			StepPositions: stepPositions,
			BestPosition:  afsa.GlobalBestPosition,
			BestValue:     afsa.GlobalBestValue,
			Iteration:     (t + 1),
		})
		if err != nil {
			fmt.Println("ERROR", err.Error())
			return afsa.GlobalBestPosition, afsa.GlobalBestValue
		}

		afsa.HistoryBest = append(afsa.HistoryBest, afsa.GlobalBestValue)

		if t > 0 && math.Abs(afsa.HistoryBest[len(afsa.HistoryBest)-1]-afsa.HistoryBest[len(afsa.HistoryBest)-2]) < afsa.Eta {
			stagnationCount++
		} else {
			stagnationCount = 0
		}

		if stagnationCount > afsa.MaxTries {
			j := afsa.Rng.Intn(afsa.PopulationSize)
			afsa.Population[j] = afsa.jumpBehavior(afsa.Population[j])
		}
	}
	return afsa.GlobalBestPosition, afsa.GlobalBestValue
}

func (afsa *AFSA) findNeighbors(index int, distanceMatrix [][]float64) []int {
	v := afsa.Visual * maxDifference(afsa.Bounds)
	row := distanceMatrix[index]

	neighbors := make([]int, 0, afsa.PopulationSize)
	for i := 0; i < afsa.PopulationSize; i++ {
		if i != index && row[i] <= v {
			neighbors = append(neighbors, i)
		}
	}
	return neighbors

}

func (afsa *AFSA) randomMove(currentPosition []float64) []float64 {
	newPosition := make([]float64, afsa.NumDimensions)
	v := afsa.Visual * maxDifference(afsa.Bounds)

	for j := range afsa.NumDimensions {
		boundMin, boundMax := afsa.Bounds[j][0], afsa.Bounds[j][1]

		delta := (afsa.Rng.Float64()*2 - 1) * math.Min(v, boundMax-boundMin)
		newVal := currentPosition[j] + delta
		newPosition[j] = math.Max(boundMin, math.Min(newVal, boundMax))
	}

	return newPosition
}

func (afsa *AFSA) chaseBehavior(i, jStar int) []float64 {
	r := afsa.Rng.Float64()
	newPosition := make([]float64, afsa.NumDimensions)
	for j := range afsa.NumDimensions {
		newPosition[j] = math.Max(afsa.Bounds[j][0], math.Min(afsa.Population[i][j]+r*(afsa.Population[jStar][j]-afsa.Population[i][j]), afsa.Bounds[j][1]))
	}
	return newPosition
}

func (afsa *AFSA) swarmBehavior(c_i, x_i []float64) []float64 {
	r := afsa.Rng.Float64()
	newPosition := make([]float64, afsa.NumDimensions)
	for j := range afsa.NumDimensions {
		newPosition[j] = math.Max(afsa.Bounds[j][0], math.Min(x_i[j]+r*(c_i[j]-x_i[j]), afsa.Bounds[j][1]))
	}
	return newPosition
}

func (afsa *AFSA) searchBehavior(i int, V_i []int) []float64 {
	j := V_i[afsa.Rng.Intn(len(V_i))]
	r := afsa.Rng.Float64()
	newPosition := make([]float64, afsa.NumDimensions)
	for k := range afsa.NumDimensions {
		newPosition[k] = math.Max(afsa.Bounds[k][0], math.Min(afsa.Population[i][k]+r*(afsa.Population[j][k]-afsa.Population[i][k]), afsa.Bounds[k][1]))
	}
	return newPosition
}

func (afsa *AFSA) jumpBehavior(x_i []float64) []float64 {
	p := afsa.Rng.Float64()
	newPosition := make([]float64, afsa.NumDimensions)

	for j := range afsa.NumDimensions {
		boundMin, boundMax := afsa.Bounds[j][0], afsa.Bounds[j][1]
		rangeJ := boundMax - boundMin

		direction := (afsa.Rng.Float64()*2 - 1) * p
		delta := direction * rangeJ

		newVal := x_i[j] + delta
		newPosition[j] = math.Max(boundMin, math.Min(newVal, boundMax))
	}

	return newPosition
}

func (afsa *AFSA) bestNeighbor(V_i []int) int {
	if len(V_i) == 0 {
		return -1
	}
	bestIndex := V_i[0]
	bestValue := afsa.Func(afsa.Population[bestIndex])
	for _, index := range V_i[1:] {
		if afsa.Func(afsa.Population[index]) < bestValue {
			bestValue = afsa.Func(afsa.Population[index])
			bestIndex = index
		}
	}
	return bestIndex
}

func (afsa *AFSA) calculateDistance(a, b []float64) float64 {
	sum := 0.0
	for i := range a {
		sum += (a[i] - b[i]) * (a[i] - b[i])
	}
	return sum
}

func (afsa *AFSA) meanPosition(neighbors []int) []float64 {
	meanPos := make([]float64, afsa.NumDimensions)
	for _, neighbor := range neighbors {
		for j := range meanPos {
			meanPos[j] += afsa.Population[neighbor][j]
		}
	}
	for j := range meanPos {
		meanPos[j] /= float64(len(neighbors))
	}
	return meanPos
}

func maxDifference(bounds [][]float64) float64 {
	maxDiff := 0.0
	for _, bound := range bounds {
		if diff := bound[1] - bound[0]; diff > maxDiff {
			maxDiff = diff
		}
	}
	return maxDiff
}
