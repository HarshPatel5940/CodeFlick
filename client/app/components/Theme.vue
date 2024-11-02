<script setup lang="ts">
const colorMode = useColorMode()
const theme = useTheme()

function toggleColorMode() {
  colorMode.value = theme.toggleTheme()
}

onMounted(() => {
  colorMode.value = theme.getTheme()
})
</script>

<template>
  <button class="p-2 bg-inherit border-none" @click="toggleColorMode">
    <Icon v-if="theme.getThemeBool()" name="solar:moon-bold" class="p-1 text-white" />
    <Icon v-if="!theme.getThemeBool()" name="solar:sun-bold" class="p-1 text-black" />
  </button>
</template>
