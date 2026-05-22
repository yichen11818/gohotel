import type { RequestOptions } from '@@/plugin-request/request';
import type { RequestConfig } from '@umijs/max';
import { history } from '@umijs/max';
import { message, notification } from 'antd';

// 错误处理方案： 错误类型
enum ErrorShowType {
  SILENT = 0,
  WARN_MESSAGE = 1,
  ERROR_MESSAGE = 2,
  NOTIFICATION = 3,
  REDIRECT = 9,
}
// 与后端约定的响应数据格式
interface ResponseStructure {
  success: boolean;
  data: any;
  errorCode?: number;
  errorMessage?: string;
  showType?: ErrorShowType;
}

interface BackendErrorResponse {
	success?: boolean;
	data?: any;
	message?: string;
	error?: {
		code?: string;
		message?: string;
	};
}

const AUTH_REDIRECT_DEBOUNCE_MS = 1500;
let lastAuthRedirectAt = 0;

const clearAuthStorage = () => {
	localStorage.removeItem('token');
	localStorage.removeItem('userInfo');
	sessionStorage.removeItem('token');
};

const getCurrentPathWithQuery = () => {
	const { pathname, search, hash } = history.location;
	return `${pathname || '/'}${search || ''}${hash || ''}`;
};

const buildLoginRedirectPath = () => {
	const currentPath = getCurrentPathWithQuery();
	if (currentPath.startsWith('/user/login')) {
		return '/user/login';
	}
	return `/user/login?redirect=${encodeURIComponent(currentPath)}`;
};

const isLoginRequest = (url?: string) => {
	if (!url) return false;
	return /\/api\/auth\/login(?:\?|$)/.test(url);
};

const AUTH_FORBIDDEN_MESSAGES = new Set([
  '账号已被封禁',
  '需要管理员权限',
  '无效的登录信息',
]);

/**
 * @name 错误处理
 * pro 自带的错误处理， 可以在这里做自己的改动
 * @doc https://umijs.org/docs/max/request#配置
 */
export const errorConfig: RequestConfig = {
  // 错误处理： umi@3 的错误处理方案。
  errorConfig: {
    // 错误抛出
    errorThrower: (res) => {
      const { success, data, errorCode, errorMessage, showType } =
        res as unknown as ResponseStructure;
      if (!success) {
        const error: any = new Error(errorMessage);
        error.name = 'BizError';
        error.info = { errorCode, errorMessage, showType, data };
        throw error; // 抛出自制的错误
      }
    },
    // 错误接收及处理
    errorHandler: (error: any, opts: any) => {
      if (opts?.skipErrorHandler) throw error;
      // 我们的 errorThrower 抛出的错误。
      if (error.name === 'BizError') {
        const errorInfo: ResponseStructure | undefined = error.info;
        if (errorInfo) {
          const { errorMessage, errorCode } = errorInfo;
          switch (errorInfo.showType) {
            case ErrorShowType.SILENT:
              // do nothing
              break;
            case ErrorShowType.WARN_MESSAGE:
              message.warning(errorMessage);
              break;
            case ErrorShowType.ERROR_MESSAGE:
              message.error(errorMessage);
              break;
            case ErrorShowType.NOTIFICATION:
              notification.open({
                description: errorMessage,
                message: errorCode,
              });
              break;
            case ErrorShowType.REDIRECT:
              // TODO: redirect
              break;
            default:
              message.error(errorMessage);
          }
        }
      } else if (error.response) {
        // Axios 的错误
        // 请求成功发出且服务器也响应了状态码，但状态代码超出了 2xx 的范围
        const backendResp = error.response.data as BackendErrorResponse | undefined;
        const backendMessage = backendResp?.error?.message || backendResp?.message;
        const requestUrl = error?.response?.config?.url as string | undefined;

        if (error.response.status === 401) {
          // 登录接口 401 交给登录页自行处理，避免被误判为会话过期
          if (isLoginRequest(requestUrl)) {
            throw error;
          }

          const now = Date.now();
          const shouldNotify = now - lastAuthRedirectAt > AUTH_REDIRECT_DEBOUNCE_MS;
          const hasToken = Boolean(localStorage.getItem('token') || sessionStorage.getItem('token'));

          clearAuthStorage();

          if (shouldNotify) {
            message.warning(backendMessage || (hasToken ? '登录状态已过期，请重新登录' : '请先登录'));
            lastAuthRedirectAt = now;
          }

          if (history.location.pathname !== '/user/login') {
            history.replace(buildLoginRedirectPath());
          }
          return;
        }

		if (error.response.status === 403) {
      const shouldForceRelogin =
        Boolean(backendMessage && AUTH_FORBIDDEN_MESSAGES.has(backendMessage)) &&
        Boolean(localStorage.getItem('token') || sessionStorage.getItem('token'));

      if (shouldForceRelogin) {
        clearAuthStorage();
        message.warning(backendMessage || '当前账号已无权限继续访问，请重新登录');
        if (history.location.pathname !== '/user/login') {
          history.replace(buildLoginRedirectPath());
        }
        return;
      }

			message.error(backendMessage || '没有权限访问该资源');
			return;
		}

		if (backendMessage) {
			message.error(backendMessage);
			return;
		}

		message.error(`Response status:${error.response.status}`);
      } else if (error.request) {
        // 请求已经成功发起，但没有收到响应
        // \`error.request\` 在浏览器中是 XMLHttpRequest 的实例，
        // 而在node.js中是 http.ClientRequest 的实例
        message.error('None response! Please retry.');
      } else {
        // 发送请求时出了点问题
        message.error('Request error, please retry.');
      }
    },
  },

  // 请求拦截器
  requestInterceptors: [
    (config: RequestOptions) => {
      // 从 localStorage 获取 token
      let token = localStorage.getItem('token') || sessionStorage.getItem('token');
      
      // 如果有 token，添加到 Authorization header
      if (token) {
        config.headers = {
          ...config.headers,
          Authorization: token.startsWith('Bearer ') ? token : `Bearer ${token}`,
        };
      }
      
      return config;
    },
  ],

  // 响应拦截器
  responseInterceptors: [
    (response) => {
      return response;
    },
  ],
};
