<script setup>
const props = defineProps({
  endVal: {
    type: Number,
    required: true,
  },
  duration: {
    type: Number,
    default: 2000,
  },
  prefix: {
    type: String,
    default: '',
  },
  suffix: {
    type: String,
    default: '',
  },
})

const counterRef = ref(null)
const displayValue = ref(0)
const started = ref(false)

onMounted(() => {
  const observer = new IntersectionObserver((entries) => {
    for (const entry of entries) {
      if (entry.isIntersecting && !started.value) {
        started.value = true
        animateCount()
      }
    }
  }, { threshold: 0.1 })

  if (counterRef.value) {
    observer.observe(counterRef.value)
  }
})

function animateCount() {
  const startTime = Date.now()
  const endTime = startTime + props.duration

  function update() {
    const now = Date.now()
    const progress = Math.min(1, (now - startTime) / props.duration)

    displayValue.value = Math.floor(progress * props.endVal)

    if (now < endTime) {
      requestAnimationFrame(update)
    }
    else {
      displayValue.value = props.endVal
    }
  }

  update()
}

const formattedValue = computed(() => {
  return props.prefix + displayValue.value + props.suffix
})
</script>

<template>
  <span ref="counterRef">{{ formattedValue }}</span>
</template>
