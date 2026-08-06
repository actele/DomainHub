<template>
  <div class="page-stack">
    <a-breadcrumb class="page-breadcrumb">
      <a-breadcrumb-item>控制台</a-breadcrumb-item>
      <a-breadcrumb-item>
        <span class="inline-flex items-center gap-1">
          域名管理
          <a-tooltip content="统一管理各服务商域名资源" position="right">
            <icon-info-circle class="text-gray-400 cursor-default" style="font-size:13px" />
          </a-tooltip>
        </span>
      </a-breadcrumb-item>
    </a-breadcrumb>

    <div class="page-heading">
      <div>
        <div class="page-eyebrow">DNS / DOMAINS</div>
        <h1 class="page-title">域名管理</h1>
        <p class="page-subtitle">统一查看和维护已接入 DNS 服务商的托管域名。</p>
      </div>
    </div>

    <div class="page-toolbar">
      <a-space class="toolbar-filters">
        <a-input v-model="filters.keyword" placeholder="搜索域名" allow-clear class="toolbar-input toolbar-input-search" />
        <a-select v-model="filters.provider" placeholder="全部服务商" allow-clear class="toolbar-input toolbar-input-provider">
          <a-option value="aliyun">阿里云</a-option>
          <a-option value="tencent">腾讯云</a-option>
          <a-option value="cloudflare">Cloudflare</a-option>
        </a-select>
      </a-space>
      <div class="toolbar-meta">
        <span class="text-sm text-gray-500">共 {{ filteredDomains.length }} 条</span>
        <a-button type="primary" @click="showAddDomainModal = true">
          <template #icon><IconPlus /></template>
          添加域名
        </a-button>
      </div>
    </div>

    <!-- 域名列表 — 桌面 / 平板 -->
    <a-card v-if="!isMobile" :bordered="false" class="general-card">
      <a-table :loading="loading" :data="filteredDomains" :pagination="tablePagination">
        <template #columns>
          <a-table-column title="域名" data-index="name" :width="320">
            <template #cell="{ record }">
              <div class="flex items-center gap-2 py-1">
                <IconPublic class="text-gray-400" />
                <span class="text-primary font-semibold">{{ record.name }}</span>
              </div>
            </template>
          </a-table-column>
          <a-table-column title="服务商" :width="180">
            <template #cell="{ record }">
              {{ providerNames[record.provider] || record.provider }}
            </template>
          </a-table-column>
          <a-table-column title="创建时间" :width="200">
            <template #cell="{ record }">
              {{ record.created_at ? new Date(record.created_at).toLocaleString() : '-' }}
            </template>
          </a-table-column>
          <a-table-column title="操作" align="right" :width="180">
            <template #cell="{ record }">
              <div class="table-actions">
                <a-button type="text" status="normal"
                  @mouseenter="preloadDomainRecords(record)"
                  @focus="preloadDomainRecords(record)"
                  @click="$router.push({ name: 'domain-records', params: { id: record.id } })">
                  解析记录
                </a-button>
                <a-button type="text" status="danger" @click="handleDeleteDomain(record)">
                  删除
                </a-button>
              </div>
            </template>
          </a-table-column>
        </template>
        <template #empty>
          <div class="text-center py-8">
            <IconPublic class="text-gray-400 text-3xl mb-2" />
            <div class="text-gray-900 font-medium">暂无数据</div>
            <div class="text-gray-500 text-sm mt-1">当前筛选条件下暂无域名，请调整筛选条件或新增数据</div>
            <a-button type="primary" class="mt-4" @click="showAddDomainModal = true">
              <template #icon>
                <IconPlus />
              </template>
              添加域名
            </a-button>
          </div>
        </template>
      </a-table>
    </a-card>

    <!-- 域名列表 — 手机（卡片栈） -->
    <div v-else class="domain-mobile-list">
      <a-spin :loading="loading" class="w-full">
        <div v-if="!filteredDomains.length" class="text-center py-10 text-gray-500 text-sm">
          暂无域名
        </div>
        <div v-for="record in pagedDomains" :key="record.id" class="list-card">
          <div class="list-card-header">
            <div>
              <div class="list-card-title flex items-center gap-2">
                <IconPublic class="text-gray-400" />
                <span class="text-primary">{{ record.name }}</span>
              </div>
              <div class="list-card-meta">
                {{ providerNames[record.provider] || record.provider }}
                · {{ record.created_at ? new Date(record.created_at).toLocaleString() : '-' }}
              </div>
            </div>
          </div>
          <div class="list-card-actions">
            <a-button type="text"
              @mouseenter="preloadDomainRecords(record)"
              @focus="preloadDomainRecords(record)"
              @click="$router.push({ name: 'domain-records', params: { id: record.id } })">
              解析记录
            </a-button>
            <a-button type="text" status="danger" @click="handleDeleteDomain(record)">
              删除
            </a-button>
          </div>
        </div>
        <a-pagination
          v-if="filteredDomains.length > tablePagination.pageSize"
          class="mt-3 justify-center"
          :total="filteredDomains.length"
          :page-size="tablePagination.pageSize"
          :current="currentPage"
          show-total
          @change="(p) => (currentPage = p)"
          @page-size-change="(s) => { tablePagination.pageSize = s; currentPage = 1 }"
        />
      </a-spin>
    </div>

    <!-- 添加域名对话框 -->
    <a-modal v-model:visible="showAddDomainModal" title="添加域名" @ok="handleAddDomain" @cancel="resetAddDomain"
      :ok-button-props="{ disabled: !newDomain.name || !newDomain.provider }" ok-text="添加" cancel-text="取消"
      :width="520">
      <a-form :model="newDomain" layout="vertical" class="modal-form">
        <!-- 步骤指示条 -->
        <div class="modal-step-bar">
          <div class="step-item" :class="{ active: !newDomain.key_id, done: !!newDomain.key_id }">
            <span class="step-dot">1</span>
            <span>选择服务商密钥</span>
          </div>
          <div class="step-connector"></div>
          <div class="step-item" :class="{ active: !!newDomain.key_id }">
            <span class="step-dot">2</span>
            <span>选择托管域名</span>
          </div>
        </div>

        <a-form-item field="key_id" label="服务商密钥">
          <a-select v-model="newDomain.key_id" placeholder="请选择已配置的密钥" @change="handleKeyChange">
            <a-option v-for="key in providerKeys" :key="key.id" :value="key.id">
              <span class="font-medium">{{ key.name }}</span>
              <span class="text-gray-400 ml-2 text-xs">{{ providerNames[key.provider] || key.provider }}</span>
            </a-option>
          </a-select>
        </a-form-item>

        <a-form-item v-if="newDomain.key_id" field="name" label="托管域名">
          <a-select v-model="newDomain.name" placeholder="正在拉取可用域名..." :loading="loadingAvailableDomains">
            <template #empty>
              <div class="text-center py-3 text-gray-400 text-sm">该密钥下暂无可用域名</div>
            </template>
            <a-option v-for="domain in availableDomains" :key="domain.name" :value="domain.name"
              :disabled="domain.added">
              <span>{{ domain.name }}</span>
              <a-tag v-if="domain.added" size="small" color="gray" class="ml-2">已添加</a-tag>
            </a-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { get, post, del } from '@/utils/api'
