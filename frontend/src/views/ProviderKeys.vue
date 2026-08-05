<template>
  <div class="page-stack">
    <a-breadcrumb class="page-breadcrumb">
      <a-breadcrumb-item>控制台</a-breadcrumb-item>
      <a-breadcrumb-item>
        <span class="inline-flex items-center gap-1">
          服务商密钥
          <a-tooltip content="集中管理 DNS 服务商访问凭据" position="right">
            <icon-info-circle class="text-gray-400 cursor-default" style="font-size:13px" />
          </a-tooltip>
        </span>
      </a-breadcrumb-item>
    </a-breadcrumb>

    <div class="page-heading">
      <div>
        <div class="page-eyebrow">CREDENTIALS</div>
        <h1 class="page-title">服务商密钥</h1>
        <p class="page-subtitle">集中管理 DNS 服务商访问凭据，保存后即可在域名管理中使用。</p>
      </div>
      <a-button type="primary" @click="openAddKeyModal">
        <template #icon><IconPlus /></template>
        添加密钥
      </a-button>
    </div>

    <a-alert v-if="saveFeedback" type="success" show-icon closable class="page-feedback"
      @close="saveFeedback = null">
      <template #title>{{ saveFeedback.title }}</template>
      {{ saveFeedback.message }}
    </a-alert>

    <div class="page-toolbar">
      <a-space wrap class="toolbar-filters">
        <a-input v-model="filters.keyword" placeholder="搜索密钥名称" allow-clear class="toolbar-input toolbar-input-search" />
        <a-select v-model="filters.provider" placeholder="全部服务商" allow-clear class="toolbar-input toolbar-input-provider">
          <a-option v-for="provider in providerOptions" :key="provider.value" :value="provider.value">
            {{ provider.label }}
          </a-option>
        </a-select>
      </a-space>
      <div class="toolbar-meta">
        <span class="count-chip">
          <IconTool />
          共 {{ filteredKeys.length }} 条
        </span>
      </div>
    </div>

    <!-- 密钥列表 -->
    <a-card :bordered="false" class="general-card">
      <a-table :loading="loading" :data="filteredKeys" :pagination="tablePagination">
        <template #columns>
          <a-table-column title="密钥名称" :width="260">
            <template #cell="{ record }">
              <div class="key-name-cell">
                <div class="key-avatar"><IconTool /></div>
                <div class="min-w-0">
                  <div class="key-name">{{ record.name }}</div>
                  <div class="key-id">ID · {{ record.id }}</div>
                </div>
              </div>
            </template>
          </a-table-column>
          <a-table-column title="服务商" :width="180">
            <template #cell="{ record }">
              <a-tag :color="providerTagColors[record.provider] || 'gray'" size="small">
                {{ providerNames[record.provider] || record.provider }}
              </a-tag>
            </template>
          </a-table-column>
          <a-table-column title="凭据摘要" :width="320">
            <template #cell="{ record }">
              <div class="credential-summary">
                <span class="credential-label">{{ getKeyFieldName(record.provider, 'access_key') }}</span>
                <code>{{ maskKey(record.access_key) }}</code>
              </div>
            </template>
          </a-table-column>
          <a-table-column title="创建时间" :width="200">
            <template #cell="{ record }">
              <span class="muted-cell">{{ new Date(record.created_at).toLocaleString() }}</span>
            </template>
          </a-table-column>
          <a-table-column title="操作" align="right" :width="180">
            <template #cell="{ record }">
              <div class="table-actions">
                <a-button type="text" status="normal" @click="handleEditKeyName(record)">
                  编辑名称
                </a-button>
                <a-button type="text" status="danger" :loading="deletingKeyId === record.id"
                  :disabled="!!deletingKeyId" @click="handleDeleteKey(record)">
                  删除
                </a-button>
              </div>
            </template>
          </a-table-column>
        </template>
        <template #empty>
          <div class="empty-state">
            <div class="empty-state-icon"><IconTool /></div>
            <div class="empty-state-title">暂无密钥</div>
            <div class="empty-state-description">当前筛选条件下暂无密钥，可以先添加一个服务商凭据。</div>
            <a-button type="primary" class="mt-4" @click="openAddKeyModal">
              <template #icon>
                <IconPlus />
              </template>
              添加密钥
            </a-button>
          </div>
        </template>
      </a-table>
    </a-card>

    <!-- 添加密钥对话框 -->
    <a-modal v-model:visible="showAddKeyModal" title="添加服务商密钥" @ok="handleAddKey" @cancel="handleCancelAdd"
      :ok-button-props="{ disabled: !canSubmitNewKey || savingKey, loading: savingKey }"
      ok-text="添加" cancel-text="取消" :width="540">
      <a-form :model="newKey" layout="vertical" class="modal-form">
        <a-alert v-if="formError" type="error" show-icon class="form-feedback mb-4">
          {{ formError }}
        </a-alert>
        <div class="modal-intro">
          <div class="modal-intro-icon"><IconTool /></div>
          <div>
            <div class="modal-intro-title">安全添加凭据</div>
            <div class="modal-intro-text">密钥仅用于访问对应服务商 API，页面不会展示完整 Secret。</div>
          </div>
        </div>

        <!-- 第一步：自定义名称 -->
        <a-form-item field="name" label="密钥名称" required>
          <a-input v-model="newKey.name" placeholder="输入一个便于识别的名称" allow-clear />
        </a-form-item>

        <!-- 第二步：下拉选择服务商 -->
        <a-form-item field="provider" label="服务商">
          <a-select v-model="newKey.provider" placeholder="请选择服务商" @change="handleProviderChange">
            <a-option v-for="provider in providerOptions" :key="provider.value" :value="provider.value">
              {{ provider.label }}
            </a-option>
          </a-select>
        </a-form-item>

        <!-- 第三步：对应服务商凭据 -->
        <template v-if="newKey.provider">
          <div class="modal-credential-block">
            <div class="block-label">API 凭据</div>
            <a-form-item field="access_key" :label="getKeyFieldName(newKey.provider, 'access_key')">
              <a-input v-model="newKey.access_key"
                :placeholder="`请输入 ${getKeyFieldName(newKey.provider, 'access_key')}`"
                allow-clear />
            </a-form-item>
            <a-form-item v-if="showSecretKey" field="secret_key" :label="getKeyFieldName(newKey.provider, 'secret_key')">
              <a-input-password v-model="newKey.secret_key"
                :placeholder="`请输入 ${getKeyFieldName(newKey.provider, 'secret_key')}`"
                allow-clear />
            </a-form-item>
            <div v-else class="credential-hint">Cloudflare 使用 API Token，无需填写 Secret Key。</div>
          </div>
        </template>
      </a-form>
    </a-modal>

    <!-- 编辑密钥名称对话框 -->
    <a-modal v-model:visible="showEditNameModal" title="修改密钥名称" @ok="handleUpdateKeyName" @cancel="resetEditForm"
      :ok-button-props="{ disabled: !canSubmitEditName || savingName, loading: savingName }" ok-text="保存" cancel-text="取消" :width="420">
      <a-form :model="editKey" layout="vertical" class="modal-form">
        <a-form-item field="newName" label="新名称" required>
          <a-input v-model="editKey.newName" placeholder="请输入新的密钥名称" allow-clear />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, computed } from 'vue'
