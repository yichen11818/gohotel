// @ts-ignore
/* eslint-disable */
import { request } from "@umijs/max";

/** 创建清洁任务 创建新的清洁任务 POST /api/admin/work-orders/cleaning */
export async function postAdminWorkOrdersCleaning(
  body: Record<string, any>,
  options?: { [key: string]: any }
) {
  return request<API.Response>("/api/admin/work-orders/cleaning", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    data: body,
    ...(options || {}),
  });
}

/** 分派员工 为工单分派负责员工 POST /api/admin/work-orders/cleaning/${param0}/assign */
export async function postAdminWorkOrdersCleaningIdAssign(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.postAdminWorkOrdersCleaningIdAssignParams,
  body: Record<string, any>,
  options?: { [key: string]: any }
) {
  const { id: param0, ...queryParams } = params;
  return request<API.Response>(
    `/api/admin/work-orders/cleaning/${param0}/assign`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      params: { ...queryParams },
      data: body,
      ...(options || {}),
    }
  );
}

/** 完成清洁任务 标记清洁任务为已完成 POST /api/admin/work-orders/cleaning/${param0}/complete */
export async function postAdminWorkOrdersCleaningIdComplete(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.postAdminWorkOrdersCleaningIdCompleteParams,
  options?: { [key: string]: any }
) {
  const { id: param0, ...queryParams } = params;
  return request<API.Response>(
    `/api/admin/work-orders/cleaning/${param0}/complete`,
    {
      method: "POST",
      params: { ...queryParams },
      ...(options || {}),
    }
  );
}

/** 获取清洁任务列表 获取所有清洁任务，支持按状态筛选 GET /api/admin/work-orders/cleanings */
export async function getAdminWorkOrdersCleanings(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getAdminWorkOrdersCleaningsParams,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: API.Housekeeping[] }>(
    "/api/admin/work-orders/cleanings",
    {
      method: "GET",
      params: {
        ...params,
      },
      ...(options || {}),
    }
  );
}

/** 创建维修请求 创建新的维修请求 POST /api/admin/work-orders/repair */
export async function postAdminWorkOrdersRepair(
  body: Record<string, any>,
  options?: { [key: string]: any }
) {
  return request<API.Response>("/api/admin/work-orders/repair", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    data: body,
    ...(options || {}),
  });
}

/** 完成维修任务 标记维修任务为已完成 POST /api/admin/work-orders/repair/${param0}/complete */
export async function postAdminWorkOrdersRepairIdComplete(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.postAdminWorkOrdersRepairIdCompleteParams,
  body: Record<string, any>,
  options?: { [key: string]: any }
) {
  const { id: param0, ...queryParams } = params;
  return request<API.Response>(
    `/api/admin/work-orders/repair/${param0}/complete`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      params: { ...queryParams },
      data: body,
      ...(options || {}),
    }
  );
}

/** 获取维修任务列表 获取所有维修任务，支持按状态筛选 GET /api/admin/work-orders/repairs */
export async function getAdminWorkOrdersRepairs(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getAdminWorkOrdersRepairsParams,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: API.Maintenance[] }>(
    "/api/admin/work-orders/repairs",
    {
      method: "GET",
      params: {
        ...params,
      },
      ...(options || {}),
    }
  );
}
