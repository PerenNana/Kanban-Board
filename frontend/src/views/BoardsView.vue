<template>
  <section>
    <h1>📋 Boards</h1>

    <ul v-if="boards.length > 0" class="board-list">
      <li v-for="board in boards" :key="board.id">
        <RouterLink :to="`/boards/${board.id}`">{{ board.name }}</RouterLink>
        <span class="meta">{{ board.columns.length }} Spalten</span>
      </li>
    </ul>

    <p v-else>Noch keine Boards vorhanden.</p>
  </section>
</template>

<script>
import { dummyBoards } from '@/data/dummy'

export default {
  name: 'BoardsView',
  data() {
    return {
      boards: structuredClone(dummyBoards),
    }
  },
}
</script>

<style scoped>
.board-list {
  list-style: none;
  padding: 0;
  display: grid;
  gap: 0.75rem;
  max-width: 480px;
}
.board-list li {
  background: #fff;
  border-radius: 8px;
  padding: 1rem;
  display: flex;
  justify-content: space-between;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.15);
}
.meta {
  color: #5e6c84;
  font-size: 0.9rem;
}
</style>