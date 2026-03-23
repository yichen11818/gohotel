/**
 * 酒店与房间相关 API（按当前后端真实能力适配为单酒店模型）
 */

import { get } from '@/utils/request.js'

export const DEFAULT_HOTEL_ID = 1

const DEFAULT_HOTEL = {
  id: DEFAULT_HOTEL_ID,
  name: '七天酒店',
  address: '欢迎入住七天酒店',
  phone: '',
  latitude: 30.56,
  longitude: 114.28,
  intro: '欢迎入住七天酒店，享受舒适便捷的入住体验。',
  checkInTime: '14:00',
  checkOutTime: '12:00',
  serviceTime: '24小时服务',
  coverImages: [],
  notices: []
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

const normalizeRoom = (room = {}) => {
  const images = ensureArray(room.images)
  const facilities = ensureArray(room.facilities)
  const mainImage = images[0] || room.image_url || room.image || 'https://dummyimage.com/720x420/f5f5f5/999999&text=Room'

  return {
    id: String(room.id || ''),
    roomNumber: room.room_number || '',
    name: room.room_type || room.name || `房间 ${room.room_number || room.id || ''}`,
    roomType: room.room_type || room.name || '',
    price: ensureNumber(room.price),
    originalPrice: ensureNumber(room.original_price),
    area: ensureNumber(room.area),
    bedType: room.bed_type || '标准床型',
    capacity: ensureNumber(room.capacity, 2),
    description: room.description || '房间整洁舒适，适合短住与商务出行。',
    facilities,
    images: images.length ? images : [mainImage],
    image: mainImage,
    status: room.status || 'available',
    cleanStatus: room.clean_status || '',
    raw: room,
  }
}

const normalizeSettings = (settings = {}, hotelId = DEFAULT_HOTEL_ID) => {
  const profile = settings.hotel_profile || {}
  const bookingRules = settings.booking_rules || {}
  const geo = profile.geo || {}
  const coverImages = ensureArray(profile.cover_images)

  return {
    ...DEFAULT_HOTEL,
    id: hotelId,
    name: profile.name || profile.brand_name || DEFAULT_HOTEL.name,
    address: profile.address || DEFAULT_HOTEL.address,
    phone: profile.front_desk_phone || DEFAULT_HOTEL.phone,
    latitude: ensureNumber(geo.lat, DEFAULT_HOTEL.latitude),
    longitude: ensureNumber(geo.lng, DEFAULT_HOTEL.longitude),
    intro: profile.intro_html || DEFAULT_HOTEL.intro,
    logo: profile.logo_url || '',
    coverImages,
    checkInTime: bookingRules.check_in_time || DEFAULT_HOTEL.checkInTime,
    checkOutTime: bookingRules.check_out_time || DEFAULT_HOTEL.checkOutTime,
    serviceTime: profile.service_time_text || DEFAULT_HOTEL.serviceTime,
  }
}

export const getPublicSettings = async (hotelId = DEFAULT_HOTEL_ID) => {
  const result = await get('/settings/public', { hotel_id: hotelId })
  return normalizeSettings(result, hotelId)
}

export const getActiveNotices = async () => {
  const result = await get('/notices/active')
  return Array.isArray(result) ? result : []
}

export const getHotelDetail = async (hotelId = DEFAULT_HOTEL_ID) => {
  const [settings, rooms, notices] = await Promise.all([
    getPublicSettings(hotelId).catch(() => normalizeSettings({}, hotelId)),
    getRoomTypes(hotelId).catch(() => []),
    getActiveNotices().catch(() => []),
  ])

  const prices = rooms.map((room) => room.price).filter((price) => price > 0)

  return {
    ...settings,
    roomCount: rooms.length,
    lowestPrice: prices.length ? Math.min(...prices) : 0,
    notices,
  }
}

export const getHotelList = async () => {
  return [await getHotelDetail(DEFAULT_HOTEL_ID)]
}

export const searchHotels = async () => {
  return [await getHotelDetail(DEFAULT_HOTEL_ID)]
}

export const getRoomTypes = async (_hotelId = DEFAULT_HOTEL_ID, params = {}) => {
  const requestParams = {
    page: params.page || 1,
    page_size: params.pageSize || params.page_size || 50,
  }

  const result = params.type
    ? await get('/rooms/search/type', {
        ...requestParams,
        type: params.type,
      })
    : await get('/rooms/available', requestParams)

  return (Array.isArray(result) ? result : []).map(normalizeRoom)
}

export const getRoomDetail = async (id) => {
  const result = await get(`/rooms/${id}`)
  return normalizeRoom(result)
}

export const getHotelFacilities = async (hotelId = DEFAULT_HOTEL_ID) => {
  const hotel = await getHotelDetail(hotelId)
  return hotel.services || []
}

export const getHotelReviews = async () => {
  return []
}

export const favoriteHotel = async () => {
  throw new Error('当前版本暂不支持收藏酒店')
}

export const unfavoriteHotel = async () => {
  throw new Error('当前版本暂不支持收藏酒店')
}
