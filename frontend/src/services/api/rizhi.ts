// @ts-ignore
/* eslint-disable */
import { request } from "@umijs/max";

/** 获取日志列表 分页获取日志列表（管理员） GET /api/admin/logs */
export async function getAdminLogs(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getAdminLogsParams,
  options?: { [key: string]: any }
) {
  return request<API.PageResponse & { data?: API.Log[] }>("/api/admin/logs", {
    method: "GET",
    params: {
      // page has a default value: 1
      page: "1",
      // page_size has a default value: 20
      page_size: "20",
      ...params,
    },
    ...(options || {}),
  });
}

/** 日志上报 客户端批量上报日志 POST /api/logs/report */
export async function postLogsReport(
  body: API.LogReportRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response>("/api/logs/report", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    data: body,
    ...(options || {}),
  });
}
