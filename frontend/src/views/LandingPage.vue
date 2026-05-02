<script setup>
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { getToken } from '../services/auth'
import star36Url from '../assets/positivus/star-36.svg'
import star100Url from '../assets/positivus/star-100.svg'
import catStartUrl from '../assets/cats/catStart.png'
import catMoodUrl from '../assets/cats/catMood.png'
import catDraftUrl from '../assets/cats/catDraft.png'
import catPhotoUrl from '../assets/cats/catPhoto.png'
import catTitleUrl from '../assets/cats/catTitle.png'

const router = useRouter()
const token = computed(() => getToken())

const processOpenIndex = ref(0)
const processSteps = [
  {
    title: '🐱 先来我家坐坐',
    desc: '注册登录，CatDiary 会给你留一个暖暖的小窝。'
  },
  {
    title: '📒 写下今天的小秘密',
    desc: '随便写点什么都可以，日记会乖乖听你说。'
  },
  {
    title: '🎭 贴个今日猫情绪',
    desc: '开心喵？委屈喵？用一个标签记住今天的你。'
  },
  {
    title: '🍃 草稿自动藏好',
    desc: '突然跑开也没关系，日记会自己藏好等你回来。'
  },
  {
    title: '📸 丢一张照片进来',
    desc: '上传照片或截图，让记忆有毛有肉～'
  },
  {
    title: '🌟 发布 & 慢慢翻',
    desc: '回头看每一篇，都是你走过的软软小爪印。'
  }
]

function goStart() {
  if (token.value) {
    router.push({ name: 'diary-new' })
    return
  }
  router.push({ name: 'login', query: { redirect: '/app/diaries/new' } })
}
</script>

