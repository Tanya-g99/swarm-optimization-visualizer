<template>
    <div class="container">
        <div class="left">
            <form class="card-shadow" @submit.prevent="handleSubmit">
                <p class="algo-name">{{ algoName }}</p>
                <div class="form-scroll">
                    <CustomField class="function-field" label="Целевая функция" type="text" v-model="targetFunction"
                        description="Функция, которую необходимо минимизировать." :error="fieldsError.targetFunction">
                        <select class="field-input" v-model="selectedFunction" @change="updateTargetFunction">
                            <option value="" disabled style="display: none;">Выберите функцию</option>
                            <option v-for="(func, key) in functionOptions" :key="key" :value="key">
                                {{ func.label }}
                            </option>
                        </select>
                    </CustomField>


                    <div class="parameter-group">
                        <BoundsField label="Границы X" v-model:left.number="minX" v-model:right.number="maxX"
                            validateRange v-model:error="fieldsError.xBounds" />
                        <BoundsField label="Границы Y" v-model:left.number="minY" v-model:right.number="maxY"
                            validateRange v-model:error="fieldsError.yBounds" />
                        <CustomField label="Количество итераций" type="number" :step="1" v-model.number="maxIter"
                            :error="fieldsError.maxIter" @update:model-value="updateMaxIter"
                            description=" Алгоритм завершит свою работу при достижении указанного количества итераций." />
                    </div>

                    <div class="parameter-group">
                        <CustomField label="Размер популяции" type="number" :step="1" v-model.number="populationSize"
                            @update:model-value="updatePopulationFields" description="Количество особей в популяции."
                            :error="fieldsError.populationSize" />

                        <button v-if="!showInitialPopulation" type="button" @click="initPopulation"
                            class="toggle-btn init-btn btn">
                            Задать начальную популяцию
                        </button>

                        <button v-if="showInitialPopulation" type="button" @click="showInitialPopulation = false"
                            class="toggle-btn remove-btn btn">
                            Убрать начальную популяцию
                        </button>
                    </div>

                    <div v-if="showInitialPopulation" style="padding-bottom: 16px;">
                        <p class="field-label">Начальная популяция</p>
                        <div class="parameter-group population">
                            <div v-for="(coords, index) in initialPopulation" :key="index">
                                <BoundsField v-model:left.number="coords[0]" v-model:right.number="coords[1]"
                                    :label="`Особь ${index + 1}`" />
                            </div>
                        </div>
                    </div>

                    <div class="parameter-group">
                        <!-- Слот для дополнительных полей -->
                        <slot name="additional-fields"></slot>

                        <CustomField label="Seed" type="number" v-model.number="seed"
                            description="Начальное значение для генератора случайных чисел." />
                    </div>
                </div>

                <div class="parameter-group">
                    <button v-if="working" type="button" @click="emit('stop')" class="end-btn btn">Завершить</button>
                    <button v-else type="submit" class="submit-btn btn">Запустить</button>
                </div>
            </form>

            <Graph2D style="height: 100%;" :targetFunction="parseFunction" :points="data?.stepPositions || []"
                :xMin="minX" :xMax="maxX" :yMin="minY" :yMax="maxY" :settings="graphSettings" />

        </div>
        <div class="right">
            <Graph3D :targetFunction="parseFunction" :points="data?.stepPositions || []" :xMin="minX" :xMax="maxX"
                :yMin="minY" :yMax="maxY" :settings="graphSettings" />
            <div class="graph-description">
                <p>Итерация: {{ data?.iteration ? data.iteration : " - " }}</p>
                <p>Лучшее решение: {{ data?.bestPosition ?
                    "(" + data.bestPosition[0].toFixed(3) + "; " + data.bestPosition[1].toFixed(3) + ")"
                    : " - " }}</p>
                <p>Лучшее значение: {{ data?.bestValue || data?.bestValue === 0 ? data.bestValue : " - " }}</p>
                <div class="segments-slider-container">
                    <label for="segments-slider">Уровень детализации:</label>
                    <input id="segments-slider" type="range" min="100" step="100" max="500"
                        v-model="graphSettings.segments" />
                </div>

            </div>
        </div>
    </div>

</template>

<script setup>
import { ref, computed, watch, onBeforeMount } from 'vue';
import CustomField from 'components/CustomField.vue';
import BoundsField from 'components/BoundsField.vue';
import Graph3D from 'src/components/Graph3D.vue';
import Graph2D from 'src/components/Graph2D.vue';
import functions from 'utils/functions';
import { Parser } from 'expr-eval';
import utils from 'utils/storage'

const props = defineProps({
    algoName: String,
    data: Object,
    working: {
        type: Boolean,
        default: false
    }
});

