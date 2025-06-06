<script setup>
// This page will automatically redirect to the login page
definePageMeta({
  layout: 'auth',
})

// Redirect to login page
onMounted(() => {
  const router = useRouter()
  router.replace('/auth/login')
})
</script>

<template>
  <div class="flex justify-center items-center h-full">
    <p>Redirecting to login...</p>
  </div>
</template>
