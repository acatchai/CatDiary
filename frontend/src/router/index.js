import { createRouter, createWebHistory } from 'vue-router'

import LandingPage from '../views/LandingPage.vue'

const routes = [
    {
        path: '/',
        name: 'landing',
        component: LandingPage,
    }
]


const router = createRouter({
    history: createWebHistory(),
    routes,
    scrollBehavior() {
        return { top: 0 }
    },
})

router.beforeEach((to) => {
    if (!to.meta?.requireAuth) return true

    const token = getToken()
    if (token) return true

    const redirect = to.fullPath || '/app/diaries'
    return { name: 'login', query: { redirect } }

})

export default router
