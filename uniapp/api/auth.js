/**
 * 用户认证相关API
 */

import { get, post } from '@/utils/request.js'

/**
 * 用户登录
 * @param {Object} data - 登录信息
 * @param {String} data.username - 用户名
 * @param {String} data.password - 密码
 */
export const login = (data) => {
  return post('/auth/login', data)
}

/**
 * 用户注册
 * @param {Object} data - 注册信息
 * @param {String} data.username - 用户名
 * @param {String} data.password - 密码
 * @param {String} data.email - 邮箱
 * @param {String} data.phone - 手机号
 */
export const register = (data) => {
  return post('/auth/register', data)
}

/**
 * 微信小程序登录
 * @param {Object} data - 微信登录信息
 * @param {String} data.code - 微信授权码
 * @param {String} data.encryptedData - 加密数据
 * @param {String} data.iv - 初始向量
 */
export const wechatLogin = (data) => {
  return post('/auth/wechat-login', data)
}

/**
 * 获取用户信息
 */
export const getUserInfo = () => {
  return get('/auth/me')
}

/**
 * 刷新token
 * @param {String} refreshToken - 刷新token
 */
export const refreshToken = (refreshToken) => {
  return post('/auth/refresh', { refreshToken })
}

/**
 * 退出登录
 */
export const logout = () => {
  return post('/auth/logout')
}
