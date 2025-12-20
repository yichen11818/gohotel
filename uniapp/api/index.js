/**
 * API统一导出
 */

import * as hotel from './hotel.js'
import * as booking from './booking.js'
import * as user from './user.js'
import * as banner from './banner.js'

export default {
  hotel,
  booking,
  user,
  banner
}

// 也可以单独导出
export { hotel, booking, user, banner }

















