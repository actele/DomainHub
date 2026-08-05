<template>
  <div class="page-stack">
    <a-breadcrumb class="page-breadcrumb">
      <a-breadcrumb-item>控制台</a-breadcrumb-item>
      <a-breadcrumb-item>
        <router-link :to="{ name: 'domains' }" class="hover:text-blue-500">域名管理</router-link>
      </a-breadcrumb-item>
      <a-breadcrumb-item v-if="domain">{{ domain.name }}</a-breadcrumb-item>
      <a-breadcrumb-item>
        <span class="inline-flex items-center gap-1">
          解析记录
          <a-tooltip :content="`管理 ${domain?.name || ''} 的 DNS 解析记录`" position="right">
            <icon-info-circle class="text-gray-400 cursor-default" style="font-size:13px" />
          </a-tooltip>
        </span>
      </a-breadcrumb-item>
    </a-breadcrumb>

    <div class="page-heading">
      <div>
        <div class="page-eyebrow">DNS / RECORDS</div>
        <h1 class="page-title">解析记录</h1>
        <p class="page-subtitle">{{ domain?.name || '当前域名' }} 的 DNS 记录与代理状态。</p>
      </div>
    </div>

    <div class="page-toolbar">
      <a-space>
        <router-link :to="{ name: 'domains' }">
          <a-button size="small">
            <template #icon><icon-left /></template>
            返回
          </a-button>
        </router-link>
        <a-input v-model="filters.keyword" placeholder="搜索主机记录/记录值" allow-clear style="width: 260px" />
        <a-select v-model="filters.type" placeholder="全部记录类型" allow-clear style="width: 180px">
          <a-option value="A">A</a-option>
          <a-option value="AAAA">AAAA</a-option>
          <a-option value="CNAME">CNAME</a-option>
          <a-option value="MX">MX</a-option>
          <a-option value="TXT">TXT</a-option>
          <a-option value="NS">NS</a-option>
          <a-option value="SRV">SRV</a-option>
          <a-option value="CAA">CAA</a-option>
        </a-select>
      </a-space>
      <div class="flex items-center gap-3">
        <span v-if="recordCacheStatus" class="text-xs text-gray-400">{{ recordCacheStatus }}</span>
        <a-button size="small" :loading="refreshing" @click="fetchRecords({ force: true })">
          <template #icon><icon-refresh /></template>
          刷新
        </a-button>
        <span class="text-sm text-gray-500">共 {{ filteredRecords.length }} 条</span>
        <a-button type="primary" @click="showAddRecordModal = true">
          <template #icon><icon-plus /></template>
          添加记录
        </a-button>
      </div>
    </div>

    <!-- 记录列表 -->
    <a-card :bordered="false" class="general-card">
      <a-table :loading="loading" :data="filteredRecords" :pagination="tablePagination">
        <template #columns>
          <a-table-column title="主机记录" data-index="Name" :width="220">
            <template #cell="{ record }">
              <span class="text-primary">{{ record.Name }}</span>
            </template>
          </a-table-column>
          <a-table-column title="记录类型" data-index="Type" :width="120" />
          <a-table-column title="线路类型" data-index="Line" :width="140">
            <template #cell="{ record }">
              {{ record.Line || '默认' }}
            </template>
          </a-table-column>
          <a-table-column title="记录值" data-index="Value" :width="280" />
          <a-table-column title="TTL" data-index="TTL" :width="100">
            <template #cell="{ record }">
              {{ record.TTL }}秒
            </template>
          </a-table-column>
          <a-table-column title="备注" data-index="Remark" :width="180">
            <template #cell="{ record }">
              {{ record.Remark || '-' }}
            </template>
          </a-table-column>
          <a-table-column title="状态" :width="190">
            <template #cell="{ record }">
              <a-space>
                <a-tag :color="record.Enabled ? 'green' : 'gray'">
                  {{ record.Enabled ? '启用' : '禁用' }}
                </a-tag>
                <a-tag v-if="domain?.provider === 'cloudflare'" :color="record.Proxied ? 'orange' : 'gray'"
                  class="cursor-pointer" @click="toggleProxy(record)">
                  <template #icon>
                    <icon-safe />
                  </template>
                  {{ record.Proxied ? '已代理' : '未代理' }}
                </a-tag>
              </a-space>
            </template>
          </a-table-column>
          <a-table-column title="操作" align="right" :width="180">
            <template #cell="{ record }">
              <div class="table-actions">
                <a-button type="text" status="normal" @click="editRecord(record)">编辑</a-button>
                <a-button type="text" status="danger" @click="deleteRecord(record)">删除</a-button>
              </div>
            </template>
          </a-table-column>
        </template>
        <template #empty>
          <div class="text-center py-8">
            <icon-file class="text-gray-400 text-3xl mb-2" />
            <div class="text-gray-900 font-medium">暂无数据</div>
            <div class="text-gray-500 text-sm mt-1">当前筛选条件下暂无解析记录，请调整筛选条件或新增数据</div>
            <a-button type="primary" class="mt-4" @click="showAddRecordModal = true">
              <template #icon>
                <icon-plus />
              </template>
              添加记录
            </a-button>
          </div>
        </template>
      </a-table>

    </a-card>

    <!-- 添加/编辑记录对话框 -->
    <a-modal v-model:visible="showAddRecordModal" :title="currentRecord ? '编辑解析记录' : '添加解析记录'"
      @ok="handleSaveRecord" @cancel="resetRecordForm"
      :ok-button-props="{ disabled: !isFormValid }" ok-text="保存" cancel-text="取消" :width="640">
      <a-form :model="newRecord" layout="vertical" class="modal-form">
        <div class="form-grid-2">
          <a-form-item field="Type" label="记录类型">
            <a-select v-model="newRecord.Type" placeholder="请选择类型">
              <a-option value="A">A &nbsp;—&nbsp; IPv4 地址</a-option>
              <a-option value="AAAA">AAAA &nbsp;—&nbsp; IPv6 地址</a-option>
              <a-option value="CNAME">CNAME &nbsp;—&nbsp; 别名</a-option>
              <a-option value="MX">MX &nbsp;—&nbsp; 邮件交换</a-option>
              <a-option value="TXT">TXT &nbsp;—&nbsp; 文本记录</a-option>
              <a-option value="NS">NS &nbsp;—&nbsp; 域名服务器</a-option>
              <a-option value="SRV">SRV &nbsp;—&nbsp; 服务定位</a-option>
              <a-option value="CAA">CAA &nbsp;—&nbsp; 证书颁发机构</a-option>
            </a-select>
          </a-form-item>
          <a-form-item field="TTL" label="TTL 缓存时间">
            <a-select v-model="newRecord.TTL">
              <a-option :value="600">10 分钟</a-option>
              <a-option :value="1800">30 分钟</a-option>
              <a-option :value="3600">1 小时</a-option>
              <a-option :value="43200">12 小时</a-option>
              <a-option :value="86400">24 小时</a-option>
            </a-select>
          </a-form-item>
        </div>
        <a-form-item field="Name" label="主机记录">
          <a-input v-model="newRecord.Name" placeholder="子域名，如 www；根域名填 @" allow-clear />
        </a-form-item>
        <a-form-item field="Value" label="记录值">
          <a-input v-model="newRecord.Value" :placeholder="getRecordValuePlaceholder(newRecord.Type)" allow-clear />
        </a-form-item>
        <a-form-item field="Remark" label="备注">
          <a-input v-model="newRecord.Remark" placeholder="可选，便于识别此记录的用途" allow-clear />
        </a-form-item>
        <a-form-item :hide-label="true">
          <div class="modal-switch-row">
            <a-switch v-model="newRecord.Enabled" />
            <span class="text-sm font-medium" :class="newRecord.Enabled ? 'text-blue-600' : 'text-gray-400'">
              {{ newRecord.Enabled ? '记录已启用' : '记录已停用' }}
            </span>
          </div>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { get, post, put, del } from '@/utils/api'
