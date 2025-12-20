/**
 * 横幅相关API
 */

import { get } from '@/utils/request.js'

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
