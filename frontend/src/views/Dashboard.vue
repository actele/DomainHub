<template>
  <a-layout class="h-screen dashboard-layout">
    <a-layout-header class="dashboard-header fixed w-full z-10 px-8 flex justify-between items-center">
      <div class="flex items-center gap-4">
        <div class="brand-badge">DM</div>
        <div>
          <div class="text-xl font-bold text-slate-900">域名管理系统</div>
          <div class="text-xs text-slate-500">企业级 DNS 管理控制台</div>
        </div>
      </div>
      <div class="flex items-center gap-4">
        <div class="header-pill">
          <span class="status-dot"></span>
          {{ currentSection }}
        </div>
          <div class="text-sm font-medium" :class="isAdmin ? 'text-blue-600' : 'text-slate-500'">
            {{ isAdmin ? '管理员' : '普通用户' }}
          </div>
        <a-button type="text" size="large" @click="handleLogout">
          <template #icon>
            <IconExport />
          </template>
          退出登录
        </a-button>
      </div>
    </a-layout-header>

    <a-layout style="margin-top: 72px">
      <a-layout-sider :width="220" class="h-[calc(100vh-72px)] bg-white">
        <a-menu :selected-keys="[$route.name]" class="dashboard-menu" :style="{ height: '100%', borderRight: 0 }"
          @menuItemClick="handleMenuClick">
          <a-menu-item key="domains">
            <template #icon>
              <IconPublic />
            </template>
            域名管理
          </a-menu-item>
          <a-menu-item key="provider-keys">
            <template #icon>
              <IconLock />
            </template>
            服务商密钥
          </a-menu-item>
            <template v-if="isAdmin">
              <a-menu-item-group title="系统管理">
                <a-menu-item key="admin-users">
                  <template #icon>
                    <IconUser />
                  </template>
                  账号管理
                </a-menu-item>
                <a-menu-item key="admin-providers">
                  <template #icon>
                    <IconTool />
                  </template>
                  服务商管理
                </a-menu-item>
              </a-menu-item-group>
            </template>
          </a-menu>
      </a-layout-sider>

      <a-layout-content class="dashboard-content p-6">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useRoute } from 'vue-router'
import { IconPublic, IconLock, IconExport, IconUser, IconTool } from '@arco-design/web-vue/es/icon'

const router = useRouter()
const route = useRoute()

const sectionNameMap = {
  domains: '域名管理',
  'domain-records': '解析记录',
    'provider-keys': '服务商密钥',
    'admin-users': '账号管理',
    'admin-providers': '服务商管理'
}

const currentSection = computed(() => {
  return sectionNameMap[route.name] || '控制台'
})

const handleLogout = () => {
  localStorage.removeItem('token')
    localStorage.removeItem('role')
  router.push('/login')
}

const handleMenuClick = (key) => {
  if (key.startsWith('admin-')) {
    router.push(`/dashboard/admin/${key.replace('admin-', '')}`)
  } else {
    router.push(`/dashboard/${key}`)
  }
}

const isAdmin = computed(() => localStorage.getItem('role') === 'admin')
</script>

<style scoped>
.dashboard-layout {
  background: #f5f7fa;
}

.dashboard-header {
  height: 72px;
  line-height: 72px;
  background: rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid #e5e7eb;
}

.brand-badge {
  width: 38px;
  height: 38px;
  border-radius: 8px;
  background: #1677ff;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 700;
  letter-spacing: 0.5px;
}

.dashboard-menu {
  padding: 12px 10px;
}

.header-pill {
  height: 30px;
  padding: 0 12px;
  border: 1px solid #dbeafe;
  border-radius: 999px;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: #1d4ed8;
  background: #eff6ff;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #10b981;
  box-shadow: 0 0 0 4px rgba(16, 185, 129, 0.15);
}

:deep(.dashboard-menu .arco-menu-item) {
  border-radius: 8px;
  margin-bottom: 4px;
}

:deep(.dashboard-menu .arco-menu-selected) {
  background: #e6f4ff;
  color: #1677ff;
  font-weight: 600;
}

.dashboard-content {
  background: transparent;
}
</style>