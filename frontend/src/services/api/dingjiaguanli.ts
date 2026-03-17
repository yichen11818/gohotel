// @ts-ignore
/* eslint-disable */
import { request } from "@umijs/max";

/** 获取定价规则列表 获取所有定价规则 GET /api/admin/pricing/rules */
export async function getAdminPricingRules(options?: { [key: string]: any }) {
  return request<API.Response>("/api/admin/pricing/rules", {
    method: "GET",
    ...(options || {}),
  });
}

/** 创建定价规则 创建新的定价规则 POST /api/admin/pricing/rules */
export async function postAdminPricingRules(
  body: Record<string, any>,
  options?: { [key: string]: any }
) {
  return request<API.Response>("/api/admin/pricing/rules", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    data: body,
    ...(options || {}),
  });
}
