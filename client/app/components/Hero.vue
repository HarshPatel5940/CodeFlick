<script setup lang="ts">
const icons = [
  {
    name: 'nuxt',
    iconName: 'logos:nuxt-icon',
    label: 'Nuxt',
  },
  {
    name: 'pinia',
    iconName: 'logos:pinia',
    label: 'Pinia',
  },
  {
    name: 'drizzle',
    iconName: 'simple-icons:drizzle',
    label: 'Drizzle ORM',
  },
  {
    name: 'biome',
    iconName: 'devicon:biome',
    label: 'Biome',
  },
  {
    name: 'docker',
    iconName: 'logos:docker-icon',
    label: 'Docker',
  },
  {
    name: 'unocss',
    iconName: 'logos:unocss',
    label: 'UnoCSS',
  },
]
</script>

<template>
  <section
    class="bg-inherit w-full h-full justify-center overflow-auto items-center flex flex-col"
  >
    <div class="max-w-6xl mx-auto flex flex-col items-center">
      <div class="flex flex-col items-center gap-4">
        <Icon name="logos:nuxt-icon" class="text-green-500 w-16 h-16" />
        <h1 class="dark:text-gray-100 text-gray-900 text-4xl font-bold">
          Nuxt Starter
        </h1>
      </div>
      <div
        class="grid grid-cols-1 text-gray-900 dark:text-gray-100 md:grid-cols-3 gap-4 w-full max-w-4xl"
      >
        <div
          v-for="icon in icons"
          class="dark:bg-neutral-800/50 bg-gray-100 p-6 rounded-xl border-2 border-neutral-700/50 hover:border-green-500"
        >
          <div :key="icon.name" class="flex flex-col items-center gap-3">
            <Icon
              :name="icon.iconName"
              :class="icon.name === 'drizzle' ? 'text-yellow-400' : ''"
              class="w-12 h-12 group-hover:scale-110"
            />
            <span class="font-medium">{{ icon.label }}</span>
          </div>
        </div>
      </div>
      <Counter />
    </div>
  </section>
</template>
