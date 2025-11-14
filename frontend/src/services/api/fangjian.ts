// @ts-ignore
/* eslint-disable */
import { request } from "@umijs/max";

/** 获取房间列表 获取所有房间列表，支持分页 GET /api/rooms */
export async function getRooms(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getRoomsParams,
  options?: { [key: string]: any }
) {
  return request<API.Room[]>("/api/rooms", {
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

/** 获取房间详情 根据房间ID获取房间详细信息 GET /api/rooms/${param0} */
export async function getRoomsId(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getRoomsIdParams,
  options?: { [key: string]: any }
) {
  const { id: param0, ...queryParams } = params;
  return request<API.Room>(`/api/rooms/${param0}`, {
    method: "GET",
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** 获取可用房间列表 获取所有可用状态的房间列表，支持分页 GET /api/rooms/available */
export async function getRoomsAvailable(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getRoomsAvailableParams,
  options?: { [key: string]: any }
) {
  return request<API.Room[]>("/api/rooms/available", {
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

/** 根据楼层获取房间 根据楼层号获取该楼层的所有房间，支持分页 GET /api/rooms/floor/${param0} */
export async function getRoomsFloorFloor(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getRoomsFloorFloorParams,
  options?: { [key: string]: any }
) {
  const { floor: param0, ...queryParams } = params;
  return request<API.Room[]>(`/api/rooms/floor/${param0}`, {
    method: "GET",
    params: {
      // page has a default value: 1
      page: "1",
      // page_size has a default value: 10
      page_size: "10",
      ...queryParams,
    },
    ...(options || {}),
  });
}

/** 按房型搜索房间 根据房型搜索房间，支持分页 GET /api/rooms/search/type */
export async function getRoomsSearchType(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getRoomsSearchTypeParams,
  options?: { [key: string]: any }
) {
  return request<API.Room[]>("/api/rooms/search/type", {
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
