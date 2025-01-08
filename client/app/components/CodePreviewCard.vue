<script setup>
import { onMounted } from 'vue'
import Prism from '~/plugins/myPrism'

// defineProps({
//   title: {
//     type: String,
//     required: true,
//   },
//   filename: {
//     type: String,
//     required: true,
//   },
//   customURL: {
//     type: String,
//     required: true,
//   },
//   lang: {
//     type: String,
//     default: 'js',
//     required: true,
//   },
// })

onMounted(() => {
  Prism.highlightAll()
})

</script>

<template>
  <div class="w-full h-full">
  <pre><code class="lang-js dark:bg-inherit"><span>const x</span> = <span>'hello world'</span></code></pre>
    <div class="w-full h-full bg-inherit dark:bg-inherit">
      <div class="flex flex-row justify-between items-center">
        <h1 class="text-2xl font-bold text-gray-800 dark:text-gray-200">{{ title }}</h1>
        <!-- <a :href="fileURL" target="_blank" class="text-blue-500 dark:text-blue-300">View</a> -->
      </div>

    </div>
  </div>

</template>