import { prefetchDomainRecords } from '@/utils/domainRecordsCache'
import { IconPublic, IconPlus, IconInfoCircle } from '@arco-design/web-vue/es/icon'
import { Message, Modal } from '@arco-design/web-vue'
import { useMediaQuery } from '@/composables/useBreakpoint'

const isMobile = useMediaQuery('(max-width: 767px)')

const domains = ref([])
const loading = ref(true)
const showAddDomainModal = ref(false)
const providerKeys = ref([])
const availableDomains = ref([])
const loadingAvailableDomains = ref(false)
const filters = ref({
  keyword: '',
  provider: ''
})
const newDomain = ref({
  name: '',
  key_id: '',
  provider: ''
})

const providerNames = {
  aliyun: '阿里云',
  tencent: '腾讯云',
  cloudflare: 'Cloudflare'
}

const tablePagination = {
  pageSize: 10,
  showTotal: true,
  showPageSize: true,
  pageSizeOptions: [10, 20, 50, 100]
}

const currentPage = ref(1)

const filteredDomains = computed(() => {
  return domains.value.filter(item => {
    const matchProvider = !filters.value.provider || item.provider === filters.value.provider
    const query = filters.value.keyword?.trim().toLowerCase()
    const matchKeyword = !query ||
      item.name?.toLowerCase().includes(query)
    return matchProvider && matchKeyword
  })
})

