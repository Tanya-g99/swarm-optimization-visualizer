package algos

import (
	"math"
)

type ABCRequest struct {
	AlgoRequest

	Limit       *int `json:"limit,omitempty"`
	ForagerSize *int `json:"foragerSize,omitempty"`
}

type ABC struct {
	Algo

	Limit        int
	ForagerSize  int
	ObserverSize int
	Trials       []int
}

func NewABC(request ABCRequest) (Algorithm, error) {
	algo, err := NewAlgo(request.AlgoRequest)
	if err != nil {
		return nil, err
	}

	foragerSize := setDefault(request.ForagerSize, algo.PopulationSize/2)
	observerSize := algo.PopulationSize - foragerSize

	return &ABC{
		Algo:         *algo,
		Limit:        setDefault(request.Limit, algo.PopulationSize*algo.NumDimensions/2),
		ForagerSize:  foragerSize,
		ObserverSize: observerSize,
		Trials:       make([]int, algo.PopulationSize),
	}, nil
}

func (abc *ABC) Run(send func(Response) error) ([]float64, float64) {
	err := send(Response{
		StepPositions: abc.Population,
		BestPosition:  abc.GlobalBestPosition,
		BestValue:     abc.GlobalBestValue,
		Iteration:     0,
	})
	if err != nil {
		return abc.GlobalBestPosition, abc.GlobalBestValue
	}

	for t := range abc.Iterations {
		abc.foragerPhase()
		abc.observerPhase()
		abc.scoutPhase()

		err = send(Response{
			StepPositions: abc.Population,
			BestPosition:  abc.GlobalBestPosition,
			BestValue:     abc.GlobalBestValue,
			Iteration:     t + 1,
		})
		if err != nil {
			return abc.GlobalBestPosition, abc.GlobalBestValue
		}
	}

	return abc.GlobalBestPosition, abc.GlobalBestValue
}

func (abc *ABC) foragerPhase() {
	for i := range abc.ForagerSize {
		k := abc.Rng.Intn(abc.ForagerSize) // другой собиратель
		for k == i {
			k = abc.Rng.Intn(abc.ForagerSize)
		}
		candidate := abc.mutate(abc.Population[i], abc.Population[k])

		if abc.Func(candidate) < abc.Func(abc.Population[i]) {
			abc.Population[i] = candidate
			abc.Trials[i] = 0
			abc.updateGlobalBest(i)
		} else {
			abc.Trials[i]++
		}
	}
}

func (abc *ABC) observerPhase() {
	for i := abc.ForagerSize; i < abc.ForagerSize+abc.ObserverSize; i++ {
		j := abc.selectForagerByFitness()  // собиратель
		k := abc.Rng.Intn(abc.ForagerSize) // другой собиратель
		for k == j {
			k = abc.Rng.Intn(abc.ForagerSize)
		}
		candidate := abc.mutate(abc.Population[j], abc.Population[k])

		if abc.Func(candidate) < abc.Func(abc.Population[j]) {
			abc.Population[j] = candidate
			abc.Trials[j] = 0
			abc.updateGlobalBest(j)
		} else {
			abc.Trials[j]++
		}
	}
}

func (abc *ABC) scoutPhase() {
	for i := range abc.ForagerSize { // Только собиратели могут стать разведчиками
		if abc.Trials[i] > abc.Limit {
			abc.Population[i] = abc.randomSolution()
			abc.Trials[i] = 0
			abc.updateGlobalBest(i)
		}
	}
}

func (abc *ABC) mutate(solution, other []float64) []float64 {
	newSolution := make([]float64, len(solution))
	copy(newSolution, solution)
	s := abc.Rng.Intn(len(solution))
	r := abc.Rng.Float64()*2 - 1
	newSolution[s] = math.Max(abc.Bounds[s][0], math.Min(solution[s]+r*(solution[s]-other[s]), abc.Bounds[s][1]))
	return newSolution
}

func (abc *ABC) selectForagerByFitness() int {
	sumFitness := 0.0
	for i := range abc.ForagerSize {
		sumFitness += abc.Func(abc.Population[i])
	}

	threshold := abc.Rng.Float64() * sumFitness
	sum := 0.0
	for i := range abc.ForagerSize {
		sum += abc.Func(abc.Population[i])
		if sum >= threshold {
			return i
		}
	}
	return 0
}

func (abc *ABC) randomSolution() []float64 {
	solution := make([]float64, abc.NumDimensions)
	for i := range solution {
		solution[i] = abc.Bounds[i][0] + abc.Rng.Float64()*(abc.Bounds[i][1]-abc.Bounds[i][0])
	}
	return solution
}

func (abc *ABC) updateGlobalBest(i int) {
	value := abc.Func(abc.Population[i])
	if value < abc.GlobalBestValue {
		abc.GlobalBestValue = value
		copy(abc.GlobalBestPosition, abc.Population[i])
	}
}
