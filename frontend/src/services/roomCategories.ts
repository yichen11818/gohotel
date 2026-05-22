import {
  getAdminRoomCategories as getAdminRoomCategoriesGenerated,
  getAdminRoomCategoriesId as getAdminRoomCategoriesIdGenerated,
  postAdminRoomCategories,
  postAdminRoomCategoriesId,
  postAdminRoomCategoriesIdOpenApiDelete,
} from '@/services/api/guanliyuan';
import { request } from '@umijs/max';

type RoomCategoryListResponse = API.PageResponse & { data?: API.RoomCategory[] };
type RoomCategoryDetailResponse = API.Response & { data?: API.RoomCategory };

const normalizeRoomCategoryListResponse = (response: any): RoomCategoryListResponse => {
  if (Array.isArray(response)) {
    return {
      data: response,
      success: true,
      page: {
        page: 1,
        page_size: response.length,
        total: response.length,
        total_pages: 1,
      },
    };
  }

  return {
    data: response?.data || [],
    success: response?.success ?? true,
    page: response?.page,
  };
};

const normalizeRoomCategoryDetailResponse = (response: any): RoomCategoryDetailResponse => {
  if (response && typeof response === 'object' && 'data' in response) {
    return response as RoomCategoryDetailResponse;
  }

  return {
    success: true,
    data: response as API.RoomCategory,
  };
};

// Keep a stable adapter over generated services so page code does not depend on
// openapi naming quirks or inconsistent response shapes after regeneration.
export async function getAdminRoomCategories(
  params?: API.getAdminRoomCategoriesParams,
  options?: { [key: string]: any },
) {
  const response = await getAdminRoomCategoriesGenerated(
    {
      page: params?.page ?? 1,
      page_size: params?.page_size ?? 50,
    },
    options,
  );

  return normalizeRoomCategoryListResponse(response);
}

export async function getPublicRoomCategories(
  params?: API.getAdminRoomCategoriesParams,
  options?: { [key: string]: any },
) {
  if (typeof window !== 'undefined' && typeof window.fetch === 'function') {
    const url = new URL('/api/room-categories', window.location.origin);
    url.searchParams.set('page', String(params?.page ?? 1));
    url.searchParams.set('page_size', String(params?.page_size ?? 50));

    const response = await window.fetch(url.toString(), {
      method: 'GET',
      credentials: 'same-origin',
      ...(options || {}),
    });

    if (!response.ok) {
      throw new Error(`获取房型分类失败: ${response.status}`);
    }

    const data = await response.json();
    return normalizeRoomCategoryListResponse(data);
  }

  const response = await request<RoomCategoryListResponse>('/api/room-categories', {
    method: 'GET',
    params: {
      page: params?.page ?? 1,
      page_size: params?.page_size ?? 50,
    },
    ...(options || {}),
  });

  return normalizeRoomCategoryListResponse(response);
}

export async function getAdminRoomCategoryById(id: number, options?: { [key: string]: any }) {
  const response = await getAdminRoomCategoriesIdGenerated({ id }, options);
  return normalizeRoomCategoryDetailResponse(response);
}

export async function createAdminRoomCategory(
  body: API.CreateRoomCategoryRequest,
  options?: { [key: string]: any },
) {
  const response = await postAdminRoomCategories(body, options);
  return normalizeRoomCategoryDetailResponse(response);
}

export async function updateAdminRoomCategory(
  id: number,
  body: API.UpdateRoomCategoryRequest,
  options?: { [key: string]: any },
) {
  const response = await postAdminRoomCategoriesId({ id }, body, options);
  return normalizeRoomCategoryDetailResponse(response);
}

export async function deleteAdminRoomCategory(id: number, options?: { [key: string]: any }) {
  return postAdminRoomCategoriesIdOpenApiDelete({ id }, options);
}
