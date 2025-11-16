// @ts-ignore
/* eslint-disable */
import { request } from "@umijs/max";

/** 获取所有预订（管理员） 管理员获取所有预订列表，支持分页 GET /api/admin/bookings */
export async function getAdminBookings(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getAdminBookingsParams,
  options?: { [key: string]: any }
) {
  return request<API.Booking[]>("/api/admin/bookings", {
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

/** 获取用户列表（管理员） 管理员获取所有用户列表，支持分页 GET /api/admin/users */
export async function getAdminUsers(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getAdminUsersParams,
  options?: { [key: string]: any }
) {
  return request<API.User[]>("/api/admin/users", {
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

/** 获取用户详情（管理员） 管理员根据用户ID获取用户详细信息 GET /api/admin/users/${param0} */
export async function getAdminUsersId(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getAdminUsersIdParams,
  options?: { [key: string]: any }
) {
  const { id: param0, ...queryParams } = params;
  return request<API.User>(`/api/admin/users/${param0}`, {
    method: "GET",
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** 添加用户 管理员添加新的用户账户，默认密码为yumi123456 POST /api/admin/users/user */
export async function postAdminUsersUser(
  body: API.AddUserRequest,
  options?: { [key: string]: any }
) {
  return request<API.User>("/api/admin/users/user", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    data: body,
    ...(options || {}),
  });
}

/** 办理入住（管理员） 管理员为已确认的预订办理入住 POST /api/bookings/${param0}/checkin */
export async function postBookingsIdCheckin(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.postBookingsIdCheckinParams,
  options?: { [key: string]: any }
) {
  const { id: param0, ...queryParams } = params;
  return request<Record<string, any>>(`/api/bookings/${param0}/checkin`, {
    method: "POST",
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** 办理退房（管理员） 管理员为入住中的预订办理退房 POST /api/bookings/${param0}/checkout */
export async function postBookingsIdCheckout(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.postBookingsIdCheckoutParams,
  options?: { [key: string]: any }
) {
  const { id: param0, ...queryParams } = params;
  return request<Record<string, any>>(`/api/bookings/${param0}/checkout`, {
    method: "POST",
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** 确认预订（管理员） 管理员确认待处理的预订 POST /api/bookings/${param0}/confirm */
export async function postBookingsIdConfirm(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.postBookingsIdConfirmParams,
  options?: { [key: string]: any }
) {
  const { id: param0, ...queryParams } = params;
  return request<Record<string, any>>(`/api/bookings/${param0}/confirm`, {
    method: "POST",
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** 创建房间（管理员） 管理员创建新房间 POST /api/rooms */
export async function postRooms(
  body: API.CreateRoomRequest,
  options?: { [key: string]: any }
) {
  return request<API.Room>("/api/rooms", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    data: body,
    ...(options || {}),
  });
}

/** 更新房间（管理员） 管理员更新房间信息 POST /api/rooms/${param0} */
export async function postRoomsId(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.postRoomsIdParams,
  body: API.UpdateRoomRequest,
  options?: { [key: string]: any }
) {
  const { id: param0, ...queryParams } = params;
  return request<API.Room>(`/api/rooms/${param0}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    params: { ...queryParams },
    data: body,
    ...(options || {}),
  });
}

/** 删除房间（管理员） 管理员删除房间 POST /api/rooms/${param0}/delete */
export async function postRoomsIdOpenApiDelete(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.postRoomsId_openAPI_deleteParams,
  options?: { [key: string]: any }
) {
  const { id: param0, ...queryParams } = params;
  return request<Record<string, any>>(`/api/rooms/${param0}/delete`, {
    method: "POST",
    params: { ...queryParams },
    ...(options || {}),
  });
}
