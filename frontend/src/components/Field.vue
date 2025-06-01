<template>
    <div ref="fieldRef" class="field-container">
        <label :for="id" class="field-label">
            {{ label }}
        </label>
        <slot :id="id" :focus="onFocus" :blur="onBlur"></slot>

        <Teleport to="body">
            <p v-if="showDescription && description && !error" class="field-message field-description" :style="descriptionStyle">
                {{ description }}
            </p>
        </Teleport>

        <Teleport to="body">
            <p v-if="error" class="error-text field-message field-error error-border px-8" @focus="updatePosition"
                :style="descriptionStyle">
                {{ " ! " + error }}
            </p>
        </Teleport>
    </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';

const props = defineProps({
    label: { type: String, required: true },
    description: { type: String, default: '' },
    error: { type: String, default: '' }
});

const showDescription = ref(false);
const descriptionStyle = ref({});
const fieldRef = ref(null);
const scrollParent = ref(null);
// Определяем ближайший скроллимый родительский элемент
const getScrollParent = (element) => {
    while (element) {
        const overflow = window.getComputedStyle(element).overflowY;
        if (overflow === 'auto' || overflow === 'scroll') {
            return element;
        }
        element = element.parentElement;
    }
    return window; // Если нет скроллимого родителя, используем `window`
};

const onFocus = () => {
    showDescription.value = true;
};

const updatePosition = () => {
    if (fieldRef.value) {
        const rect = fieldRef.value.getBoundingClientRect();
        descriptionStyle.value = {
            top: `${rect.bottom + window.scrollY}px`,
            left: `${rect.left}px`,
        };
    }
};

const onBlur = () => {
    showDescription.value = false;
};

const id = `field-${Math.random().toString(36).substr(2, 9)}`;

onMounted(() => {
    scrollParent.value = getScrollParent(fieldRef.value);
    scrollParent.value.addEventListener("scroll", updatePosition);
    window.addEventListener("scroll", updatePosition);
    updatePosition();
});

onUnmounted(() => {
    if (scrollParent.value) {
        scrollParent.value.removeEventListener("scroll", updatePosition);
    }
    window.removeEventListener("scroll", updatePosition);
});
</script>

<style scoped lang="scss">
.field-container {
    display: flex;
    flex-direction: column;
    position: relative;
    margin-bottom: 1rem;


    .field-label {
        font-weight: bold;
        display: flex;
    }

    .field-error {
        border-top-left-radius: 0px;
    }

    .field-description {
        color: var(--color-card-text);
    }
}
</style>
