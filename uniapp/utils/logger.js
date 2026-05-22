/**
 * 日志工具类
 * 支持：日志分级、本地存储、上报服务器、性能监控
 */

import { getStorage, setStorage } from './storage.js'
import { post } from './request.js'
import { IS_PROD } from './env.js'

// 日志级别
export const LogLevel = {
  DEBUG: 0,
  INFO: 1,
  WARN: 2,
  ERROR: 3,
  FATAL: 4
}

// 日志级别名称
const LogLevelNames = {
  0: 'DEBUG',
  1: 'INFO',
  2: 'WARN',
  3: 'ERROR',
  4: 'FATAL'
}

// 日志配置
const config = {
  // 当前环境的日志级别
  level: IS_PROD ? LogLevel.INFO : LogLevel.DEBUG,
  
  // 是否启用控制台输出
  enableConsole: !IS_PROD,
  
  // 是否启用本地存储
  enableStorage: true,
  
  // 是否启用远程上报
  enableReport: IS_PROD,
  
  // 本地存储的最大日志条数
  maxStorageSize: 100,
  
  // 日志上报的API地址
  reportUrl: '/logs/report',
  
  // 自动上报的日志级别（ERROR及以上）
  autoReportLevel: LogLevel.ERROR,
  
  // 本地存储的key
  storageKey: 'gohotel_logs'
}

/**
 * 日志管理类
 */
class Logger {
  constructor() {
    this.logs = []
    this.loadLogsFromStorage()
    this.initErrorCapture()
  }

  /**
   * 从本地存储加载日志
   */
  loadLogsFromStorage() {
    if (config.enableStorage) {
      const logs = getStorage(config.storageKey, [])
      this.logs = Array.isArray(logs) ? logs : []
    }
  }

  /**
   * 保存日志到本地存储
   */
  saveLogsToStorage() {
    if (config.enableStorage) {
      // 只保留最新的日志
      const logsToSave = this.logs.slice(-config.maxStorageSize)
      setStorage(config.storageKey, logsToSave)
    }
  }

  /**
   * 初始化全局错误捕获
   */
  initErrorCapture() {
    // 捕获 Vue 错误
    if (uni && uni.onError) {
      uni.onError((error) => {
        this.error('全局错误', { error: error.toString() })
      })
    }

    // 捕获未处理的 Promise 错误
    if (uni && uni.onUnhandledRejection) {
      uni.onUnhandledRejection((event) => {
        this.error('未处理的Promise错误', { 
          reason: event.reason?.toString() || 'Unknown reason'
        })
      })
    }
  }

  /**
   * 记录日志
   * @param {Number} level - 日志级别
   * @param {String} message - 日志消息
   * @param {Object} data - 附加数据
   */
  log(level, message, data = {}) {
    // 级别过滤
    if (level < config.level) {
      return
    }

    // 创建日志对象
    const logEntry = {
      level: LogLevelNames[level],
      levelValue: level,
      message,
      data,
      timestamp: new Date().toISOString(),
      platform: this.getPlatform(),
      page: this.getCurrentPage(),
      userAgent: this.getUserAgent()
    }

    // 控制台输出
    if (config.enableConsole) {
      this.consoleLog(level, logEntry)
    }

    // 添加到日志列表
    this.logs.push(logEntry)

    // 保存到本地存储
    this.saveLogsToStorage()

    // 自动上报错误日志
    if (config.enableReport && level >= config.autoReportLevel) {
      this.reportLog(logEntry)
    }

    return logEntry
  }

