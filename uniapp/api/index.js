/**
 * API统一导出
 */

import * as hotel from './hotel.js'
import * as booking from './booking.js'
import * as user from './user.js'
import * as banner from './banner.js'
import * as auth from './auth.js'
import * as upload from './upload.js'

export default {
  hotel,
  booking,
  user,
  banner,
  auth,
  upload
}

// 也可以单独导出
export { hotel, booking, user, banner, auth, upload }

















