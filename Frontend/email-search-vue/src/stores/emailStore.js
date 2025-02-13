// Frontend/email-search-vue/src/stores/emailStore.js
import { defineStore } from 'pinia';
import emailService from '../services/emailService';

export const useEmailStore = defineStore('emailStore', {
  state: () => ({
    emails: [],
    selectedEmail: null,
    page: 0,
    pageSize: 5,
    totalEmails: 0,
    totalPages: 0,
    folderFilter: '',
    textFilter: '',
    selectedField: '_all', // Este valor se actualiza desde el selector en SearchBar
  }),
  actions: {
    async fetchEmails() {
      let combinedTerm = '';
      if (this.folderFilter && this.textFilter) {
        combinedTerm = `${this.folderFilter} ${this.textFilter}`;
      } else if (this.folderFilter) {
        combinedTerm = `${this.folderFilter}`;
      } else if (this.textFilter) {
        combinedTerm = `${this.textFilter}`;
      }
      
      // Usamos una variable local para definir el field a usar en la búsqueda.
      let searchField = this.selectedField;
      if (this.folderFilter && this.textFilter) {
        console.log("Los dos filtros están activos");
        // Si aplicas ambos filtros, puedes forzar _all sin modificar selectedField
        searchField = '_all';
      }
      try {
        console.log("La búsqueda es:" + combinedTerm + " y el campo es:" + searchField);
        // Usa la función fetchEmails del servicio
        const data = await emailService.fetchEmails({
          term: combinedTerm,
          field: searchField,
          from: this.page * this.pageSize,
          size: this.pageSize,
        });
        console.log(combinedTerm);
        this.emails = data.hits.hits.map(hit => hit._source);
        this.totalEmails = data.hits.total ? data.hits.total.value : 0;
        this.totalPages = Math.ceil(this.totalEmails / this.pageSize);
      } catch (error) {
        console.error('Error fetching emails:', error);
      }
    },
    removeFilter(filterType) {
      if (filterType === 'folder') {
        this.folderFilter = '';
        // Al remover el filtro de carpeta, selectedField se mantiene con el valor seleccionado en el SearchBar.
      } else if (filterType === 'text') {
        this.textFilter = '';
      }
      this.fetchEmails();
    },
    prevPage() {
      if (this.page > 0) {
        this.page--;
        this.fetchEmails();
      }
    },
    nextPage() {
      if ((this.page + 1) * this.pageSize < this.totalEmails) {
        this.page++;
        this.fetchEmails();
      }
    },
    setFolderFilter(main, sub) {
      const folderPath = `${main}/${sub}`;
      this.folderFilter = 'enron_mail_20110402/maildir/' + folderPath;
      this.fetchEmails();
    },
  },
});