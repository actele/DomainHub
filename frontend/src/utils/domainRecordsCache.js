import { get } from './api'

// DNS 服务商接口通常比本地接口慢，短时间内返回页面时优先复用缓存。
// TTL 较短，避免用户长时间看到过期的解析记录；页面仍提供手动刷新。
export const RECORDS_CACHE_TTL = 30 * 1000
export const RECORDS_CACHE_STALE_TTL = 5 * 60 * 1000

const STORAGE_PREFIX = 'domain-manager:records:'
const memoryCache = new Map()
const pendingRequests = new Map()

const cloneRecords = (records) => (Array.isArray(records)
  ? records.map(record => ({ ...record }))
  : [])

const getAuthScope = () => {
  const token = localStorage.getItem('token')
  if (!token) return 'guest'

  // 不把 JWT 放进 storage key，只取 user_id 和 exp 作为当前登录会话的隔离标识。
  try {
    const payloadPart = token.split('.')[1]
    if (!payloadPart || typeof atob !== 'function') return 'session'
    const normalized = payloadPart.replace(/-/g, '+').replace(/_/g, '/')
    const payload = JSON.parse(atob(normalized + '='.repeat((4 - normalized.length % 4) % 4)))
    return `${payload.user_id || 'user'}:${payload.exp || 'session'}`
  } catch {
    return 'session'
  }
}

const getCacheKey = (domainId) => `${getAuthScope()}:${domainId}`
const getStorageKey = (domainId) => `${STORAGE_PREFIX}${getCacheKey(domainId)}`

const isValidEntry = (entry) => Boolean(
  entry &&
  Array.isArray(entry.records) &&
  Number.isFinite(entry.fetchedAt)
)

const readCacheEntry = (domainId) => {
  const cacheKey = getCacheKey(domainId)
  const memoryEntry = memoryCache.get(cacheKey)
  if (isValidEntry(memoryEntry)) return memoryEntry

  try {
    const stored = sessionStorage.getItem(getStorageKey(domainId))
    if (!stored) return null
    const entry = JSON.parse(stored)
    if (!isValidEntry(entry)) return null
    memoryCache.set(cacheKey, entry)
    return entry
  } catch {
    // 隐私模式或禁用 storage 时仍可使用内存缓存。
    return null
  }
}

const writeCacheEntry = (domainId, records) => {
  const entry = {
    records: cloneRecords(records),
    fetchedAt: Date.now()
  }
  const cacheKey = getCacheKey(domainId)
  memoryCache.set(cacheKey, entry)

  try {
    sessionStorage.setItem(getStorageKey(domainId), JSON.stringify(entry))
  } catch {
    // storage 写入失败不影响本次请求结果。
  }

  return entry
}

const cloneEntry = (entry) => entry && ({
  records: cloneRecords(entry.records),
  fetchedAt: entry.fetchedAt
})

export const getCachedDomainRecords = (domainId) => cloneEntry(readCacheEntry(domainId))

export const isDomainRecordsCacheFresh = (entry, now = Date.now()) => Boolean(
  isValidEntry(entry) &&
  now - entry.fetchedAt < RECORDS_CACHE_TTL
)

export const isDomainRecordsCacheUsable = (entry, now = Date.now()) => Boolean(
  isValidEntry(entry) &&
  now - entry.fetchedAt < RECORDS_CACHE_STALE_TTL
)

export const fetchDomainRecords = async (domainId, { force = false } = {}) => {
  const cacheKey = getCacheKey(domainId)
  const cached = readCacheEntry(domainId)

  if (!force && isDomainRecordsCacheFresh(cached)) {
    return {
      records: cloneRecords(cached.records),
      fetchedAt: cached.fetchedAt,
      fromCache: true
    }
  }

  // 预加载和页面加载同时发生时共用一个请求，避免重复调用服务商 API。
  if (pendingRequests.has(cacheKey)) return pendingRequests.get(cacheKey)

  const request = get(`/api/v1/domains/${domainId}/records`)
    .then(records => {
      const entry = writeCacheEntry(domainId, records)
      return {
        records: cloneRecords(entry.records),
        fetchedAt: entry.fetchedAt,
        fromCache: false
      }
    })
    .finally(() => {
      pendingRequests.delete(cacheKey)
    })

  pendingRequests.set(cacheKey, request)
  return request
}

export const prefetchDomainRecords = (domainId) => {
  const cached = readCacheEntry(domainId)
  if (isDomainRecordsCacheFresh(cached)) return Promise.resolve(cloneEntry(cached))
  return fetchDomainRecords(domainId)
}

export const invalidateDomainRecords = (domainId) => {
  const cacheKey = getCacheKey(domainId)
  memoryCache.delete(cacheKey)
  try {
    sessionStorage.removeItem(getStorageKey(domainId))
  } catch {
    // storage 不可用时只清理内存缓存。
  }
}

export const clearDomainRecordsCache = () => {
  memoryCache.clear()
  try {
    const keysToRemove = []
    for (let index = 0; index < sessionStorage.length; index += 1) {
      const key = sessionStorage.key(index)
      if (key?.startsWith(STORAGE_PREFIX)) keysToRemove.push(key)
    }
    keysToRemove.forEach(key => sessionStorage.removeItem(key))
  } catch {
    // storage 不可用时只清理内存缓存。
  }
}
