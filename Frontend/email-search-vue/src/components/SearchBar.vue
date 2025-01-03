<template>
  <div class="search-container">
    <select v-model="selectedField" class="field-selector" @change="emitSearch">
      <option value="body">Body</option>
      <option value="from">From</option>
      <option value="to">To</option>
      <option value="subject">Subject</option>
    </select>
    <input 
      v-model="modelValue" 
      placeholder="Search..." 
      class="search-box"
      @input="handleInput" 
    />
  </div>
</template>

<script setup>
import { ref } from 'vue';

const modelValue = ref('');
const selectedField = ref('body'); // Campo predeterminado para buscar
const emit = defineEmits(['update:modelValue', 'update:field', 'search']);

let timeoutId = null;

const handleInput = (event) => {
  emit('update:modelValue', event.target.value);
  clearTimeout(timeoutId);
  timeoutId = setTimeout(() => {
    emitSearch();
  }, 500);
};

const emitSearch = () => {
  emit('update:field', selectedField.value); // Emite el campo seleccionado
  emit('search');
};
</script>

<style scoped>
.search-container {
  display: flex;
  gap: 10px;
}

.field-selector {
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 5px;
}

.search-box {
  flex: 1;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 5px;
}
</style>