const emit = defineEmits(['submit', 'stop', 'update:data', 'update:populationSize']);

const fieldsError = ref({
    targetFunction: utils.useSessionStorageField('targetFunctionError', '').value,
    populationSize: utils.useSessionStorageField('populationSizeError', '').value,
    maxIter: utils.useSessionStorageField('maxIterError', '').value,
    xBounds: utils.useSessionStorageField('xBoundsError', '').value,
    yBounds: utils.useSessionStorageField('yBoundsError', '').value,
});

const graphSettings = ref({
    minZ: Infinity,
    maxZ: -Infinity,
    xScale: 1,
    yScale: 1,
    zScale: 1,
    segments: 300,
});

const updateSettings = () => {
    const func = parseFunction.value;
    if (!(func === null || fieldsError.value.xBounds || fieldsError.value.yBounds)) {
        graphSettings.value.minZ = Infinity;
        graphSettings.value.maxZ = -Infinity;

        const xstep = (maxX.value - minX.value) / graphSettings.value.segments, ystep = (maxY.value - minY.value) / graphSettings.value.segments

        for (let x = minX.value; x <= maxX.value; x += xstep) {
            for (let y = minY.value; y <= maxY.value; y += ystep) {
                const z = func.evaluate({ x: x, y: y });
                if (!isNaN(z)) {
                    graphSettings.value.minZ = Math.min(graphSettings.value.minZ, z);
                    graphSettings.value.maxZ = Math.max(graphSettings.value.maxZ, z);
                }
            }
        }
        graphSettings.value.xScale = Math.max(100, Math.abs(minX.value), Math.abs(maxX.value)) / 100;
        graphSettings.value.yScale = Math.max(100, Math.abs(minY.value), Math.abs(maxY.value)) / 100;
        graphSettings.value.zScale = Math.max(100, Math.abs(graphSettings.value.minZ), Math.abs(graphSettings.value.maxZ)) / 100;
    }
};


const targetFunction = utils.useSessionStorageField('targetFunction', 'cos(x+y)');
const selectedFunction = utils.useSessionStorageField('selectedFunction', '');

const functionOptions = {
    rastrigin: functions.rastrigin,
    rosenbrock: functions.rosenbrock,
    ackley: functions.ackley,
    levi: functions.levi,
    schwefel: functions.schwefel
};

const updateTargetFunction = () => {
    if (selectedFunction.value && functionOptions[selectedFunction.value]) {
        targetFunction.value = functionOptions[selectedFunction.value].formula;
    }
};

const parseFunction = computed(() => {
    try {
        const result = Parser.parse(targetFunction.value);
        if (isNaN(result.evaluate({ x: minX.value, y: minY.value })) && isNaN(result.evaluate({ x: maxX.value, y: maxY.value }))) {
            fieldsError.value.targetFunction = 'Ошибка парсинга функции';
            return null;
        }
        fieldsError.value.targetFunction = '';
        return result;
    } catch (error) {
        fieldsError.value.targetFunction = 'Ошибка парсинга функции';
        return null;
    }
});

watch(targetFunction, () => {
    updateSettings();
    emit('update:data', undefined);
    if (selectedFunction.value && targetFunction.value != functionOptions[selectedFunction.value].formula) {
        selectedFunction.value = "";
    }

})


const maxIter = utils.useSessionStorageField('maxIter', 100);

const updateMaxIter = () => fieldsError.value.maxIter =
    maxIter.value < 0 ? 'Количество итераций должно быть неотрицательным' :
        maxIter.value % 1 != 0 ? 'Количество итераций должно быть целым' :
            maxIter.value > 1000 ? 'Количество итераций должно быть не больше 1000' : '';


const minX = utils.useSessionStorageField('minX', 10);
const maxX = utils.useSessionStorageField('maxX', 50);
const minY = utils.useSessionStorageField('minY', 5);
const maxY = utils.useSessionStorageField('maxY', 50);

watch([minX, maxX, minY, maxY], () => {
    updateSettings();
    if (maxX.value - minX.value > 500) {
        fieldsError.value.xBounds = 'Слишком большой диапазон. Максимально допустимое значение: 500.';
    }
    if (maxY.value - minY.value > 500) {
        fieldsError.value.yBounds = 'Слишком большой диапазон. Максимально допустимое значение: 500.';
    }
});


const populationSize = utils.useSessionStorageField('populationSize', 50);
const showInitialPopulation = utils.useSessionStorageField('showInitialPopulation', false);
const initialPopulation = utils.useSessionStorageField('initialPopulation', []);

const updatePopulationFields = () => {
    const newSize = populationSize.value;
    checkPopulationSize();
    if (!fieldsError.value.populationSize) {
        while (initialPopulation.value.length < newSize) {
            initialPopulation.value.push([0, 0]);
        }
        while (initialPopulation.value.length > newSize) {
            initialPopulation.value.pop();
        }
    }
};

