// @ts-ignore
/* eslint-disable */
import { request } from "@umijs/max";

/** 获取酒店列表 获取所有酒店列表，支持分页 GET /api/admin/hotels */
export async function getAdminHotels(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getAdminHotelsParams,
  options?: { [key: string]: any }
) {
  return request<Record<string, any>>("/api/admin/hotels", {
    method: "GET",
    params: {
      // page has a default value: 1
      page: "1",
      // page_size has a default value: 10
      page_size: "10",
      ...params,
    },
    ...(options || {}),
  });
}

/** 创建酒店 创建新的酒店 POST /api/admin/hotels */
export async function postAdminHotels(
  body: API.CreateHotelRequest,
  options?: { [key: string]: any }
) {
  return request<API.Hotel>("/api/admin/hotels", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    data: body,
    ...(options || {}),
  });
}

/** 获取酒店详情 根据ID获取酒店详情 GET /api/admin/hotels/${param0} */
export async function getAdminHotelsId(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getAdminHotelsIdParams,
  options?: { [key: string]: any }
) {
  const { id: param0, ...queryParams } = params;
  return request<API.Hotel>(`/api/admin/hotels/${param0}`, {
    method: "GET",
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** 更新酒店信息 更新指定酒店的信息 POST /api/admin/hotels/${param0} */
export async function postAdminHotelsId(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.postAdminHotelsIdParams,
  body: API.UpdateHotelRequest,
  options?: { [key: string]: any }
) {
  const { id: param0, ...queryParams } = params;
  return request<API.Hotel>(`/api/admin/hotels/${param0}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    params: { ...queryParams },
    data: body,
    ...(options || {}),
  });
}

/** 删除酒店 删除指定酒店 POST /api/admin/hotels/${param0}/delete */
export async function postAdminHotelsIdOpenApiDelete(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.postAdminHotelsId_openAPI_deleteParams,
  options?: { [key: string]: any }
) {
  const { id: param0, ...queryParams } = params;
  return request<Record<string, any>>(`/api/admin/hotels/${param0}/delete`, {
    method: "POST",
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** 获取酒店设置 获取指定酒店的配置信息 GET /api/admin/settings */
export async function getAdminSettings(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getAdminSettingsParams,
  options?: { [key: string]: any }
) {
  return request<API.HotelSettings>("/api/admin/settings", {
    method: "GET",
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 保存酒店设置 保存或更新酒店配置 POST /api/admin/settings/save */
export async function postAdminSettingsSave(
  body: API.SaveHotelSettingsRequest,
  options?: { [key: string]: any }
) {
  return request<Record<string, any>>("/api/admin/settings/save", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    data: body,
    ...(options || {}),
  });
}
