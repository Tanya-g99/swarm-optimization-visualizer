import { createRouter, createWebHistory } from 'vue-router'

const links = [
  {
    path: '/info',
    name: 'Информация',
    component: () => import('views/InfoPage.vue')
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/info'
    },
    {
      path: '/',
      name: 'root',
      component: () => import('views/MainLayout.vue'),
      children: [
        ...links,
        {
          path: '/algo/:algoName',
          name: 'Algo',
          component: () => import('views/Algo.vue')
        },]
    }
  ]
})

export default { router, links }
