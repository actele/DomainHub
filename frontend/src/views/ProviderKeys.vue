<template>
  <div class="space-y-6">
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

    <div class="page-toolbar">
      <a-space>
        <a-input v-model="filters.keyword" placeholder="搜索密钥名称" allow-clear style="width: 240px" />
        <a-select v-model="filters.provider" placeholder="全部服务商" allow-clear style="width: 180px">
          <a-option value="aliyun">阿里云</a-option>
          <a-option value="tencent">腾讯云</a-option>
          <a-option value="cloudflare">Cloudflare</a-option>
        </a-select>
      </a-space>
      <div class="flex items-center gap-3">
        <span class="text-sm text-gray-500">共 {{ filteredKeys.length }} 条</span>
        <a-button type="primary" @click="showAddKeyModal = true">
          <template #icon><IconPlus /></template>
          添加密钥
        </a-button>
      </div>
    </div>

    <!-- 密钥列表 -->
    <a-card :bordered="false" class="general-card">
      <a-table :loading="loading" :data="filteredKeys" :pagination="tablePagination">
        <template #columns>
          <a-table-column title="密钥名称" :width="260">
            <template #cell="{ record }">
              <div class="flex items-center gap-2 py-1">
                <IconTool class="text-gray-400" />
                <span class="text-primary font-semibold">{{ record.name }}</span>
              </div>
            </template>
          </a-table-column>
          <a-table-column title="服务商" :width="180">
            <template #cell="{ record }">
              {{ providerNames[record.provider] || record.provider }}
            </template>
          </a-table-column>
          <a-table-column title="凭据摘要" :width="320">
            <template #cell="{ record }">
              {{ getKeyFieldName(record.provider, 'access_key') }}: {{ maskKey(record.access_key) }}
            </template>
          </a-table-column>
          <a-table-column title="创建时间" :width="200">
            <template #cell="{ record }">
              {{ new Date(record.created_at).toLocaleString() }}
            </template>
          </a-table-column>
          <a-table-column title="操作" align="right" :width="180">
            <template #cell="{ record }">
              <div class="table-actions">
                <a-button type="text" status="normal" @click="handleEditKeyName(record)">
                  编辑名称
                </a-button>
                <a-button type="text" status="danger" @click="handleDeleteKey(record)">
                  删除
                </a-button>
              </div>
            </template>
          </a-table-column>
        </template>
        <template #empty>
          <div class="text-center py-8">
            <IconTool class="text-gray-400 text-3xl mb-2" />
            <div class="text-gray-900 font-medium">暂无数据</div>
            <div class="text-gray-500 text-sm mt-1">当前筛选条件下暂无密钥，请调整筛选条件或新增数据</div>
            <a-button type="primary" class="mt-4" @click="showAddKeyModal = true">
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
    <a-modal v-model:visible="showAddKeyModal" title="添加服务商密钥" @ok="handleAddKey" @cancel="resetForm"
      :ok-button-props="{ disabled: !canSubmitNewKey }"
      ok-text="添加" cancel-text="取消" :width="540">
      <a-form :model="newKey" layout="vertical" class="modal-form">
        <!-- 第一步：自定义名称 -->
        <a-form-item field="name" label="密钥名称" required>
          <a-input v-model="newKey.name" placeholder="输入一个便于识别的名称" allow-clear />
        </a-form-item>

        <!-- 第二步：下拉选择服务商 -->
        <a-form-item field="provider" label="服务商">
          <a-select v-model="newKey.provider" placeholder="请选择服务商" @change="() => { newKey.access_key=''; newKey.secret_key='' }">
            <a-option value="aliyun">阿里云</a-option>
            <a-option value="tencent">腾讯云</a-option>
            <a-option value="cloudflare">Cloudflare</a-option>
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
          </div>
        </template>
      </a-form>
    </a-modal>

    <!-- 编辑密钥名称对话框 -->
    <a-modal v-model:visible="showEditNameModal" title="修改密钥名称" @ok="handleUpdateKeyName" @cancel="resetEditForm"
      :ok-button-props="{ disabled: !editKey.newName }" ok-text="保存" cancel-text="取消" :width="420">
      <a-form :model="editKey" layout="vertical" class="modal-form">
        <a-form-item field="newName" label="新名称" required>
          <a-input v-model="editKey.newName" placeholder="请输入新的密钥名称" allow-clear />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { IconPlus, IconTool, IconInfoCircle } from '@arco-design/web-vue/es/icon'
import { get, post, del, put } from '@/utils/api'
import { Message, Modal } from '@arco-design/web-vue'

const keys = ref([])
const loading = ref(true)
const showAddKeyModal = ref(false)
const showEditNameModal = ref(false)
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
}

const handleAddKey = async () => {
  try {
    await post('/api/v1/provider-keys', newKey.value)
    Message.success('添加密钥成功')
    showAddKeyModal.value = false
    resetForm()
    await fetchKeys()
  } catch (e) {
    Message.error('添加密钥失败：' + (e.message || '未知错误'))
  }
}

// 删除密钥
const handleDeleteKey = (key) => {
  Modal.warning({
    title: '确认删除',
    content: `确定要删除密鑰 "${key.name}" 吗？关联的域名将无法使用。`,
    okText: '删除',
    cancelText: '取消',
    okButtonProps: { status: 'danger' },
    async onOk() {
      try {
        await del('/api/v1/provider-keys', { provider: key.provider })
        Message.success('删除密鑰成功')
        await fetchKeys()
      } catch (e) {
        Message.error('删除密鑰失败：' + (e.message || '未知错误'))
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
  try {
    await put(`/api/v1/provider-keys/${editKey.value.id}/name`, {
      new_name: editKey.value.newName
    })
    Message.success('更新密钥名称成功')
    showEditNameModal.value = false
    resetEditForm()
    await fetchKeys()
  } catch (e) {
    Message.error('更新密钥名称失败：' + (e.message || '未知错误'))
  }
}

// 掩码显示密钥
const maskKey = (key) => {
  if (!key) return ''
  return key.slice(0, 4) + '*'.repeat(key.length - 8) + key.slice(-4)
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
</script>

<style scoped></style>
