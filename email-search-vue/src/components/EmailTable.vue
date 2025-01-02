<template>
  <div class="flex w-full h-screen">
    <!-- Panel izquierdo -->
    <div class="w-full p-8 overflow-y-auto">
      <SearchBar 
        v-model="searchQuery" 
        @search="fetchEmails" 
        class="mb-6"
      />

      <!-- Tabla -->
      <table class="min-w-full table-auto border-collapse">
        <thead>
          <tr>
            <th class="px-4 py-2 border-b text-left bg-gray-100">Message ID</th>
            <th class="px-4 py-2 border-b text-left bg-gray-100">Subject</th>
            <th class="px-4 py-2 border-b text-left bg-gray-100">From</th>
            <th class="px-4 py-2 border-b text-left bg-gray-100">To</th>
            <th class="px-4 py-2 border-b text-left bg-gray-100">Date</th>
          </tr>
        </thead>
        <tbody>
          <tr 
            v-for="email in emails" 
            :key="email._id" 
            @click="selectEmail(email)"
            :class="{ 'bg-blue-100': selectedEmail && selectedEmail._id === email._id, 'hover:bg-gray-200 cursor-pointer': true }"
          >
            <td class="px-4 py-2 border-b">{{ email.message_id }}</td>
            <td class="px-4 py-2 border-b">{{ email.subject }}</td>
            <td class="px-4 py-2 border-b">{{ email.from }}</td>
            <td class="px-4 py-2 border-b">{{ email.to }}</td>
            <td class="px-4 py-2 border-b">{{ email.date }}</td>
          </tr>
        </tbody>
      </table>

      <!-- Paginación -->
      <div class="flex justify-center items-center mt-4 space-x-4">
        <button 
          @click="prevPage" 
          :disabled="page <= 0"
          class="px-4 py-2 bg-gray-300 rounded disabled:opacity-50"
        >
          Anterior
        </button>

        <span class="text-lg">Página {{ page + 1 }} de {{ totalPages }}</span>

        <button 
          @click="nextPage" 
          :disabled="(page.value + 1) * pageSize >= totalEmails"
          class="px-4 py-2 bg-gray-300 rounded disabled:opacity-50"
        >
          Siguiente
        </button>
      </div>
    </div>

    <!-- Modal de detalles de email -->
    <EmailDetail 
      :email="selectedEmail"
      :isVisible="isModalVisible"
      @close="closeModal"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import SearchBar from './SearchBar.vue';
import EmailDetail from './EmailDetail.vue'; // Importar el componente de detalle de correo

const emails = ref([]);
const selectedEmail = ref(null);
const searchQuery = ref('');
const page = ref(0);  
const pageSize = 10;  
const totalEmails = ref(0);
const totalPages = ref(1);
const isModalVisible = ref(false); // Estado para mostrar el modal

const fetchEmails = async () => {
  try {
    const response = await axios.post(
      'http://localhost:8080/api/search',
      {
        term: searchQuery.value || "*", // Si no hay búsqueda, buscamos todos
        from: page.value * pageSize,
        size: pageSize
      }
    );

    emails.value = response.data.hits.hits.map(hit => hit._source);
    totalEmails.value = response.data.hits.total ? response.data.hits.total.value : 0;
    totalPages.value = Math.ceil(totalEmails.value / pageSize);
  } catch (error) {
    console.error('Error fetching emails:', error);
  }
};

const selectEmail = (email) => {
  selectedEmail.value = email;
  isModalVisible.value = true; // Mostrar el modal cuando se seleccione un correo
};

const closeModal = () => {
  isModalVisible.value = false; // Cerrar el modal
};

const prevPage = () => {
  if (page.value > 0) {
    page.value--;
    fetchEmails();
  }
};

const nextPage = () => {
  if ((page.value + 1) * pageSize < totalEmails.value) {
    page.value++;
    fetchEmails();
  }
};

onMounted(() => fetchEmails());
</script>

<style scoped>
/* Estilos adicionales si es necesario */
</style>
