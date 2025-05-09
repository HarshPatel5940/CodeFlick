<script setup lang="ts">
const props = defineProps({
  initialData: {
    type: Object,
    default: () => ({
      title: '',
      fileName: '',
      content: '',
      customUrl: '',
    }),
  },
  isEdit: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['gist-created', 'gist-updated'])
const toast = useToast()
const BE = import.meta.env.VITE_BE_URL

const gistTitle = ref(props.initialData.title)
const gistFileName = ref(props.initialData.fileName)
const gistContent = ref(props.initialData.content)
const customURL = ref(props.initialData.customUrl)
const isSubmitting = ref(false)

async function createGist(gistPublic: boolean) {
  isSubmitting.value = true

  try {
    const multiPartData = new FormData()

    multiPartData.append('gist_title', gistTitle.value)
    multiPartData.append('is_public', String(gistPublic))
    customURL.value && multiPartData.append('custom_url', customURL.value)
    gistFileName.value && multiPartData.append('file_name', gistFileName.value)

    const file = new File([gistContent.value], gistFileName.value, { type: 'text/plain' })
    multiPartData.append('file', file)

    const endpoint = props.isEdit
      ? `${BE}/api/gists/update/${props.initialData.id}`
      : `${BE}/api/gists/new`

    const response = await fetch(endpoint, {
      method: props.isEdit ? 'PUT' : 'POST',
      body: multiPartData,
      credentials: 'include',
    })

    const body = await response.json()

    return {
      status: response.status,
      bodyData: body,
    }
  }
  catch (error) {
    console.error('Error creating gist:', error)
    return {
      status: 500,
      bodyData: { success: false, message: 'An unexpected error occurred' },
    }
  }
  finally {
    isSubmitting.value = false
  }
}

function validateForm() {
  const fields = [
    { id: 'gistTitle', value: gistTitle.value, minLength: 5 },
    { id: 'gistFilename', value: gistFileName.value, minLength: 5 },
    { id: 'gistContent', value: gistContent.value, minLength: 1 },
  ]

  let safeCheck = true

  for (const field of fields) {
    const isValid = field.value && field.value.length >= field.minLength
    safeCheck = Boolean(safeCheck && isValid)

    const element = document.getElementById(field.id)
    if (element) {
      element.classList.toggle('border-red-500', !isValid)
    }
  }

  return safeCheck
}

async function handleCreateGist(gistPublic: boolean) {
  if (!validateForm()) {
    toast.add({
      title: 'Invalid Input',
      description: 'Please fill all the required fields correctly.',
      color: 'red',
      timeout: 3000,
    })
    return
  }

  const { status, bodyData } = await createGist(gistPublic)

  if (status !== 201) {
    if (status === 409) {
      toast.add({
        title: 'Conflict',
        description: 'Gist with the same short URL already exists',
        color: 'red',
        timeout: 3000,
      })
      return
    }

    toast.add({
      title: props.isEdit ? 'Failed to Update Gist' : 'Failed to Create Gist',
      description: 'Something went wrong',
      color: 'red',
      timeout: 3000,
    })
    return
  }

  toast.add({
    title: props.isEdit ? 'Gist Updated!' : 'Gist Created!',
    color: 'green',
    icon: 'heroicons:check-circle',
    timeout: 2000,
  })

  if (!bodyData.data.updatedAt) {
    bodyData.data.updatedAt = new Date().toISOString()
  }
  if (!bodyData.data.viewCount) {
    bodyData.data.viewCount = 0
  }

  emit(props.isEdit ? 'gist-updated' : 'gist-created', bodyData.data)

  if (!props.isEdit) {
    gistTitle.value = ''
    gistFileName.value = ''
    gistContent.value = ''
    customURL.value = ''

    navigateTo(gistPublic
      ? `/${bodyData.data.shortUrl}`
      : `/${bodyData.data.shortUrl}?gid=${bodyData.data.fileId}`)
  }
}

function handleFileUpload(event: Event): void {
  const target = event.target as HTMLInputElement
  const files = target.files

  if (files && files.length > 0) {
    const file = files[0]
    if (file) {
      const reader = new FileReader()
      reader.onload = (e: ProgressEvent<FileReader>) => {
        gistContent.value = e.target?.result as string
        gistFileName.value = file.name
      }
      reader.readAsText(file)
    }
  }
}
</script>

<template>
  <div class="flex flex-col gap-y-4 p-5 w-full max-h-[85vh] overflow-y-auto bg-white/10 dark:bg-gray-800/20 rounded-lg">
    <h2 class="text-lg font-mono tracking-wider underline-offset-2 dark:no-underline pb-2 border-b border-gray-200 dark:border-gray-700">
      {{ isEdit ? 'Edit Gist' : 'Create New Gist' }}
    </h2>
    <div class="flex flex-col justify-between gap-4">
      <UInput
        id="gistTitle"
        v-model="gistTitle"
        class="w-full"
        size="xl"
        placeholder="Gist Title*"
        :required="true"
        :disabled="isSubmitting"
      />
      <div class="flex flex-col md:flex-row justify-between gap-2">
        <UInput
          id="gistFilename"
          v-model="gistFileName"
          class="w-full"
          size="xl"
          placeholder="File Name*"
          :required="true"
          :disabled="isSubmitting"
        />
        <UInput
          id="customURL"
          v-model="customURL"
          class="w-full"
          size="xl"
          placeholder="CustomURL (Optional)"
          :required="false"
          :disabled="isSubmitting"
        />
      </div>

      <UTextarea
        id="gistContent"
        v-model="gistContent"
        resize
        size="xl"
        placeholder="Type or paste your code here..."
        :required="true"
        :disabled="isSubmitting"
        class="min-h-[100px]"
      />
      <div class="flex flex-row justify-between gap-2">
        <div class="relative">
          <input
            id="fileInput"
            type="file"
            accept=".txt"
            class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
            :disabled="isSubmitting"
            @change="handleFileUpload"
          >
          <UButton
            icon="i-heroicons-arrow-up-on-square"
            size="sm"
            class="rounded-full"
            color="gray"
            variant="outline"
            label="Upload File"
            :disabled="isSubmitting"
          />
        </div>

        <div class="flex flex-row justify-end gap-2">
          <UButton
            icon="i-heroicons-lock-closed"
            size="sm"
            class="rounded-full"
            color="gray"
            variant="solid"
            :label="isEdit ? 'Update Secret' : 'Create Secret'"
            :trailing="false"
            :loading="isSubmitting"
            :disabled="isSubmitting"
            @click="handleCreateGist(false)"
          />
          <UButton
            icon="i-heroicons-globe-alt"
            size="sm"
            class="rounded-full"
            color="primary"
            variant="solid"
            :label="isEdit ? 'Update Public' : 'Create Public'"
            :trailing="false"
            :loading="isSubmitting"
            :disabled="isSubmitting"
            @click="handleCreateGist(true)"
          />
        </div>
      </div>
    </div>
  </div>
</template>
