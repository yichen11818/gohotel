/**
 * API配置文件
 */

// 开发环境API地址
const DEV_API_URL = 'http://192.168.1.10:19999'

// 生产环境API地址
const PROD_API_URL = 'http://192.168.1.10:19999'

// 根据环境选择API地址
export const API_BASE_URL = process.env.NODE_ENV === 'production' ? PROD_API_URL : DEV_API_URL

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

// 是否显示请求日志
export const SHOW_REQUEST_LOG = process.env.NODE_ENV === 'development'