const checkPopulationSize = () => fieldsError.value.populationSize =
    populationSize.value <= 0 ? 'Размер популяции должен быть больше 0' :
        populationSize.value % 1 != 0 ? 'Размер популяции должен быть целым числом' :
            populationSize.value > 1000 ? 'Размер популяции должен быть не больше 1000' : ''

const initPopulation = () => {
    updatePopulationFields();
    showInitialPopulation.value = true;
}

const seed = utils.useSessionStorageField('seed', null);


const handleSubmit = () => {
    if (Object.values(fieldsError.value).every(value => value === ''))
        emit('submit', {
            targetFunction: targetFunction.value,
            maxIter: maxIter.value,
            bounds: [
                [minX.value, maxX.value],
                [minY.value, maxY.value]
            ],
            populationSize: populationSize.value,
            initialPopulation: showInitialPopulation.value ? initialPopulation.value : undefined,
            seed: seed.value
        });
};

onBeforeMount(() => { updateSettings(); })

watch(populationSize, (newValue) => {
    emit('update:populationSize', newValue);
});

defineExpose({ fieldsError, populationSize, maxIter, checkPopulationSize });
</script>

<style scoped lang="scss">
@media all and (orientation: landscape) {
    .container {
        grid-template-columns: 1fr 1.5fr;
    }
}

.container {
    height: 100%;
    min-height: min-content;
    width: 100%;
    display: grid;
    gap: 2rem;
    padding: 0 2rem;

    .left {
        display: flex;
        flex-direction: column;
        gap: 1rem;

        form {
            display: flex;
            flex-direction: column;
            // gap: 1rem;
            background-color: var(--color-background-soft);
            border-radius: var(--button-radius);

            ::-webkit-scrollbar {
                width: 8px;
            }

            ::-webkit-scrollbar-track {
                background: var(--color-background-soft);
                border-radius: 10px;
            }

            ::-webkit-scrollbar-thumb {
                background: var(--color-scrollbar);
                border-radius: 10px;
            }

            ::-webkit-scrollbar-thumb:hover {
                background: #555;
            }

            .algo-name {
                text-align: center;
                font-size: 2rem;
                color: var(--color-heading);
                font-weight: bolder;
                box-shadow: 0 7px 10px rgba(0, 0, 0, 0.1);
                z-index: 5;
            }

            .form-scroll {
                height: 400px;
                max-height: 400px;
                overflow-y: scroll;
                padding: var(--card-padding);
                border-radius: var(--field-border);

                .field-label {
                    font-weight: bold;
                    margin-bottom: 0.5rem;
                }

                .parameter-group {
                    display: flex;
                    flex-wrap: wrap;
                    gap: 2rem;
                    row-gap: 1rem;
                }

                .population {
                    max-height: 310px;
                    overflow: auto;
                    border: thick double var(--color-card-border);
                    padding: 4px var(--card-large-padding) var(--card-large-padding);
                    background-color: var(--color-card-bg);
                    border-radius: 8px;
                }

                .toggle-btn {
                    margin-top: auto;
                    margin-bottom: 16px;
                    color: var(--color-button-text);
                }

                .init-btn {
                    background: var(--color-success);
                }

                .init-btn:hover {
                    background: var(--color-success-bright);
                }

                .remove-btn {
                    background: var(--color-error);
                }

                .remove-btn:hover {
                    background: var(--color-error-bright);
                }
            }

            .submit-btn {
                background: var(--color-button);
                color: var(--color-button-text);
                width: 100%;
            }

            .submit-btn:hover {
                background-color: var(--color-button-hover);
            }


            .end-btn {
                background: var(--color-error);
                color: var(--color-button-text);
                width: 100%;
            }

            .end-btn:hover {
                background: var(--color-error-bright);
            }
        }

    }

    .right {
        position: relative;

        .graph-description {
            background-color: var(--color-card-bg-opacity);
            border-top-right-radius: var(--cadr-large-radius);
            border-right: 1px solid var(--color-card-border-opacity);
            border-top: 1px solid var(--color-card-border-opacity);
            border-left: 1px solid var(--color-card-border);
            border-right: 1px solid var(--color-card-border);
            padding: var(--card-large-padding);
            font-size: 1rem;
            position: absolute;
            bottom: 0;
            left: 0;


            .segments-slider-container {
                display: flex;
                flex-wrap: wrap;
                align-items: center;
                gap: 1rem;

                #segments-slider {
                    width: 200px;
                    cursor: pointer;
                }
            }
        }
    }

}

@media (max-width: 850px) {
    .container {
        grid-template-columns: 1fr;
    }
}
</style>