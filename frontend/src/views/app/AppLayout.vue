<script setup>
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { apiRequest } from "../../services/api"
import { setToken } from "../../services/auth"

const router = useRouter()
const loggingOut = ref(false)

const navItems = computed(() => [
    { to: { name: 'diary-list' }, label: '日记' },
    { to: { name: 'diary-new' }, label: '写日记' },
    { to: { name: 'diaft-list' }, label: '草稿' },
    { to: { name: 'profile' }, label: '个人资料' },
    { to: { name: 'security' }, label: '安全设置' },
])

async function logout() {
    loggingOut.value = true
    try {
        await apiRequest('/auth/logout', { method: 'POST' }).catch(() => null)
    } finally {
        setToken('')
        loggingOut.value = false
        router.replace({ name: 'login' })
    }
}
</script>

<template>
    <div class="min-h-screen bg-base-100">
        <div class="cd-container py-[30px]">
            <div class="grid gap-[20px] md:grid-cols-[260px_1fr]">
                <aside class="cd-card bg-[#F3F3F3] p-[25px]">
                    <div class="flex items-center gap-[12px]">
                        <div class="h-[36px] w-[36px] rounded-full bg-white"></div>
                        <div class="cd-h4">CatDiary</div>
                    </div>

                    <nav class="mt-[25px] grid gap-[10px]">
                        <RouterLink v-for="item in navItems" :key="item.label" :to="item.to"
                            class="rounded-[14px] border border-[#191A23] bg-white px-[16px] py-[12px] cd-p">
                            {{ item.label }}
                        </RouterLink>
                    </nav>

                    <div class="mt-[25px]">
                        <button
                            class="w-full rounded-[14px] border border-[#191A23] bg-transparent px-[16px] py-[12px] cd-p disabled:opacity-60"
                            :disabled="loggingOut" type="button" @click="logout">
                            {{ loggingOut ? '退出中...' : '退出登录' }}
                        </button>
                    </div>
                </aside>

                <section class="cd-card bg-[#F3F3F3] p-[25px]">
                    <RouterView />
                </section>
            </div>
        </div>
    </div>
</template>
