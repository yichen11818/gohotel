/**
 * 性能监控工具
 */

import logger from './logger.js'

/**
 * 性能监控类
 */
class Performance {
  constructor() {
    this.marks = new Map()
    this.measures = []
  }

  /**
   * 标记时间点
   * @param {String} name - 标记名称
   */
  mark(name) {
    this.marks.set(name, Date.now())
  }

  /**
   * 测量时间差
   * @param {String} name - 测量名称
   * @param {String} startMark - 起始标记
   * @param {String} endMark - 结束标记（可选，默认为当前时间）
   */
  measure(name, startMark, endMark = null) {
    const startTime = this.marks.get(startMark)
    if (!startTime) {
      console.warn(`找不到起始标记: ${startMark}`)
      return null
    }

    const endTime = endMark ? this.marks.get(endMark) : Date.now()
    if (endMark && !endTime) {
      console.warn(`找不到结束标记: ${endMark}`)
      return null
    }

    const duration = endTime - startTime
    const measure = {
      name,
      startMark,
      endMark: endMark || 'now',
      duration,
      timestamp: new Date().toISOString()
    }

    this.measures.push(measure)
    
    // 记录到日志
    logger.logPerformance(name, duration, 'ms')

    return measure
  }

  /**
   * 清除标记
   * @param {String} name - 标记名称，不传则清除所有
   */
  clearMarks(name = null) {
    if (name) {
      this.marks.delete(name)
    } else {
      this.marks.clear()
    }
  }

  /**
   * 清除测量
   */
  clearMeasures() {
    this.measures = []
  }

  /**
   * 获取所有测量结果
   */
  getMeasures() {
    return [...this.measures]
  }

  /**
   * 获取测量统计
   */
  getStatistics() {
    const stats = {}
    
    this.measures.forEach(measure => {
      if (!stats[measure.name]) {
        stats[measure.name] = {
          count: 0,
          total: 0,
          min: Infinity,
          max: -Infinity,
          avg: 0
        }
      }

      const stat = stats[measure.name]
      stat.count++
      stat.total += measure.duration
      stat.min = Math.min(stat.min, measure.duration)
      stat.max = Math.max(stat.max, measure.duration)
      stat.avg = stat.total / stat.count
    })

    return stats
  }

  /**
   * 监控函数执行时间
   * @param {Function} fn - 要监控的函数
   * @param {String} name - 函数名称
   */
  async monitorFunction(fn, name) {
    const startMark = `${name}_start`
    const endMark = `${name}_end`

    this.mark(startMark)
    
    try {
      const result = await fn()
      this.mark(endMark)
      this.measure(name, startMark, endMark)
      return result
    } catch (error) {
      this.mark(endMark)
      this.measure(`${name}_error`, startMark, endMark)
      throw error
    }
  }

  /**
   * 创建一个计时器
   * @param {String} name - 计时器名称
   */
  timer(name) {
    const startTime = Date.now()
    
    return {
      stop: () => {
        const duration = Date.now() - startTime
        logger.logPerformance(name, duration, 'ms')
        return duration
      }
    }
  }
}

// 创建单例
const performance = new Performance()

/**
 * 装饰器：监控方法执行时间
 */
export function measureTime(name) {
  return function(target, propertyKey, descriptor) {
    const originalMethod = descriptor.value

    descriptor.value = async function(...args) {
      const timer = performance.timer(name || propertyKey)
      try {
        const result = await originalMethod.apply(this, args)
        timer.stop()
        return result
      } catch (error) {
        timer.stop()
        throw error
      }
    }

    return descriptor
  }
}

/**
 * 监控API请求性能
 */
export class ApiPerformanceMonitor {
  constructor() {
    this.requests = []
  }

  /**
   * 开始监控请求
   */
  start(url, method) {
    const id = `${method}_${url}_${Date.now()}`
    const startTime = Date.now()
    
    return {
      id,
      url,
      method,
      startTime,
      end: (statusCode, dataSize) => {
        const duration = Date.now() - startTime
        const request = {
          id,
          url,
          method,
          statusCode,
          dataSize,
          duration,
          timestamp: new Date().toISOString()
        }
        
        this.requests.push(request)
        logger.logPerformance(`API_${method}_${url}`, duration, 'ms')
        
        return request
      }
    }
  }

  /**
   * 获取慢请求
   */
  getSlowRequests(threshold = 1000) {
    return this.requests.filter(req => req.duration > threshold)
  }

  /**
   * 获取请求统计
   */
  getStatistics() {
    const stats = {
      total: this.requests.length,
      avgDuration: 0,
      slowRequests: 0,
      byStatus: {}
    }

    let totalDuration = 0
    this.requests.forEach(req => {
      totalDuration += req.duration
      
      if (req.duration > 1000) {
        stats.slowRequests++
      }

      stats.byStatus[req.statusCode] = (stats.byStatus[req.statusCode] || 0) + 1
    })

    stats.avgDuration = this.requests.length > 0 ? totalDuration / this.requests.length : 0

    return stats
  }

  /**
   * 清空记录
   */
  clear() {
    this.requests = []
  }
}

export default performance

// 导出监控器实例
export const apiMonitor = new ApiPerformanceMonitor()

















