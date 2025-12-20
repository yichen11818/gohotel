/**
 * 网络请求封装
 * 基于 uni.request 的二次封装
 * 集成日志记录和性能监控
 */

import { BASE_URL, REQUEST_TIMEOUT, TOKEN_KEY, HEADERS, SHOW_REQUEST_LOG } from '@/config/api.config.js'
import logger from './logger.js'
import { apiMonitor } from './performance.js'
import errorHandler from './error-handler.js'

/**
 * 请求拦截器
 */
const requestInterceptor = (config) => {
  // 添加token
  const token = uni.getStorageSync(TOKEN_KEY)
  if (token) {
    config.header = {
      ...config.header,
      'Authorization': `Bearer ${token}`
    }
  }

  // 打印请求日志
  if (SHOW_REQUEST_LOG) {
    console.log('=== 请求开始 ===')
    console.log('URL:', config.url)
    console.log('Method:', config.method)
    console.log('Data:', config.data)
    console.log('Header:', config.header)
  }

  return config
}

/**
 * 响应拦截器
 */
const responseInterceptor = (response) => {
  const { statusCode, data } = response

  // 打印响应日志
  if (SHOW_REQUEST_LOG) {
    console.log('=== 响应结束 ===')
    console.log('Status:', statusCode)
    console.log('Data:', data)
  }

  // HTTP状态码判断
  if (statusCode >= 200 && statusCode < 300) {
    // 业务状态码判断（根据您的后端约定调整）
    if (data.code === 0 || data.code === 200 || data.success ) {
      return Promise.resolve(data.data || data)
    } else {
      // 业务错误
      handleBusinessError(data)
      return Promise.reject(data)
    }
  } else {
    // HTTP错误
    handleHttpError(statusCode)
    return Promise.reject(response)
  }
}

/**
 * 处理业务错误
 */
const handleBusinessError = (data) => {
  const errorMsg = data.message || data.msg || '请求失败'

  // 特殊错误码处理
  switch (data.code) {
    case 401:
      // 未授权，清除token并跳转登录
      uni.removeStorageSync(TOKEN_KEY)
      uni.showToast({
        title: '请先登录',
        icon: 'none'
      })
      setTimeout(() => {
        uni.reLaunch({
          url: '/pages/login/login'
        })
      }, 1500)
      break
    case 403:
      uni.showToast({
        title: '没有权限',
        icon: 'none'
      })
      break
    case 404:
      uni.showToast({
        title: '请求的资源不存在',
        icon: 'none'
      })
      break
    case 500:
      uni.showToast({
        title: '服务器错误',
        icon: 'none'
      })
      break
    default:
      uni.showToast({
        title: errorMsg,
        icon: 'none'
      })
  }
}

/**
 * 处理HTTP错误
 */
const handleHttpError = (statusCode) => {
  let message = '请求失败'

  switch (statusCode) {
    case 400:
      message = '请求参数错误'
      break
    case 401:
      message = '未授权，请登录'
      uni.removeStorageSync(TOKEN_KEY)
      setTimeout(() => {
        uni.reLaunch({
          url: '/pages/login/login'
        })
      }, 1500)
      break
    case 403:
      message = '拒绝访问'
      break
    case 404:
      message = '请求地址不存在'
      break
    case 500:
      message = '服务器内部错误'
      break
    case 502:
      message = '网关错误'
      break
    case 503:
      message = '服务不可用'
      break
    case 504:
      message = '网关超时'
      break
    default:
      message = `请求失败(${statusCode})`
  }

  uni.showToast({
    title: message,
    icon: 'none'
  })
}

/**
 * 核心请求方法
 */
const request = (options) => {
  // 默认配置
  const defaultConfig = {
    url: BASE_URL + options.url,
    method: options.method || 'GET',
    data: options.data || {},
    header: {
      ...HEADERS,
      ...options.header
    },
    timeout: options.timeout || REQUEST_TIMEOUT,
    dataType: 'json',
    responseType: 'text'
  }

  // 执行请求拦截器
  const config = requestInterceptor(defaultConfig)

  // 开始性能监控
  const perfMonitor = apiMonitor.start(config.url, config.method)
  const startTime = Date.now()

  // 返回Promise
  return new Promise((resolve, reject) => {
    uni.request({
      ...config,
      success: (response) => {
        const duration = Date.now() - startTime

        // 记录性能
        const dataSize = JSON.stringify(response.data).length
        perfMonitor.end(response.statusCode, dataSize)

        // 记录响应日志
        logger.logResponse(config.url, config.method, response.statusCode, response.data, duration)

        // 执行响应拦截器
        responseInterceptor(response)
          .then(data => resolve(data))
          .catch(error => reject(error))
      },
      fail: (error) => {
        const duration = Date.now() - startTime

        // 记录错误日志
        logger.error('请求失败', {
          url: config.url,
          method: config.method,
          error: error.errMsg,
          duration: `${duration}ms`
        })

        // 统一错误处理
        errorHandler.handleError(error, 'network', {
          url: config.url,
          method: config.method
        })

        // 网络错误处理
        if (error.errMsg.includes('timeout')) {
          uni.showToast({
            title: '请求超时',
            icon: 'none'
          })
        } else if (error.errMsg.includes('fail')) {
          uni.showToast({
            title: '网络连接失败',
            icon: 'none'
          })
        }

        reject(error)
      }
    })
  })
}

/**
 * GET请求
 */
export const get = (url, params = {}, options = {}) => {
  return request({
    url,
    method: 'GET',
    data: params,
    ...options
  })
}

/**
 * POST请求
 */
export const post = (url, data = {}, options = {}) => {
  return request({
    url,
    method: 'POST',
    data,
    ...options
  })
}

/**
 * PUT请求
 */
export const put = (url, data = {}, options = {}) => {
  return request({
    url,
    method: 'PUT',
    data,
    ...options
  })
}

/**
 * DELETE请求
 */
export const del = (url, data = {}, options = {}) => {
  return request({
    url,
    method: 'DELETE',
    data,
    ...options
  })
}

/**
 * 文件上传
 */
export const upload = (url, filePath, formData = {}, options = {}) => {
  const token = uni.getStorageSync(TOKEN_KEY)

  return new Promise((resolve, reject) => {
    uni.uploadFile({
      url: BASE_URL + url,
      filePath,
      name: options.name || 'file',
      formData,
      header: {
        'Authorization': token ? `Bearer ${token}` : ''
      },
      success: (response) => {
        const data = JSON.parse(response.data)
        if (data.code === 0 || data.code === 200) {
          resolve(data.data || data)
        } else {
          handleBusinessError(data)
          reject(data)
        }
      },
      fail: (error) => {
        uni.showToast({
          title: '上传失败',
          icon: 'none'
        })
        reject(error)
      }
    })
  })
}

export default {
  get,
  post,
  put,
  del,
  upload
}

