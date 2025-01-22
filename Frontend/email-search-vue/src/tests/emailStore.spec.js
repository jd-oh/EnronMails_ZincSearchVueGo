import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useEmailStore } from '../../src/stores/emailStore'
import emailService from '../../src/services/emailService'

vi.mock('../../src/services/emailService', () => ({
    default: {
      fetchEmails: vi.fn()
    }
  }))
  
  describe('emailStore', () => {
    beforeEach(() => {
      setActivePinia(createPinia())
      vi.resetAllMocks()
    })
  
    it('fetchEmails actualiza el estado de emails', async () => {
      const mockResponse = {
        hits: {
          hits: [
            { _source: { subject: 'Test Email 1' } },
            { _source: { subject: 'Test Email 2' } }
          ],
          total: { value: 2 }
        }
      }
      // Accede directamente a fetchEmails
      emailService.fetchEmails.mockResolvedValueOnce(mockResponse)
  
      const store = useEmailStore()
      await store.fetchEmails()
  
      expect(emailService.fetchEmails).toHaveBeenCalled()
      expect(store.emails).toHaveLength(2)
      expect(store.totalEmails).toBe(2)
      expect(store.totalPages).toBe(1)
    })
  
    it('prevPage no hace nada si está en la página 0', () => {
      const store = useEmailStore()
      store.page = 0
      store.prevPage()
      expect(store.page).toBe(0) // no retrocede más
    })
  
    it('prevPage disminuye la página si es mayor a 0', async () => {
      const mockResponse = {
        hits: {
          hits: [
            { _source: { subject: 'Test Email 1' } },
            { _source: { subject: 'Test Email 2' } }
          ],
          total: { value: 2 }
        }
      }
      emailService.fetchEmails.mockResolvedValueOnce(mockResponse)
  
      const store = useEmailStore()
      store.page = 1
      await store.prevPage()
      expect(store.page).toBe(0)
    })
  
    it('nextPage avanza la página si hay más emails', async () => {
      const mockResponse = {
        hits: {
          hits: [
            { _source: { subject: 'Test Email 3' } },
            { _source: { subject: 'Test Email 4' } }
          ],
          total: { value: 10 }
        }
      }
      emailService.fetchEmails.mockResolvedValueOnce(mockResponse)
  
      const store = useEmailStore()
      store.page = 0
      store.totalEmails = 10
      store.pageSize = 5
      await store.nextPage()
      expect(store.page).toBe(1)
    })
  
    it('nextPage no avanza la página si no hay más emails', () => {
      const store = useEmailStore()
      store.page = 1
      store.totalEmails = 10
      store.pageSize = 10
      store.nextPage()
      expect(store.page).toBe(1) // se mantiene, no avanza
    })
  })