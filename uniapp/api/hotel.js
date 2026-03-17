/**
 * 酒店相关API
 */

import { get, post, put, del } from '@/utils/request.js'

// === 房间相关API ===

/**
 * 获取所有房间列表
 * @param {Object} params - 查询参数
 * @param {Number} params.page - 页码，默认1
 * @param {Number} params.page_size - 每页数量，默认10
 */
export const getRoomList = (params = {}) => {
  return get('/rooms', params)
}

/**
 * 获取可用房间列表
 * @param {Object} params - 查询参数
 * @param {Number} params.page - 页码，默认1
 * @param {Number} params.page_size - 每页数量，默认10
 */
export const getAvailableRooms = (params = {}) => {
  return get('/rooms/available', params)
}

/**
 * 根据楼层获取房间
 * @param {Number} floor - 楼层号
 * @param {Object} params - 查询参数
 * @param {Number} params.page - 页码，默认1
 * @param {Number} params.page_size - 每页数量，默认10
 */
export const getRoomsByFloor = (floor, params = {}) => {
  return get(`/rooms/floor/${floor}`, params)
}

/**
 * 根据房型搜索房间
 * @param {Object} params - 搜索参数
 * @param {String} params.room_type - 房型
 * @param {Number} params.page - 页码，默认1
 * @param {Number} params.page_size - 每页数量，默认10
 */
export const searchRoomsByType = (params) => {
  return get('/rooms/search/type', params)
}

/**
 * 获取房间详情
 * @param {Number} id - 房间ID
 */
export const getRoomDetail = (id) => {
  return get(`/rooms/${id}`)
}

/**
 * 创建房间（管理员）
 * @param {Object} data - 房间信息
 */
export const createRoom = (data) => {
  return post('/rooms', data)
}

/**
 * 更新房间（管理员）
 * @param {Number} id - 房间ID
 * @param {Object} data - 房间信息
 */
export const updateRoom = (id, data) => {
  return post(`/rooms/${id}`, data)
}

/**
 * 删除房间（管理员）
 * @param {Number} id - 房间ID
 */
export const deleteRoom = (id) => {
  return post(`/rooms/${id}/delete`)
}

/**
 * 批量创建房间（管理员）
 * @param {Object} data - 批量房间信息
 */
export const batchCreateRooms = (data) => {
  return post('/rooms/batch', data)
}

// === 设施相关API ===

/**
 * 获取所有设施（管理员）
 * @param {Object} params - 查询参数
 * @param {Number} params.page - 页码，默认1
 * @param {Number} params.page_size - 每页数量，默认10
 */
export const getFacilities = (params = {}) => {
  return get('/admin/facilities', params)
}

/**
 * 根据楼层获取设施（管理员）
 * @param {Number} floor - 楼层号
 */
export const getFacilitiesByFloor = (floor) => {
  return get(`/admin/facilities/floor/${floor}`)
}

/**
 * 获取设施详情（管理员）
 * @param {Number} id - 设施ID
 */
export const getFacilityDetail = (id) => {
  return get(`/admin/facilities/${id}`)
}

/**
 * 创建设施（管理员）
 * @param {Object} data - 设施信息
 */
export const createFacility = (data) => {
  return post('/admin/facilities', data)
}

/**
 * 更新设施（管理员）
 * @param {Number} id - 设施ID
 * @param {Object} data - 设施信息
 */
export const updateFacility = (id, data) => {
  return post(`/admin/facilities/${id}`, data)
}

/**
 * 批量更新设施位置（管理员）
 * @param {Object} data - 批量设施位置信息
 */
export const batchUpdateFacilities = (data) => {
  return post('/admin/facilities/batch', data)
}

















