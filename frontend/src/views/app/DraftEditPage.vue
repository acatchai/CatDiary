<script setup>
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { apiRequest } from '../../services/api';
import { getToken } from '../../services/auth';

const route = useRoute()
const router = useRouter()

const id = computed(() => String(route.params.id || ''))

const loading = ref(false)
const saving = ref(false)
const flushing = ref(false)
const deleting = ref(false)
const errorText = ref('')
const hintText = ref('')

const draft = ref(null)
const version = ref(0)

const form = ref({
    title: '',
    content: '',
    mood: '',
    weather: '',
    location: '',
    occurredAtLocal: '',
})

let autosaveTimer = null
let lastSavedSnapshot = ''

function toLocalDatetimeValue(rfc3339) {
    if (!rfc3339) return ''
    const d = new Date(rfc3339)
    if (Number.isNaN(d.getTime())) return ''
    const pad = (n) => String(n).padStart(2, '0')
    const yyyy = d.getFullYear()
    const mm = pad(d.getMonth() + 1)
    const dd = pad(d.getDate())
    const hh = pad(d.getHours())
    const mi = pad(d.getMinutes())
    return `${yyyy}-${mm}-${dd}T${hh}:${mi}`
}

function snapshot() {
    return JSON.stringify({
        title: form.value.title,
        content: form.value.content,
        mood: form.value.mood,
        weather: form.value.weather,
        location: form.value.location,
        occurredAtLocal: form.value.occurredAtLocal,
        version: version.value,
    })
}

async function load() {
    loading.value = true
    errorText.value = ''
    hintText.value = ''
    try {
        const res = await apiRequest(`/drafts/${id.value}`)
        draft.value = res?.data || null
        version.value = Number(draft.value?.version) || 0
        form.value = {
            title: draft.value?.title || '',
            content: draft.value?.content || '',
            mood: draft.value?.mood || '',
            weather: draft.value?.weather || '',
            location: draft.value?.location || '',
            occurredAtLocal: toLocalDatetimeValue(draft.value?.occured_at),
        }
        lastSavedSnapshot = snapshot()
    } catch (e) {
        errorText.value = e?.message || '加载失败'
    } finally {
        loading.value = false
    }
}

function scheduleAutosave() {
    if (loading.value) return
    if (autosaveTimer) clearTimeout(autosaveTimer)
    autosaveTimer = setTimeout(() => {
        autosaveTimer = null
        autosave()
    }, 900)
}

async function autosave() {
    if (saving.value || loading.value) return
    const nowSnapshot = snapshot()
    if (nowSnapshot === lastSavedSnapshot) return
    if (!form.value.title.trim() || !form.value.content.trim()) return

    saving.value = true
    errorText.value = ''
    hintText.value = '自动保存中...'

    try {
        const body = {
            expected_version: version.value ? version.value : undefined,
            title: form.value.title,
            content: form.value.content,
            mood: form.value.mood,
            weather: form.value.weather,
            location: form.value.location,
        }
        if (form.value.occurredAtLocal) {
            const d = new Date(form.value.occurredAtLocal)
            if (!Number.isNaN(d.getTime())) body.occured_at = d.toISOString()
        }

        const res = await apiRequest(`/drafts/${id.value}`, { method: 'PUT', body: JSON.stringify(body) })
        draft.value = res?.data || draft.value
        version.value = Number(draft.value?.version) || version.value
        lastSavedSnapshot = snapshot()
        hintText.value = '已自动保存'
    } catch (e) {
        if (e?.status === 409 && e?.data?.current_version != null) {
            errorText.value = '草稿已在别处更新，已为你重新加载最新版本喵！'
            await load()
        } else {
            errorText.value = e?.message || '自动保存失败'
        }
        hintText.value = ''
    } finally {
        saving.value = false
    }
}

async function manualSave() {
    saving.value = true
    errorText.value = ''
    hintText.value = ''
    try {
        const body = {
            expected_version: version.value ? version.value : undefined,
            title: form.value.title,
            content: form.value.content,
            mood: form.value.mood,
            weather: form.value.weather,
            location: form.value.location,
        }
        if (form.value.occurredAtLocal) {
            const d = new Date(form.value.occurredAtLocal)
            if (!Number.isNaN(d.getTime())) body.occured_at = d.toISOString()
        }
        const res = await apiRequest(`/drafts/${id.value}`, { method: 'PUT', body: JSON.stringify(body) })
        draft.value = res?.data || draft.value
        version.value = Number(draft.value?.version) || version.value
        lastSavedSnapshot = snapshot()
        hintText.value = '已保存'
    } catch (e) {
        if (e?.status === 409 && e?.data?.current_version != null) {
            errorText.value = '草稿已在别处更新，已为你重新加载最新版本喵！'
            await load()
        } else {
            errorText.value = e?.message || '保存失败'
        }
    } finally {
        saving.value = false
    }
}

