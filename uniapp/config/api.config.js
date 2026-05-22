/**
 * API配置文件
 */

import { getEnvString, IS_DEV } from '@/utils/env.js'

// 默认 API 地址留空，H5 开发环境走同源 /api + Vite 代理。
// 生产/真机环境可通过 UNI_APP_API_BASE_URL 或本地存储覆盖为完整域名。
const DEFAULT_API_URL = ''
export const API_BASE_URL_STORAGE_KEY = 'gohotel_api_base_url'

const normalizeApiBaseUrl = (value) => {
  if (typeof value !== 'string') return ''
  return value.trim().replace(/\/+$/, '')
}

const getStoredApiBaseUrl = () => {
  try {
    if (typeof uni === 'undefined' || typeof uni.getStorageSync !== 'function') {
      return ''
    }
    return normalizeApiBaseUrl(uni.getStorageSync(API_BASE_URL_STORAGE_KEY))
  } catch (_error) {
    return ''
  }
}

const ENV_API_URL = normalizeApiBaseUrl(
  getEnvString(['UNI_APP_API_BASE_URL', 'VITE_UNI_APP_API_BASE_URL'], '')
)

// 根据环境选择 API 地址，优先使用可配置值
export const API_BASE_URL = getStoredApiBaseUrl() || ENV_API_URL || DEFAULT_API_URL

// API版本
export const API_VERSION = '/api'

// 完整的API地址
export const BASE_URL = API_BASE_URL + API_VERSION

// 请求超时时间（毫秒）
export const REQUEST_TIMEOUT = 10000

// Token存储的key
export const TOKEN_KEY = 'gohotel_token'

// 用户信息存储的key
export const USER_INFO_KEY = 'gohotel_user_info'

// 请求头配置
export const HEADERS = {
  'Content-Type': 'application/json',
}

// 开发环境默认关闭逐条请求日志，避免后端不可达时刷屏影响调试。
export const SHOW_REQUEST_LOG = false
