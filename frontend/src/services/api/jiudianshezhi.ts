// @ts-ignore
/* eslint-disable */
import { request } from "@umijs/max";

/** 获取酒店公开设置 获取指定酒店的公开配置信息 GET /api/settings/public */
export async function getSettingsPublic(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getSettingsPublicParams,
  options?: { [key: string]: any }
) {
  return request<API.HotelSettings>("/api/settings/public", {
    method: "GET",
    params: {
      ...params,
    },
    ...(options || {}),
  });
}
