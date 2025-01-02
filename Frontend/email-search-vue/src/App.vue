<template>
  <div class="container">
    <h1>Email Search</h1>
    <EmailTable />
  </div>
</template>

<script setup>
import EmailTable from './components/EmailTable.vue';
</script>

<style>
.container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
  font-family: Arial, sans-serif;
}
</style>
