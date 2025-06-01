import { ref, watch } from 'vue';

function useSessionStorageField(key, defaultValue) {
    const saved = sessionStorage.getItem(key);
    const data = ref(saved !== null && saved !== "undefined" ? JSON.parse(saved) : defaultValue);

    watch(data, (newValue) => {
        sessionStorage.setItem(key, JSON.stringify(newValue));
    }, { deep: true });

    return data;
}

export default {
    useSessionStorageField,
};