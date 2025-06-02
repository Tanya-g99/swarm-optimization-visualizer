<template>
    <Form :algoName="config.title" ref="form" v-model:data="data" @submit="submitForm" @stop="stopAlgo"
        :working="working" @update:populationSize="handlePopulationSizeUpdate(form.value)">
        <template #additional-fields>
            <template v-for="([key, field]) in fieldsEntries" :key="key">
                <!-- Простое поле -->
                <component v-if="!Array.isArray(field.model)" :is="field.component"
                    v-model.number="config.state[field.model].value" v-bind="field.props" :error="fieldsErrors[key]"
                    @update:model-value="onFieldUpdate(key)" />

                <component v-else :is="field.component" v-bind="field.props"
                    v-model:left.number="config.state[field.model[0]].value"
                    v-model:right.number="config.state[field.model[1]].value" v-model:error="fieldsErrors[key]"
                    @update:left="onFieldUpdate(key)" @update:right="onFieldUpdate(key)" />
            </template>
            {{ fieldsErrors['visual'] }}
        </template>
    </Form>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import Form from 'components/Form.vue'
import { useRoute } from 'vue-router';

import configs from 'configs/algorithms'

const route = useRoute();
const algoName = computed(() => route.params.algoName);

let config = configs[algoName.value];

let fieldsEntries = Object.entries(config.fields);
let fieldsErrors = config.fieldsErrors;

let socket;

watch(algoName, (newAlgo) => {
    config = configs[newAlgo];
    fieldsEntries = Object.entries(config.fields);
    fieldsErrors = config.fieldsErrors;

    data.value = null;
})

const data = ref(null);
const form = ref(null);

const validateField = (key) => {
    const field = config.fields[key];
    if (field && typeof field.validate === 'function') {
        field.validate(form.value);
    }
}

const onFieldUpdate = (key) => {
    const field = config.fields[key];
    if (field && typeof field.update === 'function') {
        field.update(form.value);
    }
    validateField(key);
}

const handlePopulationSizeUpdate = (newVal) => {
    if (typeof config.onPopulationChange === 'function') {
        config.onPopulationChange(form.value);
    }
    if (config.fields.foragerSize) validateField('foragerSize');
    if (config.fields.observerSize) validateField('observerSize');
}


const working = ref(false)

const submitForm = (formData) => {
    if (!Object.values(fieldsErrors.value).every(err => err === "")) {
        return;
    }
    socket = new WebSocket(config.socketUrl);
    socket.onopen = () => {
        working.value = true;
        socket.send(JSON.stringify(config.submit(formData)));
    }
    socket.onmessage = (event) => {
        data.value = JSON.parse(event.data);
    }

    socket.onclose = () => {
        working.value = false;
    }
    socket.onerror = (err) => {
        console.error('WebSocket error:', err);
        working.value = false;
    }
}

const stopAlgo = () => {
    if (socket) {
        socket.close();
        working.value = false;
        socket = null;
    }
}

onMounted(() => {
    if (typeof config.init === 'function') {
        config.init(form.value);
    }
    for (const key in config.fields) {
        validateField(key)
    }
})
</script>

<style scoped lang="scss"></style>