import { createRouter, createWebHistory } from 'vue-router'

import LandingPage from '../views/LandingPage.vue'
import loginPage from '../views/LoginPage.vue'
import RegisterPage from '../views/RegisterPage.vue'

const routes = [
    {
        path: '/',
        name: 'landing',
        component: LandingPage,
    },
    {
        path: '/login',
        name: 'login',
        component: loginPage,
    },
    {
        path: '/register',
        name: 'register',
        component: RegisterPage,
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
