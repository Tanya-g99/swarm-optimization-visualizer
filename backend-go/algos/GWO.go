package algos

import (
	"fmt"
	"math"
)

type GWORequest struct {
	AlgoRequest
	InitialA *float64 `json:"InitialA"`
	InitialC *float64 `json:"InitialC"`
}

type GWO struct {
	Algo
	a float64
	C float64

	beta       []float64
	betaValue  float64
	delta      []float64
	deltaValue float64
}

func NewGWO(request GWORequest) (Algorithm, error) {
	algo, err := NewAlgo(request.AlgoRequest)
	if err != nil {
		return nil, err
	}

	return &GWO{
		Algo:       *algo,
		a:          setDefault(request.InitialA, 2.0),
		C:          setDefault(request.InitialC, 2.0),
		beta:       nil,
		betaValue:  math.Inf(1),
		delta:      nil,
		deltaValue: math.Inf(1),
	}, nil
}

func (gwo *GWO) Run(send func(Response) error) ([]float64, float64) {
	gwo.updateBestWolves()
	err := send(Response{
		StepPositions: gwo.Population,
		BestPosition:  gwo.GlobalBestPosition,
		BestValue:     gwo.GlobalBestValue,
		Iteration:     0,
	})

	if err != nil {
		fmt.Println("CONN ERR", err.Error())
		return gwo.GlobalBestPosition, gwo.GlobalBestValue
	}
	for t := range gwo.Iterations {
		a := gwo.a - (gwo.a*float64(t))/float64(gwo.Iterations)

		for i, w := range gwo.Population {
			Xnew := make([]float64, gwo.NumDimensions)
			copy(Xnew, w)
			for j := range gwo.NumDimensions {
				X1 := gwo.hunting(gwo.GlobalBestPosition[j], w[j], a)
				X2 := gwo.hunting(gwo.beta[j], w[j], a)
				X3 := gwo.hunting(gwo.delta[j], w[j], a)

				Xnew[j] = (X1 + X2 + X3) / 3
				Xnew[j] = math.Max(gwo.Bounds[j][0], math.Min(Xnew[j], gwo.Bounds[j][1]))
			}
			if gwo.Func(Xnew) < gwo.Func(w) {
				gwo.Population[i] = Xnew
			}

		}
		gwo.updateBestWolves()

		err := send(Response{
			StepPositions: gwo.Population,
			BestPosition:  gwo.GlobalBestPosition,
			BestValue:     gwo.GlobalBestValue,
			Iteration:     t + 1,
		})
		if err != nil {
			fmt.Println("CONN ERR", err.Error())
			return gwo.GlobalBestPosition, gwo.GlobalBestValue
		}
	}

	return gwo.GlobalBestPosition, gwo.GlobalBestValue
}

func (gwo *GWO) updateBestWolves() {
	for _, wolf := range gwo.Population {
		value := gwo.Func(wolf)
		if value < gwo.GlobalBestValue {
			gwo.delta, gwo.deltaValue = gwo.beta, gwo.betaValue
			gwo.beta, gwo.betaValue = gwo.GlobalBestPosition, gwo.GlobalBestValue
			gwo.GlobalBestPosition, gwo.GlobalBestValue = wolf, value
		} else if value < gwo.betaValue {
			gwo.delta, gwo.deltaValue = gwo.beta, gwo.betaValue
			gwo.beta, gwo.betaValue = gwo.GlobalBestPosition, gwo.GlobalBestValue
			gwo.GlobalBestPosition, gwo.GlobalBestValue = wolf, value
		} else if value < gwo.deltaValue {
			gwo.delta, gwo.deltaValue = wolf, value
		}
	}
}

func (gwo *GWO) hunting(prey, wolf, a float64) float64 {
	r1, r2 := gwo.Rng.Float64(), gwo.Rng.Float64()
	A := a * (2*r1 - 1)
	C := gwo.C * r2
	D := math.Abs(C*prey - wolf)
	return prey - A*D
}
