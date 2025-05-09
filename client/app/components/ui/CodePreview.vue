<script setup>
import hljs from 'highlight.js/lib/core'
import bash from 'highlight.js/lib/languages/bash'
import css from 'highlight.js/lib/languages/css'
import go from 'highlight.js/lib/languages/go'
import java from 'highlight.js/lib/languages/java'
import javascript from 'highlight.js/lib/languages/javascript'
import json from 'highlight.js/lib/languages/json'
import markdown from 'highlight.js/lib/languages/markdown'
import python from 'highlight.js/lib/languages/python'
import ruby from 'highlight.js/lib/languages/ruby'
import typescript from 'highlight.js/lib/languages/typescript'
import xml from 'highlight.js/lib/languages/xml'
import yaml from 'highlight.js/lib/languages/yaml'

import 'highlight.js/styles/atom-one-dark.css'

const props = defineProps({
  title: {
    type: String,
    default: '',
  },
  filename: {
    type: String,
    default: 'example.go',
  },
  customURL: {
    type: String,
    default: '',
  },
  lang: {
    type: String,
    default: 'javascript',
  },
  code: {
    type: String,
    default: '',
  },
  showEditButton: {
    type: Boolean,
    default: false,
  },
  showCopyButton: {
    type: Boolean,
    default: true,
  },
  allowEditing: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:code'])

hljs.registerLanguage('javascript', javascript)
hljs.registerLanguage('typescript', typescript)
hljs.registerLanguage('python', python)
hljs.registerLanguage('java', java)
hljs.registerLanguage('xml', xml)
hljs.registerLanguage('html', xml)
hljs.registerLanguage('css', css)
hljs.registerLanguage('bash', bash)
hljs.registerLanguage('ruby', ruby)
hljs.registerLanguage('go', go)
hljs.registerLanguage('json', json)
hljs.registerLanguage('yaml', yaml)
hljs.registerLanguage('markdown', markdown)

const codeContent = ref(props.code)
const isEditing = ref(props.allowEditing)
const isCopied = ref(false)
const codeBlock = ref(null)
const textarea = ref(null)
const editorHeight = ref('200px')

const languageMap = {
  js: { display: 'JavaScript', highlight: 'javascript' },
  javascript: { display: 'JavaScript', highlight: 'javascript' },
  jsx: { display: 'JSX', highlight: 'javascript' },
  ts: { display: 'TypeScript', highlight: 'typescript' },
  typescript: { display: 'TypeScript', highlight: 'typescript' },
  tsx: { display: 'TSX', highlight: 'typescript' },
  html: { display: 'HTML', highlight: 'html' },
  xml: { display: 'XML', highlight: 'xml' },
  css: { display: 'CSS', highlight: 'css' },
  py: { display: 'Python', highlight: 'python' },
  python: { display: 'Python', highlight: 'python' },
  java: { display: 'Java', highlight: 'java' },
  rb: { display: 'Ruby', highlight: 'ruby' },
  ruby: { display: 'Ruby', highlight: 'ruby' },
  go: { display: 'Go', highlight: 'go' },
  json: { display: 'JSON', highlight: 'json' },
  yaml: { display: 'YAML', highlight: 'yaml' },
  yml: { display: 'YAML', highlight: 'yaml' },
  sh: { display: 'Shell', highlight: 'bash' },
  bash: { display: 'Bash', highlight: 'bash' },
  md: { display: 'Markdown', highlight: 'markdown' },
}

const displayFilename = computed(() => {
  return props.filename.length > 30
    ? `${props.filename.substring(0, 27)}...`
    : props.filename
})

const highlightLang = computed(() => {
  const lang = detectLanguage.value.toLowerCase()
  return languageMap[lang]?.highlight || 'javascript'
})

const displayLang = computed(() => {
  const lang = props.lang.toLowerCase()
  return languageMap[lang]?.display || props.lang
})

const detectLanguage = computed(() => {
  const extension = props.filename.split('.').pop()?.toLowerCase()
  return extension && languageMap[extension] ? extension : props.lang
})

const highlightedCode = computed(() => {
  try {
    if (!codeContent.value)
      return ''
    const result = hljs.highlight(codeContent.value, { language: highlightLang.value })
    return result.value
  }
  catch (error) {
    console.error('Highlighting error:', error)
    return codeContent.value
  }
})

function toggleEdit() {
  isEditing.value = !isEditing.value

  if (isEditing.value) {
    nextTick(() => {
      adjustTextareaHeight()
    })
  }
}

function copyToClipboard() {
  navigator.clipboard.writeText(codeContent.value)
  isCopied.value = true

  setTimeout(() => {
    isCopied.value = false
  }, 2000)
}

function adjustTextareaHeight() {
  if (textarea.value) {
    textarea.value.style.height = 'auto'
    textarea.value.style.height = `${textarea.value.scrollHeight}px`
    editorHeight.value = `${textarea.value.scrollHeight}px`
  }
}

watch(codeContent, (newValue) => {
  emit('update:code', newValue)
})

watch(() => props.code, (newValue) => {
  codeContent.value = newValue
})

onMounted(() => {
  if (isEditing.value && textarea.value) {
    adjustTextareaHeight()
  }
})
</script>

<template>
  <div class="code-preview-container rounded-lg shadow-xl overflow-hidden border border-myborder dark:border-myborder bg-gray-800 dark:bg-gray-900">
    <div class="flex items-center justify-between p-3 bg-mybg dark:bg-mybg border-b border-myborder">
      <div class="flex items-center">
        <div class="flex space-x-2 mr-4">
          <div class="w-3 h-3 bg-red-500 rounded-full" />
          <div class="w-3 h-3 bg-yellow-500 rounded-full" />
          <div class="w-3 h-3 bg-green-500 rounded-full" />
        </div>
        <div class="text-gray-300 text-sm font-mono truncate">
          {{ displayFilename }}
        </div>
      </div>

      <div class="flex items-center space-x-2">
        <button v-if="showEditButton" class="text-xs px-2 py-1 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded transition-colors" @click="toggleEdit">
          {{ isEditing ? 'Preview' : 'Edit' }}
        </button>
        <a
          v-if="customURL"
          :href="customURL"
          target="_blank"
          class="text-xs px-2 py-1 bg-blue-600 hover:bg-blue-500 text-white rounded flex items-center transition-colors"
        >
          <Icon name="tabler:external-link" class="mr-1 w-3 h-3" />
          View
        </a>
        <button
          v-if="showCopyButton"
          class="text-xs px-2 py-1 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded flex items-center transition-colors"
          :class="{ 'bg-green-700 hover:bg-green-600': isCopied }"
          @click="copyToClipboard"
        >
          <Icon :name="isCopied ? 'tabler:check' : 'tabler:copy'" class="mr-1 w-3 h-3" />
          {{ isCopied ? 'Copied!' : 'Copy' }}
        </button>
      </div>
    </div>

    <div class="relative">
      <textarea
        v-if="isEditing"
        ref="textarea"
        v-model="codeContent"
        class="w-full font-mono text-sm p-4 bg-[#282c34] text-gray-100 outline-none resize-none"
        :style="{ height: editorHeight }"
        @input="adjustTextareaHeight"
      />

      <div v-else class="code-block-wrapper bg-[#282c34] p-4 overflow-auto">
        <pre ref="codeBlock" class="hljs p-0 m-0 bg-transparent"><code v-html="highlightedCode" /></pre>
      </div>

      <div class="absolute bottom-2 right-2 text-xs text-gray-500 bg-[#282c34]/90 px-2 py-1 rounded-md">
        {{ languageMap[detectLanguage.value]?.display || displayLang }}
      </div>
    </div>
  </div>
</template>

<style>
/* Code block wrapper */
.code-block-wrapper {
  max-height: 500px;
}

/* Custom scrollbar */
.code-preview-container ::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

.code-preview-container ::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 4px;
}

.code-preview-container ::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 4px;
}

.code-preview-container ::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.3);
}

/* Ensure highlight.js styling is clean */
.hljs {
  background: transparent !important;
  padding: 0 !important;
  overflow: visible !important;
}
</style>
