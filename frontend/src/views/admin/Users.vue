<template>
  <div class="space-y-4">
    <a-breadcrumb class="page-breadcrumb">
      <a-breadcrumb-item>控制台</a-breadcrumb-item>
      <a-breadcrumb-item>
        <span class="inline-flex items-center gap-1">
          账号管理
          <a-tooltip content="管理系统内所有用户账号、角色与状态" position="right">
            <icon-info-circle class="text-gray-400 cursor-default" style="font-size:13px" />
          </a-tooltip>
        </span>
      </a-breadcrumb-item>
    </a-breadcrumb>

    <div class="page-toolbar">
      <a-space>
        <a-input v-model="filters.keyword" placeholder="搜索用户名" allow-clear style="width: 240px" />
        <a-select v-model="filters.role" placeholder="全部角色" allow-clear style="width: 160px">
          <a-option value="admin">管理员</a-option>
          <a-option value="user">普通用户</a-option>
        </a-select>
        <a-select v-model="filters.status" placeholder="全部状态" allow-clear style="width: 160px">
          <a-option value="active">正常</a-option>
          <a-option value="disabled">已禁用</a-option>
        </a-select>
      </a-space>
      <div class="flex items-center gap-3">
        <span class="text-xs text-slate-400">root 为系统内置账号</span>
        <span class="text-sm text-gray-500">共 {{ filteredUsers.length }} 个账号</span>
        <a-button type="primary" @click="showCreateModal = true">
          <template #icon><icon-plus /></template>
          添加账号
        </a-button>
      </div>
    </div>

    <a-card :bordered="false" class="general-card">
      <a-table :loading="loading" :data="filteredUsers" :pagination="tablePagination">
        <template #columns>
          <a-table-column title="用户名" :width="200">
            <template #cell="{ record }">
              <div class="flex items-center gap-2">
                <div class="w-7 h-7 rounded-full bg-blue-100 flex items-center justify-center text-blue-600 font-bold text-xs">
                  {{ record.username[0].toUpperCase() }}
                </div>
                <span class="font-medium">{{ record.username }}</span>
              </div>
            </template>
          </a-table-column>
          <a-table-column title="角色" :width="140">
            <template #cell="{ record }">
              <a-tag :color="record.role === 'admin' ? 'blue' : 'gray'">
                {{ record.role === 'admin' ? '管理员' : '普通用户' }}
              </a-tag>
            </template>
          </a-table-column>
          <a-table-column title="状态" :width="120">
            <template #cell="{ record }">
              <a-tag :color="record.status === 'active' ? 'green' : 'red'">
                {{ record.status === 'active' ? '正常' : '已禁用' }}
              </a-tag>
            </template>
          </a-table-column>
          <a-table-column title="创建时间" data-index="created_at" :width="200" />
          <a-table-column title="操作" align="right" :width="240">
            <template #cell="{ record }">
              <div class="table-actions">
                <a-button type="text" size="small" @click="openResetPwd(record)">重置密码</a-button>
                <a-button type="text" size="small" @click="toggleRole(record)">
                  {{ record.role === 'admin' ? '降为用户' : '升为管理员' }}
                </a-button>
                <a-button type="text" size="small"
                  :status="record.status === 'active' ? 'danger' : 'normal'"
                  @click="toggleStatus(record)">
                  {{ record.status === 'active' ? '禁用' : '启用' }}
                </a-button>
              </div>
            </template>
          </a-table-column>
        </template>
        <template #empty>
          <div class="text-center py-8 text-gray-400">暂无账号数据</div>
        </template>
      </a-table>
    </a-card>

    <!-- 添加账号 -->
    <a-modal v-model:visible="showCreateModal" title="添加账号" @ok="handleCreate" @cancel="resetCreate"
      :ok-button-props="{ disabled: !newUser.username || !newUser.password }"
      ok-text="创建" cancel-text="取消" :width="460">
      <a-form :model="newUser" layout="vertical" class="modal-form">
        <a-form-item field="username" label="用户名" required>
          <a-input v-model="newUser.username" placeholder="请输入用户名" allow-clear />
        </a-form-item>
        <a-form-item field="password" label="初始密码" required>
          <a-input-password v-model="newUser.password" placeholder="至少 6 位" allow-clear />
        </a-form-item>
        <a-form-item field="role" label="角色">
          <a-radio-group v-model="newUser.role" type="button" style="width:100%">
            <a-radio value="user" style="flex:1;text-align:center">普通用户</a-radio>
            <a-radio value="admin" style="flex:1;text-align:center">管理员</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 重置密码 -->
    <a-modal v-model:visible="showResetPwdModal" title="重置密码" @ok="handleResetPwd" @cancel="resetPwdForm.password = ''"
      :ok-button-props="{ disabled: !resetPwdForm.password || resetPwdForm.password.length < 6 }"
      ok-text="确认重置" cancel-text="取消" :width="420">
      <a-form layout="vertical" class="modal-form">
        <div class="mb-4 text-sm text-gray-500">
          将重置账号 <span class="font-semibold text-gray-800">{{ resetPwdForm.username }}</span> 的密码
        </div>
        <a-form-item label="新密码" required>
          <a-input-password v-model="resetPwdForm.password" placeholder="至少 6 位" allow-clear />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { get, post, put } from '@/utils/api'
