<script setup>
import { onMounted } from 'vue'
import Prism from '~/plugins/myPrism'

const gistTitle = ref('')
const gistDescription = ref('')
const gistFileName = ref('')
const gistContent = ref('')


onMounted(() => {
  Prism.highlightAll()
})
</script>

<template>
  <div class="flex flex-col gap-y-4 p-5 md:px-10 w-full max-h-[85vh] overflow-y-scroll mr-8 my-8 md:my-14">
    <h2 class="text-md md:text-xl font-mono tracking-wider underline-offset-2 dark:no-underline">
      Create New Gists
    </h2>
    <div class="flex flex-col justify-between gap-4 h-11">
      <div class="flex flex-row justify-between gap-2">
        <UInput
          v-model="gistTitle"
          class="w-full"
          size="xl"
          placeholder="Gist Title"
        />
        <UInput
          v-model="gistFileName"
          class="w-full"
          size="xl"
          placeholder="FileName"
        />
      </div>
      <UInput
        v-model="gistDescription"
        size="lg"
        placeholder="Gist Description"
        :ui="{ icon: { trailing: { pointer: '' } } }"
      >
        <template #trailing>
          <UButton
            v-show="gistDescription !== ''"
            color="gray"
            variant="link"
            icon="i-heroicons-x-mark-20-solid"
            :padded="false"
            @click="gistDescription = ''"
          />
        </template>
      </UInput>

      <UTextarea v-model="gistContent" resize size="lg" placeholder="Zaabang your gist content here!!!" />
      <div class="flex flex-row justify-end pt-3 gap-2 h-11">
        <UButton
          icon="i-heroicons-pencil-square"
          size="sm"
          class="rounded-full"
          color="primary"
          variant="solid"
          label="Create Public Gists"
          :trailing="false"
        />
        <UButton
          icon="i-heroicons-pencil-square"
          size="sm"
          class="rounded-full"
          color="primary"
          variant="solid"
          label="Create Secret Gists"
          :trailing="false"
        />
      </div>
      <UDivider class="pt-3 shadow-lg  dark:shadow-myborder/75" />
    </div>
  </div>
</template>
