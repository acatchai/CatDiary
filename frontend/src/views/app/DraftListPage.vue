<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { apiRequest } from '../../services/api'

const router = useRouter()

const loading = ref(false)
const errorText = ref('')

const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const items = ref([])
const creating = ref(false)

const totalPages = computed(() => {
    const t = Number(total.value) || 0
    const s = Number(pageSize.value) || 20
    return Math.max(1, Math.ceil(t / s))
})

async function load() {
    loading.value = true
    errorText.value = ''
    try {
        const res = await apiRequest(`/drafts?page=${page.value}&page_size=${pageSize.value}`)
        const data = res?.data || {}
        items.value = Array.isArray(data?.items) ? data.items : []
        total.value = Number(data?.total) || 0
    } catch (e) {
        errorText.value = e?.message || '加载失败'
    } finally {
        loading.value = false
    }
}

async function createDraft() {
    creating.value = true
    errorText.value = ''
    try {
        const res = await apiRequest('/drafts', {
            method: 'POST',
            body: JSON.stringify({
                title: '未命名草稿',
                content: '今天想记录点什么？',
                mood: '',
                weather: '',
                location: '',
            })
        })
        const id = res?.data?.id
        if (id != null) {
            router.push({ name: 'draft-edit', params: { id: String(id) } })
            return
        }
        await load()
    } catch (e) {
        errorText.value = e?.message || '创建失败'
    } finally {
        creating.value = false
    }
}

function prevPage() {
    if (page.value <= 1) return
    page.value -= 1
    load()
}

function nextPage() {
    if (page.value >= totalPages.value) return
    page.value += 1
    load()
}

onMounted(load)
</script>

<template>
    <div>
        <div class="flex flex-col gap-[10px] md:flex-row md:items-center md:justify-between">
            <div>
                <div class="cd-h2">草稿</div>
                <div class="cd-p mt-[8px]">别担心写不完哦，草稿会替你留着一个小尾巴，等你回来再续上～</div>
            </div>
            <button class="cd-btn cd-btn-primary cd-p disabled:opacity-60" :disabled="creating" type="button"
                @click="createDraft">
                {{ creating ? '创建中...' : '新建草稿' }}
            </button>
        </div>
        <div v-if="errorText" class="mt-[20px] rounded-[14px] border border-[#191A23] bg-white p-[14px] cd-p">
            {{ errorText }}
        </div>

        <div class="mt-[20px] grid gap-[15px]">
            <div v-if="loading" class="rounded-[14px] border border-[#191A23] bg-white p-[18px] cd-p">
                加载中...
            </div>

            <div v-else-if="items.length === 0" class="rounded-[14px] border border-[#191A23] bg-white p-[18px] cd-p">
                还没有草稿，先建一个吧。
            </div>

            <article v-else v-for="d in items" :key="d.id"
                class="rounded-[20px] border border-[#191A23] bg-white p-[20px]">
                <div class="flex items-start justify-between gap-[20px]">
                    <div class="min-w-0">
                        <RouterLink class="cd-h4 block truncate underline" :to="`/app/drafts/${d.id}`">
                            {{ d.title }}
                        </RouterLink>
                        <div class="cd-p mt-[8px] flex flex-wrap items-center gap-[10px] text-[#191A23]">
                            <span v-if="d.mood" class="cd-chip cd-p px-[10px] py-[2px]">{{ d.mood }}</span>
                            <span v-if="d.version != null">版本：{{ d.version }}</span>
                            <span v-if="d.updated_at">更新时间：{{ d.updated_at }}</span>
                        </div>
                    </div>
                </div>
            </article>
        </div>

        <div class="mt-[25px] flex items-center justify-between">
            <button class="cd-btn cd-btn-outline cd-p disabled:opacity-60" :disabled="page <= 1" type="button"
                @click="prevPage">
                上一页
            </button>
            <div class="cd-p">第{{ page }} / {{ totalPages }} 页</div>
            <button class="cd-btn cd-btn-outline cd-p disabled:opacity-60" :disabled="page >= totalPages" type="button"
                @click="nextPage">
                下一页
            </button>
        </div>
    </div>
</template>
