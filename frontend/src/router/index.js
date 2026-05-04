import { createRouter, createWebHistory } from 'vue-router'

import LandingPage from '../views/LandingPage.vue'
import loginPage from '../views/LoginPage.vue'
import RegisterPage from '../views/RegisterPage.vue'
import AppLayout from '../views/app/AppLayout.vue'
import DiaryListPage from '../views/app/DiaryListPage.vue'
import DiaryNewPage from '../views/app/DiaryNewPage.vue'
import DiaryDetailPage from '../views/app/DiaryDetailPage.vue'
import DraftListPage from '../views/app/DraftListPage.vue'
import DraftEditPage from '../views/app/DraftEditPage.vue'
import ProfilePage from '../views/app/ProfilePage.vue'

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
    },
    {
        path: '/app',
        component: AppLayout,
        meta: { requiresAuth: true },
        children: [
            { path: '', redirect: { name: 'diary-list' } },
            {
                path: 'diaries',
                name: 'diary-list',
                component: DiaryListPage,
                meta: { requiresAuth: true }
            },
            {
                path: 'diaries/new',
                name: 'diary-new',
                component: DiaryNewPage,
                meta: { requiresAuth: true }
            },
            {
                path: 'diaries/:id',
                name: 'diary-detail',
                component: DraftDetailPage,
                meta: { requiresAuth: true },
            },
            {
                path: 'drafts',
                name: 'draft-list',
                component: DraftListPage,
                meta: { requiresAuth: true },
            },
            {
                path: 'drafts/:id',
                name: 'draft-edit',
                component: DraftEditPage,
                meta: { requiresAuth: true }
            },
            {
                path: 'profile',
                name: 'profile',
                component: ProfilePage,
                meta: { requiresAuth: true }
            }
        ]
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