async function flushToDiary() {
    flushing.value = true
    errorText.value = ''
    hintText.value = ''
    try {
        await apiRequest(`/drafts/${id.value}/flush`, { method: 'POST' })
        router.replace({ name: 'diary-list' })
    } catch (e) {
        errorText.value = e?.message || '发布失败'
    } finally {
        flushing.value = false
    }
}

async function removeDraft() {
    deleting.value = true
    errorText.value = ''
    hintText.value = ''
    try {
        await apiRequest(`/drafts/${id.value}`, { method: 'DELETE' })
        router.replace({ name: 'draft-list' })
    } catch (e) {
        errorText.value = e?.message || '删除失败'
    } finally {
        deleting.value = false
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
            headers: token ? { Authorization: `Bearer ${token}` } : {},
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
            hintText.value = '上传接口没有返回可用 URL喵~'
            return
        }
        form.value.content = `${form.value.content}\n\n![](${url})\n`
        scheduleAutosave()
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

watch(
    () => [form.value.title, form.value.content, form.value.mood, form.value.weather, form.value.location, form.value.occurredAtLocal],
    () => scheduleAutosave(),
)

onMounted(load)
onBeforeUnmount(() => {
    if (autosaveTimer) clearTimeout(autosaveTimer)
})
</script>

<template>
    <div>
        <div class="flex flex-col gap-[10px] md:flex-row md:items-center md:justify-between">
            <div>
                <div class="cd-h2">编辑草稿</div>
                <div class="cd-p mt-[8px]">
                    <span v-if="version">版本：{{ version }}</span>
                    <span v-if="hintText" class="ml-[10px]">{{ hintText }}</span>
                </div>
            </div>
            <div class="flex flex-wrap gap-[10px]">
                <RouterLink class="cd-btn cd-btn-outline cd-p inline-flex items-center justify-center" to="/app/drafts">
                    返回草稿
                </RouterLink>
                <button class="cd-btn cd-btn-outline cd-p disabled:opacity-60" :disabled="saving || loading"
                    type="button" @click="manualSave">
                    {{ saving ? '保存中...' : '保存' }}
                </button>
                <button class="cd-btn cd-btn-primary cd-p disabled:opacity-60" :disabled="flushing || saving || loading"
                    type="button" @click="flushToDiary">
                    {{ flushing ? '发布中...' : '发布为日记' }}
                </button>
            </div>
        </div>

        <div v-if="errorText" class="mt-[20px] rounded-[14px] border border-[#191A23] bg-white p-[14px] cd-p">
            {{ errorText }}
        </div>

        <div v-if="loading" class="mt-[20px] rounded-[14px] border border-[#191A23] bg-white p-[18px] cd-p">
            加载中...
        </div>

        <form v-else class="mt-[20px] grid gap-[20px]" @submit.prevent="manualSave">
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
                <div class="cd-p text-[#191A23]">上传接口后端目前是 TODO，启用后可直接插入图片链接。</div>
            </div>

            <div class="grid gap-[20px] md:grid-cols-2">
                <label class="grid gap-[10px]">
                    <span class="cd-p">心情</span>
                    <input v-model="form.mood"
                        class="h-[54px] rounded-[14px] border border-[#191A23] bg-white px-[18px] outline-none"
                        maxlength="20" placeholder="比如：开心 / 平静 / 低落" />
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
                        maxlength="100" placeholder="比如：家里 / 公司 / 路上" />
                </label>
                <label class="grid gap-[10px]">
                    <span class="cd-p">发生时间（可选）</span>
                    <input v-model="form.occurredAtLocal"
                        class="h-[54px] rounded-[14px] border border-[#191A23] bg-white px-[18px] outline-none"
                        type="datetime-local" />
                </label>
            </div>

            <div class="flex flex-wrap gap-[10px]">
                <button class="cd-btn cd-btn-primary cd-p disabled:opacity-60" :disabled="saving" type="submit">
                    {{ saving ? '保存中...' : '保存' }}
                </button>
                <button class="cd-btn cd-btn-outline cd-p disabled:opacity-60" :disabled="deleting || saving"
                    type="button" @click="removeDraft">
                    {{ deleting ? '删除中...' : '删除草稿' }}
                </button>
            </div>
        </form>
    </div>
</template>
