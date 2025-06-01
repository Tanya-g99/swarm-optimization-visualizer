const math = require('mathjs');
const seedrandom = require('seedrandom');
const { parentPort, workerData } = require('worker_threads');

class AFSA {
    constructor(func, nDim, maxIter, populationSize = 50, maxTryNum = 5, eta = 1e-4, visual = [1, 8], bounds = null, initialPopulation = null, teta = 1, seed = null) {
        // Инициализация генератора случайных чисел с заданным seed
        if (seed !== null) {
            this.rng = seedrandom(seed); // Используем seedrandom для генерации случайных чисел
        } else {
            this.rng = Math.random; // В противном случае используем стандартный Math.random
        }

        this.func = func;
        this.nDim = nDim;
        this.maxIter = maxIter;
        this.maxTryNum = maxTryNum;
        this.eta = eta;
        this.minVisual = visual[0];
        this.initialVisual = visual[1];
        this.teta = teta;

        // Установка границ
        this.bounds = bounds || Array(nDim).fill([-100, 100]);

        if (initialPopulation === null) {
            this.populationSize = populationSize;
            this.population = this.createRandomPopulation();
        } else {
            this.population = initialPopulation;
            this.populationSize = initialPopulation.length;
        }

        [this.globalBestPosition, this.globalBestValue] = this.population.reduce(
            (best, position) => {
                const value = this.func(...position);
                return value < best[1] ? [position, value] : best;
            },
            [this.population[0], this.func(...this.population[0])]
        );
        this.historyBest = [this.globalBestValue];
    }

    createRandomPopulation() {
        return Array.from({ length: this.populationSize }, () =>
            Array.from({ length: this.nDim }, (_, j) => {
                const [lowerBound, upperBound] = this.bounds[j];
                return this.rng() * (upperBound - lowerBound) + lowerBound;
            })
        );
    }

    run() {
        parentPort.postMessage({ iteration: 0, stepPositions: this.population.slice(), bestPosition: this.globalBestPosition, bestValue: this.globalBestValue });
        
        let stagnationCount = 0;

        // Предварительные вычисления (векторизация)
        const distanceMatrix = this.calculateDistanceMatrix();

        for (let t = 0; t < this.maxIter; t++) {
            const stepPositions = [];  // Сохраняем позиции на каждом шаге

            // Уменьшение визуального диапазона до минимального значения
            this.visual = Math.max(this.minVisual, this.initialVisual * (1 - t / this.maxIter));

            for (let i = 0; i < this.populationSize; i++) {
                // Поиск соседей для рыбы i
                const V_i = this.findNeighbors(i, distanceMatrix);

                let y_i;
                if (V_i.length === 0) {  // Случайное поведение
                    y_i = this.randomMove(this.population[i]);
                } else {
                    if (V_i.length / this.populationSize > this.teta) {  // Заполненный визуальный диапазон
                        y_i = this.searchBehavior(i, V_i);
                    } else {
                        const c_i = this.meanPosition(V_i);
                        if (this.func(...c_i) < this.func(...this.population[i])) {  // Роение
                            y_i = this.swarmBehavior(c_i, this.population[i]);
                        } else {  // Поиск
                            y_i = this.searchBehavior(i, V_i);
                        }
                    }

                    // Поведение преследования
                    const j_star = this.bestNeighbor(V_i);
                    if (j_star !== null) {
                        if (this.func(...this.population[j_star]) < this.func(...this.population[i])) {
                            y_i = this.chaseBehavior(i, j_star);
                        } else {
                            y_i = this.searchBehavior(i, V_i);
                        }
                    }
                }

                // Обновление
                const f_y_i = this.func(...y_i);
                const f_population_i = this.func(...this.population[i]);

                if (f_y_i < f_population_i) {
                    this.population[i] = y_i;
                }

                // Обновление лучшего решения
                if (f_y_i < this.globalBestValue) {
                    this.globalBestValue = f_y_i;
                    this.globalBestPosition = y_i;
                }

                // Добавляем текущие позиции для рисования траекторий
                stepPositions.push(this.population[i]);
            }

            this.historyBest.push(this.globalBestValue);

            parentPort.postMessage({ iteration: t + 1, stepPositions: stepPositions, bestPosition: this.globalBestPosition, bestValue: this.globalBestValue });

            // Проверка на стагнацию
            if (t > 0 && Math.abs(this.historyBest[this.historyBest.length - 1] - this.historyBest[this.historyBest.length - 2]) < this.eta) {
                stagnationCount++;
            } else {
                stagnationCount = 0;
            }

            if (stagnationCount > this.maxTryNum) {  // Условие остановки
                const j = Math.floor(this.rng() * this.populationSize);
                this.population[j] = this.jumpBehavior(this.population[j]);
            }
        }

        return { bestPosition: this.globalBestPosition, bestValue: this.globalBestValue };
    }

    calculateDistanceMatrix() {
        const matrix = [];
        for (let i = 0; i < this.populationSize; i++) {
            matrix[i] = [];
            for (let j = 0; j < this.populationSize; j++) {
                matrix[i][j] = this.distance(this.population[i], this.population[j]);
            }
        }
        return matrix;
    }

