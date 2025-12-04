/**
 * 全局错误处理器
 */

import logger from './logger.js'

/**
 * 错误处理类
 */
class ErrorHandler {
  constructor() {
    this.errorQueue = []
    this.maxQueueSize = 50
    this.errorHandlers = new Map()
    
    this.init()
  }

  /**
   * 初始化错误捕获
   */
  init() {
    // 注册全局错误处理
    if (typeof uni !== 'undefined') {
      // App全局错误
      uni.onError((error) => {
        this.handleError(error, 'global')
      })

      // 未处理的Promise错误
      uni.onUnhandledRejection((event) => {
        this.handleError(event.reason, 'unhandledRejection')
      })
    }
  }

  /**
   * 处理错误
   * @param {Error|String} error - 错误对象或错误信息
   * @param {String} type - 错误类型
   * @param {Object} context - 上下文信息
   */
  handleError(error, type = 'unknown', context = {}) {
    const errorInfo = this.parseError(error, type, context)
    
    // 添加到队列
    this.errorQueue.push(errorInfo)
    if (this.errorQueue.length > this.maxQueueSize) {
      this.errorQueue.shift()
    }

    // 记录日志
    logger.error(`[${type}] ${errorInfo.message}`, {
      ...errorInfo,
      stack: errorInfo.stack
    })

    // 执行自定义错误处理器
    this.executeHandlers(errorInfo)

    // 根据错误类型显示提示
    this.showErrorTip(errorInfo)

    return errorInfo
  }

  /**
   * 解析错误信息
   */
  parseError(error, type, context) {
    const errorInfo = {
      type,
      message: '',
      stack: '',
      timestamp: new Date().toISOString(),
      page: this.getCurrentPage(),
      context
    }

    if (error instanceof Error) {
      errorInfo.message = error.message
      errorInfo.stack = error.stack || ''
      errorInfo.name = error.name
    } else if (typeof error === 'string') {
      errorInfo.message = error
    } else {
      errorInfo.message = JSON.stringify(error)
    }

    return errorInfo
  }

  /**
   * 获取当前页面
   */
  getCurrentPage() {
    try {
      const pages = getCurrentPages()
      if (pages.length > 0) {
        const currentPage = pages[pages.length - 1]
        return currentPage.route || 'unknown'
      }
    } catch (e) {
      return 'unknown'
    }
    return 'unknown'
  }

  /**
   * 执行自定义错误处理器
   */
  executeHandlers(errorInfo) {
    this.errorHandlers.forEach((handler, name) => {
      try {
        handler(errorInfo)
      } catch (e) {
        console.error(`错误处理器 ${name} 执行失败:`, e)
      }
    })
  }

  /**
   * 显示错误提示
   */
  showErrorTip(errorInfo) {
    // 生产环境显示友好提示
    if (process.env.NODE_ENV === 'production') {
      const userFriendlyMessages = {
        'network': '网络连接失败，请检查网络设置',
        'timeout': '请求超时，请稍后重试',
        'server': '服务器错误，请稍后重试',
        'auth': '登录已过期，请重新登录',
        'permission': '没有权限执行此操作',
        'default': '操作失败，请稍后重试'
      }

      const message = userFriendlyMessages[errorInfo.type] || userFriendlyMessages.default

      uni.showToast({
        title: message,
        icon: 'none',
        duration: 2000
      })
    } else {
      // 开发环境显示详细错误
      console.error('错误详情:', errorInfo)
    }
  }

  /**
   * 注册错误处理器
   * @param {String} name - 处理器名称
   * @param {Function} handler - 处理函数
   */
  registerHandler(name, handler) {
    this.errorHandlers.set(name, handler)
  }

  /**
   * 移除错误处理器
   * @param {String} name - 处理器名称
   */
  removeHandler(name) {
    this.errorHandlers.delete(name)
  }

  /**
   * 获取错误队列
   */
  getErrors(filter = {}) {
    let errors = [...this.errorQueue]

    if (filter.type) {
      errors = errors.filter(err => err.type === filter.type)
    }

    if (filter.page) {
      errors = errors.filter(err => err.page === filter.page)
    }

    if (filter.startTime) {
      errors = errors.filter(err => 
        new Date(err.timestamp) >= new Date(filter.startTime)
      )
    }

    return errors
  }

  /**
   * 清空错误队列
   */
  clearErrors() {
    this.errorQueue = []
  }

  /**
   * 手动上报错误
   */
  async reportError(error, context = {}) {
    const errorInfo = this.parseError(error, 'manual', context)
    
    try {
      // 上报到服务器
      await logger.reportLog({
        level: 'ERROR',
        message: errorInfo.message,
        data: errorInfo
      })
      
      return true
    } catch (e) {
      console.error('错误上报失败:', e)
      return false
    }
  }

  /**
   * Try-Catch包装器
   */
  async tryCatch(fn, errorType = 'tryCatch', context = {}) {
    try {
      return await fn()
    } catch (error) {
      this.handleError(error, errorType, context)
      throw error
    }
  }

  /**
   * 安全执行函数（捕获错误但不抛出）
   */
  async safeExecute(fn, defaultValue = null, errorType = 'safeExecute') {
    try {
      return await fn()
    } catch (error) {
      this.handleError(error, errorType)
      return defaultValue
    }
  }
}

// 创建单例
const errorHandler = new ErrorHandler()

// 导出便捷方法
export default errorHandler

export const handleError = (error, type, context) => 
  errorHandler.handleError(error, type, context)

export const registerErrorHandler = (name, handler) => 
  errorHandler.registerHandler(name, handler)

export const tryCatch = (fn, errorType, context) => 
  errorHandler.tryCatch(fn, errorType, context)

export const safeExecute = (fn, defaultValue, errorType) => 
  errorHandler.safeExecute(fn, defaultValue, errorType)

/**
 * 错误边界装饰器
 */
export function errorBoundary(errorType = 'decorated') {
  return function(target, propertyKey, descriptor) {
    const originalMethod = descriptor.value

    descriptor.value = async function(...args) {
      try {
        return await originalMethod.apply(this, args)
      } catch (error) {
        errorHandler.handleError(error, errorType, {
          method: propertyKey,
          args
        })
        throw error
      }
    }

    return descriptor
  }
}

















