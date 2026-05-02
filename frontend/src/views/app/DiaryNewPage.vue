<script setup>
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { apiRequest } from '../../services/api'
import { getToken } from '../../services/auth'

const router = useRouter()

const form = ref({
    title: '',
    content: '',
    mood: '',
    weather: '',
    location: '',
    occurredAtLocal: '',
})

const saving = ref(false)
const errorText = ref('')
const hintText = ref('')

const occurredAtRFC3339 = computed(() => {
    const v = form.value.occurredAtLocal
    if (!v) return ''
    const d = new Date(v)
    if (Number.isNaN(d.getTime())) return ''
    return d.toISOString()
})

async function submit() {
    errorText.value = ''
    hintText.value = ''
    saving.value = true
    try {
        const body = {
            title: form.value.title,
            content: form.value.content,
            mood: form.value.mood,
            weather: form.value.weather,
            location: form.value.location,
        }
        if (occurredAtRFC3339.value) body.occurredAt = occurredAtRFC3339.value

        await apiRequest('/diaries', { method: 'POST', body: JSON.stringify(body) })
        router.replace({ name: 'diary-list' })
    } catch (e) {
        errorText.value = e?.message || '保存失败'
    } finally {
        saving.value = false
    }
}

async function insertImageFromFile(file) {
    errorText.value = ''
    hintText.value = ''

    const token = getToken()
    const formData = new FormData()
    formData.append('file', file)

    try {
        const res = await fetch('/api/v1/uploads', {
            method: 'POST',
            headers: token ? { 'Authorization': `Bearer ${token}` } : {},
            body: formData,
        })
        if (!res.ok) {
            const t = await res.text().catch(() => '')
            throw new Error(t || `上传失败 (${res.status})`)
        }
        const text = await res.text().catch(() => '')
        let data = null
        try {
            data = JSON.parse(text)
        } catch {
            data = null
        }
        const url = data?.url || data?.data?.url
        if (!url) {
            hintText.value = '上传接口未返回可用 URL'
            return
        }
        form.value.content = `${form.value.content}\n\n!(${url})\n`
    } catch (e) {
        errorText.value = e?.message || '上传失败'
    }
}

function onPickImage(e) {
    const file = e?.target?.files?.[0]
    e.target.value = ''
    if (!file) return
    insertImageFromFile(file)
}
</script>

<template>
    <div>
        <div class="flex flex-col gap-[10px] md:flex-row md:items-center md:justify-between">
            <div>
                <div class="cd-h2">写日记</div>
                <div class="cd-p mt-[8px]">把今天交给文字与猫。</div>
            </div>
            <RouterLink class="cd-btn cd-btn-outline cd-p inline-flex items-center justify-center" to="/app/diaries">
                返回列表
            </RouterLink>
        </div>

        <form class="mt-[25px] grid gap-[20px]" @submit.prevent="submit">
            <label class="grid gap-[10px]">
                <span class="cd-p">标题</span>
                <input v-model="form.title"
                    class="h-[54px] rounded-[14px] border border-[#191A23] bg-white px-[18px] outline-none"
                    maxlength="100" required />
            </label>

            <label class="grid gap-[10px]">
                <span class="cd-p">正文</span>
                <textarea v-model="form.content"
                    class="min-h-[220px] rounded-[14px] border border-[#191A23] bg-white px-[18px] py-[14px] outline-none"
                    required></textarea>
            </label>

            <div class="flex flex-wrap items-center gap-[10px]">
                <label class="cd-btn cd-btn-outline cd-p inline-flex cursor-pointer items-center justify-center">
                    插入图片
                    <input class="hidden" accept="image/*" type="file" @change="onPickImage" />
                </label>
                <div class="cd-p text-[#191A23]">
                    <span v-if="hintText">{{ hintText }}</span>
                    <span v-else>上传接口后端目前是 TODO，启用后可直接插入图片链接。</span>
                </div>
            </div>

            <div class="grid gap-[20px] md:grid-cols-2">
                <label class="grid gap-[10px]">
                    <span class="cd-p">Mood</span>
                    <input v-model="form.mood"
                        class="h-[54px] rounded-[14px] border border-[#191A23] bg-white px-[18px] outline-none"
                        maxlength="20" placeholder="开心 / 平静 / 低落" />
                </label>
                <label class="grid gap-[10px]">
                    <span class="cd-p">天气</span>
                    <input v-model="form.weather"
                        class="h-[54px] rounded-[14px] border border-[#191A23] bg-white px-[18px] outline-none"
                        maxlength="20" placeholder="晴 / 阴 / 雨" />
                </label>
                <label class="grid gap-[10px]">
                    <span class="cd-p">地点</span>
                    <input v-model="form.location"
                        class="h-[54px] rounded-[14px] border border-[#191A23] bg-white px-[18px] outline-none"
                        maxlength="100" placeholder="北京 / 上海 / 深圳">
                </label>
                <label class="grid gap-[10px]">
                    <span class="cd-p">记录时间(可选)</span>
                    <input v-model="form.occurredAtLocal"
                        class="h-[54px] rounded-[14px] border border-[#191A23] bg-white px-[18px] outline-none"
                        type="datetime-local" />
                </label>
            </div>

            <div v-if="errorText" class="rounded-[14px] border border-[#191A23] bg-white p-[14px] cd-p">
                {{ errorText }}
            </div>

            <div class="flex gap-[15px]">
                <button class="cd-btn cd-btn-primary cd-p disabled:opacity-60" :disabled="saving" type="submit">
                    {{ saving ? '保存中...' : '保存' }}
                </button>
                <button class="cd-btn cd-btn-outline cd-p" type="button"
                    @click="form.value = { title: '', content: '', mood: '', weather: '', location: '', occurredAtLocal: '' }">
                    清空
                </button>
            </div>
        </form>
    </div>
</template>
