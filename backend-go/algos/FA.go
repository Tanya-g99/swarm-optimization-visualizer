package algos

import (
	"math"
)

type FARequest struct {
	AlgoRequest

	Beta0 *float64 `json:"beta0,omitempty"`
	Gamma *float64 `json:"gamma,omitempty"`
	Alpha *float64 `json:"alpha,omitempty"`
}

type FA struct {
	Algo

	Beta0 float64
	Gamma float64
	Alpha float64
}

func NewFA(request FARequest) (Algorithm, error) {
	algo, err := NewAlgo(request.AlgoRequest)
	if err != nil {
		return nil, err
	}

	return &FA{
		Algo:  *algo,
		Alpha: setDefault(request.Alpha, 0.01),
		Beta0: setDefault(request.Beta0, 1.0),
		Gamma: setDefault(request.Gamma, 0.8),
	}, nil
}

func (fa *FA) Run(send func(Response) error) ([]float64, float64) {
	err := send(Response{
		StepPositions: fa.Population,
		BestPosition:  fa.GlobalBestPosition,
		BestValue:     fa.GlobalBestValue,
		Iteration:     0,
	})
	if err != nil {
		return fa.GlobalBestPosition, fa.GlobalBestValue
	}

	for t := range fa.Iterations {
		for i := range fa.PopulationSize {
			maxDistance := 0.0
			for j := range fa.PopulationSize {
				if i != j {
					// расстояние между светлячками i и j
					rij := 0.0
					for k := range fa.NumDimensions {
						rij += math.Pow(fa.Population[i][k]-fa.Population[j][k], 2)
					}
					rij = math.Sqrt(rij)
					if rij > maxDistance {
						maxDistance = rij
					}
				}
			}

			for j := range fa.PopulationSize {
				if i == j {
					continue
				}
				if fa.Func(fa.Population[j]) < fa.Func(fa.Population[i]) {
					fa.Population[i] = fa.UpdatePosition(fa.Population[i], fa.Population[j], maxDistance)
				}
			}
			if fa.Func(fa.Population[i]) < fa.GlobalBestValue {
				fa.GlobalBestValue = fa.Func(fa.Population[i])
				fa.GlobalBestPosition = fa.Population[i]
			}
		}

		err = send(Response{
			StepPositions: fa.Population,
			BestPosition:  fa.GlobalBestPosition,
			BestValue:     fa.GlobalBestValue,
			Iteration:     t + 1,
		})
		if err != nil {
			return fa.GlobalBestPosition, fa.GlobalBestValue
		}
	}
	return fa.GlobalBestPosition, fa.GlobalBestValue
}

func (fa *FA) UpdatePosition(xi, xj []float64, maxDistance float64) []float64 {
	newPosition := make([]float64, len(xi))

	rij := fa.calculateDistance(xi, xj)
	gamma_i := fa.Gamma / maxDistance

	for k := range len(xi) {
		randValue := fa.Rng.Float64()
		directedComponent := fa.Beta0 * math.Exp(-gamma_i*rij*rij) * (xj[k] - xi[k])

		newPosition[k] = xi[k] + directedComponent + fa.Alpha*(randValue-0.5)
		newPosition[k] = math.Max(fa.Bounds[k][0], math.Min(fa.Bounds[k][1], newPosition[k]))
	}

	return newPosition
}

func (fa *FA) calculateDistance(a, b []float64) float64 {
	sum := 0.0
	for i := range a {
		sum += (a[i] - b[i]) * (a[i] - b[i])
	}
	return math.Sqrt(sum)
}
