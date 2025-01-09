// stores/emailStore.js
import { defineStore } from 'pinia';
import emailService from '../services/emailService';

export const useEmailStore = defineStore('email', {
  state: () => ({
    emails: [],
    selectedEmail: null,
    searchQuery: '',
    selectedField: 'body',
    page: 0,
    pageSize: 5,
    totalEmails: 0,
    totalPages: 1,
  }),
  actions: {
    async fetchEmails() {
        try {
          const data = await emailService.fetchEmails({
            term: this.searchQuery,
            field: this.selectedField,
            from: this.page * this.pageSize,
            size: this.pageSize,
          });
          console.log('Fetched data:', data); // Verifica la estructura de los datos
          this.emails = data.hits.hits.map(hit => hit._source);
          this.totalEmails = data.hits.total ? data.hits.total.value : 0;
          this.totalPages = Math.ceil(this.totalEmails / this.pageSize);
        } catch (error) {
          console.error('Error fetching emails:', error);
        }
      }
    },      
});