import { IconPlus, IconTool, IconInfoCircle } from '@arco-design/web-vue/es/icon'
import { get, post, del, put } from '@/utils/api'
import { Message, Modal } from '@arco-design/web-vue'

const keys = ref([])
const loading = ref(true)
const showAddKeyModal = ref(false)
const showEditNameModal = ref(false)
const savingKey = ref(false)
const savingName = ref(false)
const deletingKeyId = ref(null)
const formError = ref('')
const saveFeedback = ref(null)
let feedbackTimer = null
const newKey = ref({
  name: '',
  provider: '',
  access_key: '',
  secret_key: ''
})
const editKey = ref({
  provider: '',
  newName: ''
})
const filters = ref({
  keyword: '',
  provider: ''
})

const providerNames = {
  aliyun: '阿里云',
  tencent: '腾讯云',
  cloudflare: 'Cloudflare'
}

const providerOptions = [
  { value: 'aliyun', label: '阿里云' },
  { value: 'tencent', label: '腾讯云' },
  { value: 'cloudflare', label: 'Cloudflare' }
]

const providerTagColors = {
  aliyun: 'orange',
  tencent: 'arcoblue',
  cloudflare: 'purple'
}

const providerKeyNames = {
  aliyun: {
    access_key: 'AccessKey ID',
    secret_key: 'AccessKey Secret'
  },
  tencent: {
    access_key: 'SecretId',
    secret_key: 'SecretKey'
  },
  cloudflare: {
    access_key: 'API Token',
    secret_key: ''
  }
}

