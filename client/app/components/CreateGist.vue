<script setup lang="ts">
const toast = useToast()

const gistTitle = ref('')
const gistFileName = ref('')
const gistContent = ref('')
const customURL = ref('')
// const gistPublic = ref(true)

async function createGist(gistPublic: boolean) {
  const BE = import.meta.env.VITE_BE_URL
  const multiPartData = new FormData()

  multiPartData.append('gist_title', gistTitle.value)
  multiPartData.append('is_public', String(gistPublic))
  customURL.value && multiPartData.append('custom_url', customURL.value)

  const file = new File([gistContent.value], gistFileName.value, { type: 'text/plain' })
  multiPartData.append('file', file)

  const response = await fetch(`${BE}/api/gists/new`, {
    method: 'POST',
    body: multiPartData,
    credentials: 'include',
  })

  const body = await response.json()

  return {
    status: response.status,
    bodyData: body,
  }
  /* example response

  {
    "created_at": "2025-01-12T17:29:02.133109+05:30",
    "fileSize": 73,
    "file_id": "01JHD63EZNXTVP81JAYDG8ND2V",
    "key": "01JBGCCRA3ZJK3WF4HC7K7WHYR/01JHD63EZNXTVP81JAYDG8ND2V",
    "message": "Gist uploaded successfully!",
    "short_url": "README.md-01JHD63EZNXTVP81JAYDG8ND2V",
    "success": true,
    "title": "CodeFLuke Readmeh"
  }

  */
}

async function handleCreateGist(gistPublic: boolean) {
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
      if (isValid) {
        element.classList.remove('border-red-500')
      }
      else {
        element.classList.add('border-red-500')
      }
    }
    else {
      console.error('Element not found')
    }
  }

  if (!safeCheck) {
    toast.add({
      title: 'Invalid Input',
      description: 'Please fill all the required fields correctly.',
      color: 'red',
      timeout: 3000,
    })
    return
  }

  const { status, bodyData } = await createGist(gistPublic)
  if (status !== 200) {
    if (status === 409) {
      return toast.add({
        title: 'Conflict',
        description: 'Gist with the same short URL already exists',
        color: 'red',
        timeout: 3000,
      })
    }

    toast.add({
      title: 'Failed to Create Gist',
      description: 'Something went wrong while creating the gist',
      color: 'red',
      timeout: 3000,
    })
  }

  if (bodyData.success) {
    toast.add({
      title: 'Gist Created',
      description: 'Your gist has been created successfully',
      color: 'green',
      timeout: 3000,
    })
  }
}
</script>

<template>
  <div class="flex flex-col gap-y-4 p-5 md:px-10 w-full max-h-[85vh] overflow-y-scroll mr-8 my-8 md:my-14">
    <h2 class="text-md md:text-xl font-mono tracking-wider underline-offset-2 dark:no-underline">
      Create New Gists
    </h2>
    <div class="flex flex-col justify-between gap-4 h-11">
      <UInput
        id="gistTitle"
        v-model="gistTitle"
        class="w-full"
        size="xl"
        placeholder="Gist Title*"
      />
      <div class="flex flex-row justify-between gap-2">
        <UInput
          id="gistFilename"
          v-model="gistFileName"
          class="w-full"
          size="xl"
          placeholder="File Name*"
        />
        <UInput
          v-model="customURL"
          class="w-full"
          size="xl"
          placeholder="CustomURL (Optional)"
        />
      </div>

      <UTextarea id="gistContent" v-model="gistContent" resize size="xl" placeholder="Zaabang your gist content here!!!" />
      <div class="flex flex-row justify-end pt-3 gap-2 h-11">
        <UButton
          icon="i-heroicons-pencil-square"
          size="sm"
          class="rounded-full"
          color="primary"
          variant="solid"
          label="Create Public Gist"
          :trailing="false"
          @click="handleCreateGist(true)"
        />
        <UButton
          icon="i-heroicons-pencil-square"
          cre size="sm"
          class="rounded-full"
          color="primary"
          variant="solid"
          label="Create Secret Gist"
          :trailing="false"
          @click="(handleCreateGist(false))"
        />
      </div>
      <UDivider class="pt-2 dark:pt-3 shadow-lg shadow-myLightBorder/30 dark:shadow-myborder/75" />
    </div>
  </div>
</template>
