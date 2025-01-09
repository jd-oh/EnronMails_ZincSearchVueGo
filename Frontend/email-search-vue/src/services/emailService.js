import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080/api',
  timeout: 5000, // Tiempo de espera para las solicitudes
});

// Interceptor para manejar errores globalmente
api.interceptors.response.use(
  (response) => response,
  (error) => {
    console.error('API Error:', error.response || error.message);
    return Promise.reject(error);
  }
);

/**
 * Buscar correos electrónicos en la API.
 * @param {Object} params - Parámetros de búsqueda.
 * @param {string} params.term - Término de búsqueda.
 * @param {string} params.field - Campo para buscar.
 * @param {number} params.from - Índice inicial para la paginación.
 * @param {number} params.size - Tamaño de la página.
 * @returns {Promise<Object>} - Resultado de la búsqueda.
 */
export const fetchEmails = async ({ term, field, from, size }) => {
  console.log('Parameters sent to API:', { term, field, from, size }); // Verifica los valores
  const response = await api.post('/search', {
    term: term || '*',
    field,
    from,
    size,
  });
  return response.data;
};

export default { fetchEmails };
