<script setup>
import { ref } from 'vue'
import { apiRequest } from '../../services/api'

const saving = ref(false)
const errorText = ref('')
const hintText = ref('')

const form = ref({
    oldPassword: '',
    newPassword: '',
    newPassword2: '',
})

async function submit() {
    errorText.value = ''
    hintText.value = ''

    if (!form.value.oldPassword || !form.value.newPassword) {
        errorText.value = '请填写旧密码与新密码'
        return
    }
    if (form.value.newPassword !== form.value.newPassword2) {
        errorText.value = '两次输入的新密码不一致'
        return
    }

    saving.value = true
    try {
        const res = await apiRequest('/users/me/password', {
            method: 'PATCH',
            body: JSON.stringify({
                old_password: form.value.oldPassword,
                new_password: form.value.newPassword,
            }),
        })
        hintText.value = res?.message || '密码修改成功'
        form.value = { oldPassword: '', newPassword: '', newPassword2: '' }
    } catch (e) {
        errorText.value = e?.message || '修改失败'
    } finally {
        saving.value = false
    }
}
</script>

<template>
    <div>
        <div class="flex flex-col gap-[10px] md:flex-row md:items-center md:justify-between">
            <div>
                <div class="cd-h2">安全设置</div>
                <div class="cd-p mt-[8px]">
                    <span v-if="hintText">{{ hintText }}</span>
                    <span v-else>定期换个密码，心情会更安稳。</span>
                </div>
            </div>
        </div>

        <div v-if="errorText" class="mt-[20px] rounded-[14px] border border-[#191A23] bg-white p-[14px] cd-p">
            {{ errorText }}
        </div>

        <div class="mt-[20px] cd-card bg-[#F3F3F3] p-[25px]">
            <form class="grid gap-[20px] max-w-[520px]" @submit.prevent="submit">
                <label class="grid gap-[10px]">
                    <span class="cd-p">旧密码</span>
                    <input v-model="form.oldPassword"
                        class="h-[54px] rounded-[14px] border border-[#191A23] bg-white px-[18px] outline-none"
                        type="password" autocomplete="current-password" />
                </label>
                <label class="grid gap-[10px]">
                    <span class="cd-p">新密码</span>
                    <input v-model="form.newPassword"
                        class="h-[54px] rounded-[14px] border border-[#191A23] bg-white px-[18px] outline-none"
                        type="password" autocomplete="new-password" />
                </label>
                <label class="grid gap-[10px]">
                    <span class="cd-p">确认新密码</span>
                    <input v-model="form.newPassword2"
                        class="h-[54px] rounded-[14px] border border-[#191A23] bg-white px-[18px] outline-none"
                        type="password" autocomplete="new-password" />
                </label>
                <button class="cd-btn cd-btn-primary cd-p disabled:opacity-60" :disabled="saving" type="submit">
                    {{ saving ? '提交中...' : '修改密码' }}
                </button>
            </form>
        </div>
    </div>
</template>
