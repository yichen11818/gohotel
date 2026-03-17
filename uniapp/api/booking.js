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
 * 获取我的预订列表
 * @param {Object} params - 查询参数
 * @param {String} params.status - 预订状态（pending/confirmed/cancelled/completed）
 * @param {Number} params.page - 页码
 * @param {Number} params.page_size - 每页数量
 */
export const getMyBookings = (params) => {
  return get('/bookings/my', params)
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
 * @param {Object} data - 取消信息
 * @param {String} data.reason - 取消原因
 */
export const cancelBooking = (id, data) => {
  return post(`/bookings/${id}/cancel`, data)
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

// === 管理员相关API ===

/**
 * 获取所有预订（管理员）
 * @param {Object} params - 查询参数
 * @param {Number} params.page - 页码，默认1
 * @param {Number} params.page_size - 每页数量，默认10
 */
export const getAllBookings = (params = {}) => {
  return get('/admin/bookings', params)
}

/**
 * 根据房间号和状态获取预订列表（管理员）
 * @param {Object} params - 查询参数
 * @param {String} params.room_number - 房间号
 * @param {String} params.status - 预订状态
 */
export const getBookingsByRoom = (params) => {
  return get('/admin/bookings/room', params)
}

/**
 * 通过客人信息搜索预订（管理员）
 * @param {Object} params - 搜索参数
 * @param {String} params.guest_name - 客人姓名
 * @param {String} params.guest_phone - 客人手机号
 * @param {String} params.status - 预订状态
 */
export const searchBookings = (params) => {
  return get('/admin/bookings/search', params)
}

/**
 * 确认预订（管理员）
 * @param {Number} id - 预订ID
 */
export const confirmBookingAdmin = (id) => {
  return post(`/admin/bookings/${id}/confirm`)
}

/**
 * 办理入住（管理员）
 * @param {Number} id - 预订ID
 */
export const checkinBooking = (id) => {
  return post(`/admin/bookings/${id}/checkin`)
}

/**
 * 办理退房（管理员）
 * @param {Number} id - 预订ID
 */
export const checkoutBooking = (id) => {
  return post(`/admin/bookings/${id}/checkout`)
}

















