<script setup>
import { ref } from 'vue'

const toast = useToast()
const myGistsRef = ref(null)
const publicGistsRef = ref(null)
const isMobile = ref(false)
const searchQuery = ref('')

onMounted(() => {
  checkScreenSize()
  window.addEventListener('resize', checkScreenSize)

  onUnmounted(() => {
    window.removeEventListener('resize', checkScreenSize)
  })
})

function checkScreenSize() {
  isMobile.value = window.innerWidth < 1024
}

function handleGistCreated(newGist) {
  toast.add({
    title: 'Success!',
    description: 'Your gist has been created.',
    color: 'green',
    timeout: 3000,
  })

  if (myGistsRef.value && typeof myGistsRef.value.addGist === 'function') {
    myGistsRef.value.addGist(newGist)
  }

  if (newGist.isPublic && publicGistsRef.value && typeof publicGistsRef.value.loadGists === 'function') {
    publicGistsRef.value.loadGists()
  }
}

function handleGistUpdated(updatedGist) {
  if (myGistsRef.value && typeof myGistsRef.value.updateGist === 'function') {
    myGistsRef.value.updateGist(updatedGist)
  }

  if (publicGistsRef.value && typeof publicGistsRef.value.loadGists === 'function') {
    publicGistsRef.value.loadGists()
  }
}

function handleSearch() {
  if (myGistsRef.value && typeof myGistsRef.value.searchGists === 'function') {
    myGistsRef.value.searchGists(searchQuery.value)
  }
}
</script>

<template>
  <main class="w-full min-h-screen dark:bg-mybg bg-myLightBg overflow-x-hidden">
    <div class="flex flex-row pt-4 md:pt-5 pb-5 md:pb-10 justify-center">
      <Navbar />
    </div>

    <div v-if="isMobile" class="container mx-auto px-4 flex flex-col">
      <div class="w-full mb-6">
        <DashGistCreateForm
          @gist-created="handleGistCreated"
          @gist-updated="handleGistUpdated"
        />
      </div>

      <div class="mb-4">
        <div class="flex">
          <UInput
            v-model="searchQuery"
            placeholder="Search my gists..."
            icon="i-heroicons-magnifying-glass"
            class="w-full"
            @keyup.enter="handleSearch"
          />
          <UButton
            color="primary"
            class="ml-2"
            icon="i-heroicons-magnifying-glass"
            @click="handleSearch"
          />
        </div>
      </div>

      <div class="w-full mb-6">
        <h2 class="text-xl font-mono tracking-wider mb-4 pl-2 border-l-4 border-green-500 dark:border-green-400">
          Public Gists
        </h2>
        <LazyDashGistList
          ref="publicGistsRef"
          fetch-url="/api/gists?fetchPublic=yes"
          empty-message="No public gists found"
          class="w-full"
        />
      </div>

      <div class="w-full mb-6">
        <h2 class="text-xl font-mono tracking-wider mb-4 pl-2 border-l-4 border-blue-500 dark:border-blue-400">
          My Gists
        </h2>
        <LazyDashGistList
          ref="myGistsRef"
          fetch-url="/api/gists"
          empty-message="You haven't created any gists yet"
          class="w-full"
        />
      </div>
    </div>

    <div v-else class="container mx-auto px-4">
      <div class="grid grid-cols-12 gap-6">
        <div class="col-span-12 lg:col-span-7 xl:col-span-8">
          <DashGistCreateForm
            @gist-created="handleGistCreated"
            @gist-updated="handleGistUpdated"
          />

          <div class="mt-8">
            <h2 class="text-xl font-mono tracking-wider mb-4 pl-2 border-l-4 border-green-500 dark:border-green-400">
              Public Gists
            </h2>
            <LazyDashGistList
              ref="publicGistsRef"
              fetch-url="/api/gists?fetchPublic=yes"
              empty-message="No public gists found"
              class="w-full"
              list-height="50vh"
            />
          </div>
        </div>

        <div class="col-span-12 lg:col-span-5 xl:col-span-4">
          <div class="mb-4">
            <div class="flex">
              <UInput
                v-model="searchQuery"
                placeholder="Search my gists..."
                icon="i-heroicons-magnifying-glass"
                class="w-full"
                @keyup.enter="handleSearch"
              />
              <UButton
                color="primary"
                class="ml-2"
                icon="i-heroicons-magnifying-glass"
                @click="handleSearch"
              />
            </div>
          </div>

          <h2 class="text-xl font-mono tracking-wider mb-4 pl-2 border-l-4 border-blue-500 dark:border-blue-400">
            My Gists
          </h2>
          <LazyDashGistList
            ref="myGistsRef"
            fetch-url="/api/gists"
            empty-message="You haven't created any gists yet"
            class="w-full"
            list-height="80vh"
          />
        </div>
      </div>
    </div>
  </main>
</template>