const tablePagination = {
  pageSize: 10,
  showTotal: true,
  showPageSize: true,
  pageSizeOptions: [10, 20, 50, 100]
}

const canSubmitNewKey = computed(() => {
  if (!newKey.value.name || !newKey.value.provider || !newKey.value.access_key) {
    return false
  }

  if (!showSecretKey.value) {
    return true
  }

  return !!newKey.value.secret_key
})

const canSubmitEditName = computed(() => !!editKey.value.newName?.trim())

const filteredKeys = computed(() => {
  return keys.value.filter(item => {
    const matchProvider = !filters.value.provider || item.provider === filters.value.provider
    const query = filters.value.keyword?.trim().toLowerCase()
    const matchKeyword = !query || (item.name || '').toLowerCase().includes(query)
    return matchProvider && matchKeyword
  })
})

// 获取密钥列表
const fetchKeys = async () => {
  loading.value = true
  try {
    keys.value = await get('/api/v1/provider-keys')
  } catch (e) {
    Message.error('获取密钥列表失败')
  } finally {
    loading.value = false
  }
}

// 添加密钥
const resetForm = () => {
  newKey.value = {
    name: '',
    provider: '',
    access_key: '',
    secret_key: ''
  }
  formError.value = ''
}

const openAddKeyModal = () => {
  resetForm()
  showAddKeyModal.value = true
}

const handleCancelAdd = () => {
  showAddKeyModal.value = false
  resetForm()
}

const handleProviderChange = () => {
  newKey.value.access_key = ''
  newKey.value.secret_key = ''
  formError.value = ''
}

const handleAddKey = async () => {
  if (!canSubmitNewKey.value || savingKey.value) return

  savingKey.value = true
  formError.value = ''
  try {
    const providerLabel = providerNames[newKey.value.provider] || '服务商'
    await post('/api/v1/provider-keys', newKey.value)
    await fetchKeys()
    showAddKeyModal.value = false
    resetForm()
    setSaveFeedback('密钥已保存', `${providerLabel}凭据已就绪，现在可以去域名管理中选择它。`)
    Message.success('密钥已保存，可立即使用')
  } catch (e) {
    formError.value = e.message || '保存失败，请检查网络或凭据后重试'
    Message.error('添加密钥失败')
  } finally {
    savingKey.value = false
  }
}

const setSaveFeedback = (title, message) => {
  saveFeedback.value = { title, message }
  if (feedbackTimer) clearTimeout(feedbackTimer)
  feedbackTimer = setTimeout(() => {
    saveFeedback.value = null
  }, 5000)
}

