/**
 * 用户相关API
 */

import { get, post } from '@/utils/request.js'
import { TOKEN_KEY, USER_INFO_KEY } from '@/config/api.config.js'

const MEMBER_LEVEL_LABELS = {
  normal: '普通会员',
  silver: '白银会员',
  gold: '黄金会员',
  platinum: '铂金会员'
}

const formatBalance = (balance) => {
  const amount = Number(balance || 0)
  return Number.isFinite(amount) ? amount.toFixed(2) : '0.00'
}

const normalizeUserProfile = (user = {}) => ({
  ...user,
  nickname: user.real_name || user.username || '用户',
  member_level_name: MEMBER_LEVEL_LABELS[user.level] || MEMBER_LEVEL_LABELS.normal,
  memberLevel: MEMBER_LEVEL_LABELS[user.level] || MEMBER_LEVEL_LABELS.normal,
  is_vip: user.level && user.level !== 'normal',
  balance: formatBalance(user.balance),
  points: Number(user.points || 0),
})

const persistAuth = (result = {}) => {
  if (result.token) {
    uni.setStorageSync(TOKEN_KEY, result.token)
  }

  if (result.user) {
    uni.setStorageSync(USER_INFO_KEY, normalizeUserProfile(result.user))
  }
}

/**
 * 用户登录
 * @param {Object} data - 登录信息
 * @param {String} data.username - 用户名
 * @param {String} data.password - 密码
 */
export const login = async (data) => {
  const result = await post('/auth/login', data)
  persistAuth(result)
  return result
}

/**
 * 微信登录
 * @param {Object} data - 微信登录信息
 * @param {String} data.code
 * @param {String} data.nickname
 * @param {String} data.avatar
 */
export const wechatLogin = async (data) => {
  const result = await post('/auth/wechat-login', data)
  persistAuth(result)
  return result
}

/**
 * 用户注册
 */
export const register = (data) => {
  return post('/auth/register', data)
}

/**
 * 退出登录
 */
export const logout = async () => {
  uni.removeStorageSync(TOKEN_KEY)
  uni.removeStorageSync(USER_INFO_KEY)
  return true
}

/**
 * 获取用户信息
 */
export const getUserInfo = async () => {
  const result = await get('/users/profile')
  const profile = normalizeUserProfile(result)
  uni.setStorageSync(USER_INFO_KEY, profile)
  return profile
}

/**
 * 更新用户信息
 */
export const updateUserInfo = async (data) => {
  const payload = {
    real_name: data.real_name || data.nickname || '',
    avatar: data.avatar || '',
    phone: data.phone || ''
  }
  const result = await post('/users/profile', payload)
  const profile = normalizeUserProfile(result)
  uni.setStorageSync(USER_INFO_KEY, profile)
  return profile
}

/**
 * 修改密码
 */
export const changePassword = (data) => {
  return post('/users/password', {
    old_password: data.old_password || data.oldPassword,
    new_password: data.new_password || data.newPassword,
  })
}

/**
 * 获取用户积分
 */
export const getUserPoints = async () => {
  const profile = await getUserInfo()
  return {
    points: Number(profile.points || 0)
  }
}

/**
 * 获取用户收藏列表
 */
export const getFavorites = async () => {
  return []
}

/**
 * 获取用户优惠券列表
 */
export const getCoupons = async () => {
  return []
}

/**
 * 获取会员信息
 */
export const getMemberInfo = async () => {
  const profile = await getUserInfo()
  return {
    level: profile.level || 'normal',
    member_level_name: profile.member_level_name,
    balance: profile.balance,
    total_spend: profile.total_spend || 0,
    points: profile.points || 0,
  }
}

/**
 * 升级会员
 */
export const upgradeMember = async () => {
  throw new Error('当前版本暂不支持会员升级')
}
