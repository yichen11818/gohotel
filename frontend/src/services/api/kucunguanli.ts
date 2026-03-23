// @ts-ignore
/* eslint-disable */
import { request } from "@umijs/max";

/** 获取库存网格 获取指定日期范围的库存网格数据 GET /api/admin/inventory/grid */
export async function getAdminInventoryGrid(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getAdminInventoryGridParams,
  options?: { [key: string]: any }
) {
  return request<API.Response & { data?: API.InventoryGridResponse }>(
    "/api/admin/inventory/grid",
    {
      method: "GET",
      params: {
        ...params,
      },
      ...(options || {}),
    }
  );
}

/** 初始化库存 初始化房型库存 POST /api/admin/inventory/init */
export async function postAdminInventoryInit(
  body: Record<string, any>,
  options?: { [key: string]: any }
) {
  return request<API.Response>("/api/admin/inventory/init", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    data: body,
    ...(options || {}),
  });
}

/** 更新库存 调整指定日期范围的库存 POST /api/admin/inventory/update */
export async function postAdminInventoryUpdate(
  body: Record<string, any>,
  options?: { [key: string]: any }
) {
  return request<API.Response>("/api/admin/inventory/update", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    data: body,
    ...(options || {}),
  });
}