// 删除密钥
const handleDeleteKey = (key) => {
  Modal.warning({
    title: '确认删除',
    content: `确定要删除密钥「${key.name}」吗？关联的域名将无法使用。`,
    okText: '删除',
    cancelText: '取消',
    okButtonProps: { status: 'danger' },
    async onOk() {
      if (deletingKeyId.value) return
      deletingKeyId.value = key.id
      try {
        await del('/api/v1/provider-keys', { provider: key.provider })
        await fetchKeys()
        Message.success('密钥已删除')
      } catch (e) {
        Message.error('删除密钥失败：' + (e.message || '未知错误'))
      } finally {
        deletingKeyId.value = null
      }
    }
  })
}

// 编辑密钥名称
const handleEditKeyName = (key) => {
  editKey.value = {
    id: key.id,
    newName: key.name
  }
  showEditNameModal.value = true
}

const resetEditForm = () => {
  editKey.value = {
    provider: '',
    newName: ''
  }
}

const handleUpdateKeyName = async () => {
  if (!canSubmitEditName.value || savingName.value) return

  savingName.value = true
  try {
    await put(`/api/v1/provider-keys/${editKey.value.id}/name`, {
      new_name: editKey.value.newName.trim()
    })
    await fetchKeys()
    showEditNameModal.value = false
    resetEditForm()
    Message.success('密钥名称已更新')
  } catch (e) {
    Message.error('更新密钥名称失败：' + (e.message || '未知错误'))
  } finally {
    savingName.value = false
  }
}

// 掩码显示密钥
const maskKey = (key) => {
  if (!key) return ''
  if (key.includes('****')) return key
  if (key.length <= 8) return '••••••••'
  return `${key.slice(0, 4)}••••${key.slice(-4)}`
}

// 获取当前服务商的字段名称
const getKeyFieldName = (provider, field) => {
  if (!provider) return field === 'access_key' ? 'Access Key' : 'Secret Key'
  return providerKeyNames[provider]?.[field] || ''
}

// 是否显示 Secret Key 字段
const showSecretKey = computed(() => {
  return !newKey.value.provider || providerKeyNames[newKey.value.provider]?.secret_key !== ''
})

onMounted(fetchKeys)
onBeforeUnmount(() => {
  if (feedbackTimer) clearTimeout(feedbackTimer)
})
</script>

<style scoped>
.page-feedback {
  border: 1px solid #bbf7d0;
  background: linear-gradient(135deg, #f0fdf4, #ecfdf5);
}

.key-name-cell {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}

.key-avatar,
.modal-intro-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex: 0 0 auto;
  width: 34px;
  height: 34px;
  border-radius: 10px;
  color: #2563eb;
  background: #eff6ff;
  border: 1px solid #dbeafe;
}

.key-name {
  overflow: hidden;
  color: #0f172a;
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.key-id,
.muted-cell {
  color: #94a3b8;
  font-size: 12px;
}

.credential-summary {
  display: flex;
  align-items: center;
  gap: 8px;
}

.credential-label {
  color: #64748b;
  font-size: 12px;
}

.credential-summary code {
  padding: 3px 7px;
  border-radius: 6px;
  color: #334155;
  background: #f8fafc;
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
  font-size: 12px;
}

.empty-state {
  padding: 48px 20px;
  text-align: center;
}

.empty-state-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  border-radius: 14px;
  color: #64748b;
  background: #f1f5f9;
  font-size: 22px;
}

.empty-state-title {
  margin-top: 14px;
  color: #0f172a;
  font-weight: 600;
}

.empty-state-description {
  margin-top: 5px;
  color: #94a3b8;
  font-size: 13px;
}

.modal-intro {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
  padding: 12px;
  border: 1px solid #dbeafe;
  border-radius: 12px;
  background: linear-gradient(135deg, #f8fbff, #eff6ff);
}

.modal-intro-title {
  color: #1e3a8a;
  font-size: 13px;
  font-weight: 600;
}

.modal-intro-text,
.credential-hint {
  margin-top: 3px;
  color: #64748b;
  font-size: 12px;
  line-height: 1.5;
}

.form-feedback {
  margin-bottom: 16px;
}
</style>