  /**
   * 控制台输出
   */
  consoleLog(level, logEntry) {
    const style = this.getConsoleStyle(level)
    const time = new Date(logEntry.timestamp).toLocaleTimeString()
    const prefix = `[${logEntry.level}] [${time}]`
    
    switch (level) {
      case LogLevel.DEBUG:
        console.debug(prefix, logEntry.message, logEntry.data)
        break
      case LogLevel.INFO:
        console.info(prefix, logEntry.message, logEntry.data)
        break
      case LogLevel.WARN:
        console.warn(prefix, logEntry.message, logEntry.data)
        break
      case LogLevel.ERROR:
      case LogLevel.FATAL:
        console.error(prefix, logEntry.message, logEntry.data)
        break
      default:
        console.log(prefix, logEntry.message, logEntry.data)
    }
  }

  /**
   * 获取控制台样式
   */
  getConsoleStyle(level) {
    const styles = {
      [LogLevel.DEBUG]: 'color: #999',
      [LogLevel.INFO]: 'color: #2196F3',
      [LogLevel.WARN]: 'color: #FF9800',
      [LogLevel.ERROR]: 'color: #F44336',
      [LogLevel.FATAL]: 'color: #F44336; font-weight: bold'
    }
    return styles[level] || ''
  }

  /**
   * 上报日志到服务器
   */
  async reportLog(logEntry) {
    try {
      await post(config.reportUrl, {
        logs: [logEntry],
        deviceInfo: this.getDeviceInfo()
      })
    } catch (error) {
      console.error('日志上报失败:', error)
    }
  }

  /**
   * 批量上报日志
   * @param {Array} logs - 日志数组，默认上报所有日志
   */
  async reportLogs(logs = null) {
    try {
      const logsToReport = logs || this.logs
      if (logsToReport.length === 0) {
        return
      }

      await post(config.reportUrl, {
        logs: logsToReport,
        deviceInfo: this.getDeviceInfo()
      })

      // 上报成功后清空日志
      if (!logs) {
        this.clearLogs()
      }

      return true
    } catch (error) {
      console.error('批量上报日志失败:', error)
      return false
    }
  }

  /**
   * 获取平台信息
   */
  getPlatform() {
    // #ifdef H5
    return 'H5'
    // #endif
    
    // #ifdef MP-WEIXIN
    return '微信小程序'
    // #endif
    
    // #ifdef MP-ALIPAY
    return '支付宝小程序'
    // #endif
    
    // #ifdef APP-PLUS
    return 'App'
    // #endif
    
    return 'Unknown'
  }

  /**
   * 获取当前页面路径
   */
  getCurrentPage() {
    try {
      const pages = getCurrentPages()
      if (pages.length > 0) {
        const currentPage = pages[pages.length - 1]
        return currentPage.route || currentPage.$page?.fullPath || 'Unknown'
      }
    } catch (error) {
      console.error('获取页面路径失败:', error)
    }
    return 'Unknown'
  }

  /**
   * 获取用户代理信息
   */
  getUserAgent() {
    try {
      const systemInfo = uni.getSystemInfoSync()
      return `${systemInfo.platform} ${systemInfo.system} ${systemInfo.model}`
    } catch (error) {
      return 'Unknown'
    }
  }

  /**
   * 获取设备信息
   */
  getDeviceInfo() {
    try {
      const systemInfo = uni.getSystemInfoSync()
      return {
        platform: systemInfo.platform,
        system: systemInfo.system,
        model: systemInfo.model,
        brand: systemInfo.brand,
        version: systemInfo.version,
        screenWidth: systemInfo.screenWidth,
        screenHeight: systemInfo.screenHeight,
        pixelRatio: systemInfo.pixelRatio,
        language: systemInfo.language
      }
    } catch (error) {
      return {}
    }
  }

  /**
   * DEBUG 级别日志
   */
  debug(message, data) {
    return this.log(LogLevel.DEBUG, message, data)
  }

  /**
   * INFO 级别日志
   */
  info(message, data) {
    return this.log(LogLevel.INFO, message, data)
  }

  /**
   * WARN 级别日志
   */
  warn(message, data) {
    return this.log(LogLevel.WARN, message, data)
  }

  /**
   * ERROR 级别日志
   */
  error(message, data) {
    return this.log(LogLevel.ERROR, message, data)
  }

