<template>
  <div class="container">
    <div class="left-panel">
      <input 
        v-model="searchQuery" 
        placeholder="Buscar por subject, from o to" 
        class="search-box"
        @input="fetchEmails(0)"
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

const emails = ref([]);
const selectedEmail = ref(null);
const searchQuery = ref('');
const page = ref(0);  
const pageSize = 10;  
const totalEmails = ref(0);
const totalPages = ref(1);

const fetchEmails = async (pageNumber = 0) => {
  try {
    const response = await axios.post(
      'http://localhost:8080/api/search',
      {
        term: searchQuery.value || "*",
        from: pageNumber * pageSize,
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
    fetchEmails(page.value);
  }
};

const nextPage = () => {
  if ((page.value + 1) * pageSize < totalEmails.value) {
    page.value++;
    fetchEmails(page.value);
  }
};

onMounted(() => fetchEmails(0));
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
  font-size: 1em;
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
.email-body {
  background-color: #fff;
  border: 1px solid #ddd;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  position: relative;
}
.body-text {
  white-space: pre-wrap;
  font-family: monospace;
}
.empty-state {
  text-align: center;
  margin-top: 50px;
  color: #888;
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
</style>
