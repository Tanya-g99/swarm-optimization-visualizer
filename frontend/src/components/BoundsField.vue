<template>
    <Field :label="label" :description="description" :error="error">
        <template v-slot="{ id, focus, blur }">
            <div class="bounds-input">
                (<input class="field-input" :id="id" type="number" :step="step" v-model.number="left" @focus="focus"
                    @blur="blur" @input="$emit('update:left', $event.target.value)" />
                <span>;</span>
                <input class="field-input" :id="id + '-2'" type="number" :step="step" v-model.number="right"
                    @focus="focus" @blur="blur" @input="$emit('update:right', $event.target.value)" />)
            </div>
        </template>
    </Field>
</template>

<script setup>
import { onMounted, ref, watch } from 'vue';
import Field from './Field.vue';

const props = defineProps({
    label: { type: String, required: true },
    description: { type: String, default: '' },
    error: { type: String, default: '' },
    left: Number,
    right: Number,
    step: {
        type: Number,
        default: 0.1
    },
    validateRange: {
        type: Boolean,
        default: false
    },
});

const emit = defineEmits(['update:left', 'update:right', 'update:error']);

const left = ref(props.left);
const right = ref(props.right);

const checkRange = () => {
    if (props.validateRange) {
        emit('update:error',
            left.value > right.value ? 'Левая граница не может быть больше правой.' :
                left.value === right.value ? 'Значения границ должны быть разные.' :
                    '')
    }

    if (!props.error) {
        emit('update:left', left.value);
        emit('update:right', right.value);
    }
}

watch([left, right], () => {
    checkRange();
});

watch(() => props.left, (newVal) => left.value = newVal);
watch(() => props.right, (newVal) => right.value = newVal);

onMounted(checkRange);
</script>

<style scoped lang="scss">
.bounds-input {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 1rem;

    input {
        width: 70px;
        text-align: center;
    }
}
</style>