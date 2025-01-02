<template>
  <div class="container">
    <div class="left-panel">
      <SearchBar 
        v-model="searchQuery" 
        @search="fetchEmails" 
      />

      <table>
        <thead>
          <tr>
            <th>Message ID</th>
            <th>Subject</th>
            <th>From</th>
            <th>To</th>
            <th>Date</th>
          </tr>
        </thead>
        <tbody>
          <tr 
            v-for="email in emails" 
            :key="email._id" 
            @click="selectEmail(email)"
            :class="{ selected: selectedEmail && selectedEmail._id === email._id }"
          >
            <td>{{ email.message_id }}</td>
            <td>{{ email.subject }}</td>
            <td>{{ email.from }}</td>
            <td>{{ email.to }}</td>
            <td>{{ email.date }}</td>
          </tr>
        </tbody>
      </table>

      <div class="pagination">
        <button 
          @click="prevPage" 
          :disabled="page <= 0"
        >
          Anterior
        </button>

        <span>Página {{ page + 1 }} de {{ totalPages }}</span>

        <button 
          @click="nextPage" 
          :disabled="(page.value + 1) * pageSize >= totalEmails"
        >
          Siguiente
        </button>
      </div>
    </div>

    <!-- Modal de Detalle del Email -->
    <div v-if="isModalVisible" class="modal-overlay">
      <div class="modal">
        <!-- Botón para cerrar el modal -->
        <button class="close-button" @click="closeModal">X</button>
        <h3>{{ selectedEmail.subject }}</h3>
        <p><strong>From:</strong> {{ selectedEmail.from }}</p>
        <p><strong>To:</strong> {{ selectedEmail.to }}</p>
        <p><strong>Date:</strong> {{ selectedEmail.date }}</p>
        <hr />
        <div class="body-text">{{ selectedEmail.body }}</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import SearchBar from './SearchBar.vue';

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
  isModalVisible.value = true; // Mostrar el modal al seleccionar un correo
};

const closeModal = () => {
  isModalVisible.value = false; // Ocultar el modal
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
.container {
  display: flex;
  height: 100vh;
}

.left-panel {
  width: 60%;
  padding: 20px;
}

.right-panel {
  width: 40%;
  padding: 20px;
  border-left: 1px solid #ddd;
  overflow-y: auto;
}

.search-box {
  width: 100%;
  padding: 10px;
  margin-bottom: 20px;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th, td {
  border: 1px solid #ddd;
  padding: 10px;
  text-align: left;
}

th {
  background-color: #f4f4f4;
}

tr {
  cursor: pointer;
}

tr:hover {
  background-color: #f0f0f0;
}

.selected {
  background-color: #d1e7ff;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

button {
  margin: 0 10px;
  padding: 10px 20px;
  cursor: pointer;
}

button:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 999;
}

.modal {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  width: 60%;
  max-width: 800px; /* Establecer un tamaño máximo para el modal */
  position: relative;
  max-height: 80vh; /* Evitar que el modal se haga demasiado grande */
  overflow-y: auto; /* Agregar scroll interno cuando el contenido sea largo */
}

.close-button {
  position: absolute;
  top: 10px;
  right: 10px;
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #999;
}

.close-button:hover {
  color: #333;
}

.body-text {
  white-space: pre-wrap;
  font-family: monospace;
  word-wrap: break-word;
}
</style>
