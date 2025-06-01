import CustomField from 'components/CustomField.vue'
import utils from 'utils/storage'
import { ref } from 'vue'

const subpopulationsCount = utils.useSessionStorageField('subpopulationsCount', 10)
const iMax = utils.useSessionStorageField('iMax', 10)
const populationSize = utils.useSessionStorageField('populationSize', 50)

const fieldsErrors = ref({
    subpopulationsCount: '',
    iMax: '',
    populationSize: ''
})

export default {
    title: 'Алгоритм перемешанных лягушачьих прыжков',
    socketUrl: '/api/SFLA',

    state: {
        subpopulationsCount,
        iMax,
        populationSize
    },

    fieldsErrors,

    fields: {
        subpopulationsCount: {
            component: CustomField,
            model: 'subpopulationsCount',
            props: {
                label: 'Количество субпопуляций',
                type: 'number',
                step: 1,
                description: 'Количество субпопуляций, на которые будет разделена общая популяция. Субпопуляции — это группы особей, которые работают независимо друг от друга, с целью локального поиска и нахождения оптимальных решений в пределах своей подгруппы.'
            },
            validate: () => {
                if (subpopulationsCount.value <= 0) {
                    fieldsErrors.value.subpopulationsCount = 'Количество субпопуляций должно быть больше 0'
                } else if (subpopulationsCount.value > populationSize.value) {
                    fieldsErrors.value.subpopulationsCount = 'Количество субпопуляций должно быть не больше размера популяции'
                } else {
                    fieldsErrors.value.subpopulationsCount = ''
                }
            },
            update: (form) => {
                if (populationSize.value < subpopulationsCount.value)
                    form.fieldsError.populationSize = 'Размер популяции должен быть не меньше количества субпопуляций';
                else form.checkPopulationSize();
            }
        },
        iMax: {
            component: CustomField,
            model: 'iMax',
            props: {
                label: 'Количество итераций локального поиска',
                type: 'number',
                step: 1,
                description: 'Максимальное количество итераций локального поиска внутри каждой субпопуляции.'
            },
            validate: () => {
                fieldsErrors.value.iMax =
                    iMax.value <= 0 ? 'Количество итераций локального поиска должно быть больше 0' : ''
            }
        }
    },

    submit(formData) {
        return {
            ...formData,
            subpopulationsCount: subpopulationsCount.value,
            iMax: iMax.value
        }
    },

    onPopulationChange(form) {
        this.fields.subpopulationsCount.validate(form)
        this.fields.subpopulationsCount.update(form)
    }
}
