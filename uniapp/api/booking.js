/**
 * 预订相关 API（按当前后端真实能力适配）
 */

import { get, post } from '@/utils/request.js'

const STATUS_TEXT_MAP = {
  pending: '待确认',
  confirmed: '待入住',
  checkin: '入住中',
  checkout: '已完成',
  cancelled: '已取消',
}

const STATUS_CLASS_MAP = {
  pending: 'pending',
  confirmed: 'confirmed',
  checkin: 'checkin',
  checkout: 'checkout',
  cancelled: 'cancelled',
}

const ensureArray = (value) => {
  if (Array.isArray(value)) {
    return value.filter(Boolean)
  }

  if (typeof value === 'string' && value.trim()) {
    try {
      const parsed = JSON.parse(value)
      return Array.isArray(parsed) ? parsed.filter(Boolean) : []
    } catch (_error) {
      return value
        .split(',')
        .map((item) => item.trim())
        .filter(Boolean)
    }
  }

  return []
}

const ensureNumber = (value, fallback = 0) => {
  const num = Number(value)
  return Number.isFinite(num) ? num : fallback
}

const formatDate = (value) => {
  if (!value) return ''
  return String(value).slice(0, 10)
}

const extractRoomImage = (room = {}) => {
  const images = ensureArray(room.images)
  return images[0] || room.image_url || room.image || 'https://dummyimage.com/720x420/f5f5f5/999999&text=Room'
}

const normalizeBooking = (booking = {}) => {
  const room = booking.room || {}
  const status = booking.status || 'pending'
  const roomId = booking.room_id || room.id || ''
  const roomImage = extractRoomImage(room)

  return {
    id: String(booking.id || ''),
    orderNo: String(booking.booking_number || booking.id || ''),
    status,
    statusText: STATUS_TEXT_MAP[status] || status,
    statusClass: STATUS_CLASS_MAP[status] || 'default',
    roomId: String(roomId),
    roomName: room.room_type || room.name || '房型待确认',
    roomType: room.room_type || room.name || '房型待确认',
    roomNumber: room.room_number || '',
    roomImage,
    hotelName: '七天酒店',
    hotelImage: roomImage,
    checkIn: formatDate(booking.check_in),
    checkOut: formatDate(booking.check_out),
    totalDays: ensureNumber(booking.total_days, 1),
    totalPrice: ensureNumber(booking.actual_price ?? booking.total_price),
    createdAt: booking.created_at || '',
    guestName: booking.guest_name || '',
    guestPhone: booking.guest_phone || '',
    guestIdCard: booking.guest_id_card || '',
    specialRequest: booking.special_request || '',
    paymentStatus: booking.payment_status || 'unpaid',
    canCancel: ['pending', 'confirmed'].includes(status),
    raw: booking,
  }
}

export const createBooking = (data) => {
  const payload = {
    room_id: Number(data.room_id || data.roomId),
    check_in: data.check_in || data.checkIn,
    check_out: data.check_out || data.checkOut,
    guest_name: data.guest_name || data.guestName,
    guest_phone: data.guest_phone || data.guestPhone,
    guest_id_card: data.guest_id_card || data.guestIdCard || '',
    special_request: data.special_request || data.specialRequest || '',
  }

  return post('/bookings', payload)
}

export const getBookingList = async (params = {}) => {
  const result = await get('/bookings/my', {
    page: params.page || 1,
    page_size: params.pageSize || params.page_size || 50,
  })

  const list = (Array.isArray(result) ? result : []).map(normalizeBooking)

  if (!params.status || params.status === 'all') {
    return list
  }

  return list.filter((item) => item.status === params.status)
}

export const getBookingDetail = async (id) => {
  const result = await get(`/bookings/${id}`)
  return normalizeBooking(result)
}

export const cancelBooking = (id, reason = '') => {
  return post(`/bookings/${id}/cancel`, { reason })
}

export const confirmBooking = async () => {
  throw new Error('当前版本暂不支持住客端确认预订')
}

export const calculatePrice = async () => {
  throw new Error('当前版本暂不支持住客端预估价格接口')
}

export const payBooking = async () => {
  throw new Error('当前版本暂不支持住客端支付')
}
