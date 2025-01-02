<template>
  <input 
    v-model="modelValue" 
    placeholder="Buscar por subject, from o to" 
    class="search-box"
    @input="handleInput" 
  />
</template>

<script setup>
import { ref } from 'vue';

const modelValue = ref(''); // Almacena el valor del campo de búsqueda
const emit = defineEmits(['update:modelValue', 'search']); // Emite los eventos necesarios

let timeoutId = null; // Para almacenar el ID del timeout

const handleInput = (event) => {
  // Emite el cambio de valor para sincronizar el campo de búsqueda
  emit('update:modelValue', event.target.value);

  // Limpiar cualquier timeout previo
  clearTimeout(timeoutId);

  // Crear un nuevo timeout para que la búsqueda ocurra después de un pequeño retraso
  timeoutId = setTimeout(() => {
    emit('search'); // Emitimos el evento de búsqueda
  }, 500); // 500 ms de retraso
};
</script>

<style scoped>
.search-box {
  width: 100%;
  padding: 10px;
  margin-bottom: 20px;
}
</style>
