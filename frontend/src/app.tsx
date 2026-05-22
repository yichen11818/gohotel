import { AvatarDropdown, AvatarName, Footer } from '@/components';
import { getUsersProfile } from '@/services/api/yonghu';
import type { Settings as LayoutSettings } from '@ant-design/pro-components';
import '@ant-design/v5-patch-for-react-19';
import type { RequestConfig, RunTimeLayoutConfig } from '@umijs/max';
import { history } from '@umijs/max';
import { message } from 'antd';
import defaultSettings from '../config/defaultSettings';
import { errorConfig } from './requestErrorConfig';
const loginPath = '/user/login';
const SESSION_EXPIRED_PARAM = 'sessionExpired';
let sessionExpiredNotified = false;
const getCurrentPathWithQuery = () => {
  const { pathname, search, hash } = history.location;
  return `${pathname || '/'}${search || ''}${hash || ''}`;
};

const buildLoginPathWithRedirect = (sessionExpired = false) => {
  const currentPath = getCurrentPathWithQuery();
  const params = new URLSearchParams();
  if (currentPath.startsWith(loginPath)) {
    if (sessionExpired) {
      params.set(SESSION_EXPIRED_PARAM, '1');
      return `${loginPath}?${params.toString()}`;
    }
    return loginPath;
  }
  params.set('redirect', currentPath);
  if (sessionExpired) {
    params.set(SESSION_EXPIRED_PARAM, '1');
  }
  return `${loginPath}?${params.toString()}`;
};

const clearAuthStorage = () => {
  localStorage.removeItem('token');
  localStorage.removeItem('userInfo');
  sessionStorage.removeItem('token');
};

const normalizeApiBaseURL = (value?: string) => {
  if (!value) return '';
  return value.trim().replace(/\/+$/, '');
};

const LEGACY_API_BASE_URL_STORAGE_KEY = 'gohotel_admin_api_base_url';

const getBrowserOriginBaseURL = () => {
  if (typeof window === 'undefined') {
    return '';
  }

  const origin = normalizeApiBaseURL(window.location.origin);
  if (!/^https?:\/\//i.test(origin)) {
    return '';
  }

  return origin;
};

const DEFAULT_API_BASE_URL = normalizeApiBaseURL(process.env.UMI_APP_API_BASE_URL);

const getApiBaseURL = () => {
  const browserOrigin = getBrowserOriginBaseURL();

  try {
    // 历史版本允许通过 localStorage 覆盖 API 地址，容易把页面请求指到旧服务。
    // 现在统一优先走环境变量或当前站点同源 API，并清理遗留值，避免页面与接口跨服务错位。
    if (typeof window !== 'undefined') {
      localStorage.removeItem(LEGACY_API_BASE_URL_STORAGE_KEY);
    }
  } catch (_e) {
    // Ignore storage cleanup failure and fall back to runtime defaults.
  }

  return DEFAULT_API_BASE_URL || browserOrigin || '';
};

const maybeShowSessionExpiredNotice = () => {
  const params = new URLSearchParams(history.location.search);
  if (params.get(SESSION_EXPIRED_PARAM) === '1' && !sessionExpiredNotified) {
    message.info('会话已过期，请重新登录以继续操作');
    sessionExpiredNotified = true;
  }
};

/**
 * @see https://umijs.org/docs/api/runtime-config#getinitialstate
 * */
