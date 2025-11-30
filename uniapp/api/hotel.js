/**
 * 酒店相关API
 */

import { get, post, put, del } from '@/utils/request.js'

/**
 * 获取酒店列表
 * @param {Object} params - 查询参数
 * @param {Number} params.page - 页码
 * @param {Number} params.pageSize - 每页数量
 * @param {String} params.keyword - 搜索关键词
 * @param {String} params.city - 城市
 */
export const getHotelList = (params) => {
  return get('/hotels', params)
}

/**
 * 获取酒店详情
 * @param {Number} id - 酒店ID
 */
export const getHotelDetail = (id) => {
  return get(`/hotels/${id}`)
}

/**
 * 搜索酒店
 * @param {Object} params - 搜索参数
 * @param {String} params.keyword - 关键词
 * @param {String} params.checkInDate - 入住日期
 * @param {String} params.checkOutDate - 离店日期
 * @param {Number} params.minPrice - 最低价格
 * @param {Number} params.maxPrice - 最高价格
 */
export const searchHotels = (params) => {
  return get('/hotels/search', params)
}

/**
 * 获取酒店房型列表
 * @param {Number} hotelId - 酒店ID
 */
export const getRoomTypes = (hotelId) => {
  return get(`/hotels/${hotelId}/rooms`)
}

/**
 * 获取酒店设施
 * @param {Number} hotelId - 酒店ID
 */
export const getHotelFacilities = (hotelId) => {
  return get(`/hotels/${hotelId}/facilities`)
}

/**
 * 获取酒店评价
 * @param {Number} hotelId - 酒店ID
 * @param {Object} params - 分页参数
 */
export const getHotelReviews = (hotelId, params) => {
  return get(`/hotels/${hotelId}/reviews`, params)
}

/**
 * 收藏酒店
 * @param {Number} hotelId - 酒店ID
 */
export const favoriteHotel = (hotelId) => {
  return post(`/hotels/${hotelId}/favorite`)
}

/**
 * 取消收藏酒店
 * @param {Number} hotelId - 酒店ID
 */
export const unfavoriteHotel = (hotelId) => {
  return del(`/hotels/${hotelId}/favorite`)
}








