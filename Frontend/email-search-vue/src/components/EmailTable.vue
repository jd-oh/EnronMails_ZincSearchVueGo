<template>
  <div class="fixed inset-0 w-full h-screen p-8 overflow-y-auto bg-gray-100 dark:bg-gray-900">
    <h1 class="mb-4 text-3xl font-extrabold text-gray-900 dark:text-white md:text-5xl lg:text-6xl text-center">
      <span class="text-transparent bg-clip-text bg-gradient-to-r to-emerald-600 from-sky-400">Enron Mails</span> List
    </h1>

    <!-- Barra de búsqueda -->
    <SearchBar 
      v-model="emailStore.searchQuery" 
      :field="emailStore.selectedField" 
      @update:field="field => emailStore.selectedField = field"
      @search="emailStore.fetchEmails" 
      class="mb-6"
    />

    <!-- Lista de correos -->
    <EmailList 
      :emails="emailStore.emails" 
      :selectedEmail="emailStore.selectedEmail"
      @select="selectEmail"
    />

    <!-- Paginación -->
    <Pagination 
      :page="emailStore.page" 
      :totalPages="emailStore.totalPages" 
      :totalEmails="emailStore.totalEmails"
      :pageSize="emailStore.pageSize"
      @prev="emailStore.prevPage"
      @next="emailStore.nextPage"
    />
    
    <!-- Modal de detalles de email -->
    <EmailDetail 
      :email="emailStore.selectedEmail"
      :isVisible="isModalVisible"
      @close="closeModal"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import SearchBar from './SearchBar.vue';
import EmailList from './EmailList.vue';
import EmailDetail from './EmailDetail.vue';
import Pagination from './Pagination.vue';
import { useEmailStore } from '../stores/emailStore';

const emailStore = useEmailStore();

const isModalVisible = ref(false);

const selectEmail = (email) => {
  emailStore.selectedEmail = email;
  isModalVisible.value = true;
};

const closeModal = () => {
  isModalVisible.value = false;
};

onMounted(() => emailStore.fetchEmails());
</script>
