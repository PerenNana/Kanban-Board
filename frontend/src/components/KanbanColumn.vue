<template>
  <div class="column">
    <h2>
      {{ column.name }}
      <span class="count">{{ column.cards.length }}</span>
    </h2>

    <VueDraggable v-model="column.cards" group="cards" :animation="150" class="cards">
      <KanbanCard v-for="card in column.cards" :key="card.id" :card="card" />
    </VueDraggable>
  </div>
</template>

<script>
import { VueDraggable } from 'vue-draggable-plus'
import KanbanCard from './KanbanCard.vue'

export default {
  name: 'KanbanColumn',
  components: { KanbanCard, VueDraggable },
  props: {
    column: { type: Object, required: true },
  },
}
</script>

<style scoped>
.column {
  background: #ebecf0;
  border-radius: 8px;
  padding: 0.75rem;
  width: 280px;
  flex-shrink: 0;
}
h2 {
  font-size: 1rem;
  margin: 0 0 0.5rem;
  display: flex;
  justify-content: space-between;
}
.cards {
  min-height: 2.5rem;
}
.card {
  cursor: grab;
}
.count {
  background: #dfe1e6;
  border-radius: 10px;
  padding: 0 0.5rem;
  font-size: 0.85rem;
}
.empty {
  color: #5e6c84;
  font-style: italic;
  font-size: 0.9rem;
}
</style>