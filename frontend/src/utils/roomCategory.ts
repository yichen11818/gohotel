import { useCallback, useEffect, useState } from 'react';

import { getPublicRoomCategories } from '@/services/roomCategories';

export const parseJsonArray = (raw?: string): string[] => {
  if (!raw) {
    return [];
  }

  try {
    const parsed = JSON.parse(raw);
    return Array.isArray(parsed)
      ? parsed.map((item) => String(item).trim()).filter(Boolean)
      : [];
  } catch (error) {
    return [];
  }
};

export const stringifyJsonArray = (values?: string[]): string =>
  JSON.stringify((values || []).map((item) => item.trim()).filter(Boolean));

export const getRoomCategoryOptions = (categories: API.RoomCategory[]) =>
  categories
    .map((category) => {
      const value = category.name?.trim();
      if (!value) {
        return null;
      }

      return {
        label: value,
        value,
      };
    })
    .filter(Boolean) as { label: string; value: string }[];

export const getRoomCategoryValueEnum = (categories: API.RoomCategory[]) =>
  Object.fromEntries(
    getRoomCategoryOptions(categories).map((option) => [
      option.value,
      {
        text: option.label,
      },
    ]),
  );

export const useRoomCategoryOptions = () => {
  const [categories, setCategories] = useState<API.RoomCategory[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<unknown>();

  const refresh = useCallback(async () => {
    setLoading(true);
    setError(undefined);

    try {
      const response = await getPublicRoomCategories({ page: 1, page_size: 200 });
      setCategories(response.data || []);
    } catch (err) {
      setCategories([]);
      setError(err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const options = getRoomCategoryOptions(categories);
  const valueEnum = getRoomCategoryValueEnum(categories);

  return {
    categories,
    error,
    loading,
    options,
    refresh,
    valueEnum,
  };
};
