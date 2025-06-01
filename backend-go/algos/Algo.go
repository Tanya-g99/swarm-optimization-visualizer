package algos

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/expr-lang/expr"
	"github.com/seehuhn/mt19937"
)

type Algorithm interface {
	Run(send func(Response) error) ([]float64, float64)
}

type AlgoRequest struct {
	Func           string      `json:"targetFunction"`
	Iterations     int         `json:"maxIter"`
	Bounds         [][]float64 `json:"bounds,omitempty"`
	Population     [][]float64 `json:"initialPopulation,omitempty"`
	PopulationSize *int        `json:"populationSize,omitempty"`
	Seed           *int        `json:"seed,omitempty"`

	NumDimensions *int `json:"numDimensions"`
}

type Algo struct {
	Func           func([]float64) float64
	Iterations     int
	Bounds         [][]float64
	Population     [][]float64
	PopulationSize int
	NumDimensions  int

	GlobalBestPosition []float64
	GlobalBestValue    float64

	Rng *rand.Rand
}

func NewAlgo(request AlgoRequest) (*Algo, error) {

	if request.Bounds != nil && len(request.Bounds) != 2 {
		return nil, errors.New("неверно заданы границы поиска")
	}

	if request.Population != nil && len(request.Population) < 1 {
		return nil, errors.New("пустая популяция")
	}

	if request.Population == nil && request.PopulationSize == nil {
		return nil, errors.New("недостаточно данных для создания популяции")
	}

	if request.NumDimensions == nil && len(request.Population) < 1 && request.Bounds == nil {
		return nil, errors.New("не задана размерность задачи")
	}

	if request.Bounds != nil && request.Population != nil && len(request.Bounds[0]) != len(request.Population[0]) {
		return nil, errors.New("несоответствие размерностей границ и популяции")
	}

	if request.Bounds != nil && request.NumDimensions != nil && len(request.Bounds[0]) != *request.NumDimensions {
		return nil, errors.New("несоответствие размерности границ и размерности задачи")
	}

	if request.NumDimensions != nil && request.Population != nil && *request.NumDimensions != len(request.Population[0]) {
		return nil, errors.New("несоответствие размерности популяции и размерности задачи")
	}

	function, err := ConvertMathExpressionToFunc(request.Func)
	if err != nil {
		return nil, errors.New("Ошибка компиляции функции:" + err.Error())
	}

	algo := &Algo{
		Func:       function,
		Iterations: request.Iterations,

		GlobalBestPosition: nil,
		GlobalBestValue:    math.Inf(1),
	}

	src := mt19937.New()
	if request.Seed != nil {
		src.Seed(int64(*request.Seed))
	} else {
		src.Seed(time.Now().UnixNano())
	}
	algo.Rng = rand.New(src)

	if request.Bounds == nil {
		if request.NumDimensions == nil {
			algo.NumDimensions = len(request.Population[0])
		}
		algo.Bounds = make([][]float64, algo.NumDimensions)
		for i := range algo.Bounds {
			algo.Bounds[i] = []float64{-100, 100}
		}
	} else if len(request.Bounds) > 0 {
		algo.Bounds = request.Bounds
		algo.NumDimensions = len(algo.Bounds[0])
	}

	if request.Population == nil {
		algo.PopulationSize = *request.PopulationSize
		algo.Population = make([][]float64, algo.PopulationSize)
		for i := range algo.PopulationSize {
			algo.Population[i] = make([]float64, algo.NumDimensions)
			for j := range algo.NumDimensions {
				algo.Population[i][j] = algo.Rng.Float64()*(algo.Bounds[j][1]-algo.Bounds[j][0]) + algo.Bounds[j][0]
			}
		}
	} else {
		algo.Population = request.Population
		algo.PopulationSize = len(algo.Population)
	}

	for i := range algo.PopulationSize {
		value := algo.Func(algo.Population[i])
		if value < algo.GlobalBestValue {
			algo.GlobalBestValue = value
			algo.GlobalBestPosition = algo.Population[i]
		}
	}

	return algo, nil
}

func ConvertMathExpressionToFunc(expression string) (func([]float64) float64, error) {
	env := map[string]interface{}{
		"x":    0.0,
		"y":    0.0,
		"sin":  math.Sin,
		"cos":  math.Cos,
		"tan":  math.Tan,
		"PI":   math.Pi,
		"log":  math.Log,
		"ln":   math.Log,
		"sqrt": math.Sqrt,
		"abs":  math.Abs,
		"pow":  math.Pow,
		"exp":  math.Exp,
	}

	program, err := expr.Compile(expression, expr.Env(env))
	if err != nil {
		return nil, err
	}

	return func(vars []float64) float64 {
		if len(vars) != 2 {
			fmt.Println("Ошибка: ожидался массив из 2 элементов [x, y]")
			return math.NaN()
		}

		env["x"] = vars[0]
		env["y"] = vars[1]

		output, err := expr.Run(program, env)
		if err != nil {
			fmt.Println("Ошибка вычисления:", err)
			return math.NaN()
		}

		if math.IsNaN(output.(float64)) {
			return math.Inf(1)
		}
		return output.(float64)
	}, nil
}
