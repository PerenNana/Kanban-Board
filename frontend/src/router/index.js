import { createRouter, createWebHistory } from 'vue-router'
import BoardsView from '../views/BoardsView.vue'
import BoardView from '../views/BoardView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', name: 'boards', component: BoardsView },
    { path: '/boards/:id', name: 'board', component: BoardView, props: true },
  ],
})

export default router