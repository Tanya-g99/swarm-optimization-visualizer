package algos

import (
	"math"
)

type SFLARequest struct {
	AlgoRequest

	SubpopulationsCount *int `json:"subpopulationsCount,omitempty"`
	IMax                *int `json:"iMax,omitempty"`
}

type SFLA struct {
	Algo

	SubpopulationsCount int
	IMax                int
}

func NewSFLA(request SFLARequest) (Algorithm, error) {

	algo, err := NewAlgo(request.AlgoRequest)
	if err != nil {
		return nil, err
	}

	return &SFLA{
		Algo:                *algo,
		SubpopulationsCount: setDefault(request.SubpopulationsCount, 1),
		IMax:                setDefault(request.IMax, 10),
	}, nil
}

func (sfla *SFLA) Run(send func(Response) error) ([]float64, float64) {
	err := send(Response{
		StepPositions: sfla.Population,
		BestPosition:  sfla.GlobalBestPosition,
		BestValue:     sfla.GlobalBestValue,
		Iteration:     0,
	})
	if err != nil {
		return sfla.GlobalBestPosition, sfla.GlobalBestValue
	}

	for t := range sfla.Iterations {
		for i := range sfla.SubpopulationsCount {
			sfla.localSearch(i)
			sfla.updateBest()
		}

		sfla.shufflePopulation()

		err = send(Response{
			StepPositions: sfla.Population,
			BestPosition:  sfla.GlobalBestPosition,
			BestValue:     sfla.GlobalBestValue,
			Iteration:     (t + 1),
		})
		if err != nil {
			return sfla.GlobalBestPosition, sfla.GlobalBestValue
		}
	}

	return sfla.GlobalBestPosition, sfla.GlobalBestValue
}

func (sfla *SFLA) localSearch(i int) {
	for range sfla.IMax {
		subpopSize := sfla.PopulationSize / sfla.SubpopulationsCount
		subpopStart := i * subpopSize
		subpopEnd := (i + 1) * subpopSize

		bestInSubpopIndex := -1
		worstInSubpopIndex := -1
		bestInSubpopValue := math.Inf(1)
		worstInSubpopValue := -math.Inf(1)

		for j := subpopStart; j < subpopEnd; j++ {
			value := sfla.Func(sfla.Population[j])
			if value < bestInSubpopValue {
				bestInSubpopValue = value
				bestInSubpopIndex = j
			}
			if value > worstInSubpopValue {
				worstInSubpopValue = value
				worstInSubpopIndex = j
			}
		}

		r := sfla.Rng.Float64()
		for d := range sfla.NumDimensions {
			sfla.Population[worstInSubpopIndex][d] += r * (sfla.Population[bestInSubpopIndex][d] - sfla.Population[worstInSubpopIndex][d])
			sfla.Population[worstInSubpopIndex][d] = math.Max(sfla.Bounds[d][0], math.Min(sfla.Population[worstInSubpopIndex][d], sfla.Bounds[d][1]))
		}

		if sfla.Func(sfla.Population[worstInSubpopIndex]) >= worstInSubpopValue {
			r = sfla.Rng.Float64()
			for d := range sfla.NumDimensions {
				sfla.Population[worstInSubpopIndex][d] += r * (sfla.GlobalBestPosition[d] - sfla.Population[worstInSubpopIndex][d])
				sfla.Population[worstInSubpopIndex][d] = math.Max(sfla.Bounds[d][0], math.Min(sfla.Population[worstInSubpopIndex][d], sfla.Bounds[d][1]))
			}

			if sfla.Func(sfla.Population[worstInSubpopIndex]) >= worstInSubpopValue {
				for d := range sfla.NumDimensions {
					sfla.Population[worstInSubpopIndex][d] = sfla.Rng.Float64()*(sfla.Bounds[d][1]-sfla.Bounds[d][0]) + sfla.Bounds[d][0]
				}
			}
		}

	}
}

func (sfla *SFLA) updateBest() {
	for _, frog := range sfla.Population {
		value := sfla.Func(frog)
		if value < sfla.GlobalBestValue {
			sfla.GlobalBestValue = value
			sfla.GlobalBestPosition = frog
		}
	}
}

func (sfla *SFLA) shufflePopulation() {
	sfla.Rng.Shuffle(sfla.PopulationSize, func(i, j int) {
		sfla.Population[i], sfla.Population[j] = sfla.Population[j], sfla.Population[i]
	})
}
