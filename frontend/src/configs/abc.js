import { ref } from 'vue'
import CustomField from 'components/CustomField.vue';
import utils from 'utils/storage';


const foragerSize = utils.useSessionStorageField('foragerSize', 25);
const observerSize = utils.useSessionStorageField('observerSize', 0);
const limit = utils.useSessionStorageField('limit', 100);

const fieldsErrors = ref({
    foragerSize: '',
    observerSize: '',
    limit: '',
});

const abc = {
    title: 'Алгоритм искусственной пчелиной колонии',
    socketUrl: '/api/ABC',
    state: {
        foragerSize,
        observerSize,
        limit,
    },
    fieldsErrors,
    fields: {
        foragerSize: {
            model: 'foragerSize',
            component: CustomField,
            props: {
                label: 'Размер популяции собирателей',
                type: 'number',
                step: 1,
                description: 'Пчелы-собиратели ищут новые решения.'
            },
            validate: (form) => {
                const val = abc.state.foragerSize.value;
                const total = form.populationSize;
                fieldsErrors.value.foragerSize =
                    val <= 0 ? 'Должен быть больше 0' : val >= total ? 'Должен быть меньше размера популяции' : '';
            },
            update: (form) => {
                abc.state.observerSize.value = form.populationSize - abc.state.foragerSize.value;
                abc.fields.observerSize.validate(form);
            }
        },
        observerSize: {
            model: 'observerSize',
            component: CustomField,
            props: {
                label: 'Размер популяции наблюдателей',
                type: 'number',
                step: 1,
                description: 'Пчелы-наблюдатели выбирают лучшие решения.'
            },
            validate: (form) => {
                const val = abc.state.observerSize.value;
                const total = form.populationSize;
                fieldsErrors.value.observerSize =
                    val <= 0 ? 'Размер популяции собирателей должен быть больше 0' :
                        val >= total ? 'Размер популяции собирателей должен быть меньше размера популяции' : '';
            },
            update: (form) => {
                abc.state.foragerSize.value = form.populationSize - abc.state.observerSize.value;
                abc.fields.foragerSize.validate(form);
            }
        },
        limit: {
            model: 'limit',
            component: CustomField,
            props: {
                label: 'Лимит стагнации',
                type: 'number',
                step: 1,
                description: 'Максимальное количество неудачных попыток улучшения перед заменой решения.'
            },
            validate: () => {
                const val = abc.state.limit.value;
                fieldsErrors.value.limit = val <= 0 ? 'Лимит стагнации должен быть больше 0' : '';
            }
        }
    },
    init: (form) => {
        abc.state.observerSize.value = form.populationSize - foragerSize.value;
    },
    submit: (formData) => {
        return {
            ...formData,
            foragerSize: abc.state.foragerSize.value,
            observerSize: abc.state.observerSize.value,
            limit: abc.state.limit.value
        }
    },
    onPopulationChange: (form) => {
        // abc.state.foragerSize.value = Math.floor(form.populationSize / 2);

        abc.state.observerSize.value = form.populationSize - abc.state.foragerSize.value;
        abc.fields.foragerSize.validate(form);
        abc.fields.observerSize.validate(form);
    }
}

export default abc;