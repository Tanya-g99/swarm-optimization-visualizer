import CustomField from 'components/CustomField.vue'
import utils from 'utils/storage'
import { ref } from 'vue'

const a = utils.useSessionStorageField('a', 2.0)
const C = utils.useSessionStorageField('C', 2.0)

const fieldsErrors = ref({
    a: '',
    C: ''
})

export default {
    title: 'Алгоритм серого волка (GWO)',
    socketUrl: '/api/GWO',

    state: {
        a,
        C
    },

    fieldsErrors,

    fields: {
        a: {
            component: CustomField,
            model: 'a',
            props: {
                label: 'Параметр a',
                type: 'number',
                step: 0.01,
                description: 'Параметр, управляющий изменением положения волков на каждой итерации.',
            },
            validate: () => {
                fieldsErrors.value.a = a.value <= 0 ? 'Параметр a должен быть положительным' : ''
            }
        },
        C: {
            component: CustomField,
            model: 'C',
            props: {
                label: 'Параметр C',
                type: 'number',
                step: 0.01,
                description: 'Параметр, управляющий влиянием лучшего решения на движение волков.',
            },
            validate: () => {
                fieldsErrors.value.C = C.value <= 0 ? 'Параметр C должен быть положительным' : ''
            }
        }
    },

    submit(formData) {
        return {
            ...formData,
            a: a.value,
            C: C.value
        }
    },
}
