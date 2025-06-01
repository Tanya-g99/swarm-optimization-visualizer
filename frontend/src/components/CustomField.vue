<template>
    <Field :label="label" :description="description" :error="error">
        <template #default="{ id, focus, blur }">
            <div class="custom-field-block">
                <input class="field-input left" :id="id" :type="type" :step="step" :value="modelValue" v-bind="$attrs"
                    @focus="focus" @blur="blur"
                    @input="$emit('update:modelValue', $event.target.value === '' ? null : $event.target.value)"
                    :class="{ 'error-border': error }" />

                <slot class="right"></slot>
            </div>
        </template>
    </Field>
</template>

<script setup>
import Field from './Field.vue';

const props = defineProps({
    label: { type: String, required: true },
    description: { type: String, default: '' },
    error: { type: String, default: '' },
    type: { type: String, default: 'text' },
    step: { type: Number, default: 1 },
    modelValue: { type: [String, Number], default: '' },

});


const emit = defineEmits(['update:modelValue']);
</script>


<style scoped lang="scss">
@media (max-width: 850px) {
    .custom-field-block {
        width: 100%;

        .left {
            width: 100%;
        }
    }
}

@media (min-width: 850px) {
    .custom-field-block {
        display: flex;

        .left {
            flex-grow: 1;
        }

        .right {
            flex-shrink: 0;
        }

    }

    .custom-field-block> :not(:first-child):not(:last-child) {
        border-radius: 0;
    }

    .custom-field-block> :first-child:not(:only-child) {
        border-top-right-radius: 0;
        border-bottom-right-radius: 0;
    }

    .custom-field-block> :last-child:not(:only-child) {
        border-top-left-radius: 0;
        border-bottom-left-radius: 0;
    }
}
</style>