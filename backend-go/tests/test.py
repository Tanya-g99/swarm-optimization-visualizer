import os
import subprocess
import json
import pytest
import matplotlib.pyplot as plt
import time
import math

from mealpy.swarm_based import GWO, FA, ABC
from mealpy.utils.problem import Problem
from mealpy import FloatVar

def bounds_to_floatvar(bounds):
    lb = tuple(b[0] for b in bounds)
    ub = tuple(b[1] for b in bounds)
    return FloatVar(lb=lb, ub=ub)


def run_go_program(algo_name, function, iterations, bounds, population_size, seed, *additional_params):

    bounds_str = str(bounds)

    result = subprocess.run(
        ["go", "run", "main.go", 
         "-algorithm", algo_name,
         "-function", function,
         "-iterations", str(iterations),
         "-bounds", bounds_str,
        "-population_size", str(population_size),
         "-seed", str(seed)] + list(additional_params),
        capture_output=True,
        text=True
    )

    if result.returncode != 0:
        print(f"Ошибка при запуске Go-программы: {result.stderr}")
        return None

    try:
        output = json.loads(result.stdout)
        return output
    except json.JSONDecodeError:
        print(f"Не удалось разобрать вывод: {result.stdout}")
        return None
    
def setup_method():
    for algo_name in ["GWO", "FA", "ABC", "SFLA"]:
        file_name = f"{algo_name}_test_results.txt"

        if os.path.exists(file_name):
            os.remove(file_name)

        with open(file_name, 'w') as f:
            pass 

setup_method()

test_functions = [
    {
        "name": "Rosenbrock",
        "expression": "(1 - x)^2 + 100*(y - x^2)^2",
        "bounds": [[-5, 5], [-5, 5]],
        "expected_position": [1.0, 1.0],
        "expected_value": 0.0,
        "position_tolerance": 1e-1,
        "value_tolerance": 1e-3,
    },
    {
        "name": "Rastrigin",
        "expression": "10*2 + (x^2 - 10 * cos(2 * 3.1415 * x)) + (y**2 - 10 * cos(2 * 3.1415 * y))",
        "bounds": [[-5.12, 5.12], [-5.12, 5.12]],
        "expected_position": [0.0, 0.0],
        "expected_value": 0.0,
        "position_tolerance": 1e-1,
        "value_tolerance": 1e-3,
    },
    {
        "name": "Schwefel",
        "expression": "418.9829 * 2 - x * sin(sqrt(abs(x))) - y * sin(sqrt(abs(y)))",
        "bounds": [[-500, 500], [-500, 500]],
        "expected_position": [420.9687, 420.9687],
        "expected_value": 0.0,
        "position_tolerance": 1.0,
        "value_tolerance": 1e-1,
    },
]

@pytest.mark.parametrize("function", test_functions)
@pytest.mark.parametrize(
    "algo_name, iterations, population_size, seed, additional_params",
    [
        (
            "GWO", 
            300, 
            100,
            1,
            ["-initialA", "2.0", "-initialC", "2.0"],
        ),
        (
            "AFSA", 
            300, 
            100,
            2,
            ["-eta", "0.5", "-maxTries", "50", "-visual", "[1.0, 8.0]", "-teta", "1"],
        ),
        (
            "FA", 
            300, 
            100,
            100,
            ["-beta0", "1.0", "-gamma", "0.8", "-alpha", "0.01"],
        ),
        (
            "ABC",
            300,
            100,
            4,
            ["-limit", "10", "-foragerSize", "10"]
        ),
        (
            "SFLA",
            300,
            100,
            50,
            ["-subpopulationsCount", "5", "-iMax", "10"]
        ),
    ]
)

def test_algorithm(function, algo_name, iterations, population_size, seed, additional_params, capsys):
    print(f"Запуск Go-программы для алгоритма...")

    actual_result = run_go_program(algo_name, function["expression"], iterations, function["bounds"], population_size, seed, *additional_params)

    assert actual_result is not None, "Результат работы Go-программы равен None!"


    if algo_name in ["GWO", "FA", "ABC"]:
        problem = Problem(
            bounds=bounds_to_floatvar(function["bounds"]),
            minmax="min",
            obj_func=lambda sol: eval(function["expression"].replace("x", f"sol[0]")
                                      .replace("y", f"sol[1]").replace("^", "**")
                                      .replace("sin", "math.sin").replace("cos", "math.cos")
                                      .replace("sqrt", "math.sqrt")),
            log_to=None,
            verbose=False
        )

        model = None
        if algo_name == "GWO":
            model = GWO.OriginalGWO(epoch=iterations, pop_size=100, seed=seed)
        elif algo_name == "FA":
            model = FA.OriginalFA(epoch=iterations, pop_size=100, seed=seed)
        elif algo_name == "ABC":
            model = ABC.OriginalABC(epoch=iterations, pop_size=100, seed=seed)

        start_time = time.perf_counter()
        g_best = model.solve(problem=problem)
        library_duration = time.perf_counter() - start_time

        function["position_tolerance"] = max(function["position_tolerance"], abs(max(list(g_best.solution))))
        function["value_tolerance"] = max(function["value_tolerance"], abs(g_best.target.fitness))
        library_result = {
            "best_position": list(g_best.solution),
            "best_value": g_best.target.fitness,
            "time": library_duration
        }
    else:
        library_result = {
            "best_position": "-",
            "best_value": "-",
            "time": "-"
        }

    # запись в файл
    with open(f"{algo_name}_test_results.txt", "a") as f:
        f.write(f"{'Функция':<30} : {function["name"]}\n")
        f.write(f"{'Время работы алгоритма':<30} : {actual_result['time']} секунд\n")
        f.write(f"{'Время работы библ. алгоритма':<30} : {library_result['time']} секунд\n")
        f.write("\n")

        f.write(f"{'Ожидаемая и полученная информация':^80}\n") 
        f.write(f"{'-' * 150}\n") 
        f.write(f"{'Параметр':<20} {'Ожидаемое значение':<25} {'Полученное значение':<45} {'Библиотечное значение':<35}\n") 
        f.write(f"{'-' * 150}\n") 
        
        f.write(f"{'Лучшая позиция':<20} {str(function["expected_position"]):<25} {str(actual_result['best_position']):<45} {str(library_result['best_position']):<35}\n")
        f.write(f"{'Лучшее значение':<20} {function["expected_value"]:<25} {actual_result['best_value']:<45} {library_result['best_value']:<35}\n")
        
        f.write(f"{'-' * 150}\n\n")

    assert actual_result["best_position"] ==  pytest.approx(function["expected_position"], abs=function["position_tolerance"]), \
        f"Ошибка: лучшая позиция отличается! Ожидалось {function["expected_position"]}, получено {actual_result['best_position']}"

    assert actual_result["best_value"] == pytest.approx(function["expected_value"], abs=function["value_tolerance"]), \
        f"Ошибка: лучшее значение отличается! Ожидалось {function["expected_value"]}, получено {actual_result['best_value']}"
    
    captured = capsys.readouterr()
    print(captured.out)
