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
          :disabled="(page + 1) * pageSize >= totalEmails"
        >
          Siguiente
        </button>
      </div>
    </div>

    <div class="right-panel">
      <h2>Detalle del Email</h2>
      <div v-if="selectedEmail" class="email-body">
        <button class="close-button" @click="closeEmail">X</button>
        <h3>{{ selectedEmail.subject }}</h3>
        <p><strong>From:</strong> {{ selectedEmail.from }}</p>
        <p><strong>To:</strong> {{ selectedEmail.to }}</p>
        <p><strong>Date:</strong> {{ selectedEmail.date }}</p>
        <hr />
        <p class="body-text">{{ selectedEmail.body }}</p>
      </div>
      <div v-else class="empty-state">
        <p>Selecciona un email para ver el detalle</p>
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
};

const closeEmail = () => {
  selectedEmail.value = null;
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
/* Estilos ya definidos */
</style>
