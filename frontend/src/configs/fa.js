import { ref } from 'vue'
import CustomField from 'components/CustomField.vue'
import utils from 'utils/storage'

const alpha = utils.useSessionStorageField('alpha', 0.01)
const beta0 = utils.useSessionStorageField('beta0', 1)
const gamma = utils.useSessionStorageField('gamma', 0.8)

const fieldsErrors = ref({
    alpha: '',
    beta0: '',
    gamma: ''
})

export default {
    title: 'Светлячковый алгоритм',
    socketUrl: '/api/firefly',

    state: {
        alpha,
        beta0,
        gamma,
    },

    fieldsErrors,

    fields: {
        alpha: {
            component: CustomField,
            model: 'alpha',
            props: {
                label: 'Параметр яркости (α)',
                type: 'number',
                step: 0.01,
                description: 'Определяет скорость уменьшения яркости света с расстоянием между светлячками.',
            },
            validate: () => {
                fieldsErrors.value.alpha = alpha.value < 0 ? 'Параметр яркости должен быть неотрицательным' : ''
            },
        },
        beta0: {
            component: CustomField,
            model: 'beta0',
            props: {
                label: 'Параметр случайности (β0)',
                type: 'number',
                step: 0.01,
                description: 'Контролирует степень случайности в движении светлячков. Высокие значения увеличивают разведку.',
            },
            validate: () => {
                fieldsErrors.value.beta0 = beta0.value < 0 ? 'Параметр случайности должен быть неотрицательным' : ''
            },
        },
        gamma: {
            component: CustomField,
            model: 'gamma',
            props: {
                label: 'Параметр плотности светлячков (𝛾)',
                type: 'number',
                step: 0.01,
                description: 'Определяет поведение светлячков в зависимости от плотности других светлячков.',
            },
            validate: () => {
                fieldsErrors.value.gamma = gamma.value < 0 ? 'Параметр плотности светлячков должен быть неотрицательным' : ''
            },
        }
    },

    submit(formData) {
        return {
            ...formData,
            alpha: alpha.value,
            beta0: beta0.value,
            gamma: gamma.value,
        }
    },

}
