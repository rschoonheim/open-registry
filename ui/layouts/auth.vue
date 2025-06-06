<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50">
    <div class="max-w-md w-full space-y-8 p-6 bg-white shadow-md rounded-lg">
      <div class="text-center">
        <h2 class="mt-6 text-3xl font-extrabold text-gray-900">
          <slot name="title">Open Registry - Authentication</slot>
        </h2>
        <p class="mt-2 text-sm text-gray-600">
          <slot name="description"></slot>
        </p>
      </div>

      <!-- Main content area for login/register forms -->
      <NuxtPage />

      <!-- Footer area for links -->
      <div class="mt-6 text-center text-sm">
      </div>
    </div>
  </div>
</template>

