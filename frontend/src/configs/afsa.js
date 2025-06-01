import { ref } from 'vue';
import CustomField from 'components/CustomField.vue';
import BoundsField from 'components/BoundsField.vue';
import utils from 'utils/storage';

const maxTryNum = utils.useSessionStorageField('maxTryNum', 5);
const eta = utils.useSessionStorageField('eta', 1e-4);
const visualMin = utils.useSessionStorageField('visualMin', 1);
const visualMax = utils.useSessionStorageField('visualMax', 8);
const teta = utils.useSessionStorageField('teta', 1);

const fieldsErrors = ref({
    maxTryNum: '',
    eta: '',
    visual: '',
    teta: ''
});

const afsa = {
    title: 'Алгоритм искусственного косяка рыб',
    socketUrl: '/api/AFSA',

    state: {
        maxTryNum,
        eta,
        visualMin,
        visualMax,
        teta,
    },

    fieldsErrors,

    fields: {
        maxTryNum: {
            model: 'maxTryNum',
            component: CustomField,
            props: {
                label: 'Параметр стагнации',
                type: 'number',
                step: 0.01,
                description: 'Количество итераций без улучшения, после которых начнется случайное движение.',
            },
            validate: (form) => {
                fieldsErrors.value.maxTryNum =
                    afsa.state.maxTryNum.value < 0 ? 'Параметр стагнации должен быть неотрицательным' :
                        afsa.state.maxTryNum.value >= form.maxIter ? 'Параметр стагнации должен быть меньше количества итераций' : '';
            },
        },
        eta: {
            model: 'eta',
            component: CustomField,
            props: {
                label: 'Контроль стагнации',
                type: 'number',
                step: 0.0001,
                description: 'Минимальное изменение лучшего значения функции, после которого считается, что алгоритм застопорился.',
            },
            validate: () => {
                fieldsErrors.value.eta =
                    afsa.state.eta.value < 0 ? 'Контроль стагнации должен быть неотрицательным' : '';
            },
        },
        visual: {
            model: ['visualMin', 'visualMax'],
            component: BoundsField,
            props: {
                label: 'Визуальный диапазон',
                validateRange: true,
                description: 'Определяет насколько далеко особь может видеть соседей. Этот параметр часто уменьшается в процессе оптимизации',
            },
            validate: () => {
                if (afsa.state.visualMin.value < 0)
                    fieldsErrors.value.visual = 'Границы визуального диапазона должны быть неотрицательными';
            },
        },
        teta: {
            model: 'teta',
            component: CustomField,
            props: {
                label: 'Параметр плотности косяка',
                type: 'number',
                step: 0.01,
                description: 'Определяет поведение рыб в зависимости от плотности соседей.',
            },
            validate: () => {
                fieldsErrors.value.teta =
                    afsa.state.teta.value < 0 ? 'Параметр плотности косяка должен быть неотрицательным' : '';
            },
        },
    },
    submit: (formData) => {
        return {
            ...formData,
            maxTryNum: afsa.state.maxTryNum.value,
            eta: afsa.state.eta.value,
            visual: [afsa.state.visualMin.value, afsa.state.visualMax.value],
            teta: afsa.state.teta.value,
        };
    },
};

export default afsa;
