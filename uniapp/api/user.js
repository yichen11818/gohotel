/**
 * 用户相关API
 */

import { get, post, put, del } from '@/utils/request.js'
import { TOKEN_KEY, USER_INFO_KEY } from '@/config/api.config.js'

/**
 * 用户登录
 * @param {Object} data - 登录信息
 * @param {String} data.phone - 手机号
 * @param {String} data.password - 密码
 */
export const login = async (data) => {
  const result = await post('/auth/login', data)
  
  // 保存token和用户信息
  if (result.token) {
    uni.setStorageSync(TOKEN_KEY, result.token)
  }
  if (result.userInfo) {
    uni.setStorageSync(USER_INFO_KEY, result.userInfo)
  }
  
  return result
}

/**
 * 用户注册
 * @param {Object} data - 注册信息
 * @param {String} data.phone - 手机号
 * @param {String} data.password - 密码
 * @param {String} data.code - 验证码
 * @param {String} data.nickname - 昵称
 */
export const register = (data) => {
  return post('/auth/register', data)
}

/**
 * 退出登录
 */
export const logout = async () => {
  try {
    await post('/auth/logout')
  } finally {
    // 清除本地存储
    uni.removeStorageSync(TOKEN_KEY)
    uni.removeStorageSync(USER_INFO_KEY)
  }
}

/**
 * 发送验证码
 * @param {String} phone - 手机号
 * @param {String} type - 类型（register/login/reset）
 */
export const sendVerifyCode = (phone, type) => {
  return post('/auth/send-code', { phone, type })
}

/**
 * 获取用户信息
 */
export const getUserInfo = () => {
  return get('/user/info')
}

/**
 * 更新用户信息
 * @param {Object} data - 用户信息
 * @param {String} data.nickname - 昵称
 * @param {String} data.avatar - 头像URL
 * @param {String} data.gender - 性别
 * @param {String} data.birthday - 生日
 */
export const updateUserInfo = (data) => {
  return put('/user/info', data)
}

/**
 * 修改密码
 * @param {Object} data - 密码信息
 * @param {String} data.oldPassword - 旧密码
 * @param {String} data.newPassword - 新密码
 */
export const changePassword = (data) => {
  return put('/user/password', data)
}

/**
 * 获取用户积分
 */
export const getUserPoints = () => {
  return get('/user/points')
}

/**
 * 获取用户收藏列表
 * @param {Object} params - 分页参数
 */
export const getFavorites = (params) => {
  return get('/user/favorites', params)
}

/**
 * 获取用户优惠券列表
 * @param {Object} params - 查询参数
 * @param {String} params.status - 状态（unused/used/expired）
 */
export const getCoupons = (params) => {
  return get('/user/coupons', params)
}

/**
 * 获取会员信息
 */
export const getMemberInfo = () => {
  return get('/user/member')
}

/**
 * 升级会员
 * @param {String} level - 会员等级（silver/gold/platinum）
 */
export const upgradeMember = (level) => {
  return post('/user/member/upgrade', { level })
}









