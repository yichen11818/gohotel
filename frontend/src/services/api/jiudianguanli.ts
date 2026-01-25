// @ts-ignore
/* eslint-disable */
import { request } from '@umijs/max';

export async function getAdminHotels(
  params: { page?: number; page_size?: number; pageSize?: number },
  options?: { [key: string]: any },
) {
  return request<Record<string, any>>('/api/admin/hotels', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

export async function postAdminHotels(
  body: { name: string; status?: string },
  options?: { [key: string]: any },
) {
  return request<API.Hotel>('/api/admin/hotels', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

export async function getAdminHotelsId(
  params: { id: string },
  options?: { [key: string]: any },
) {
  const { id: param0 } = params;
  return request<API.Hotel>(`/api/admin/hotels/${param0}`, {
    method: 'GET',
    ...(options || {}),
  });
}

export async function postAdminHotelsId(
  params: { id: string },
  body: { name?: string; status?: string },
  options?: { [key: string]: any },
) {
  const { id: param0 } = params;
  return request<API.Hotel>(`/api/admin/hotels/${param0}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

export async function postAdminHotelsIdDelete(
  params: { id: string },
  options?: { [key: string]: any },
) {
  const { id: param0 } = params;
  return request<Record<string, any>>(`/api/admin/hotels/${param0}/delete`, {
    method: 'POST',
    ...(options || {}),
  });
}

export async function getAdminSettings(
  params: { hotel_id: number },
  options?: { [key: string]: any },
) {
  return request<API.Response & { data?: Record<string, any> }>('/api/admin/settings', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

export async function postAdminSettingsSave(
  body: { hotel_id: number; settings: Record<string, any> },
  options?: { [key: string]: any },
) {
  return request<API.Response>('/api/admin/settings/save', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

export async function getSettingsPublic(
  params: { hotel_id: number },
  options?: { [key: string]: any },
) {
  return request<API.Response & { data?: Record<string, any> }>('/api/settings/public', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}
