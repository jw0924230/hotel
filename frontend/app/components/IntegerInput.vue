<template>
  <div class="integer-input">
    <input
      type="text"
      inputmode="numeric"
      pattern="[0-9]*"
      autocomplete="off"
      :value="displayValue"
      :placeholder="placeholder"
      @beforeinput="preventInvalidInput"
      @input="handleInput"
      @paste="handlePaste"
    />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  modelValue: number
  placeholder?: string
}>(), {
  placeholder: '0'
})

const emit = defineEmits<{
  'update:modelValue': [value: number]
}>()

const displayValue = computed(() =>
  Number.isFinite(props.modelValue) && props.modelValue > 0
    ? String(Math.trunc(props.modelValue))
    : ''
)

const digitsOnly = (value: string) => value.replace(/\D/g, '')

const preventInvalidInput = (event: InputEvent) => {
  if (event.inputType.startsWith('delete') || event.data === null) return
  if (!/^\d+$/.test(event.data)) event.preventDefault()
}

const updateValue = (input: HTMLInputElement, rawValue: string) => {
  const sanitized = digitsOnly(rawValue)
  input.value = sanitized
  emit('update:modelValue', sanitized ? Number.parseInt(sanitized, 10) : 0)
}

const handleInput = (event: Event) => {
  const input = event.target as HTMLInputElement
  updateValue(input, input.value)
}

const handlePaste = (event: ClipboardEvent) => {
  event.preventDefault()
  const input = event.target as HTMLInputElement
  const pasted = digitsOnly(event.clipboardData?.getData('text') || '')
  const start = input.selectionStart || 0
  const end = input.selectionEnd || 0
  updateValue(input, input.value.slice(0, start) + pasted + input.value.slice(end))
}
</script>

<style scoped>
.integer-input {
  width: 100%;
}

.integer-input input {
  width: 100%;
  padding: 10px 14px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  font-size: 14px;
  background-color: white;
  color: #1e293b;
  transition: all 0.2s;
  box-sizing: border-box;
}

.integer-input input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.15);
}

:global(.dark-mode) .integer-input input {
  background-color: #0f172a;
  border-color: #334155;
  color: #ffffff;
}

:global(.dark-mode) .integer-input input:focus {
  border-color: #38bdf8;
  box-shadow: 0 0 0 3px rgba(56, 189, 248, 0.15);
}
</style>
