// @ts-ignore
/* eslint-disable */
import { request } from "@umijs/max";

/** 用户登录 用户登录接口，返回用户信息和 JWT token POST /api/auth/login */
export async function postAuthLogin(
  body: API.LoginRequest,
  options?: { [key: string]: any }
) {
  return request<API.Response>("/api/auth/login", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    data: body,
    ...(options || {}),
  });
}

/** 用户注册 新用户注册接口 POST /api/auth/register */
export async function postAuthRegister(
  body: API.RegisterRequest,
  options?: { [key: string]: any }
) {
  return request<API.User>("/api/auth/register", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    data: body,
    ...(options || {}),
  });
}
