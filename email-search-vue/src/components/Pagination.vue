<template>
  <div class="flex flex-col items-center">
    <!-- Texto de ayuda -->
    <span class="text-sm text-gray-700 dark:text-gray-400">
      Showing 
      <span class="font-semibold text-gray-900 dark:text-white">
        {{ page * pageSize + 1 }}
      </span> 
      to 
      <span class="font-semibold text-gray-900 dark:text-white">
        {{ Math.min((page + 1) * pageSize, totalEmails) }}
      </span> 
      of 
      <span class="font-semibold text-gray-900 dark:text-white">
        {{ totalEmails }}
      </span> 
      Entries
    </span>

    <!-- Botones -->
    <div class="inline-flex mt-2 xs:mt-0">
      <button 
        @click="$emit('prev')" 
        :disabled="page <= 0" 
        class="flex items-center justify-center px-4 h-10 text-base font-medium text-white bg-gray-800 rounded-s hover:bg-gray-900 disabled:opacity-50 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white"
      >
        Prev
      </button>
      <button 
        @click="$emit('next')" 
        :disabled="(page + 1) * pageSize >= totalEmails" 
        class="flex items-center justify-center px-4 h-10 text-base font-medium text-white bg-gray-800 border-0 border-s border-gray-700 rounded-e hover:bg-gray-900 disabled:opacity-50 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white"
      >
        Next
      </button>
    </div>
  </div>
</template>

<script setup>
defineProps(['page', 'totalPages', 'totalEmails', 'pageSize']);
defineEmits(['prev', 'next']);
</script>
