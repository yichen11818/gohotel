// @ts-ignore
/* eslint-disable */
import { request } from "@umijs/max";

/** 创建预订 创建新的房间预订，需要登录 POST /api/bookings */
export async function postBookings(
  body: API.CreateBookingRequest,
  options?: { [key: string]: any }
) {
  return request<API.Booking>("/api/bookings", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    data: body,
    ...(options || {}),
  });
}

/** 获取预订详情 根据预订ID获取预订详细信息，只能查看自己的预订 GET /api/bookings/${param0} */
export async function getBookingsId(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getBookingsIdParams,
  options?: { [key: string]: any }
) {
  const { id: param0, ...queryParams } = params;
  return request<API.Booking>(`/api/bookings/${param0}`, {
    method: "GET",
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** 取消预订 取消指定的预订，只能取消自己的预订 POST /api/bookings/${param0}/cancel */
export async function postBookingsIdCancel(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.postBookingsIdCancelParams,
  body: Record<string, any>,
  options?: { [key: string]: any }
) {
  const { id: param0, ...queryParams } = params;
  return request<Record<string, any>>(`/api/bookings/${param0}/cancel`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    params: { ...queryParams },
    data: body,
    ...(options || {}),
  });
}

/** 获取我的预订列表 获取当前登录用户的所有预订列表，支持分页 GET /api/bookings/my */
export async function getBookingsMy(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getBookingsMyParams,
  options?: { [key: string]: any }
) {
  return request<API.Booking[]>("/api/bookings/my", {
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
