<script setup>
import { onMounted, ref } from 'vue'
import { apiRequest } from '../../services/api'

const loading = ref(false)
const saving = ref(false)
const errorText = ref('')
const hintText = ref('')

const me = ref(null)
const form = ref({
    username: '',
    email: '',
    avatar: '',
})

async function load() {
    loading.value = true
    errorText.value = ''
    hintText.value = ''
    try {
        const res = await apiRequest('/users/me')
        me.value = res?.data || null
        form.value = {
            username: me.value?.username || '',
            email: me.value?.email || '',
            avatar: me.value?.avatar || '',
        }
    } catch (e) {
        errorText.value = e?.message || '加载失败'
    } finally {
        loading.value = false
    }
}

async function save() {
    saving.value = true
    errorText.value = ""
    hintText.value = ""
    try {
        const body = {}
        if (form.value.username && form.value.username !== me.value?.username) body.username = form.value.username
        if (form.value.email !== (me.value?.email || '')) body.email = form.value.email
        if (form.value.avatar !== (me.value?.avatar || '')) body.avatar = form.value.avatar

        const res = await apiRequest('/users/me', { method: 'PATCH', body: JSON.stringify(body) })
        me.value = res?.data || me.value
        hintText.value = res?.message || '已更新'
    } catch (e) {
        errorText.value = e?.message || '保存失败'
    } finally {
        saving.value = false
    }
}

onMounted(load)
</script>

<template>
    <div>
        <div class="flex flex-col gap-[10px] md:flex-row md:items-center md:justify-between">
            <div>
                <div class="cd-h2">个人资料</div>
                <div class="cd-p mt-[8px]">
                    <span v-if="hintText">{{ hintText }}</span>
                    <span v-else>将你的名字与头像变可爱~</span>
                </div>
            </div>
            <button class="cd-btn cd-btn-primary cd-p disabled:opacity-60" :disabled="saving || loading" type="button"
                @click="save">
                {{ saving ? '保存中...' : '保存' }}
            </button>
        </div>

        <div v-if="errorText" class="mt-[20px] rounded-[14px] border border-[#191A23] bg-white p-[14px] cd-p">
            {{ errorText }}
        </div>

        <div v-if="loading" class="mt-[20px] rounded-[14px] border border-[#191A23] bg-white p-[18px] cd-p">
            加载中...
        </div>

        <div v-else class="mt-[20px] cd-card bg-[#F3F3F3] p-[25px]">
            <div class="flex flex-col gap-[20px] md:flex-row md:items-start md:gap-[30px]">
                <div class="w-full md:w-[220px]">
                    <div class="rounded-[20px] border border-[#191A23] bg-white p-[18px]">
                        <div class="cd-p">头像预览</div>
                        <div
                            class="mt-[12px] h-[120px] w-[120px] overflow-hidden rounded-full border border-[#191A23] bg-[#F3F3F3]">
                            <img v-if="form.avatar" class="h-full w-full object-cover" :src="form.avatar" alt="" />
                        </div>
                    </div>
                </div>

                <form class="flex-1 grid gap-[20px]" @submit.prevent="save">
                    <label class="grid gap-[10px]">
                        <span class="cd-p">用户名</span>
                        <input v-model="form.username"
                            class="h-[54px] rounded-[14px] border border-[#191A23] bg-white px-[18px] outline-none"
                            maxlength="50" />
                    </label>
                    <label class="grid gap-[10px]">
                        <span class="cd-p">邮箱</span>
                        <input v-model="form.email"
                            class="h-[54px] rounded-[14px] border border-[#191A23] bg-white px-[18px] outline-none" />
                    </label>
                    <label class="grid gap-[10px]">
                        <span class="cd-p">头像 URL</span>
                        <input v-model="form.avatar"
                            class="h-[54px] rounded-[14px] border border-[#191A23] bg-white px-[18px] outline-none"
                            placeholder="https://..." />
                    </label>
                    <div class="cd-p text-[#191A23]">
                        <span v-if="me?.created_at">创建时间：{{ me.created_at }}</span>
                    </div>
                </form>
            </div>
        </div>
    </div>
</template>
