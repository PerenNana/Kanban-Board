<template>
  <div class="card">
    <strong>{{ card.title }}</strong>
    <p v-if="card.description">{{ card.description }}</p>
  </div>
</template>

<script>
export default {
  name: 'KanbanCard',
  props: {
    card: { type: Object, required: true },
  },
}
</script>

<style scoped>
.card {
  background: #fff;
  border-radius: 6px;
  padding: 0.6rem 0.75rem;
  margin: 0.5rem 0;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.15);
}
.card:hover {
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}
.card:active {
  cursor: grabbing;
}
p {
  margin: 0.25rem 0 0;
  color: #5e6c84;
  font-size: 0.9rem;
}
</style>