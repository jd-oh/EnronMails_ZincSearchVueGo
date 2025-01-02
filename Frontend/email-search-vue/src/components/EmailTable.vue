<template>
  <div class="fixed inset-0 w-full h-screen p-8 overflow-y-auto bg-gray-100">
  <h1 class="mb-4 text-3xl font-extrabold text-gray-900 dark:text-white md:text-5xl lg:text-6xl"><span class="text-transparent bg-clip-text bg-gradient-to-r to-emerald-600 from-sky-400">Enron Mails</span> List</h1>
   <!-- Barra de búsqueda -->
    <SearchBar 
      v-model="searchQuery" 
      @search="fetchEmails" 
      class="mb-6"
    />

    <!-- Lista de correos -->
    <EmailList 
      :emails="emails" 
      :selectedEmail="selectedEmail"
      @select="selectEmail"
    />

    <!-- Paginación -->
    <Pagination 
      :page="page" 
      :totalPages="totalPages" 
      :totalEmails="totalEmails"
      :pageSize="pageSize"
      @prev="prevPage"
      @next="nextPage"
    />
    
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
import EmailList from './EmailList.vue';  
import EmailDetail from './EmailDetail.vue'; 
import Pagination from './Pagination.vue'; 

const emails = ref([]);
const selectedEmail = ref(null);
const searchQuery = ref('');
const page = ref(0);  
const pageSize = 5;  
const totalEmails = ref(0);
const totalPages = ref(1);
const isModalVisible = ref(false); 


const fetchEmails = async () => {
  try {
    const response = await axios.post(
      'http://localhost:8080/api/search',
      {
        term: searchQuery.value || "*", 
        from: page.value * pageSize,
        size: pageSize
      }
    );

    emails.value = response.data.hits.hits.map(hit => hit._source);
    totalEmails.value = response.data.hits.total ? response.data.hits.total.value : 0;
    console.log(totalEmails.value)
    totalPages.value = Math.ceil(totalEmails.value / pageSize);
  } catch (error) {
    console.error('Error fetching emails:', error);
  }
};

const selectEmail = (email) => {
  selectedEmail.value = email;
  isModalVisible.value = true; 
};

const closeModal = () => {
  isModalVisible.value = false; 
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