<template>
  <div class="min-h-screen bg-base-100">
    <header class="cd-container pt-[20px]">
      <nav class="flex items-center justify-between">
        <RouterLink to="/" class="flex items-center gap-3">
          <div class="h-[36px] w-[36px] rounded-full bg-base-200"></div>
          <div class="cd-h4">CatDiary</div>
        </RouterLink>

        <div class="hidden items-center gap-[40px] md:flex">
          <a class="cd-p" href="#about">关于</a>
          <a class="cd-p" href="#features">功能</a>
          <a class="cd-p" href="#example">示例</a>
          <a class="cd-p" href="#plan">计划</a>
          <a class="cd-p" href="#contact">反馈</a>
        </div>

        <button class="cd-btn cd-btn-outline cd-p" type="button" @click="goStart">
          把今天写下来，是送给未来自己的一份小礼物。
        </button>
      </nav>
    </header>

    <main class="cd-container pt-[70px]">
      <section class="flex flex-col gap-[40px] md:flex-row md:items-start">
        <div class="flex-1">
          <h1 class="cd-h1">
            让我们
            <br />
            陪你记录每一天
          </h1>
          <p class="cd-p max-w-[520px] text-[#191A23]">
            CatDiary 是一个日记与草稿系统: 支持心情、位置、图片上传与自动保存，让记录变得轻松、可爱、可靠。
          </p>

          <div class="mt-[35px] flex flex-wrap gap-[15px]">
            <button class="cd-btn cd-btn-primary cd-p" type="button" @click="goStart">
              立即开始
            </button>
            <RouterLink class="cd-btn cd-btn-outline cd-p inline-flex items-center justify-center" to="/app/diaries">
              进入日记
            </RouterLink>
          </div>
        </div>

        <div class="flex-1">
          <div class="relative h-[420px] w-full overflow-hidden rounded-[45px] border border-[#191A23] bg-[#F3F3F3]">
            <img :src="catTitleUrl" alt="" class="absolute inset-0 h-full w-full object-cover select-none">
          </div>
        </div>
      </section>

      <section id="about" class="mt-[70px]">
        <div class="flex flex-col gap-[20px] md:flex-row md:items-center md:justify-between">
          <div class="flex items-center gap-[10px]">
            <span class="cd-chip cd-h4">关于</span>
            <span class="cd-h2">CatDiary</span>
          </div>
          <p class="cd-p max-w-[560px]">
            记录不是任务，而是一种柔软的自我照顾。随时开始，随时停下，回头看时，全是自己走过的光。
          </p>
        </div>
      </section>

      <section id="features" class="mt-[70px]">
        <div class="flex flex-col gap-[20px] md:flex-row md:items-center md:justify-between">
          <div class="flex items-center gap-[10px]">
            <span class="flex items-center gap-[10px]"></span>
            <span class="cd-chip cd-h4">功能</span>
            <span class="cd-h2">你会用到的</span>
          </div>
          <p class="cd-p max-w-[560px]">
            日记、草稿、Mood、图片与安全设置，先将最重要的体验做扎实。
          </p>
        </div>

        <div class="mt-[40px] grid gap-[40px] md:grid-cols-2">
          <div class="cd-card flex items-center justify-between p-[50px]">
            <div class="flex flex-col justify-between gap-[40px]">
              <div class="flex flex-col gap-[8px]">
                <span class="cd-chip cd-h4">快速记录</span>
                <span class="cd-chip cd-h4">新建日记</span>
              </div>
              <RouterLink class="cd-p inline-flex items-center gap-[15px]" to="/app/diaries/new">
                <span>开始</span>
                <span
                  class="inline-flex h-[32px] w-[32px] items-center justify-center rounded-full border border-[#191A23] bg-white">
                  <svg width="14" height="14" viewBox="0 0 14 14" fill="none">
                    <path d="M4 7H10" stroke="#191A23" stroke-width="1.5" stroke-linecap="round" />
                    <path d="M8 5L10 7L8 9" stroke="#191A23" stroke-width="1.5" stroke-linecap="round"
                      stroke-linejoin="round" />
                  </svg>
                </span>
              </RouterLink>
            </div>
            <img :src="catStartUrl" alt="" class="h-[170px] w-[210px] rounded-[16px] bg-white object-contain">
          </div>

          <div class="cd-card flex items-center justify-between p-[50px]">
            <div class="flex flex-col justify-between gap-[40px]">
              <div class="flex flex-col gap-[8px]">
                <span class="cd-chip cd-h4">心情 Mood</span>
                <span class="cd-chip cd-h4">一眼回顾</span>
              </div>
              <RouterLink class="cd-p inline-flex items-center gap-[15px]" to="/app/diaries">
                <span>查看</span>
                <span
                  class="inline-flex h-[32px] w-[32px] items-center justify-center rounded-full border border-[#191A23] bg-white">
                  <svg width="14" height="14" viewBox="0 0 14 14" fill="none">
                    <path d="M4 7H10" stroke="#191A23" stroke-width="1.5" stroke-linecap="round" />
                    <path d="M8 5L10 7L8 9" stroke="#191A23" stroke-width="1.5" stroke-linecap="round"
                      stroke-linejoin="round" />
                  </svg>
                </span>
              </RouterLink>
            </div>
            <img :src="catMoodUrl" alt="" class="h-[170px] w-[210px] rounded-[16px] bg-white object-contain">
          </div>

          <div class="cd-card flex items-center justify-between p-[50px]">
            <div class="flex flex-col justify-between gap-[40px]">
              <div class="flex flex-col gap-[8px]">
                <span class="cd-chip cd-h4">草稿</span>
                <span class="cd-chip cd-h4">随时继续</span>
              </div>
              <RouterLink class="cd-p inline-flex items-center gap-[15px]" to="/app/diaries">
                <span>进入</span>
                <span
                  class="inline-flex h-[32px] w-[32px] items-center justify-center rounded-full border border-[#191A23] bg-white">
                  <svg width="14" height="14" viewBox="0 0 14 14" fill="none">
                    <path d="M4 7H10" stroke="#191A23" stroke-width="1.5" stroke-linecap="round" />
                    <path d="M8 5L10 7L8 9" stroke="#191A23" stroke-width="1.5" stroke-linecap="round"
                      stroke-linejoin="round" />
                  </svg>
                </span>
              </RouterLink>
            </div>
            <img :src="catDraftUrl" alt="" class="h-[170px] w-[210px] rounded-[16px] bg-white object-contain">
          </div>

          <div class="cd-card flex items-center justify-between p-[50px]">
            <div class="flex flex-col justify-between gap-[40px]">
              <div class="flex flex-col gap-[8px]">
                <span class="cd-chip cd-h4">图片上传</span>
                <span class="cd-chip cd-h4">让回忆更清晰</span>
              </div>
              <RouterLink class="cd-p inline-flex items-center gap-[15px]" to="/app/diaries/new">
                <span>试试</span>
                <span
                  class="inline-flex h-[32px] w-[32px] items-center justify-center rounded-full border border-[#191A23] bg-white">
                  <svg width="14" height="14" viewBox="0 0 14 14" fill="none">
                    <path d="M4 7H10" stroke="#191A23" stroke-width="1.5" stroke-linecap="round" />
                    <path d="M8 5L10 7L8 9" stroke="#191A23" stroke-width="1.5" stroke-linecap="round"
                      stroke-linejoin="round" />
                  </svg>
                </span>
              </RouterLink>
            </div>
            <img :src="catPhotoUrl" alt="" class="h-[170px] w-[210px] rounded-[16px] bg-white object-contain">
          </div>
        </div>
      </section>

      <section id="examples" class="mt-[70px]">
        <div class="cd-card !bg-[#191A23] p-[50px] text-white">
          <div class="flex flex-col gap-[20px] md:flex-row md:items-center md:justify-between">
            <div class="flex items-center gap-[10px]">
              <span class="cd-chip cd-h4">示例</span>
              <span class="cd-h2 text-white">三种写法</span>
            </div>
            <p class="cd-p max-w-[560px] text-white">
              不知道怎么开始？试试这三个暖暖的小模板～ 🐱
            </p>
          </div>
          <div class="mt-[40px] grid gap-[30px] md:grid-cols-3">
            <div class="rounded-[20px] border border-white/20 p-[25px]">
              <div class="cd-h4 text-white">✨ 三件好事</div>
              <p class="cd-p mt-[10px] text-white/80">今天让我开心的三件小事，再小也值得被记住。</p>
            </div>
            <div class="rounded-[20px] border border-white/20 p-[25px]">
              <div class="cd-h4 text-white">💌 给自己的一封信</div>
              <p class="cd-p mt-[10px] text-white/80">低落的时候，像对待朋友一样，温柔地对自己说说话。</p>
            </div>
            <div class="rounded-[20px] border border-white/20 p-[25px]">
              <div class="cd-h4 text-white">🐾 猫咪日常</div>
              <p class="cd-p mt-[10px] text-white/80">记录被猫咪治愈的小瞬间，那些毛茸茸的温暖。</p>
            </div>
          </div>
        </div>
      </section>

      <section id="plan" class="mt-[70px]">
        <div class="flex flex-col gap-[20px] md:flex-row md:items-center md:justify-between">
          <div class="flex items-center gap-[10px]">
            <span class="cd-chip cd-h4">步骤</span>
            <span class="cd-h2">开始很简单</span>
          </div>
          <p class="cd-p max-w-[560px]">别担心，不用写很长。跟着这三步，1分钟就能记录下今天的心情。</p>
        </div>

        <div class="mt-[40px] grid gap-[15px]">
          <button v-for="(s, idx) in processSteps" :key="s.title"
            class="cd-card flex w-full items-center justify-between p-[30px] text-left"
            :class="processOpenIndex === idx ? '!bg-[#191A23] text-white !border-white/20' : '!bg-[#F3F3F3] text-[#191A23]'"
            type="button" @click="processOpenIndex = processOpenIndex === idx ? -1 : idx">
            <div>
              <div class="cd-h4">{{ `${idx + 1}. ${s.title}` }}</div>
              <div v-if="processOpenIndex === idx" class="cd-p mt-[10px] max-w-[860px]">
                {{ s.desc }}
              </div>
            </div>
            <div class="flex h-[40px] w-[40px] items-center justify-center rounded-full border "
              :class="processOpenIndex === idx ? 'border-white/40 bg-white/10 text-white' : 'border-[#191A23] bg-white text-[#191A23]'">
              <svg v-if="processOpenIndex !== idx" width="16" height="16" viewBox="0 0 16 16" fill="none">
                <path d="M8 3V13" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
                <path d="M3 8H13" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
              </svg>
              <svg v-else width="16" height="16" viewBox="0 0 16 16" fill="none">
                <path d="M3 8H13" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
              </svg>
            </div>
          </button>
        </div>
      </section>

      <section id="contact" class="mt-[70px] pb-[90px]">
        <div class="flex flex-col gap-[20px] md:flex-row md:items-center md:justify-between">
          <div class="flex items-center gap-[10px]">
            <span class="cd-chip cd-h4">反馈</span>
            <span class="cd-h2">🐾 喵～想先要什么功能？</span>
          </div>
          <p class="cd-p max-w-[560px]">
            你来说，我来做。最有价值的那个，我会优先安排上！
          </p>
        </div>

        <div class="mt-[40px] cd-card bg-[#F3F3F3] p-[50px]">
          <form class="grid gap-[20px] md:grid-cols-2" @submit.prevent>
            <label class="grid gap-[10px]">
              <span class="cd-p">称呼</span>
              <input class="h-[54px] rounded-[14px] border border-[#191A23] bg-white px-[18px] outline-none"
                placeholder="怎么称呼你喵？" />
            </label>
            <label class="grid gap-[10px]">
              <span class="cd-p">邮箱</span>
              <input class="h-[54px] rounded-[14px] border border-[#191A23] bg-white px-[18px] outline-none"
                placeholder="方便留个邮箱喵？" />
            </label>
            <label class="grid gap-[10px] md:col-span-2">
              <span class="cd-p">想要的功能</span>
              <textarea
                class="min-h-[150px] rounded-[14px] border border-[#191A23] bg-white px-[18px] py-[14px] outline-none"
                placeholder="比如：搜索、标签、统计、主题、同步、导出等等喵~"></textarea>
            </label>
            <div class="md:col-span-2">
              <button class="cd-btn cd-btn-primary cd-p" type="submit">提交</button>
            </div>
          </form>
        </div>
      </section>
    </main>

    <footer class="cd-container">
      <div class="rounded-t-[45px] bg-[#191A23] px-[50px] py-[40px] text-white">
        <div class="flex flex-col gap-[20px] md:flex-row md:items-center md:justify-between">
          <div class="cd-h4 text-white">CatDiary</div>
          <div class="cd-p text-white/80">开始 · 探索 · 分享</div>
        </div>
      </div>
    </footer>
  </div>
</template>
