// Frontend/email-search-vue/src/stores/emailStore.js

/**
 * @module emailStore
 * @description Pinia store para gestionar el estado de búsqueda de correos electrónicos y paginación.
 */
import { defineStore } from 'pinia';
import emailService from '../services/emailService';

/**
 * Pinia store para el manejo de emails.
 *
 * @property {Array} state.emails - Lista de correos electrónicos obtenidos de la búsqueda.
 * @property {Object|null} state.selectedEmail - Correo actualmente seleccionado para ver sus detalles.
 * @property {number} state.page - Página actual de la paginación.
 * @property {number} state.pageSize - Número de elementos por página.
 * @property {number} state.totalEmails - Total de correos electrónicos encontrados.
 * @property {number} state.totalPages - Total de páginas calculadas.
 * @property {string} state.folderFilter - Filtro de carpeta aplicado a la búsqueda.
 * @property {string} state.textFilter - Filtro de texto aplicado a la búsqueda.
 * @property {string} state.selectedField - Campo seleccionado para la búsqueda (_all por defecto).
 */
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
    /**
     * Realiza la búsqueda de emails utilizando el servicio emailService.
     * Combina los filtros de carpeta y texto, y calcula la paginación.
     *
     * @async
     * @function fetchEmails
     * @returns {Promise<void>}
     */
    async fetchEmails() {
      let combinedTerm = '';
      if (this.folderFilter && this.textFilter) {
        combinedTerm = `${this.folderFilter} ${this.textFilter}`;
      } else if (this.folderFilter) {
        combinedTerm = `${this.folderFilter}`;
      } else if (this.textFilter) {
        combinedTerm = `${this.textFilter}`;
      }
      
      // Definir el campo a utilizar para la búsqueda
      let searchField = this.selectedField;
      if (this.folderFilter && this.textFilter) {
        console.log("Los dos filtros están activos");
        // Si se aplican ambos filtros, se fuerza el campo _all
        searchField = '_all';
      }
      try {
        console.log("La búsqueda es:" + combinedTerm + " y el campo es:" + searchField);
        // Llamada al servicio para obtener emails
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
    /**
     * Remueve un filtro específico (carpeta o texto) y vuelve a realizar la búsqueda.
     *
     * @function removeFilter
     * @param {string} filterType - Tipo de filtro a remover ('folder' o 'text').
     */
    removeFilter(filterType) {
      if (filterType === 'folder') {
        this.folderFilter = '';
        // Al remover el filtro de carpeta, selectedField se mantiene con el valor seleccionado en el SearchBar.
      } else if (filterType === 'text') {
        this.textFilter = '';
      }
      this.fetchEmails();
    },
    /**
     * Navega a la página anterior, decrementando el índice de página y actualizando los emails.
     *
     * @function prevPage
     */
    prevPage() {
      if (this.page > 0) {
        this.page--;
        this.fetchEmails();
      }
    },
    /**
     * Navega a la página siguiente, incrementando el índice de página si hay más resultados.
     *
     * @function nextPage
     */
    nextPage() {
      if ((this.page + 1) * this.pageSize < this.totalEmails) {
        this.page++;
        this.fetchEmails();
      }
    },
    /**
     * Establece el filtro de carpeta basado en la carpeta principal y subcarpeta seleccionadas,
     * actualizando el estado y reejecutando la búsqueda.
     *
     * @function setFolderFilter
     * @param {string} main - Carpeta principal seleccionada.
     * @param {string} sub - Subcarpeta seleccionada.
     */
    setFolderFilter(main, sub) {
      const folderPath = `${main}/${sub}`;
      this.folderFilter = 'enron_mail_20110402/maildir/' + folderPath;
      this.fetchEmails();
    },
  },
});