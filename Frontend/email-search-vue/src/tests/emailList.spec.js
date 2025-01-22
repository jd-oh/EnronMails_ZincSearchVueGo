import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import EmailList from '../../src/components/EmailList.vue'

describe('EmailList.vue', () => {
  it('muestra la tabla si hay correos', () => {
    const wrapper = mount(EmailList, {
      props: {
        emails: [{ _id: 1, subject: 'Test', from: 'test@example.com' }]
      }
    })
    expect(wrapper.find('table').exists()).toBe(true)
  })

  it('muestra mensaje de “No results found” si no hay correos', () => {
    const wrapper = mount(EmailList, {
      props: {
        emails: []
      }
    })
    expect(wrapper.text()).toContain('No results found')
  })

  it('emite "select" al hacer click en una fila', async () => {
    const email = { _id: 1, subject: 'Test', from: 'test@example.com' }
    const wrapper = mount(EmailList, {
      props: {
        emails: [email]
      }
    })
    await wrapper.find('tbody tr').trigger('click')
    expect(wrapper.emitted().select[0][0]).toEqual(email)
  })
})