<template>
  <div class="register-container flex items-center justify-center min-h-screen px-4 py-10">
    <div class="register-card w-full max-w-5xl">
      <div class="auth-left">
        <div class="brand-chip">DomainHub</div>
        <h1 class="auth-title">创建团队入口账号</h1>
        <p class="auth-desc">注册后可维护域名、解析记录和服务商密钥，建议使用团队统一命名规范。</p>
        <div class="auth-points">
          <div>密码最少 6 位</div>
          <div>注册后可立即登录系统</div>
          <div>支持后续修改密码</div>
        </div>
      </div>

      <div class="auth-right">
        <div class="mb-8 text-center">
          <h2 class="register-title-text text-2xl font-bold">创建账号</h2>
          <p class="text-sm text-slate-500 mt-2">请输入账号与密码</p>
        </div>

        <a-form :model="formData" @submit="handleRegister">
          <a-form-item field="username" label="用户名">
            <a-input v-model="formData.username" placeholder="请输入用户名" allow-clear />
          </a-form-item>

          <a-form-item field="password" label="密码">
            <a-input-password v-model="formData.password" placeholder="请输入密码" allow-clear />
          </a-form-item>

          <a-button type="primary" html-type="submit" class="register-button w-full mb-4" :loading="loading">
            {{ loading ? '注册中...' : '注册' }}
          </a-button>
        </a-form>

        <div class="text-center">
          <a class="login-link text-primary" @click="router.push('/login')">返回登录</a>
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
const formData = reactive({
  username: '',
  password: ''
})

const handleRegister = async () => {
  if (!formData.username || !formData.password) {
    Message.warning('请填写完整的注册信息')
    return
  }

  if (formData.password.length < 6) {
    Message.warning('密码长度至少6位')
    return
  }

  try {
    loading.value = true

    await api.post('/api/v1/users/register', {
      username: formData.username.trim(),
      password: formData.password
    })

    Message.success({
      content: '注册成功，即将跳转到登录页面',
      duration: 2000
    })

    setTimeout(() => {
      router.push('/login')
    }, 1500)
  } catch (error) {
    const errorMsg = error.response?.data?.message || error.message || '注册失败，请稍后重试'
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
.register-container {
  background:
    radial-gradient(circle at 10% 12%, rgba(22, 119, 255, 0.12), transparent 34%),
    radial-gradient(circle at 88% 90%, rgba(56, 189, 248, 0.1), transparent 30%),
    #f5f7fa;
}

.register-card {
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
  font-size: 30px;
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

.register-title-text {
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

.register-button {
  height: 42px;
  border-radius: 8px;
}

.login-link {
  font-weight: 500;
}

@media (max-width: 640px) {
  .register-card {
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

  .register-card {
    border-radius: 12px;
  }
}
</style>