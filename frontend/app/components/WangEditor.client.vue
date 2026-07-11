<template>
  <div class="wang-editor-shell">
    <Toolbar
      class="wang-editor-toolbar"
      :editor="editorRef"
      :default-config="toolbarConfig"
      mode="default"
    />
    <Editor
      class="wang-editor-content"
      :style="{ height }"
      v-model="htmlValue"
      :default-config="editorConfig"
      mode="default"
      @on-created="handleCreated"
      @custom-paste="handlePaste"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, shallowRef } from 'vue'
import {
  i18nAddResources,
  i18nChangeLanguage,
  i18nGetResources,
  type IDomEditor,
  type IEditorConfig,
  type IToolbarConfig
} from '@wangeditor/editor'
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'
import '@wangeditor/editor/dist/css/style.css'

const props = withDefaults(defineProps<{
  modelValue: string
  placeholder?: string
  height?: string
}>(), {
  placeholder: '請輸入內容…',
  height: '280px'
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const traditionalCharacters: Record<string, string> = {
  '体': '體', '简': '簡', '标': '標', '题': '題', '链': '鏈', '接': '接',
  '图': '圖', '片': '片', '视': '視', '频': '頻', '颜': '顏', '色': '色',
  '背': '背', '景': '景', '字': '字', '号': '號', '对': '對', '齐': '齊',
  '列': '列', '表': '表', '有': '有', '无': '無', '序': '序', '清': '清',
  '除': '除', '格': '格', '式': '式', '撤': '撤', '销': '銷', '恢': '恢',
  '复': '復', '粗': '粗', '斜': '斜', '下': '下', '划': '劃', '线': '線',
  '删': '刪', '选': '選', '择': '擇', '输': '輸', '入': '入', '全': '全',
  '屏': '螢', '查': '查', '看': '看', '代': '代', '码': '碼', '编': '編',
  '辑': '輯', '增': '增', '加': '加', '取': '取', '消': '消', '确': '確',
  '定': '定', '设': '設', '置': '置', '网': '網', '页': '頁', '预': '預',
  '览': '覽', '文': '文', '本': '本', '块': '塊', '引': '引', '用': '用'
}

const toTraditional = (value: unknown): unknown => {
  if (typeof value === 'string') {
    return [...value].map(character => traditionalCharacters[character] || character).join('')
  }
  if (Array.isArray(value)) return value.map(toTraditional)
  if (value && typeof value === 'object') {
    return Object.fromEntries(
      Object.entries(value as Record<string, unknown>)
        .map(([key, item]) => [key, toTraditional(item)])
    )
  }
  return value
}

const traditionalResources = toTraditional(
  i18nGetResources('zh-CN')
) as Record<string, any>
i18nAddResources('zh-TW', traditionalResources)
i18nChangeLanguage('zh-TW')

const editorRef = shallowRef<IDomEditor>()
const htmlValue = computed({
  get: () => props.modelValue || '',
  set: value => emit('update:modelValue', value)
})

const toolbarConfig: Partial<IToolbarConfig> = {
  toolbarKeys: [
    'headerSelect',
    '|',
    'bulletedList',
    'numberedList',
    '|',
    'bold',
    'italic',
    'through',
    '|',
    'insertLink'
  ]
}

const editorConfig: Partial<IEditorConfig> = {
  placeholder: props.placeholder,
  MENU_CONF: {
    headerSelect: {
      levels: ['h1', 'h2', 'h3']
    }
  },
  hoverbarKeys: {
    text: {
      menuKeys: [
        'headerSelect',
        'bulletedList',
        'numberedList',
        'bold',
        'italic',
        'through',
        'insertLink'
      ]
    },
    link: {
      menuKeys: ['editLink', 'unLink', 'viewLink']
    }
  }
}

const handleCreated = (editor: IDomEditor) => {
  editorRef.value = editor
}

const handlePaste = (
  _editor: IDomEditor,
  event: ClipboardEvent,
  callback: (allowDefault: boolean) => void
) => {
  const clipboard = event.clipboardData
  const containsImageFile = Array.from(clipboard?.items || [])
    .some(item => item.kind === 'file' && item.type.startsWith('image/'))
  const containsImageHtml = /<(img|picture|source)\b/i.test(
    clipboard?.getData('text/html') || ''
  )

  if (containsImageFile || containsImageHtml) {
    event.preventDefault()
    callback(false)
    return
  }

  callback(true)
}

onBeforeUnmount(() => {
  editorRef.value?.destroy()
})
</script>

<style scoped>
.wang-editor-shell {
  overflow: hidden;
  border: 1px solid #cbd5e1;
  border-radius: 9px;
  background: #fff;
}

.wang-editor-toolbar {
  border-bottom: 1px solid #e2e8f0;
  background: #f8fafc;
}

.wang-editor-content {
  overflow-y: auto;
  color: #1e293b;
  font-size: 14px;
  line-height: 1.7;
}
</style>
