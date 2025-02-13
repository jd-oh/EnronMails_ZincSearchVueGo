<template>
  <div class="search-container flex flex-col sm:flex-row gap-4 max-w-lg mx-auto">
    <div class="relative w-full sm:w-1/3">
      <select 
        v-model="selectedField" 
        @change="emitSearch" 
        class="block py-2.5 px-0 w-full text-sm text-gray-500 bg-transparent border-0 border-b-2 border-gray-200 appearance-none focus:outline-none focus:ring-0 focus:border-gray-500">
        <option value="_all">Todos los campos</option>
        <option value="body">Body</option>
        <option value="from">From</option>
        <option value="to">To</option>
        <option value="subject">Subject</option>
      </select>
    </div>
    
    <input 
      v-model="searchInput" 
      placeholder="Search..."
      class="search-box w-full py-2 px-3 border-b-2 border-gray-200 focus:outline-none focus:border-gray-500 text-sm text-gray-700" 
      @input="handleInput"
    />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useEmailStore } from '../stores/emailStore'

const emailStore = useEmailStore()

// Inicialmente _all para que se muestre "Todos los campos"
const searchInput = ref('')
const selectedField = ref('_all')

const emit = defineEmits(['update:field', 'search'])

let timeoutId = null

const handleInput = (event) => {
  searchInput.value = event.target.value
  clearTimeout(timeoutId)
  timeoutId = setTimeout(() => {
    const trimmed = searchInput.value.trim()
    if(trimmed !== '') {
      if(emailStore.textFilter) {
        emailStore.textFilter += ' ' + trimmed
      } else {
        emailStore.textFilter = trimmed
      }
    }
    searchInput.value = ''
    emailStore.fetchEmails()
    emit('search')
  }, 500)
}

const emitSearch = () => {
  emailStore.selectedField = selectedField.value || '_all'
  emit('update:field', emailStore.selectedField)
  emailStore.fetchEmails()
  emit('search')
}
</script>