import { Message } from '@arco-design/web-vue'
import { IconPlus, IconInfoCircle } from '@arco-design/web-vue/es/icon'

const users = ref([])
const loading = ref(false)
const showCreateModal = ref(false)
const showResetPwdModal = ref(false)

const filters = ref({ keyword: '', role: '', status: '' })
const newUser = ref({ username: '', password: '', role: 'user' })
const resetPwdForm = ref({ id: null, username: '', password: '' })

const tablePagination = { pageSize: 10, showTotal: true, showPageSize: true, pageSizeOptions: [10, 20, 50] }

const filteredUsers = computed(() => {
  return users.value.filter(u => {
    const q = filters.value.keyword?.trim().toLowerCase()
    return (!q || u.username.toLowerCase().includes(q))
      && (!filters.value.role || u.role === filters.value.role)
      && (!filters.value.status || u.status === filters.value.status)
  })
})

const fetchUsers = async () => {
  loading.value = true
  try {
    const result = await get('/api/v1/admin/users')
    // 后端已经过滤内置账号，这里再做一次防御性过滤，避免旧版本接口泄露 root。
    users.value = (Array.isArray(result) ? result : []).filter(user => user.username?.toLowerCase() !== 'root')
  } catch (e) {
    Message.error('获取用户列表失败：' + (e.message || ''))
  } finally {
    loading.value = false
  }
}

const handleCreate = async () => {
  try {
    await post('/api/v1/admin/users', newUser.value)
    Message.success('账号创建成功')
    showCreateModal.value = false
    resetCreate()
    await fetchUsers()
  } catch (e) {
    Message.error('创建账号失败：' + (e.message || ''))
  }
}

const resetCreate = () => {
  newUser.value = { username: '', password: '', role: 'user' }
}

const toggleStatus = async (record) => {
  const newStatus = record.status === 'active' ? 'disabled' : 'active'
  try {
    await put(`/api/v1/admin/users/${record.id}/status`, { status: newStatus })
    Message.success(newStatus === 'active' ? '已启用账号' : '已禁用账号')
    await fetchUsers()
  } catch (e) {
    Message.error('操作失败：' + (e.message || ''))
  }
}

const toggleRole = async (record) => {
  const newRole = record.role === 'admin' ? 'user' : 'admin'
  try {
    await put(`/api/v1/admin/users/${record.id}/role`, { role: newRole })
    Message.success('角色已更新')
    await fetchUsers()
  } catch (e) {
    Message.error('操作失败：' + (e.message || ''))
  }
}

const openResetPwd = (record) => {
  resetPwdForm.value = { id: record.id, username: record.username, password: '' }
  showResetPwdModal.value = true
}

const handleResetPwd = async () => {
  try {
    await put(`/api/v1/admin/users/${resetPwdForm.value.id}/password`, { password: resetPwdForm.value.password })
    Message.success('密码已重置')
    showResetPwdModal.value = false
    resetPwdForm.value.password = ''
  } catch (e) {
    Message.error('重置密码失败：' + (e.message || ''))
  }
}

onMounted(fetchUsers)
</script>

<style scoped></style>
