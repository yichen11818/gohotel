/**
 * 文件上传相关API
 */

import { upload } from '@/utils/request.js'

/**
 * 上传图片
 * @param {String} filePath - 本地文件路径
 * @param {Object} options - 上传选项
 */
export const uploadImage = (filePath, options = {}) => {
  return upload('/upload/image', filePath, {}, {
    name: 'image',
    ...options
  })
}

/**
 * 批量上传图片
 * @param {Array} filePaths - 本地文件路径数组
 */
export const batchUploadImages = async (filePaths) => {
  const uploadPromises = filePaths.map(filePath => uploadImage(filePath))
  return Promise.all(uploadPromises)
}
