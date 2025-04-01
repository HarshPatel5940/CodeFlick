<script setup>
defineProps({
  icon: {
    type: String,
    default: 'heroicons:document',
  },
  title: {
    type: String,
    required: true,
  },
  description: {
    type: String,
    default: '',
  },
  actionLabel: {
    type: String,
    default: '',
  },
})

const emit = defineEmits(['action'])

function handleAction() {
  emit('action')
}
</script>

<template>
  <div class="flex flex-col items-center justify-center py-12 px-4 text-center">
    <Icon :name="icon" class="w-16 h-16 text-gray-400 dark:text-gray-500 mb-4" />
    <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">
      {{ title }}
    </h3>
    <p v-if="description" class="text-sm text-gray-500 dark:text-gray-400 max-w-md mb-6">
      {{ description }}
    </p>
    <UButton
      v-if="actionLabel"
      size="sm"
      color="primary"
      @click="handleAction"
    >
      {{ actionLabel }}
    </UButton>
  </div>
</template>