  /**
   * FATAL 级别日志
   */
  fatal(message, data) {
    return this.log(LogLevel.FATAL, message, data)
  }

  /**
   * 记录API请求
   */
  logRequest(url, method, params, duration) {
    this.info('API请求', {
      url,
      method,
      params,
      duration: `${duration}ms`
    })
  }

  /**
   * 记录API响应
   */
  logResponse(url, method, statusCode, data, duration) {
    if (statusCode >= 200 && statusCode < 300) {
      this.debug('API响应成功', {
        url,
        method,
        statusCode,
        duration: `${duration}ms`,
        dataSize: JSON.stringify(data).length
      })
    } else {
      this.error('API响应失败', {
        url,
        method,
        statusCode,
        data,
        duration: `${duration}ms`
      })
    }
  }

  /**
   * 记录页面访问
   */
  logPageView(pagePath, query = {}) {
    this.info('页面访问', {
      path: pagePath,
      query,
      referrer: this.getCurrentPage()
    })
  }

  /**
   * 记录用户行为
   */
  logUserAction(action, target, data = {}) {
    this.info('用户行为', {
      action,
      target,
      ...data
    })
  }

  /**
   * 记录性能指标
   */
  logPerformance(metric, value, unit = 'ms') {
    this.info('性能指标', {
      metric,
      value,
      unit
    })
  }

  /**
   * 获取所有日志
   */
  getLogs(filter = {}) {
    let logs = [...this.logs]

    // 按级别过滤
    if (filter.level !== undefined) {
      logs = logs.filter(log => log.levelValue >= filter.level)
    }

    // 按时间过滤
    if (filter.startTime) {
      logs = logs.filter(log => new Date(log.timestamp) >= new Date(filter.startTime))
    }
    if (filter.endTime) {
      logs = logs.filter(log => new Date(log.timestamp) <= new Date(filter.endTime))
    }

    // 按页面过滤
    if (filter.page) {
      logs = logs.filter(log => log.page === filter.page)
    }

    return logs
  }

  /**
   * 清空日志
   */
  clearLogs() {
    this.logs = []
    this.saveLogsToStorage()
  }

  /**
   * 获取日志统计
   */
  getStatistics() {
    const stats = {
      total: this.logs.length,
      byLevel: {},
      byPage: {},
      recentErrors: []
    }

    // 按级别统计
    this.logs.forEach(log => {
      const level = log.level
      stats.byLevel[level] = (stats.byLevel[level] || 0) + 1

      // 按页面统计
      const page = log.page
      stats.byPage[page] = (stats.byPage[page] || 0) + 1

      // 收集最近的错误
      if (log.levelValue >= LogLevel.ERROR && stats.recentErrors.length < 10) {
        stats.recentErrors.push(log)
      }
    })

    return stats
  }

  /**
   * 设置配置
   */
  setConfig(options) {
    Object.assign(config, options)
  }

  /**
   * 获取配置
   */
  getConfig() {
    return { ...config }
  }
}

// 创建单例
const logger = new Logger()

// 导出实例和方法
export default logger

export const debug = (message, data) => logger.debug(message, data)
export const info = (message, data) => logger.info(message, data)
export const warn = (message, data) => logger.warn(message, data)
export const error = (message, data) => logger.error(message, data)
export const fatal = (message, data) => logger.fatal(message, data)

export const logRequest = (url, method, params, duration) => 
  logger.logRequest(url, method, params, duration)

export const logResponse = (url, method, statusCode, data, duration) => 
  logger.logResponse(url, method, statusCode, data, duration)

export const logPageView = (pagePath, query) => 
  logger.logPageView(pagePath, query)

export const logUserAction = (action, target, data) => 
  logger.logUserAction(action, target, data)

export const logPerformance = (metric, value, unit) => 
  logger.logPerformance(metric, value, unit)

