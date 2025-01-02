<template>
  <div class="pagination">
    <button @click="$emit('prev')" :disabled="page <= 0">Anterior</button>
    <span>Página {{ page + 1 }} de {{ totalPages }}</span>
    <button 
      @click="$emit('next')" 
      :disabled="(page + 1) * pageSize >= totalEmails"
    >Siguiente</button>
  </div>
</template>

<script setup>
defineProps(['page', 'totalPages', 'totalEmails', 'pageSize']);
defineEmits(['prev', 'next']);
</script>

<style scoped>
.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>