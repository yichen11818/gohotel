// @ts-ignore
/* eslint-disable */
import { request } from "@umijs/max";

/** 修改密码 修改当前登录用户的密码 POST /api/users/password */
export async function postUsersPassword(
  body: Record<string, any>,
  options?: { [key: string]: any }
) {
  return request<Record<string, any>>("/api/users/password", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    data: body,
    ...(options || {}),
  });
}

/** 获取个人信息 获取当前登录用户的个人信息 GET /api/users/profile */
export async function getUsersProfile(options?: { [key: string]: any }) {
  return request<API.User>("/api/users/profile", {
    method: "GET",
    ...(options || {}),
  });
}

/** 更新个人信息 更新当前登录用户的个人信息 POST /api/users/profile */
export async function postUsersProfile(
  body: Record<string, any>,
  options?: { [key: string]: any }
) {
  return request<API.User>("/api/users/profile", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    data: body,
    ...(options || {}),
  });
}
