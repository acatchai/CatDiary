<script setup>
import { computed, onMounted, ref } from 'vue'
import { apiRequest } from '../../services/api'

const loading = ref(false)
const errorText = ref('')

const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const items = ref([])

const totalPages = computed(() => {
    const t = Number(total.value) || 0
    const s = Number(pageSize.value) || 20
    return Math.max(1, Math.ceil(t / s))
})

async function load() {
    loading.value = true
    errorText.value = ''
    try {
        const res = await apiRequest(`/diaries?page=${page.value}&page_size=${pageSize.value}`)
        items.value = Array.isArray(res?.items) ? res.items : []
        total.value = Number(res?.total) || 0
    } catch (e) {
        errorText.value = e?.message || '加载失败'
    } finally {
        loading.value = false
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
                <div class="cd-h2">日记</div>
                <div class="cd-p mt-[8px]">把你的每一天收进小盒子里。</div>
            </div>
            <RouterLink class="cd-btn cd-btn-primary cd-p inline-flex items-center justify-center"
                to="/app/diaries/new">
                写一条
            </RouterLink>
        </div>

        <div v-if="errorText" class="mt-[20px] rounded-[14px] border border-[#191A23] bg-white p-[14px] cd-p">
            {{ errorText }}
        </div>

        <div class="mt-[20px] grid gap-[15px]">
            <div v-if="loading" class="rounded-[14px] border border-[#191A23] bg-white p-[18px] cd-p">
                加载中...
            </div>

            <div v-else-if="items.length === 0" class="rounded-[14px] border border-[#191A23] bg-white p-[18px] cd-p">
                还没有日记，先写一条吧。
            </div>

            <article v-else v-for="d in items" :key="d.id"
                class="rounded-[20px] border border-[#191A23] bg-white p-[20px]">
                <div class="flex items-start justify-between gap-[20px]">
                    <div class="min-w-0">
                        <RouterLink class="cd-h4 block truncate underline" :to="`/app/diaries/${d.id}`">
                            {{ d.title }}
                        </RouterLink>
                        <div class="cd-p mt-[8px] text-[#191A23]">
                            <span v-if="d.mood" class="inline-flex items-center gap-[8px]">
                                <span class="cd-chip cd-p px-[10px] py-[2px]">{{ d.mood }}</span>
                            </span>
                            <span v-if="d.occurred_at" class="ml-[10px]">发生时间: {{ d.occurred_at }}</span>
                        </div>
                    </div>
                </div>
            </article>
        </div>
    </div>

    <div class="mt-[25px] flex items-center justify-between">
        <button class="cd-btn cd-btn-outline cd-p disabled:opacity-60" :disabled="page <= 1" type="button"
            @click="prevPage">
            上一页
        </button>
        <div class="cd-p">第 {{ page }} / {{ totalPages }} 页</div>
        <button class="cd-btn cd-btn-outline cd-p disabled:opacity-60" :disabled="page >= totalPages" type="button"
            @click="nextPage">
            下一页
        </button>
    </div>
</template>
