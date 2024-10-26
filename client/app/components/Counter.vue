<script setup lang="ts">
import { useCount } from '~/store/count'

const counter = useCount()

const buttons = [
  { icon: 'ic:baseline-plus', fn: counter.increment },
  { icon: 'ic:baseline-minus', fn: counter.decrement },
]
</script>

<template>
  <div class="mt-8">
    <div class="flex flex-col items-center gap-4 dark:text-gray-100 text-gray-900">
      <p class="text-3xl font-medium">
        Counter: {{ counter.count }}
      </p>
      <div class="flex items-center gap-4">
        <button
          v-for="button in buttons" class="bg-neutral-800 border-none hover:bg-neutral-700 text-gray-100 px-4 py-2 rounded-lg transition-colors duration-200 flex items-center justify-center group"
          @click="button.fn"
        >
          <Icon :name="button.icon" class="w-6 h-6" />
        </button>
      </div>
    </div>
  </div>
</template>
