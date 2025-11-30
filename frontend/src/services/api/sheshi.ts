// @ts-ignore
/* eslint-disable */
import { request } from "@umijs/max";

/** 查询所有设施 查询所有设施（分页） GET /api/facilities */
export async function getFacilities(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getFacilitiesParams,
  options?: { [key: string]: any }
) {
  return request<API.Facility[]>("/api/facilities", {
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

/** 创建设施 创建设施 POST /api/facilities */
export async function postFacilities(
  body: API.CreateFacilityRequest,
  options?: { [key: string]: any }
) {
  return request<API.Facility>("/api/facilities", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    data: body,
    ...(options || {}),
  });
}

/** 根据 ID 查找设施 根据 ID 查找设施 GET /api/facilities/${param0} */
export async function getFacilitiesId(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getFacilitiesIdParams,
  options?: { [key: string]: any }
) {
  const { id: param0, ...queryParams } = params;
  return request<API.Facility>(`/api/facilities/${param0}`, {
    method: "GET",
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** 更新设施 更新设施 PUT /api/facilities/${param0} */
export async function putFacilitiesId(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.putFacilitiesIdParams,
  body: API.UpdateFacilityRequest,
  options?: { [key: string]: any }
) {
  const { id: param0, ...queryParams } = params;
  return request<API.Facility>(`/api/facilities/${param0}`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    params: { ...queryParams },
    data: body,
    ...(options || {}),
  });
}

/** 删除设施 删除设施 DELETE /api/facilities/${param0} */
export async function deleteFacilitiesId(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.deleteFacilitiesIdParams,
  options?: { [key: string]: any }
) {
  const { id: param0, ...queryParams } = params;
  return request<Record<string, any>>(`/api/facilities/${param0}`, {
    method: "DELETE",
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** 批量更新设施位置 批量更新设施的位置和尺寸信息 PUT /api/facilities/batch */
export async function putFacilitiesBatch(
  body: API.BatchUpdateFacilitiesRequest,
  options?: { [key: string]: any }
) {
  return request<Record<string, any>>("/api/facilities/batch", {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    data: body,
    ...(options || {}),
  });
}

/** 按楼层查询设施 获取指定楼层的所有设施 GET /api/facilities/floor/${param0} */
export async function getFacilitiesFloorFloor(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getFacilitiesFloorFloorParams,
  options?: { [key: string]: any }
) {
  const { floor: param0, ...queryParams } = params;
  return request<API.Facility[]>(`/api/facilities/floor/${param0}`, {
    method: "GET",
    params: { ...queryParams },
    ...(options || {}),
  });
}
