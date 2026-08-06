<template>
  <div class="space-y-4">
    <a-breadcrumb class="page-breadcrumb">
      <a-breadcrumb-item>控制台</a-breadcrumb-item>
      <a-breadcrumb-item>
        <span class="inline-flex items-center gap-1">
          服务商管理
          <a-tooltip content="管理系统支持的 DNS 服务商，禁用后用户无法使用该服务商添加密钥" position="right">
            <icon-info-circle class="text-gray-400 cursor-default" style="font-size:13px" />
          </a-tooltip>
        </span>
      </a-breadcrumb-item>
    </a-breadcrumb>

    <div class="page-toolbar">
      <a-space class="toolbar-filters">
        <a-input v-model="filters.keyword" placeholder="搜索服务商名称" allow-clear class="toolbar-input toolbar-input-search" />
        <a-select v-model="filters.status" placeholder="全部状态" allow-clear class="toolbar-input toolbar-input-provider">
          <a-option :value="true">已启用</a-option>
          <a-option :value="false">已禁用</a-option>
        </a-select>
      </a-space>
      <div class="toolbar-meta">
        <span class="text-sm text-gray-500">共 {{ filteredProviders.length }} 个服务商</span>
        <a-button type="primary" @click="showCreateModal = true">
          <template #icon><icon-plus /></template>
          添加服务商
        </a-button>
      </div>
    </div>

    <a-card v-if="!isMobile" :bordered="false" class="general-card">
      <a-table :loading="loading" :data="filteredProviders" :pagination="false">
        <template #columns>
          <a-table-column title="服务商名称" :width="200">
            <template #cell="{ record }">
              <div class="flex items-center gap-2">
                <div class="w-7 h-7 rounded-lg bg-blue-50 flex items-center justify-center text-blue-600">
                  <icon-cloud />
                </div>
                <span class="font-semibold">{{ record.name }}</span>
              </div>
            </template>
          </a-table-column>
          <a-table-column title="标识符" data-index="type" :width="160">
            <template #cell="{ record }">
              <a-tag color="arcoblue" size="small">{{ record.type }}</a-tag>
            </template>
          </a-table-column>
          <a-table-column title="描述" data-index="description" />
          <a-table-column title="状态" :width="120">
            <template #cell="{ record }">
              <a-switch :model-value="record.enabled" size="small"
                @change="(val) => toggleStatus(record, val)" />
            </template>
          </a-table-column>
          <a-table-column title="操作" align="right" :width="120">
            <template #cell="{ record }">
              <a-button type="text" status="danger" size="small" @click="handleDelete(record)">
                删除
              </a-button>
            </template>
          </a-table-column>
        </template>
        <template #empty>
          <div class="text-center py-8 text-gray-400">暂无服务商</div>
        </template>
      </a-table>
    </a-card>

    <div v-else class="provider-mobile-list">
      <a-spin :loading="loading" class="w-full">
        <div v-if="!filteredProviders.length" class="text-center py-10 text-gray-500 text-sm">
          暂无服务商
        </div>
        <div v-for="record in filteredProviders" :key="record.id" class="list-card">
          <div class="list-card-header">
            <div class="flex items-center gap-3 min-w-0">
              <div class="w-8 h-8 rounded-lg bg-blue-50 flex items-center justify-center text-blue-600">
                <icon-cloud />
              </div>
              <div class="min-w-0">
                <div class="list-card-title">{{ record.name }}</div>
                <div class="list-card-meta">
                  <a-tag color="arcoblue" size="small">{{ record.type }}</a-tag>
                  <span v-if="record.description">{{ record.description }}</span>
                </div>
              </div>
            </div>
            <a-switch :model-value="record.enabled" size="small"
              @change="(val) => toggleStatus(record, val)" />
          </div>
          <div class="list-card-actions">
            <a-button type="text" size="small" status="danger" @click="handleDelete(record)">
              删除
            </a-button>
          </div>
        </div>
      </a-spin>
    </div>

    <!-- 添加服务商 -->
    <a-modal v-model:visible="showCreateModal" title="添加服务商" @ok="handleCreate" @cancel="resetForm"
      :ok-button-props="{ disabled: !newProvider.name || !newProvider.type }"
      ok-text="创建" cancel-text="取消" :width="480">
      <a-form :model="newProvider" layout="vertical" class="modal-form">
        <a-form-item field="name" label="显示名称" required>
          <a-input v-model="newProvider.name" placeholder="如：阿里云" allow-clear />
        </a-form-item>
        <a-form-item field="type" label="标识符" required>
          <a-input v-model="newProvider.type" placeholder="英文小写，如：aliyun" allow-clear />
          <template #extra>
            <span class="text-xs text-gray-400">与后端 provider 实现对应，不可重复</span>
          </template>
        </a-form-item>
        <a-form-item field="description" label="描述">
          <a-input v-model="newProvider.description" placeholder="可选" allow-clear />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { get, post, put, del } from '@/utils/api'
import { Message, Modal } from '@arco-design/web-vue'
import { IconPlus, IconCloud, IconInfoCircle } from '@arco-design/web-vue/es/icon'
import { useMediaQuery } from '@/composables/useBreakpoint'

const isMobile = useMediaQuery('(max-width: 767px)')

const providers = ref([])
const loading = ref(false)
const showCreateModal = ref(false)
const filters = ref({ keyword: '', status: '' })
const newProvider = ref({ name: '', type: '', description: '' })

const filteredProviders = computed(() => {
  return providers.value.filter(p => {
    const q = filters.value.keyword?.trim().toLowerCase()
    const matchKw = !q || p.name.toLowerCase().includes(q) || p.type.toLowerCase().includes(q)
    const matchStatus = filters.value.status === '' || filters.value.status === null
      ? true
      : p.enabled === filters.value.status
    return matchKw && matchStatus
  })
})

const fetchProviders = async () => {
  loading.value = true
  try {
    providers.value = await get('/api/v1/admin/providers')
  } catch (e) {
    Message.error('获取服务商列表失败：' + (e.message || ''))
  } finally {
    loading.value = false
  }
}

const toggleStatus = async (record, enabled) => {
  try {
    await put(`/api/v1/admin/providers/${record.id}/status`, { enabled })
    record.enabled = enabled
    Message.success(enabled ? '已启用' : '已禁用')
  } catch (e) {
    Message.error('操作失败：' + (e.message || ''))
  }
}

const handleCreate = async () => {
  try {
    await post('/api/v1/admin/providers', newProvider.value)
    Message.success('服务商已添加')
    showCreateModal.value = false
    resetForm()
    await fetchProviders()
  } catch (e) {
    Message.error('添加失败：' + (e.message || ''))
  }
}

const resetForm = () => {
  newProvider.value = { name: '', type: '', description: '' }
}

const handleDelete = (record) => {
  Modal.warning({
    title: '确认删除',
    content: `删除服务商「${record.name}」后，已使用该服务商的密钥仍保留，但无法再新增。`,
    okText: '删除',
    cancelText: '取消',
    okButtonProps: { status: 'danger' },
    async onOk() {
      try {
        await del(`/api/v1/admin/providers/${record.id}`)
        Message.success('已删除')
        await fetchProviders()
      } catch (e) {
        Message.error('删除失败：' + (e.message || ''))
      }
    }
  })
}

onMounted(fetchProviders)
</script>

<style scoped></style>
