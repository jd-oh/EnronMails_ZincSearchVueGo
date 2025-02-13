<template>
  <div class="p-4">
    <label class="block font-bold mb-2">Carpetas</label>
    <!-- Dropdown de carpetas principales -->
    <div class="relative inline-block text-left">
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
        class="origin-top-right absolute z-10 mt-2 w-64 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 overflow-y-auto"
        style="max-height: 10rem;"
      >
        <div class="py-1">
          <button
            v-for="(subFolders, mainFolder) in folders"
            :key="mainFolder"
            @click="selectMainFolder(mainFolder)"
            class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
          >
            {{ mainFolder }}
          </button>
        </div>
      </div>
    </div>

    <!-- Dropdown de subcarpetas -->
    <div v-if="selectedMainFolder" class="mt-4">
      <label class="block font-bold mb-2">
        Subcarpetas de {{ selectedMainFolder }}
      </label>
      <div class="relative inline-block text-left">
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
          class="origin-top-right absolute z-10 mt-2 w-64 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 overflow-y-auto"
          style="max-height: 10rem;"
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
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useEmailStore } from '../stores/emailStore'
import { fetchFolders } from '../services/folderService'

const emailStore = useEmailStore()
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

const selectMainFolder = (mainFolder) => {
  selectedMainFolder.value = mainFolder
  selectedSubFolder.value = ''
  dropdownOpen.value = false
  dropdownSubOpen.value = false
}

const selectFolder = (main, sub) => {
  selectedSubFolder.value = sub
  const folderPath = `${main}/${sub}`
  emailStore.folderFilter = 'enron_mail_20110402/maildir/' + folderPath
  emailStore.fetchEmails()
  dropdownSubOpen.value = false
}

onMounted(async () => {
  try {
    folders.value = await fetchFolders()
  } catch (error) {
    console.error("Error al cargar las carpetas:", error)
  }
})
</script>
  