    distance(x1, x2) {
        return Math.sqrt(x1.reduce((sum, xi, index) => sum + Math.pow(xi - x2[index], 2), 0));
    }

    findNeighbors(i, distanceMatrix) {
        const v = this.visual * Math.max(...this.bounds.map(b => b[1] - b[0]));
        return distanceMatrix[i].reduce((neighbors, dist, j) => {
            if (dist <= v && j !== i) neighbors.push(j);
            return neighbors;
        }, []);
    }

    randomMove(x_i) {
        const v = this.visual * Math.max(...this.bounds.map(b => b[1] - b[0]));
        const move = Array.from({ length: this.nDim }, () => (this.rng() * 2 - 1) * Math.min(v, this.bounds[0][1] - this.bounds[0][0]));
        return x_i.map((xi, index) => Math.min(Math.max(xi + move[index], this.bounds[index][0]), this.bounds[index][1]));
    }

    chaseBehavior(i, j_star) {
        const r = this.rng();
        return this.population[i].map((xi, index) => Math.min(Math.max(xi + r * (this.population[j_star][index] - xi), this.bounds[index][0]), this.bounds[index][1]));
    }

    swarmBehavior(c_i, x_i) {
        const r = this.rng();
        return x_i.map((xi, index) => Math.min(Math.max(xi + r * (c_i[index] - xi), this.bounds[index][0]), this.bounds[index][1]));
    }

    searchBehavior(i, V_i) {
        const j = V_i[Math.floor(this.rng() * V_i.length)];
        const r = this.rng();
        return this.population[i].map((xi, index) => Math.min(Math.max(xi + r * (this.population[j][index] - xi), this.bounds[index][0]), this.bounds[index][1]));
    }

    jumpBehavior(x_i) {
        const p = this.rng();
        const jumpDirection = Array.from({ length: this.nDim }, () => (this.rng() * 2 - 1) * p);
        return x_i.map((xi, index) => Math.min(Math.max(xi + jumpDirection[index] * (this.bounds[index][1] - this.bounds[index][0]), this.bounds[index][0]), this.bounds[index][1]));
    }

    bestNeighbor(V_i) {
        if (V_i.length === 0) return null;
        const funcValues = V_i.map(v => this.func(...this.population[v]));
        return V_i[funcValues.indexOf(Math.min(...funcValues))];
    }

    meanPosition(V_i) {
        const positions = V_i.map(i => this.population[i]);
        const summed = positions.reduce((sum, pos) => sum.map((v, index) => v + pos[index]), Array(this.nDim).fill(0));
        return summed.map(v => v / positions.length);
    }
}

// // Функция Растригина
// function rastrigin(x) {
//     return 10 * x.length + x.reduce((sum, xi) => sum + xi ** 2 - 10 * Math.cos(2 * Math.PI * xi), 0);
// }

// const initialPopulation = [
//     [-4.5, -1.2],
//     [3.5, 1.5],
//     [-2.0, -8.0],
//     [-1.0, 5.5],
//     [3.0, -7.0],
//     [-9.5, -6.2],
//     [8.5, 6.5],
//     [2.0, 8.0],
//     [-9.0, 9.5],
//     [-8.0, -9.0],
// ];

// const afsa = new AFSA(
//     rastrigin,                 // func
//     2,                         // nDim (размерность задачи, например, 2 для задачи Растригина)
//     500,                        // maxIter (максимальное количество итераций)
//     undefined,                        // populationSize (размер популяции)
//     undefined,                      // maxTryNum (по умолчанию 5, можно не передавать)
//     undefined,                      // eta (по умолчанию 1e-4, можно не передавать)
//     undefined,                    // visual (по умолчанию [1, 8], можно не передавать)
//     [[-10, 10], [-10, 10]],    // bounds (границы для каждой размерности)
//     initialPopulation,         // initialPopulation (если передается, используйте вашу начальную популяцию)
//     undefined,                         // teta (по умолчанию 1, можно не передавать)
//     101010                         // seed (если нужен seed для воспроизводимости)
// );
// const result = afsa.run();
// console.log("Лучшее положение:", result.bestPosition);
// console.log("Лучшее значение функции:", result.bestValue);
// console.log(workerData.func);
const afsa = new AFSA(
    (...position) => {
        const code = math.compile(workerData.targetFunction);
        if (position.length === 2) {
            const result = code.evaluate({ x: position[0], y: position[1] });
            return result;
        } else {
            console.error("Invalid position array length. Expected 2 values.");
            return null;
        }
    },
    2,                            // Размерность
    workerData.maxIter,           // Максимальное количество итераций
    workerData.populationSize,    // Размер популяции
    workerData.maxTryNum,         // Максимальное количество попыток
    workerData.eta,               // Порог для стагнации
    workerData.visual,            // Диапазон визуализации
    workerData.bounds,            // Границы функции
    workerData.initialPopulation, // Начальная популяция
    workerData.teta,              // Параметр тета
    workerData.seed               // seed для генератора случайных чисел
);

// Запуск алгоритма
console.log(afsa.run());

