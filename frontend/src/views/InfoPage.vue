<template>
    <section class="usage-page">
        <div class="usage-card">
            <h1 class="heading">Руководство пользователя</h1>

            <p class="text">
                Это веб-приложение предназначено для визуализации работы популяционных алгоритмов оптимизации, включая:
            </p>
            <ul class="text">
                <li><strong>SFLA</strong> — Shuffled Frog Leaping Algorithm</li>
                <li><strong>AFSA</strong> — Artificial Fish Swarm Algorithm</li>
                <li><strong>FA</strong> — Firefly Algorithm</li>
                <li><strong>ABC</strong> — Artificial Bee Colony</li>
                <li><strong>GWO</strong> — Grey Wolf Optimizer</li>
            </ul>

            <p class="text">
                Приложение позволяет наглядно проследить, как алгоритмы минимизируют заданную целевую функцию.
            </p>

            <h2 class="subheading">Форма параметров (левая панель)</h2>

            <ul class="text">
                <li><strong>Целевая функция</strong>
                    <ul>
                        <li>Может быть выбрана из предложенного списка (например, Rastrigin, Ackley) или задана вручную.
                        </li>
                        <li>Формат: выражение, содержащее переменные и математические функции.</li>
                        <li>Допустимые переменные: <code>x</code> и <code>y</code>.</li>
                        <li>Допустимые операторы: <code>+</code>, <code>-</code>, <code>*</code>, <code>/</code>,
                            <code>^</code> (возведение в степень).</li>
                        <li>Допустимые функции:
                            <ul>
                                <li><code>sin(x)</code> — синус</li>
                                <li><code>cos(x)</code> — косинус</li>
                                <li><code>tan(x)</code> — тангенс</li>
                                <li><code>ln(x)</code> — натуральный логарифм</li>
                                <li><code>sqrt(x)</code> — квадратный корень</li>
                                <li><code>abs(x)</code> — модуль</li>
                                <li><code>exp(x)</code> — экспонента</li>
                            </ul>
                        </li>
                        <li>Примеры допустимых выражений:
                            <ul>
                                <li><code>sin(x) + y^2</code></li>
                                <li><code>abs(x - y) + log(x^2 + y^2)</code></li>
                                <li><code>exp(-sqrt(x^2 + y^2)) * cos(x)</code></li>
                            </ul>
                        </li>
                        <li>Ошибки синтаксиса будут отображены под полем ввода.</li>
                    </ul>
                </li>

                <!--  -->

                <li><strong>Границы X и Y</strong>
                    <ul>
                        <li>Формат: два значения — минимум и максимум.</li>
                        <li>Диапазон не должен превышать 500 единиц.</li>
                        <li><code>min</code> всегда должен быть меньше <code>max</code>.</li>
                    </ul>
                </li>

                <li><strong>Количество итераций</strong>
                    <ul>
                        <li>Целое число от 1 до 1000.</li>
                    </ul>
                </li>

                <li><strong>Размер популяции</strong>
                    <ul>
                        <li>Целое число от 1 до 1000.</li>
                    </ul>
                </li>

                <li><strong>Начальная популяция (опционально)</strong>
                    <ul>
                        <li>Можно задать вручную координаты X и Y для каждой особи.</li>
                        <li>Если не задано — координаты генерируются случайно.</li>
                    </ul>
                </li>

                <li><strong>Seed</strong>
                    <ul>
                        <li>Необязательное поле для задания числа генератора случайных чисел (для воспроизводимости).
                        </li>
                    </ul>
                </li>
            </ul>

            <h2 class="subheading">Визуализация результатов</h2>
            <p class="text">
                После запуска алгоритма отображаются два графика:
            </p>
            <ul class="text">
                <li><strong>2D-график (вид сверху)</strong> — текущие позиции агентов на плоскости XY.</li>
                <li><strong>3D-график</strong> — поверхность целевой функции и расположение агентов в пространстве
                    XYZ.</li>
            </ul>

            <p class="text">В левом нижнем углу 3D-графика отображаются:</p>
            <ul class="text">
                <li><strong>Итерация:</strong> номер текущего шага.</li>
                <li><strong>Лучшее решение:</strong> координаты X и Y.</li>
                <li><strong>Лучшее значение:</strong> значение целевой функции в этой точке.</li>
            </ul>

            <h2 class="subheading">Уровень детализации графика</h2>
            <p class="text">
                Под графиком доступен ползунок <strong>"Уровень детализации"</strong>, управляющий числом сегментов
                сетки:
            </p>
            <ul class="text">
                <li><strong>Диапазон:</strong> от 100 до 500 (с шагом 100).</li>
                <li>Больше сегментов — выше точность, но ниже производительность.</li>
            </ul>

            <h2 class="subheading">Общие рекомендации</h2>
            <ul class="text">
                <li>Избегайте слишком больших диапазонов по X и Y — это замедляет визуализацию.</li>
                <li>Убедитесь, что формула функции корректна — особенно при ручном вводе.</li>
                <li>Для сложных функций начинайте с небольшого числа итераций.</li>
                <li>3D-график — для оценки рельефа и положения агентов, 2D — для анализа распределения.</li>
            </ul>
        </div>
    </section>
</template>

<style scoped>
.usage-page {
    display: flex;
    justify-content: center;
    padding: 0 2rem;
    background-color: var(--color-background);
    color: var(--color-text);
}

.usage-card {
    background-color: var(--color-card-bg);
    padding: 2rem;
    max-width: 1200px;
    width: 100%;
    border-radius: 1rem;
    border: 1px solid var(--color-border);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

.heading {
    font-size: 2rem;
    color: var(--color-heading);
    margin-bottom: 1rem;
}

.subheading {
    font-size: 1.4rem;
    color: var(--color-subheading);
    margin-top: 1.5rem;
    margin-bottom: 0.5rem;
}

.text {
    font-size: 1rem;
    line-height: 1.6;
    margin-bottom: 1rem;
}

ul.text {
    list-style-type: disc;
    padding-left: 1.5rem;
}

ul.text>li {
    margin-bottom: 0.5rem;
}

ul.text ul {
    list-style-type: circle;
    padding-left: 1.2rem;
}

code {
    background-color: var(--color-background-soft);
    padding: 0.2em 0.4em;
    border-radius: 4px;
    font-family: monospace;
}
</style>