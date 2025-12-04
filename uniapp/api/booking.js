/**
 * 预订相关API
 */

import { get, post, put, del } from '@/utils/request.js'

/**
 * 创建预订
 * @param {Object} data - 预订信息
 * @param {Number} data.hotelId - 酒店ID
 * @param {Number} data.roomTypeId - 房型ID
 * @param {String} data.checkInDate - 入住日期
 * @param {String} data.checkOutDate - 离店日期
 * @param {Number} data.roomCount - 房间数量
 * @param {String} data.guestName - 客人姓名
 * @param {String} data.guestPhone - 客人电话
 */
export const createBooking = (data) => {
  return post('/bookings', data)
}

/**
 * 获取预订列表
 * @param {Object} params - 查询参数
 * @param {String} params.status - 预订状态（pending/confirmed/cancelled/completed）
 * @param {Number} params.page - 页码
 * @param {Number} params.pageSize - 每页数量
 */
export const getBookingList = (params) => {
  return get('/bookings', params)
}

/**
 * 获取预订详情
 * @param {Number} id - 预订ID
 */
export const getBookingDetail = (id) => {
  return get(`/bookings/${id}`)
}

/**
 * 取消预订
 * @param {Number} id - 预订ID
 * @param {String} reason - 取消原因
 */
export const cancelBooking = (id, reason) => {
  return put(`/bookings/${id}/cancel`, { reason })
}

/**
 * 确认预订
 * @param {Number} id - 预订ID
 */
export const confirmBooking = (id) => {
  return put(`/bookings/${id}/confirm`)
}

/**
 * 计算预订价格
 * @param {Object} params - 计算参数
 * @param {Number} params.hotelId - 酒店ID
 * @param {Number} params.roomTypeId - 房型ID
 * @param {String} params.checkInDate - 入住日期
 * @param {String} params.checkOutDate - 离店日期
 * @param {Number} params.roomCount - 房间数量
 */
export const calculatePrice = (params) => {
  return post('/bookings/calculate-price', params)
}

/**
 * 支付预订
 * @param {Number} id - 预订ID
 * @param {String} paymentMethod - 支付方式（wechat/alipay/card）
 */
export const payBooking = (id, paymentMethod) => {
  return post(`/bookings/${id}/pay`, { paymentMethod })
}

















