const APP_MODE = typeof __APP_MODE__ === 'string' ? __APP_MODE__ : 'production'
const APP_API_BASE_URL = typeof __APP_API_BASE_URL__ === 'string' ? __APP_API_BASE_URL__ : ''

export const MODE = APP_MODE
export const IS_DEV = MODE === 'development'
export const IS_PROD = MODE === 'production'

export const getEnvString = (key, fallback = '') => {
  const keys = Array.isArray(key) ? key : [key]
  if (keys.includes('UNI_APP_API_BASE_URL') || keys.includes('VITE_UNI_APP_API_BASE_URL')) {
    return APP_API_BASE_URL || fallback
  }
  return fallback
}
