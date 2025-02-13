<template>
  <div class="fixed inset-0 w-full h-screen p-8 overflow-y-auto bg-gray-100 dark:bg-gray-900">
    <h1 class="mb-4 text-3xl font-extrabold text-gray-900 dark:text-white md:text-5xl lg:text-6xl text-center">
      <span class="text-transparent bg-clip-text bg-gradient-to-r from-blue-700 to-blue-900">Enron Mails</span> List
    </h1>

    <!-- Contenedor: SearchBar -->
    <div class="flex items-center mb-6">
      <SearchBar
        v-model="emailStore.textFilter"
        :field="emailStore.selectedField"
        @update:field="field => emailStore.selectedField = field"
        @search="emailStore.fetchEmails"
        class="flex-grow"
      />
      <!-- Botón para mostrar/ocultar FolderDropDown -->
      <button
        class="ml-4 px-4 py-2 bg-blue-500 text-white rounded-md flex items-center justify-center"
        @click="toggleFolderDropDown"
      >
        <span v-if="showFolderDropDown">
          <!-- Folder open icon -->
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="currentColor" viewBox="0 0 20 20">
            <path d="M2 4a2 2 0 012-2h4l2 2h6a2 2 0 012 2v1H2V4z" />
            <path d="M2 9h16v7a2 2 0 01-2 2H4a2 2 0 01-2-2V9z" opacity="0.5"/>
          </svg>
        </span>
        <span v-else>
          <!-- Folder closed icon -->
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="currentColor" viewBox="0 0 20 20">
            <path d="M2 4a2 2 0 012-2h4l2 2h6a2 2 0 012 2v1H2V4z" />
            <path d="M2 9h16v7a2 2 0 01-2 2H4a2 2 0 01-2-2V9z" />
          </svg>
        </span>
      </button>
    </div>

    <!-- Folder Navigation Toggleable Centered -->
    <div v-if="showFolderDropDown" class="mb-6 flex flex-col items-center">
      <!-- Dropdown para carpeta principal -->
      <div class="relative inline-block text-left mb-2">
        <button
          @click="toggleDropdown"
          type="button"
          class="inline-flex justify-between w-64 rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-sm font-medium text-gray-700 hover:bg-gray-50"
        >
          {{ selectedMainFolder || 'Selecciona una carpeta' }}
          <svg class="w-5 h-5 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              v-if="dropdownOpen"
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M5 15l7-7 7 7"
            />
            <path
              v-else
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M19 9l-7 7-7-7"
            />
          </svg>
        </button>
        <div
          v-if="dropdownOpen"
          class="origin-top-right absolute z-10 mt-2 w-64 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5"
        >
          <div class="py-1">
            <button
              v-for="(subs, main) in folders"
              :key="main"
              @click="selectMainFolder(main)"
              class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
            >
              {{ main }}
            </button>
          </div>
        </div>
      </div>

      <!-- Dropdown para subcarpetas -->
      <div v-if="selectedMainFolder" class="relative inline-block text-left">
        <button
          @click="toggleDropdownSub"
          type="button"
          class="inline-flex justify-between w-64 rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-sm font-medium text-gray-700 hover:bg-gray-50"
        >
          {{ selectedSubFolder || 'Selecciona una subcarpeta' }}
          <svg class="w-5 h-5 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              v-if="dropdownSubOpen"
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M5 15l7-7 7 7"
            />
            <path
              v-else
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M19 9l-7 7-7-7"
            />
          </svg>
        </button>
        <div
          v-if="dropdownSubOpen"
          class="origin-top-right absolute z-10 mt-2 w-64 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5"
        >
          <div class="py-1">
            <button
              v-for="sub in folders[selectedMainFolder]"
              :key="sub"
              @click="selectFolder(selectedMainFolder, sub)"
              class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
            >
              {{ sub }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Filtros activos -->
    <div class="mb-4 flex flex-wrap">
      <!-- Badge para folderFilter -->
      <span
        v-if="emailStore.folderFilter"
        class="inline-flex items-center px-2 py-1 me-2 text-sm font-medium text-blue-800 bg-blue-100 rounded-sm dark:bg-blue-900 dark:text-blue-300"
      >
        {{ emailStore.folderFilter }}
        <button
          @click="emailStore.removeFilter('folder')"
          type="button"
          class="inline-flex items-center p-1 ms-2 text-sm text-blue-400 bg-transparent rounded-xs hover:bg-blue-200 hover:text-blue-900 dark:hover:bg-blue-800 dark:hover:text-blue-300"
          aria-label="Remove"
        >
          <svg class="w-2 h-2" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 14 14">
            <path
              stroke="currentColor"
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6"
            />
          </svg>
          <span class="sr-only">Remove badge</span>
        </button>
      </span>
      <!-- Badge para textFilter -->
      <span
        v-if="emailStore.textFilter"
        class="inline-flex items-center px-2 py-1 me-2 text-sm font-medium text-blue-800 bg-blue-100 rounded-sm dark:bg-blue-900 dark:text-blue-300"
      >
        {{ emailStore.selectedField === '_all' ? emailStore.textFilter : (emailStore.selectedField + ': ' + emailStore.textFilter) }}
        <button
          type="button"
          @click="emailStore.removeFilter('text')"
          class="inline-flex items-center p-1 ms-2 text-sm text-blue-400 bg-transparent rounded-xs hover:bg-blue-200 hover:text-blue-900 dark:hover:bg-blue-800 dark:hover:text-blue-300"
          aria-label="Remove"
        >
          <svg class="w-2 h-2" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 14 14">
            <path
              stroke="currentColor"
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6"
            />
          </svg>
          <span class="sr-only">Remove badge</span>
        </button>
      </span>
    </div>

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

    <!-- Modal de detalles -->
    <EmailDetail
      :email="emailStore.selectedEmail"
      :isVisible="isModalVisible"
      @close="closeModal"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import SearchBar from './SearchBar.vue'
import EmailList from './EmailList.vue'
import EmailDetail from './EmailDetail.vue'
import Pagination from './Pagination.vue'
import { useEmailStore } from '../stores/emailStore'
import { fetchFolders } from '../services/folderService';


const emailStore = useEmailStore()
const isModalVisible = ref(false)
const showFolderDropDown = ref(false)

// Variables y métodos para los dropdowns
const folders = ref({})
const dropdownOpen = ref(false)
const dropdownSubOpen = ref(false)
const selectedMainFolder = ref('')
const selectedSubFolder = ref('')

const toggleDropdown = () => {
  dropdownOpen.value = !dropdownOpen.value
}
const toggleDropdownSub = () => {
  dropdownSubOpen.value = !dropdownSubOpen.value
}
const toggleFolderDropDown = () => {
  showFolderDropDown.value = !showFolderDropDown.value
}

const selectMainFolder = (main) => {
  selectedMainFolder.value = main
  selectedSubFolder.value = ''
  dropdownOpen.value = false
  dropdownSubOpen.value = false
  emailStore.folderFilter = ''
}

const selectFolder = (main, sub) => {
  selectedSubFolder.value = sub
  emailStore.setFolderFilter(main, sub)
  dropdownSubOpen.value = false
}

onMounted(async () => {
  try {
    folders.value = await fetchFolders();
  } catch (error) {
    console.error("Error al cargar las carpetas:", error);
  }
});

const selectEmail = (email) => {
  emailStore.selectedEmail = email
  isModalVisible.value = true
}

const closeModal = () => {
  isModalVisible.value = false
}
</script>

