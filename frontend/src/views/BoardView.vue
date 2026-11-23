<template>
  <section v-if="board">
    <RouterLink to="/">← Alle Boards</RouterLink>
    <h1>{{ board.name }}</h1>

    <div class="board">
      <KanbanColumn
        v-for="column in board.columns"
        :key="column.id"
        :column="column"
      />
    </div>
  </section>

  <section v-else>
    <h1>Board nicht gefunden 🤷</h1>
    <RouterLink to="/">← Zurück zur Übersicht</RouterLink>
  </section>
</template>

<script>
import KanbanColumn from '@/components/KanbanColumn.vue'
import { dummyBoards } from '@/data/dummy'

export default {
  name: 'BoardView',
  components: { KanbanColumn },
  props: {
    id: { type: String, required: true },
  },
  data() {
    return {
      boards: structuredClone(dummyBoards),
    }
  },
  computed: {
    board() {
      return this.boards.find((b) => b.id === Number(this.id))
    },
  },
}
</script>

<style scoped>
.board {
  display: flex;
  gap: 1rem;
  align-items: flex-start;
  overflow-x: auto;
  padding-bottom: 1rem;
}
</style>