<template>
  <div class="search-container flex flex-col sm:flex-row gap-4 max-w-lg mx-auto">
    <div class="relative w-full sm:w-1/3">
      <select 
        v-model="selectedField" 
        @change="emitSearch" 
        class="block py-2.5 px-0 w-full text-sm text-gray-500 bg-transparent border-0 border-b-2 border-gray-200 appearance-none focus:outline-none focus:ring-0 focus:border-gray-500">
        <option value="">Todos los campos</option>
        <option value="body">Body</option>
        <option value="from">From</option>
        <option value="to">To</option>
        <option value="subject">Subject</option>
        <option value="folder">Folder</option>
      </select>
    </div>
    
    <input 
      v-model="modelValue" 
      placeholder="Search..." 
      class="search-box w-full py-2 px-3 border-b-2 border-gray-200 focus:outline-none focus:border-gray-500 text-sm text-gray-700" 
      @input="handleInput" 
    />
  </div>
</template>

<script setup>
import { ref } from 'vue';

const modelValue = ref('');
const selectedField = ref(''); // Filtro por defecto desactivado
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
  const fieldToEmit = selectedField.value || '_all'; // Si no hay filtro, busca en todos los campos
  emit('update:field', fieldToEmit);
  emit('search');
};
</script>