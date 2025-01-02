<template>
  <div v-if="isVisible" class="fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center z-50">
    <div class="relative p-4 w-full max-w-2xl max-h-full">
      <!-- Modal content -->
      <div class="relative bg-white rounded-lg shadow dark:bg-gray-700">
        <!-- Modal header -->
        <div class="flex items-center justify-between p-4 md:p-5 border-b rounded-t dark:border-gray-600">
          <h3 class="text-xl font-semibold text-gray-900 dark:text-white">
            {{ email.subject }}
          </h3>
          <!-- Botón de cerrar -->
          <button @click="closeModal" class="text-gray-400 hover:text-gray-900 text-xl" aria-label="Close modal">
            ×
          </button>
        </div>

        <!-- Modal body -->
        <div class="p-4 md:p-5 space-y-4 max-h-[70vh] overflow-y-auto">
          <p><strong>From:</strong> {{ email.from }}</p>
          <p><strong>To:</strong> {{ email.to }}</p>
          <p><strong>Date:</strong> {{ email.date }}</p>
          <hr class="my-4" />
          <div class="whitespace-pre-wrap font-mono">{{ email.body }}</div>
        </div>

        <!-- Modal footer -->
        <div class="flex items-center p-4 md:p-5 border-t border-gray-200 rounded-b dark:border-gray-600">
          <button @click="closeModal" class="px-5 py-2.5 bg-blue-700 text-white rounded-lg hover:bg-blue-800">
            Close
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { defineProps, defineEmits } from 'vue';

const props = defineProps({
  email: Object, // El objeto de correo con detalles
  isVisible: Boolean, // Para controlar la visibilidad del modal
});

const emit = defineEmits(['close']); // Evento para cerrar el modal

const closeModal = () => {
  emit('close'); // Emitir un evento de cierre
};
</script>

<style scoped>
/* Aquí puedes incluir estilos específicos para el modal */
</style>
