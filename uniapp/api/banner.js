/**
 * 横幅相关API
 */

import { get, post } from '@/utils/request.js'

/**
 * 获取激活的横幅列表
 * @description 获取激活状态的横幅，用于前端展示
 */
export const getActiveBanners = () => {
  return get('/banners/active')
}

/**
 * 获取横幅详情
 * @param {Number} id - 横幅ID
 */
export const getBannerById = (id) => {
  return get(`/banners/${id}`)
}

/**
 * 获取所有横幅（管理员）
 * @param {Object} params - 查询参数
 * @param {Number} params.page - 页码，默认1
 * @param {Number} params.pageSize - 每页条数，默认10
 */
export const getAllBanners = (params = {}) => {
  return get('/admin/banners', params)
}

/**
 * 创建横幅（管理员）
 * @param {Object} data - 横幅信息
 * @param {String} data.title - 活动横幅标题
 * @param {String} data.subtitle - 活动横幅副标题
 * @param {String} data.temp_url - 临时图片URL
 * @param {String} data.link_url - 点击跳转链接
 * @param {Number} data.sort - 展示顺序
 * @param {String} data.start_time - 活动开始时间
 * @param {String} data.end_time - 活动结束时间
 */
export const createBanner = (data) => {
  return post('/admin/banners', data)
}

/**
 * 更新横幅（管理员）
 * @param {Number} id - 横幅ID
 * @param {Object} data - 横幅信息
 */
export const updateBanner = (id, data) => {
  return post(`/admin/banners/${id}`, data)
}

/**
 * 删除横幅（管理员）
 * @param {Number} id - 横幅ID
 */
export const deleteBanner = (id) => {
  return post(`/admin/banners/${id}/delete`)
}