// Mobile card view paginates client-side (table component already handles it for desktop).
const pagedDomains = computed(() => {
  const start = (currentPage.value - 1) * tablePagination.pageSize
  return filteredDomains.value.slice(start, start + tablePagination.pageSize)
})

// 用户准备进入记录页时提前请求，点击后通常可以直接从缓存渲染。
const preloadDomainRecords = (record) => {
  if (!record?.id) return
  prefetchDomainRecords(record.id).catch(() => {
    // 预加载失败不打断域名列表，进入记录页时仍会正常重试。
  })
}

// 获取域名列表
const fetchDomains = async () => {
  loading.value = true
  try {
    const response = await get('/api/v1/domains')
    domains.value = Array.isArray(response) ? response : []
  } catch (e) {
    Message.error(e.message || '获取域名列表失败')
  } finally {
    loading.value = false
  }
}

// 获取服务商密钥列表
const fetchProviderKeys = async () => {
  try {
    providerKeys.value = await get('/api/v1/provider-keys')
  } catch (e) {
    Message.error(e.message || '获取密钥列表失败')
  }
}

// 获取可用域名列表
const handleKeyChange = async () => {
  newDomain.value.name = ''
  availableDomains.value = []

  if (!newDomain.value.key_id) return

  const selectedKey = providerKeys.value.find(key => key.id === newDomain.value.key_id)
  if (!selectedKey) return

  newDomain.value.provider = selectedKey.provider
  loadingAvailableDomains.value = true

  try {
    const availableDomainList = await get('/api/v1/domains/available', {
      provider: selectedKey.provider
    })
    availableDomains.value = availableDomainList.map(domain => ({
      name: domain,
      added: domains.value.some(d => d.name == domain)
    }))
  } catch (e) {
    Message.error(e.message || '获取可用域名失败')
  } finally {
    loadingAvailableDomains.value = false
  }
}

// 重置添加域名表单
const resetAddDomain = () => {
  newDomain.value = { name: '', key_id: '', provider: '' }
  availableDomains.value = []
}

// 添加域名
const handleAddDomain = async () => {
  try {
    await post('/api/v1/domains', {
      name: newDomain.value.name,
      provider: newDomain.value.provider
    })
    showAddDomainModal.value = false
    newDomain.value = {
      name: '',
      key_id: '',
      provider: ''
    }
    availableDomains.value = []
    await fetchDomains()
  } catch (e) {
    Message.error(e.message || '添加域名失败')
  }
}

// 删除域名
const handleDeleteDomain = (domain) => {
  Modal.warning({
    title: '确认删除',
    content: `确定要删除域名 "${domain.name}" 吗？此操作不可撤销。`,
    okText: '删除',
    cancelText: '取消',
    okButtonProps: { status: 'danger' },
    async onOk() {
      try {
        await del(`/api/v1/domains/${domain.id}`)
        Message.success('删除域名成功')
        await fetchDomains()
      } catch (e) {
        Message.error(e.message || '删除域名失败')
      }
    }
  })
}

onMounted(async () => {
  await Promise.all([fetchDomains(), fetchProviderKeys()])
})
</script>

<style scoped></style>
