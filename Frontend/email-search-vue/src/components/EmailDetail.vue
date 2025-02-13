<template>
  <div
    v-if="isVisible"
    class="fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center z-50"
  >
    <div class="relative p-4 w-full max-w-2xl">
      <!-- Modal content -->
      <div
        class="relative bg-white rounded-lg shadow dark:bg-gray-700 flex flex-col"
        style="max-height: 70vh;"
      >
        <!-- Modal header -->
        <div class="flex items-center justify-between p-4 md:p-5 border-b rounded-t dark:border-gray-600">
          <h3 class="text-xl font-semibold text-gray-900 dark:text-white">
            {{ email.subject }}
          </h3>
          <button
            @click="closeModal"
            class="text-gray-400 hover:text-gray-900 text-xl"
            aria-label="Close modal"
          >
            ×
          </button>
        </div>

        <!-- Fixed email details with new initials icon -->
        <div class="flex items-center space-x-3 p-4 md:p-5 border-b dark:border-gray-600">
          <!-- New Initials Icon -->
          <div class="relative inline-flex items-center justify-center w-10 h-10 overflow-hidden bg-gray-100 rounded-full dark:bg-gray-600">
            <span class="font-medium text-gray-600 dark:text-gray-300">{{ initials }}</span>
          </div>
          <!-- Email details -->
          <div>
            <p><strong>From:</strong> {{ email.from }}</p>
            <p><strong>To:</strong> {{ email.to }}</p>
            <p><strong>Date:</strong> {{ email.date }}</p>
          </div>
        </div>

        <!-- Scrollable email body -->
        <div class="p-4 md:p-5 flex-1 overflow-y-auto space-y-4">
          <div class="whitespace-pre-wrap font-mono">{{ email.body }}</div>
        </div>

        <!-- Modal footer -->
        <div class="flex items-center p-4 md:p-5 border-t border-gray-200 rounded-b dark:border-gray-600">
          <button
            @click="closeModal"
            class="px-5 py-2.5 bg-blue-700 text-white rounded-lg hover:bg-blue-800"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { defineProps, defineEmits, computed } from 'vue';

const props = defineProps({
  email: Object, // El objeto de correo con detalles
  isVisible: Boolean, // Para controlar la visibilidad del modal
});

const emit = defineEmits(['close']); // Evento para cerrar el modal

const closeModal = () => {
  emit('close'); // Emitir un evento de cierre
};

const initials = computed(() => {
  const from = props.email?.from || "";
  // Obtener parte local del correo (antes del @)
  const localPart = from.split('@')[0];
  // Separamos por punto
  const parts = localPart.split('.');
  if (parts.length >= 2) {
    // Tomar la primera letra del primer segmento y la primera letra del segundo segmento
    return (parts[0][0] + parts[1][0]).toUpperCase();
  }
  // Si no hay punto, retornar solo la primera letra
  return localPart[0]?.toUpperCase() || "";
});
</script>