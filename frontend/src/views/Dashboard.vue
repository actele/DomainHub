<template>
  <a-layout class="h-screen dashboard-layout">
    <a-layout-header class="dashboard-header fixed w-full z-10 px-8 flex justify-between items-center">
      <div class="flex items-center gap-3">
        <a-button
          v-if="!isLg"
          class="menu-trigger"
          type="text"
          size="large"
          aria-label="打开菜单"
          @click="drawerOpen = true"
        >
          <template #icon>
            <IconMenu />
          </template>
        </a-button>
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

    <a-layout class="layout-body">
      <a-layout-sider
        v-if="isLg"
        :width="220"
        class="h-[calc(100vh-68px)] bg-white dashboard-sider"
      >
        <a-menu
          :selected-keys="[$route.name]"
          class="dashboard-menu"
          :style="{ height: '100%', borderRight: 0 }"
          @menuItemClick="handleMenuSelect"
        >
          <a-menu-item key="domains">
            <template #icon><IconPublic /></template>
            域名管理
          </a-menu-item>
          <a-menu-item key="provider-keys">
            <template #icon><IconLock /></template>
            服务商密钥
          </a-menu-item>
          <template v-if="isAdmin">
            <a-menu-item-group title="系统管理">
              <a-menu-item key="admin-users">
                <template #icon><IconUser /></template>
                账号管理
              </a-menu-item>
              <a-menu-item key="admin-providers">
                <template #icon><IconTool /></template>
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

    <a-drawer
      v-if="!isLg"
      :visible="drawerOpen"
      placement="left"
      :width="280"
      :footer="false"
      :mask-closable="true"
      :closable="true"
      class="dashboard-drawer"
      @cancel="drawerOpen = false"
      @ok="drawerOpen = false"
    >
      <a-menu
        :selected-keys="[$route.name]"
        class="dashboard-menu"
        :style="{ height: '100%', borderRight: 0 }"
        @menuItemClick="handleMenuSelect"
      >
        <a-menu-item key="domains">
          <template #icon><IconPublic /></template>
          域名管理
        </a-menu-item>
        <a-menu-item key="provider-keys">
          <template #icon><IconLock /></template>
          服务商密钥
        </a-menu-item>
        <template v-if="isAdmin">
          <a-menu-item-group title="系统管理">
            <a-menu-item key="admin-users">
              <template #icon><IconUser /></template>
              账号管理
            </a-menu-item>
            <a-menu-item key="admin-providers">
              <template #icon><IconTool /></template>
              服务商管理
            </a-menu-item>
          </a-menu-item-group>
        </template>
      </a-menu>
    </a-drawer>
  </a-layout>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useRoute } from 'vue-router'
import { IconPublic, IconLock, IconExport, IconUser, IconTool, IconMenu } from '@arco-design/web-vue/es/icon'
import { clearDomainRecordsCache } from '@/utils/domainRecordsCache'
import { useMediaQuery } from '@/composables/useBreakpoint'

const router = useRouter()
const route = useRoute()

// ≥1024px → keep sider; below that → hide sider, show drawer-backed hamburger
const isLg = useMediaQuery('(min-width: 1024px)')
const drawerOpen = ref(false)

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
  clearDomainRecordsCache()
  router.push('/login')
}

const handleMenuSelect = (key) => {
  if (key.startsWith('admin-')) {
    router.push(`/dashboard/admin/${key.replace('admin-', '')}`)
  } else {
    router.push(`/dashboard/${key}`)
  }
  // close drawer on mobile after navigation
  drawerOpen.value = false
}

const isAdmin = computed(() => localStorage.getItem('role') === 'admin')
</script>

<style scoped>
.dashboard-layout {
	min-height: 100vh;
	background: transparent;
}

.dashboard-header {
	height: 68px;
	line-height: 68px;
	background: rgba(255, 255, 255, 0.84);
	backdrop-filter: blur(16px);
	border-bottom: 1px solid rgba(226, 232, 240, 0.9);
	box-shadow: 0 4px 18px rgba(15, 23, 42, 0.04);
}

.brand-badge {
	width: 38px;
	height: 38px;
	border-radius: 12px;
	background: linear-gradient(145deg, #2563eb, #4f46e5);
	color: #fff;
	display: flex;
	align-items: center;
	justify-content: center;
	font-size: 14px;
	font-weight: 700;
	letter-spacing: 0.5px;
	box-shadow: 0 8px 18px rgba(37, 99, 235, 0.24);
}

.dashboard-menu {
	padding: 16px 12px;
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
	width: 100%;
	max-width: 1680px;
	margin: 0 auto;
	background: transparent;
}

.layout-body {
	margin-top: 68px;
}

@media (max-width: 1023px) {
	.layout-body {
		margin-top: 0;
		/* On smaller viewports the header is fixed but sider is hidden, so the
		   layout body still needs to clear it. Padding (not margin) keeps the
		   background gradient from the body itself uninterrupted. */
		padding-top: 68px;
	}
}

:deep(.dashboard-menu .arco-menu-item-group-title) {
	padding: 18px 12px 8px;
	color: #94a3b8;
	font-size: 11px;
	font-weight: 700;
	letter-spacing: 0.08em;
	text-transform: uppercase;
}

:deep(.dashboard-menu .arco-menu-item) {
	min-height: 42px;
	color: #64748b;
	font-size: 13px;
	transition: color 0.2s ease, background 0.2s ease, transform 0.2s ease;
}

:deep(.dashboard-menu .arco-menu-item:hover) {
	color: #1d4ed8;
	background: #f1f5ff;
	transform: translateX(2px);
}

@media (max-width: 760px) {
	.dashboard-header {
		padding-right: 16px !important;
		padding-left: 16px !important;
	}

	.dashboard-header > div:first-child > div:last-child {
		display: none;
	}

	.dashboard-header .header-pill {
		display: none;
	}

	.dashboard-content {
		padding: 16px !important;
	}
}
</style>
