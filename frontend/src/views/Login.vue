<template>
  <div class="login-container flex items-center justify-center min-h-screen px-4 py-10">
    <div class="login-card w-full max-w-5xl">
      <div class="auth-left">
        <div class="brand-chip">DomainHub</div>
        <h1 class="auth-title">稳定管理多云 DNS</h1>
        <p class="auth-desc">统一接入阿里云、腾讯云、Cloudflare，面向企业日常运维场景设计。</p>
        <div class="auth-points">
          <div>统一账号权限</div>
          <div>集中管理域名与记录</div>
          <div>安全存储服务商密钥</div>
        </div>
      </div>

      <div class="auth-right">
        <div class="mb-8 text-center">
          <h2 class="login-title-text text-2xl font-bold">登录系统</h2>
          <p class="text-sm text-slate-500 mt-2">请输入账号与密码</p>
        </div>

        <a-form :model="form" @submit="handleSubmit">
          <a-form-item field="username" label="用户名">
            <a-input v-model="form.username" placeholder="请输入用户名" allow-clear />
          </a-form-item>

          <a-form-item field="password" label="密码">
            <a-input-password v-model="form.password" placeholder="请输入密码" allow-clear />
          </a-form-item>

          <a-button type="primary" html-type="submit" class="login-button w-full mb-4" :loading="loading">
            {{ loading ? '登录中...' : '登录' }}
          </a-button>
        </a-form>

        <div class="text-center">
          <a class="register-link text-primary" @click="router.push('/register')">注册新账号</a>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { Message } from '@arco-design/web-vue'
import { useRouter } from 'vue-router'
import api from '@/utils/api'

const router = useRouter()
const loading = ref(false)
const form = reactive({
  username: '',
  password: ''
})

const handleSubmit = async () => {
  if (!form.username || !form.password) {
    Message.warning('请填写完整的登录信息')
    return
  }

  try {
    loading.value = true
    const data = await api.post('/api/v1/users/login', {
      username: form.username.trim(),
      password: form.password
    })

    Message.success({
      content: '登录成功，正在跳转...',
      duration: 2000
    })

    localStorage.setItem('token', data.token)
    localStorage.setItem('role', data.role || 'user')

    setTimeout(() => {
      router.push('/dashboard')
    }, 1500)
  } catch (error) {
    const errorMsg = error.response?.data?.message || error.message || '登录失败，请稍后重试'
    Message.error({
      content: errorMsg,
      duration: 3000
    })
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  background:
    radial-gradient(circle at 5% 10%, rgba(22, 119, 255, 0.12), transparent 35%),
    radial-gradient(circle at 90% 90%, rgba(56, 189, 248, 0.1), transparent 30%),
    #f5f7fa;
}

.login-card {
  border-radius: 8px;
  border: 1px solid #e5e7eb;
  background: #fff;
  box-shadow: 0 4px 12px rgba(15, 23, 42, 0.08);
  display: grid;
  grid-template-columns: 1.15fr 1fr;
  overflow: hidden;
}

.auth-left {
  padding: 48px;
  background: linear-gradient(155deg, #f0f7ff 0%, #eef2ff 100%);
  border-right: 1px solid #e5e7eb;
}

.brand-chip {
  display: inline-flex;
  height: 28px;
  align-items: center;
  padding: 0 12px;
  border-radius: 999px;
  background: #dbeafe;
  color: #1d4ed8;
  font-size: 12px;
  font-weight: 600;
}

.auth-title {
  margin-top: 20px;
  font-size: 32px;
  font-weight: 700;
  line-height: 1.25;
  color: #0f172a;
}

.auth-desc {
  margin-top: 12px;
  color: #64748b;
  font-size: 14px;
}

.auth-points {
  margin-top: 24px;
  display: grid;
  gap: 10px;
  color: #334155;
  font-size: 14px;
}

.auth-right {
  padding: 40px 36px;
}

.login-title-text {
  color: #0f172a;
  letter-spacing: 0.2px;
}

:deep(.a-input) {
  transition: all 0.2s ease;
  border-radius: 8px;
  background: #f8fafc;
  border-color: #d7e2ee;
}

:deep(.a-input:focus-within) {
  background: #fff;
  box-shadow: 0 0 0 2px rgba(15, 76, 129, 0.14);
}

.login-button {
  height: 42px;
  border-radius: 8px;
}

.register-link {
  font-weight: 500;
}

@media (max-width: 640px) {
  .login-card {
    grid-template-columns: 1fr;
  }

  .auth-left {
    padding: 24px;
    border-right: 0;
    border-bottom: 1px solid #e5e7eb;
  }

  .auth-right {
    padding: 24px;
  }

  .login-card {
    border-radius: 12px;
  }
}
</style>