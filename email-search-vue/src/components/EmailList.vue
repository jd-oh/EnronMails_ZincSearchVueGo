<template>
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
        @click="$emit('select', email)"
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
</template>

<script setup>
defineProps(['emails', 'selectedEmail']);
defineEmits(['select']);
</script>

<style scoped>
.selected {
  background-color: #d1e7ff;
}
</style>