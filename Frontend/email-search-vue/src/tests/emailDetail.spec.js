import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import EmailDetail from '../../src/components/EmailDetail.vue'

describe('EmailDetail.vue', () => {
    it('no muestra contenido si isVisible es false', () => {
      const wrapper = mount(EmailDetail, {
        props: {
          email: {},
          isVisible: false
        }
      })
      expect(wrapper.find('.fixed').exists()).toBe(false)
    })

  it('muestra contenido si isVisible es true', () => {
    const wrapper = mount(EmailDetail, {
      props: {
        email: {
          subject: 'Prueba',
          from: 'test@example.com',
          to: 'dest@example.com',
          date: '2023-01-01',
          body: 'Contenido del correo'
        },
        isVisible: true
      }
    })
    expect(wrapper.text()).toContain('Prueba')
    expect(wrapper.text()).toContain('test@example.com')
  })

  it('emite "close" al hacer click en el botón de cerrar', async () => {
    const wrapper = mount(EmailDetail, {
      props: {
        email: {},
        isVisible: true
      }
    })
    await wrapper.find('button[aria-label="Close modal"]').trigger('click')
    expect(wrapper.emitted().close).toBeTruthy()
  })
})