<template>
  <header class="topbar">
    <RouterLink to="/" class="logo">🗂️ Kanban</RouterLink>
  </header>

  <main class="content">
    <RouterView />
  </main>
</template>

<script>
export default {
  name: 'App',
}
</script>

<style scoped>
.topbar {
  background: #0b57d0;
  padding: 0.75rem 1.5rem;
}
.logo {
  color: #fff;
  font-weight: 700;
  font-size: 1.2rem;
  text-decoration: none;
}
.content {
  padding: 1.5rem;
}
</style>