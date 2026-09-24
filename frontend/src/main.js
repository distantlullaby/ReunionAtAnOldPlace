import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import Feed from './views/Feed.vue'
import Profile from './views/Profile.vue'
import './style.css'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'feed', component: Feed },
    { path: '/profile', name: 'profile', component: Profile }
  ]
})

createApp(App).use(router).mount('#app')