import {
  fetchDomainRecords,
  getCachedDomainRecords,
  isDomainRecordsCacheFresh,
  isDomainRecordsCacheUsable
} from '@/utils/domainRecordsCache'
import { Message, Modal } from '@arco-design/web-vue'
import { IconLeft, IconPlus, IconFile, IconSafe, IconInfoCircle, IconRefresh } from '@arco-design/web-vue/es/icon'

const route = useRoute()
const router = useRouter()

const domain = ref(null)
const records = ref([])
const loading = ref(true)
const refreshing = ref(false)
const lastFetchedAt = ref(null)
const recordsFromCache = ref(false)
const showAddRecordModal = ref(false)
const currentRecord = ref(null)
const filters = ref({
  keyword: '',
  type: ''
})

const newRecord = ref({
  Type: 'A',
  Name: '',
  Value: '',
  TTL: 600,
  Remark: '',
  Enabled: true
})

const isFormValid = computed(() => {
  return newRecord.value.Type && newRecord.value.Name && newRecord.value.Value
})

const tablePagination = {
  pageSize: 10,
  showTotal: true,
  showPageSize: true,
  pageSizeOptions: [10, 20, 50, 100]
}

const filteredRecords = computed(() => {
  return records.value.filter(item => {
    const matchType = !filters.value.type || item.Type === filters.value.type
    const query = filters.value.keyword?.trim().toLowerCase()
    const matchKeyword = !query ||
      (item.Name || '').toLowerCase().includes(query) ||
      (item.Value || '').toLowerCase().includes(query)
    return matchType && matchKeyword
  })
})

const recordCacheStatus = computed(() => {
  if (refreshing.value) {
    return records.value.length ? '正在同步最新记录…' : '正在加载记录…'
  }
  if (!lastFetchedAt.value) return ''

  const time = new Date(lastFetchedAt.value).toLocaleTimeString([], {
    hour: '2-digit',
    minute: '2-digit'
  })
  return `${recordsFromCache.value ? '缓存于' : '更新于'} ${time}`
})

