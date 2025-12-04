// @ts-ignore
/* eslint-disable */
import { request } from "@umijs/max";

/** 获取所有活动横幅 获取所有活动横幅，支持分页 GET /api/admin/banners */
export async function getAdminBanners(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getAdminBannersParams,
  options?: { [key: string]: any }
) {
  return request<Record<string, any>>("/api/admin/banners", {
    method: "GET",
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 创建活动横幅 创建新的活动横幅，使用临时图片URL POST /api/admin/banners */
export async function postAdminBanners(
  body: {
    /** 活动横幅标题 */
    title: string;
    /** 活动横幅副标题 */
    subtitle?: string;
    /** 临时图片URL（通过通用上传接口获取） */
    temp_url: string;
    /** 点击跳转链接 */
    link_url?: string;
    /** 展示顺序 */
    sort?: number;
    /** 活动开始时间 */
    start_time?: string;
    /** 活动结束时间 */
    end_time?: string;
  },
  options?: { [key: string]: any }
) {
  const formData = new FormData();

  Object.keys(body).forEach((ele) => {
    const item = (body as any)[ele];

    if (item !== undefined && item !== null) {
      if (typeof item === "object" && !(item instanceof File)) {
        if (item instanceof Array) {
          item.forEach((f) => formData.append(ele, f || ""));
        } else {
          formData.append(
            ele,
            new Blob([JSON.stringify(item)], { type: "application/json" })
          );
        }
      } else {
        formData.append(ele, item);
      }
    }
  });

  return request<API.Banner>("/api/admin/banners", {
    method: "POST",
    data: formData,
    requestType: "form",
    ...(options || {}),
  });
}

/** 根据ID获取活动横幅 根据ID获取活动横幅详情 GET /api/admin/banners/${param0} */
export async function getAdminBannersId(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getAdminBannersIdParams,
  options?: { [key: string]: any }
) {
  const { id: param0, ...queryParams } = params;
  return request<API.Banner>(`/api/admin/banners/${param0}`, {
    method: "GET",
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** 更新活动横幅 更新活动横幅信息，使用临时图片URL POST /api/admin/banners/${param0} */
export async function postAdminBannersId(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.postAdminBannersIdParams,
  body: {
    /** 活动横幅标题 */
    title?: string;
    /** 活动横幅副标题 */
    subtitle?: string;
    /** 临时图片URL（通过通用上传接口获取） */
    temp_url: string;
    /** 点击跳转链接 */
    link_url?: string;
    /** 展示顺序 */
    sort?: number;
    /** 活动开始时间 */
    start_time?: string;
    /** 活动结束时间 */
    end_time?: string;
  },
  options?: { [key: string]: any }
) {
  const { id: param0, ...queryParams } = params;
  const formData = new FormData();

  Object.keys(body).forEach((ele) => {
    const item = (body as any)[ele];

    if (item !== undefined && item !== null) {
      if (typeof item === "object" && !(item instanceof File)) {
        if (item instanceof Array) {
          item.forEach((f) => formData.append(ele, f || ""));
        } else {
          formData.append(
            ele,
            new Blob([JSON.stringify(item)], { type: "application/json" })
          );
        }
      } else {
        formData.append(ele, item);
      }
    }
  });

  return request<API.Banner>(`/api/admin/banners/${param0}`, {
    method: "POST",
    params: { ...queryParams },
    data: formData,
    requestType: "form",
    ...(options || {}),
  });
}

/** 删除活动横幅 删除指定ID的活动横幅 POST /api/admin/banners/${param0}/delete */
export async function postAdminBannersIdOpenApiDelete(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.postAdminBannersId_openAPI_deleteParams,
  options?: { [key: string]: any }
) {
  const { id: param0, ...queryParams } = params;
  return request<Record<string, any>>(`/api/admin/banners/${param0}/delete`, {
    method: "POST",
    params: { ...queryParams },
    ...(options || {}),
  });
}

/** 获取激活的活动横幅 获取激活状态的活动横幅，用于前端展示 GET /api/banners/active */
export async function getBannersActive(options?: { [key: string]: any }) {
  return request<API.Banner[]>("/api/banners/active", {
    method: "GET",
    ...(options || {}),
  });
}
