package main

import (
	"encoding/json"
	"flag"
	"fmt"
	test "graduate_work/algos"
	"strings"
	"time"
)

func CopySlice(slice [][]float64) [][]float64 {
	copySlice := make([][]float64, len(slice))
	for i := range slice {
		copySlice[i] = make([]float64, len(slice[i]))
		copy(copySlice[i], slice[i])
	}
	return copySlice
}

type RequestType any

func TestAlgo[T RequestType](request T, constructor func(T) (test.Algorithm, error)) {
	algorithm, err := constructor(request)
	if err != nil {
		fmt.Println("Ошибка инициализации алгоритма:", err)
		return
	}

	var history [][][]float64

	start := time.Now()
	bestPos, bestVal := algorithm.Run(func(resp test.Response) error {
		history = append(history, CopySlice(resp.StepPositions))
		return nil
	})
	elapsed := time.Since(start)

	result := map[string]any{
		"history":       history,
		"best_position": bestPos,
		"best_value":    bestVal,
		"time":          elapsed.Seconds(),
	}

	jsonOutput, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsonOutput))
}

func main() {
	// флаги для командной строки
	algoName := flag.String("algorithm", "GWO", "Название алгоритма (например, GWO, AFSA, SFLA)")
	function := flag.String("function", "x+y", "Целевая функция")
	iterations := flag.Int("iterations", 100, "Количество итераций")
	bounds := flag.String("bounds", "[[-5,5],[-5,5]]", "Границы для каждого измерения (например, [-5,5],[-5,5])")
	population := flag.String("population", "", "Начальная популяция")
	population_size := flag.Int("population_size", 50, "Размер начальной популяции")
	seed := flag.Int("seed", 1, "Seed")

	// GWO
	initialA := flag.Float64("initialA", 2.0, "Начальное значение A")
	initialC := flag.Float64("initialC", 2.0, "Начальное значение C")

	// ABC
	limit := flag.Int("limit", 100, "Параметр limit для алгоритма ABC")
	foragerSize := flag.Int("foragerSize", 10, "Параметр foragerSize для алгоритма ABC")

	// AFSA
	eta := flag.Float64("eta", 1e-4, "Параметр eta для алгоритма AFSA")
	maxTries := flag.Int("maxTries", 5, "Параметр maxTries для алгоритма AFSA")
	visual := flag.String("visual", "[1, 8]", "Визуальный диапазон для алгоритма AFSA")
	teta := flag.Float64("teta", 1, "Параметр teta для алгоритма AFSA")

	// FA
	beta0 := flag.Float64("beta0", 1.0, "Параметр beta0 для алгоритма FA")
	gamma := flag.Float64("gamma", 1.0, "Параметр gamma для алгоритма FA")
	alpha := flag.Float64("alpha", 0.1, "Параметр alpha для алгоритма FA")

	// SFLA
	subpopulationsCount := flag.Int("subpopulationsCount", 3, "Количество подгрупп для алгоритма SFLA")
	iMax := flag.Int("iMax", 100, "Количество итераций для SFLA")

	flag.Parse()

	parsedBounds, err := parseMatrix(*bounds)
	if err != nil {
		fmt.Println("Ошибка при разборе границ:", err)
		return
	}

	// общий запрос
	algoRequest := test.AlgoRequest{
		Func:       *function,
		Iterations: *iterations,
		Bounds:     parsedBounds,
		Seed:       seed,
	}

	if *population == "" {
		algoRequest.PopulationSize = population_size
	} else {
		parsedPopulation, err := parseMatrix(*population)
		if err != nil {
			fmt.Println("Ошибка при разборе популяции:", err)
			return
		}

		algoRequest.Population = parsedPopulation
	}

	var request RequestType
	var constructor func(RequestType) (test.Algorithm, error)

	switch strings.ToUpper(*algoName) {
	case "GWO":
		gwoRequest := test.GWORequest{
			AlgoRequest: algoRequest,
			InitialA:    initialA,
			InitialC:    initialC,
		}
		request = gwoRequest
		constructor = func(req RequestType) (test.Algorithm, error) {
			return test.NewGWO(req.(test.GWORequest))
		}

	case "AFSA":
		parsedVisual, err := parseSlice(*visual)
		if err != nil {
			fmt.Println("Ошибка при разборе визуального диапазона AFSA:", err)
			return
		}
		afsaRequest := test.AFSARequest{
			AlgoRequest: algoRequest,
			Eta:         eta,
			MaxTries:    maxTries,
			Visual:      parsedVisual,
			Teta:        teta,
		}
		request = afsaRequest
		constructor = func(req RequestType) (test.Algorithm, error) {
			return test.NewAFSA(req.(test.AFSARequest))
		}

	case "ABC":
		abcRequest := test.ABCRequest{
			AlgoRequest: algoRequest,
			Limit:       limit,
			ForagerSize: foragerSize,
		}
		request = abcRequest
		constructor = func(req RequestType) (test.Algorithm, error) {
			return test.NewABC(req.(test.ABCRequest))
		}

	case "FA":
		fireflyRequest := test.FARequest{
			AlgoRequest: algoRequest,
			Beta0:       beta0,
			Gamma:       gamma,
			Alpha:       alpha,
		}
		request = fireflyRequest
		constructor = func(req RequestType) (test.Algorithm, error) {
			return test.NewFA(req.(test.FARequest))
		}

	case "SFLA":
		sflaRequest := test.SFLARequest{
			AlgoRequest:         algoRequest,
			SubpopulationsCount: subpopulationsCount,
			IMax:                iMax,
		}
		request = sflaRequest
		constructor = func(req RequestType) (test.Algorithm, error) {
			return test.NewSFLA(req.(test.SFLARequest))
		}

	default:
		fmt.Println("Неизвестный алгоритм:", *algoName)
		return
	}

	TestAlgo(request, constructor)
}

// парсинг строки в срез
func parseSlice(str string) ([]float64, error) {
	var slice []float64
	err := json.Unmarshal([]byte(str), &slice)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return nil, err
	}

	return slice, nil
}

// парсинг строки в матрицу
func parseMatrix(str string) ([][]float64, error) {
	var matrix [][]float64
	err := json.Unmarshal([]byte(str), &matrix)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return nil, err
	}

	return matrix, nil
}
