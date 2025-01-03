<template>
  <div class="relative overflow-x-auto">
    <table 
      v-if="emails.length > 0" 
      class="w-full text-sm text-left text-gray-500 dark:text-gray-400"
    >
      <thead class="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
        <tr>
          <th class="px-6 py-3">Subject</th>
          <th class="px-6 py-3">From</th>
          <!-- <th class="px-6 py-3">Message Id</th> -->
          <th class="px-6 py-3">To</th>
          <th class="px-6 py-3">Date</th>
        </tr>
      </thead>
      <tbody>
        <tr 
          v-for="email in emails" 
          :key="email._id" 
          @click="$emit('select', email)"
          :class="[
            'bg-white border-b dark:bg-gray-800 dark:border-gray-700 cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-600',
            selectedEmail && selectedEmail._id === email._id ? 'bg-blue-100 dark:bg-blue-900' : ''
          ]"
        >
          <td class="px-6 py-4 h-16 font-medium text-gray-900 whitespace-nowrap dark:text-white truncate">{{ email.subject }}</td>
          <td class="px-6 py-4 h-16 truncate">{{ email.from }}</td>
          <!-- <td class="px-6 py-4 h-16 truncate">{{ email.message_id}}</td> -->
          <td class="px-6 py-4 h-16 truncate">{{ email.to }}</td>
          <td class="px-6 py-4 h-16 truncate">{{ email.date }}</td>
        </tr>
      </tbody>
    </table>
    <p v-else class="text-center text-gray-500 mt-4">No results found. Please search for emails.</p>
  </div>
</template>


<script setup>
defineProps({
  emails: {
    type: Array,
    default: () => [], // Valor predeterminado vacío para evitar errores
  },
  selectedEmail: {
    type: Object,
    default: () => null,
  },
});

defineEmits(['select']);
</script>

<style scoped>
/* Estilo adicional para resaltar la fila seleccionada */
.selected {
  background-color: #d1e7ff;
}

/* Añadir reglas para truncar el texto con puntos suspensivos */
.truncate {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