export async function getInitialState(): Promise<{
  settings?: Partial<LayoutSettings>;
  currentUser?: API.CurrentUser;
  loading?: boolean;
  fetchUserInfo?: () => Promise<API.CurrentUser | undefined>;
}> {
  const fetchUserInfo = async () => {
    const token = localStorage.getItem('token') || sessionStorage.getItem('token');

    if (!token) {
      return undefined;
    }

    try {
      const response = await getUsersProfile();
      const userInfo = response?.data;

      if (!response?.success || !userInfo) {
        return undefined;
      }

      localStorage.setItem('userInfo', JSON.stringify(userInfo));

      return {
        name: userInfo.username || userInfo.real_name || '用户',
        avatar:
          userInfo.avatar ||
          'https://gw.alipayobjects.com/zos/antfincdn/XAosXuNZyF/BiazfanxmamNRoxxVxka.png',
        userid: userInfo.id?.toString() || '',
        email: userInfo.email,
        phone: userInfo.phone,
        role: userInfo.role,
      } as API.CurrentUser;
    } catch (error) {
      const statusCode = (error as any)?.response?.status;
      if (statusCode !== 401) {
        console.error('获取当前用户信息失败:', error);
      }
      return undefined;
    }
  };
  
  // 如果不是登录页面，执行
  const { location } = history;
  if (![loginPath, '/user/register', '/user/register-result'].includes(location.pathname)) {
    const currentUser = await fetchUserInfo();

	// 仅允许管理员登录：如果本地有用户信息但角色不是管理员，直接清理并跳转登录
		if (currentUser && (currentUser as any).role !== 'admin') {
		clearAuthStorage();
		history.replace(buildLoginPathWithRedirect(true));
		return {
			fetchUserInfo,
			currentUser: undefined,
			settings: defaultSettings as Partial<LayoutSettings>,
		};
	}
    
    // 如果没有用户信息且有 token，说明 token 可能过期，跳转登录页
    if (!currentUser && (localStorage.getItem('token') || sessionStorage.getItem('token'))) {
      clearAuthStorage();
      history.replace(buildLoginPathWithRedirect(true));
    }
    
    return {
      fetchUserInfo,
      currentUser,
      settings: defaultSettings as Partial<LayoutSettings>,
    };
  }
  return {
    fetchUserInfo,
    settings: defaultSettings as Partial<LayoutSettings>,
  };
}

// ProLayout 支持的api https://procomponents.ant.design/components/layout
export const layout: RunTimeLayoutConfig = ({ initialState }) => {
  maybeShowSessionExpiredNotice();
  return {
    avatarProps: {
      src: initialState?.currentUser?.avatar,
      title: <AvatarName />,
      render: (_, avatarChildren) => {
        return <AvatarDropdown>{avatarChildren}</AvatarDropdown>;
      },
    },
    waterMarkProps: {
      content: initialState?.currentUser?.name,
    },
    footerRender: () => <Footer />,
    onPageChange: () => {
      const { location } = history;
      const hasToken = Boolean(localStorage.getItem('token') || sessionStorage.getItem('token'));
      // 如果没有登录，重定向到 login
      if (!initialState?.currentUser && location.pathname !== loginPath) {
        history.replace(buildLoginPathWithRedirect(hasToken));
      }

		// 仅允许管理员登录：如果已登录但不是管理员，同样强制回到登录页
		if (initialState?.currentUser && (initialState as any).currentUser?.role !== 'admin' && location.pathname !== loginPath) {
			clearAuthStorage();
			history.replace(buildLoginPathWithRedirect(true));
		}

      maybeShowSessionExpiredNotice();
    },
    layoutBgStyle: {
      background:
        'linear-gradient(180deg, #f8fbff 0%, #f1f6fb 42%, #eef3f8 100%)',
    },
    bgLayoutImgList: [],
    menuHeaderRender: (logoDom, titleDom) => (
      <div
        style={{
          display: 'flex',
          flexDirection: 'column',
          padding: '12px 0 10px',
          fontWeight: 600,
          letterSpacing: '0.04em',
          color: '#16324a',
        }}
      >
        <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
          {logoDom}
          <span style={{ fontSize: 18 }}>GoHotel</span>
        </div>
        <span style={{ fontSize: 12, color: '#6d8093' }}>酒店运营中台</span>
      </div>
    ),
    headerTitleRender: () => (
      <span style={{ fontWeight: 600, color: '#16324a' }}>GoHotel 管理控制台</span>
    ),
    // 自定义 403 页面
    // unAccessible: <div>unAccessible</div>,
    // 增加一个 loading 的状态
    childrenRender: (children) => {
      return <>{children}</>;
    },
    ...initialState?.settings,
  };
};

/**
 * @name request 配置，可以配置错误处理
 * 它基于 axios 和 ahooks 的 useRequest 提供了一套统一的网络请求和错误处理方案。
 * @doc https://umijs.org/docs/max/request#配置
 */
export const request: RequestConfig = {
  baseURL: getApiBaseURL(),
  ...errorConfig,
};
