import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import SearchBar from '../../src/components/SearchBar.vue'

describe('SearchBar.vue', () => {
  let wrapper

  beforeEach(() => {
    wrapper = mount(SearchBar)
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.runOnlyPendingTimers()
    vi.useRealTimers()
  })

  it('emits "update:modelValue" when input value changes', async () => {
    const input = wrapper.find('input')
    await input.setValue('test search')
    expect(wrapper.emitted()['update:modelValue']).toBeTruthy()
    expect(wrapper.emitted()['update:modelValue'][0]).toEqual(['test search'])
  })

  it('emits "update:field" and "search" when select value changes', async () => {
    const select = wrapper.find('select')
    await select.setValue('subject')
    expect(wrapper.emitted()['update:field']).toBeTruthy()
    expect(wrapper.emitted()['update:field'][0]).toEqual(['subject'])
    expect(wrapper.emitted().search).toBeTruthy()
    expect(wrapper.emitted().search.length).toBe(1)
  })

  it('emits "search" after debounce when typing in input', async () => {
    const input = wrapper.find('input')
    await input.setValue('debounce test')
    expect(wrapper.emitted().search).toBeFalsy()
    
    vi.advanceTimersByTime(500)
    await wrapper.vm.$nextTick()

    expect(wrapper.emitted().search).toBeTruthy()
    expect(wrapper.emitted().search.length).toBe(1)
  })

  it('resets debounce timer on consecutive inputs', async () => {
    const input = wrapper.find('input')
    await input.setValue('first input')
    vi.advanceTimersByTime(300)
    await input.setValue('second input')
    vi.advanceTimersByTime(300)
    expect(wrapper.emitted().search).toBeFalsy()

    vi.advanceTimersByTime(200)
    await wrapper.vm.$nextTick()

    expect(wrapper.emitted().search).toBeTruthy()
    expect(wrapper.emitted().search.length).toBe(1)
  })

  it('emits "search" immediately when select field changes', async () => {
    const select = wrapper.find('select')
    await select.setValue('from')
    expect(wrapper.emitted().search.length).toBe(1)
  })
})