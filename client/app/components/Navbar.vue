<script setup lang="ts">
const toast = useToast()
const profile = useProfile()

async function signOut() {
  const BE = import.meta.env.VITE_BE_URL

  const response = await fetch(`${BE}/api/auth/logout`, {
    method: 'POST',
    credentials: 'include',
  })

  const body = await response.json()

  return {
    status: response.status,
    data: body,
  }
}

async function handleSignOut() {
  toast.add({
    title: 'Signing out!',
    color: 'orange',
    icon: 'mingcute:alert-line',
    timeout: 2000,
  })
  const res = await signOut()
  if (res.status !== 200) {
    toast.add({
      title: 'Failed to sign out!',
      color: 'red',
      icon: 'mingcute:alert-line',
      timeout: 2000,
    })
    return
  }
  profile.$reset()
  navigateTo('/')
  toast.add({
    title: 'Signed out!',
    color: 'green',
    icon: 'heroicons:check-circle',
    timeout: 2000,
  })
}

const items = [
  [{
    label: `${profile.data.name}`,
    slot: 'username',
    disabled: true,
  }],
  [{
    label: 'Your Gists',
    icon: 'heroicons:folder-open',
  }, {
    label: 'Your Stars',
    icon: 'heroicons:folder-open',
  }],
  [{
    label: 'Status',
    icon: 'i-heroicons-signal',
  }, {
    label: 'Settings',
    icon: 'i-heroicons-cog-8-tooth',
  }],
  [{
    label: 'Sign out',
    icon: 'i-heroicons-arrow-left-on-rectangle',
    onClick: handleSignOut,
  }],
]
</script>

<template>
  <div class="h-14 w-4/6 lg:w-5/6 mx-5 border-4 shadow-myLightBorder/50 dark:border-myborder dark:shadow-myborder/75 shadow-2xl rounded-full dark:bg-[#1F2938] flex flex-row items-center justify-between px-5 overflow-x-hidden">
    <div>
      <h1 class="text-lg md:text-2xl font-bold font-mono text-gray-900 dark:text-white">
        <a href="/">
          CodeFlick
        </a>
      </h1>
    </div>
    <div class="flex flex-row items-center h-full ">
      <UDropdown :items="items" :ui="{ item: { disabled: 'cursor-text select-text' } }" :popper="{ placement: 'bottom-start' }" class="bg-inherit">
        <Icon name="iconoir:profile-circle" class="text-green-500 w-7 h-7" />
        <template #username="{ item }">
          <div class="text-left">
            Signed in as <strong>{{ item.label }}</strong>
          </div>
        </template>

        <template #item="{ item }">
          <div :onclick="item.onClick" class="w-full flex  flex-row justify-between ">
            <span class="truncate">{{ item.label }}</span>

            <UIcon :name="item.icon" class="flex-shrink-0 h-4 w-4  text-gray-400 dark:text-gray-500 ms-auto" />
          </div>
        </template>
      </UDropdown>
      <div class="pt-1">
        <Theme />
      </div>
    </div>
  </div>
</template>
