type BackendErrorPayload = {
  error?: {
    message?: string;
  } | string;
  message?: string;
  error_message?: string;
  errorMessage?: string;
};

const extractBackendPayloadMessage = (payload: BackendErrorPayload | string | undefined): string | undefined => {
  if (!payload) {
    return undefined;
  }

  if (typeof payload === 'string') {
    return payload;
  }

  if (payload.error) {
    if (typeof payload.error === 'string') {
      return payload.error;
    }
    if (payload.error.message) {
      return payload.error.message;
    }
  }

  return payload.message || payload.error_message || payload.errorMessage;
};

export const getBackendErrorMessage = (error: unknown, fallback?: string): string => {
  const defaultMessage = fallback ?? '操作失败，请稍后再试';
  if (!error) {
    return defaultMessage;
  }

  const axiosError = (error as any)?.response;
  const backendData = axiosError?.data;
  const backendMessage = extractBackendPayloadMessage(backendData);

  if (backendMessage) {
    return backendMessage;
  }

  if (axiosError?.statusText) {
    return axiosError.statusText;
  }

  const rawMessage = (error as any)?.message;
  if (rawMessage) {
    return rawMessage;
  }

  if (typeof axiosError?.status === 'number') {
    return `请求失败（${axiosError.status}）`;
  }

  return defaultMessage;
};
