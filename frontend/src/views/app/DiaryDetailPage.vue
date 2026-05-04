<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { apiRequest } from '../../services/api'

const route = useRoute()
const router = useRouter()

const id = computed(() => String(route.params.id || ''))

const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const errorText = ref('')

const diary = ref(null)
const editMode = ref(false)

const form = ref({
    title: '',
    content: '',
    mood: '',
    weather: '',
    location: '',
    occurredAtLocal: '',
})

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

async function load() {
    loading.value = true
    errorText.value = ''
    try {
        const res = await apiRequest(`/diaries/${id.value}`)
        diary.value = res?.data || null
        form.value = {
            title: diary.value?.title || '',
            content: diary.value?.content || '',
            mood: diary.value?.mood || '',
            weather: diary.value?.weather || '',
            location: diary.value?.location || '',
            occurredAtLocal: toLocalDatetimeValue(diary.value?.occurred_at || ''),
        }
    } catch (e) {
        errorText.value = e?.message || '加载失败'
    } finally {
        loading.value = false
    }
}

async function save() {
    saving.value = true
    errorText.value = ''
    try {
        const body = {
            title: form.value.title,
            content: form.value.content,
            mood: form.value.mood,
            weather: form.value.weather,
            location: form.value.location,
        }
        if (form.value.occurredAtLocal) {
            const d = new Date(form.value.occurredAtLocal)
            if (!Number.isNaN(d.getTime())) body.occurred_at = d.toISOString()
        }
        const res = await apiRequest(`/diaries/${id.value}`, { method: 'PUT', body: JSON.stringify(body) })
        diary.value = res?.data || diary.value
        editMode.value = false
    } catch (e) {
        errorText.value = e?.message || '保存失败'
    } finally {
        saving.value = false
    }
}

async function removeDiary() {
    deleting.value = true
    errorText.value = ''
    try {
        await apiRequest(`/diaries/${id.value}`, { method: 'DELETE' })
        router.replace({ name: 'diary-list' })
    } catch (e) {
        errorText.value = e?.message || '删除失败'
    } finally {
        deleting.value = false
    }
}

onMounted(load)
</script>

<template>
    <div>
        <div class="flex flex-col gap-[10px] md:flex-row md:items-center md:justify-between">
            <div>
                <div class="cd-h2">日记详情</div>
                <div class="cd-p mt-[8px]">将它打磨成你想留住的样子。</div>
            </div>
            <div class="flex flex-wrap gap-[10px]">
                <RouterLink class="cd-btn cd-btn-outline cd-p inline-flex items-center justify-center"
                    to="/app/diaries">返回列表
                </RouterLink>
                <button v-if="!editMode" class="cd-btn cd-btn-primary cd-p" type="button" :disabled="loading"
                    @click="editMode = true">
                    编辑
                </button>
            </div>
        </div>

        <div v-if="errorText" class="mt-[20px] rounded-[14px] border border-[#191A23] bg-white p-[14px] cd-p">
            {{ errorText }}
        </div>

        <div v-if="loading" class="mt-[20px] rounded-[14px] border border-[#191A23] bg-white p-[18px] cd-p">
            加载中...
        </div>

        <div v-else class="mt-[20px] rounded-[20px] border border-[#191A23] bg-white p-[20px]">
            <div v-if="!editMode">
                <div class="cd-h3">{{ diary?.title }}</div>
                <div class="cd-p mt-[12px] whitespace-pre-wrap">{{ diary?.content }}</div>
                <div class="cd-p mt-[18px] flex flex-wrap gap-[10px] text-[#191A23]">
                    <span v-if="diary?.mood" class="cd-chip cd-p px-[10px] py-[2px]">{{ diary?.mood }}</span>
                    <span v-if="diary?.weather" class="cd-chip cd-p px-[10px] py-[2px]">{{ diary?.weather }}</span>
                    <span v-if="diary?.location" class="cd-chip cd-p px-[10px] py-[2px]">{{ diary?.location }}</span>
                    <span v-if="diary?.occurred_at">发生时间：{{ diary?.occurred_at }}</span>
                </div>
            </div>

            <form v-else class="grid gap-[20px]" @submit.prevent="save">
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

                <div class="grid gap-[20px] md:grid-cols-2">
                    <label class="grid gap-[10px]">
                        <span class="cd-p">心情</span>
                        <input v-model="form.mood"
                            class="h-[54px] rounded-[14px] border border-[#191A23] bg-white px-[18px] outline-none"
                            maxlength="20" />
                    </label>
                    <label class="grid gap-[10px]">
                        <span class="cd-p">天气</span>
                        <input v-model="form.weather"
                            class="h-[54px] rounded-[14px] border border-[#191A23] bg-white px-[18px] outline-none"
                            maxlength="20" />
                    </label>
                    <label class="grid gap-[10px]">
                        <span class="cd-p">地点</span>
                        <input v-model="form.location"
                            class="h-[54px] rounded-[14px] border border-[#191A23] bg-white px-[18px] outline-none"
                            maxlength="100" />
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
                    <button class="cd-btn cd-btn-outline cd-p" type="button" :disabled="saving"
                        @click="editMode = false">
                        取消
                    </button>
                    <button class="cd-btn cd-btn-outline cd-p disabled:opacity-60" type="button"
                        :disabled="deleting || saving" @click="removeDiary">
                        {{ deleting ? '删除中...' : '删除' }}
                    </button>
                </div>
            </form>
        </div>
    </div>
</template>
