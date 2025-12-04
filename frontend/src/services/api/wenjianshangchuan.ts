// @ts-ignore
/* eslint-disable */
import { request } from "@umijs/max";

/** 通用图片上传 上传图片到腾讯云对象存储，返回临时URL POST /api/upload/image */
export async function postUploadImage(
  body: {
    /** 资源类型：banner/room/user等 */
    type: string;
  },
  image?: File,
  options?: { [key: string]: any }
) {
  const formData = new FormData();

  if (image) {
    formData.append("image", image);
  }

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

  return request<Record<string, any>>("/api/upload/image", {
    method: "POST",
    data: formData,
    requestType: "form",
    ...(options || {}),
  });
}