const toApiRecordPayload = (record, recordID = '') => {
  return {
    id: recordID || record.Id || '',
    type: record.Type,
    name: record.Name,
    value: record.Value,
    line: record.Line || '默认',
    ttl: record.TTL,
    priority: record.Priority || 0
  }
}

// 获取域名信息
const fetchDomain = async () => {
  try {
    domain.value = await get(`/api/v1/domains/${route.params.id}`)
  } catch (e) {
    Message.error('获取域名信息失败')
  }
}

// 获取记录列表
const fetchRecords = async ({ force = false } = {}) => {
  const domainId = route.params.id
  const cached = getCachedDomainRecords(domainId)
  const canShowCache = !force && isDomainRecordsCacheUsable(cached)

  if (canShowCache) {
    records.value = cached.records
    lastFetchedAt.value = cached.fetchedAt
    recordsFromCache.value = true
    loading.value = false
  } else if (!records.value.length) {
    loading.value = true
  }

  // 新页面有可用缓存时直接显示；缓存已过期则在页面上保留旧数据并后台同步。
  if (!force && canShowCache && isDomainRecordsCacheFresh(cached)) return

  refreshing.value = true
  try {
    const result = await fetchDomainRecords(domainId, { force })
    records.value = result.records
    lastFetchedAt.value = result.fetchedAt
    recordsFromCache.value = result.fromCache
  } catch (e) {
    if (!cached || !records.value.length) {
      Message.error(e.message || '获取记录列表失败')
    } else {
      Message.warning('服务商暂时无法连接，当前显示的是上次缓存记录')
    }
  } finally {
    loading.value = false
    refreshing.value = false
  }
}

// 获取记录值的占位符
const getRecordValuePlaceholder = (type) => {
  const placeholders = {
    A: '例如：192.168.1.1',
    AAAA: '例如：2001:db8::1',
    CNAME: '例如：example.com',
    MX: '例如：mail.example.com',
    TXT: '例如：v=spf1 include:spf.example.com ~all',
    NS: '例如：ns1.example.com',
    SRV: '例如：0 5 5060 sip.example.com',
    CAA: '例如：0 issue "letsencrypt.org"'
  }
  return placeholders[type] || ''
}

// 编辑记录
const editRecord = (record) => {
  currentRecord.value = record
  newRecord.value = { ...record }
  showAddRecordModal.value = true
}

// 重置表单
const resetRecordForm = () => {
  currentRecord.value = null
  newRecord.value = {
    Type: 'A',
    Name: '',
    Value: '',
    TTL: 600,
    Remark: '',
    Enabled: true
  }
}

// 保存记录
const handleSaveRecord = async () => {
  const isEditing = Boolean(currentRecord.value)
  try {
    if (isEditing) {
      await put(
        `/api/v1/domains/${route.params.id}/records/${currentRecord.value.Id}`,
        toApiRecordPayload(newRecord.value, currentRecord.value.Id)
      )
      Message.success('更新记录成功')
    } else {
      await post(`/api/v1/domains/${route.params.id}/records`, toApiRecordPayload(newRecord.value))
      Message.success('添加记录成功')
    }
    showAddRecordModal.value = false
    currentRecord.value = null
    newRecord.value = {
      Type: 'A',
      Name: '',
      Value: '',
      TTL: 600,
      Remark: '',
      Enabled: true
    }
    await fetchRecords({ force: true })
  } catch (e) {
    Message.error(isEditing ? ('更新记录失败：' + (e.message || '未知错误')) : ('添加记录失败：' + (e.message || '未知错误')))
  }
}

// 删除记录
const deleteRecord = (record) => {
  Modal.warning({
    title: '确认删除',
    content: '确定要删除这条记录吗？',
    okText: '确定',
    cancelText: '取消',
    async onOk() {
      try {
        await del(`/api/v1/domains/${route.params.id}/records/${record.Id}`)
        Message.success('删除记录成功')
        await fetchRecords({ force: true })
      } catch (e) {
        Message.error('删除记录失败')
      }
    }
  })
}

// 切换代理状态（仅Cloudflare）
const toggleProxy = async (record) => {
  try {
    await put(`/api/v1/domains/${route.params.id}/records/${record.Id}`, {
      ...toApiRecordPayload(record, record.Id),
      proxied: !record.Proxied
    })
    await fetchRecords({ force: true })
  } catch (e) {
    Message.error('更新代理状态失败')
  }
}

const loadPage = async () => {
  domain.value = null
  records.value = []
  lastFetchedAt.value = null
  recordsFromCache.value = false
  await Promise.all([fetchDomain(), fetchRecords()])
}

onMounted(loadPage)
watch(() => route.params.id, loadPage)
</script>

<style scoped></style>
