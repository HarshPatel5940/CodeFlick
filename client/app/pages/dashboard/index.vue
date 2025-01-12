<template>
  <main class="w-screen h-screen dark:bg-mybg bg-myLightBg">
    <div class="flex flex-row pt-4 md:pt-5 justify-center">
      <Navbar />
    </div>
    <div class="w-screen h-screen bg-inherit flex flex-col md:flex-row dark:bg-inherit">
      <MyGists />
      <UDivider orientation="vertical" class="h-[75vh] self-center hidden md:block" size="sm" />
      <GlobalGists />
      <CreateGist />
    </div>
  </main>
</template>
