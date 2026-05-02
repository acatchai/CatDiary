<script setup>
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { apiRequest } from '../services/api'
import { setToken } from '../services/auth'

const route = useRoute()
const router = useRouter()

const form = ref({ username: '', password: '' })
const errorText = ref('')
const loading = ref(false)

const redirectTo = computed(() => {
    const q = route.query?.redirect
    return typeof q === 'string' && q.trim() !== '' ? q : '/app/diaries'
})

async function submit() {
    errorText.value = ''
    loading.value = true
    try {
        const res = await apiRequest('/auth/login', {
            method: 'POST',
            body: JSON.stringify({
                username: form.value.username,
                password: form.value.password,
            })
        })
        setToken(res?.token || '')
        await router.replace(redirectTo.value)
    } catch (e) {
        errorText.value = e?.message || '登录失败'
    } finally {
        loading.value = false
    }
}
</script>

<template>
    <div class="min-h-screen bg-base-100">
        <div class="cd-container py-[60px]">
            <div class="mx-auto max-w-[520px] cd-card bg-[#F3F3F3] p-[50px]">
                <div class="flex items-center justify-between">
                    <div class="cd-h2">登录</div>
                    <RouterLink class="cd-p underline" to="/">返回首页</RouterLink>
                </div>

                <form class="mt-[30px] grid gap-[20px]" @submit.prevent="submit">
                    <label class="grid gap-[10px]">
                        <span class="cd-p">用户名</span>
                        <input v-model="form.username"
                            class="h-[54px] rounded-[14px] border border-[#191A23] bg-white px-[18px] outline-none"
                            autocomplete="username" />
                    </label>

                    <label class="grid gap-[10px]">
                        <span class="cd-p">密码</span>
                        <input v-model="form.password"
                            class="h-[54px] rounded-[14px] border border-[#191A23] bg-white px-[18px] outline-none"
                            type="password" autocomplete="current-password" />
                    </label>

                    <div v-if="errorText" class="rounded-[14px] border border-[#191A23] bg-white p-[14px] cd-p">
                        {{ errorText }}
                    </div>

                    <button class="cd-btn cd-btn-primary cd-p disabled:opacity-60" :disabled="loading" type="submit">
                        {{ loading ? '登录中...' : '登录' }}
                    </button>

                    <div class="cd-p">
                        还没有账号？
                        <RouterLink class="underline" to="/register">去注册</RouterLink>
                    </div>

                </form>
            </div>
        </div>

    </div>

</template>
