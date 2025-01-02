<template>
  <div v-if="email" class="email-body">
    <button class="close-button" @click="$emit('close')">X</button>
    <h3>{{ email.subject }}</h3>
    <p><strong>From:</strong> {{ email.from }}</p>
    <p><strong>To:</strong> {{ email.to }}</p>
    <p><strong>Date:</strong> {{ email.date }}</p>
    <hr />
    <p class="body-text">{{ email.body }}</p>
  </div>
  <div v-else class="empty-state">
    <p>Selecciona un email para ver el detalle</p>
  </div>
</template>

<script setup>
defineProps(['email']);
defineEmits(['close']);
</script>

<style scoped>
.email-body {
  border: 1px solid #ddd;
  padding: 20px;
}
</style>