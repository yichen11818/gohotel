import { AvatarDropdown, AvatarName, Footer, Question } from '@/components';
import { getUsersProfile } from '@/services/api/yonghu';
import { LinkOutlined } from '@ant-design/icons';
import type { Settings as LayoutSettings } from '@ant-design/pro-components';
import { SettingDrawer } from '@ant-design/pro-components';
import '@ant-design/v5-patch-for-react-19';
import type { RequestConfig, RunTimeLayoutConfig } from '@umijs/max';
import { history, Link } from '@umijs/max';
import defaultSettings from '../config/defaultSettings';
import { errorConfig } from './requestErrorConfig';
const isDev = process.env.NODE_ENV === 'development';
const loginPath = '/user/login';

const API_BASE_URL_STORAGE_KEY = 'gohotel_admin_api_base_url';
const DEFAULT_API_BASE_URL = 'http://127.0.0.1:19999';

const getApiBaseURL = () => {
	try {
		const raw = localStorage.getItem(API_BASE_URL_STORAGE_KEY);
		if (!raw) return DEFAULT_API_BASE_URL;
		const trimmed = raw.trim().replace(/\/+$/, '');
		if (!/^https?:\/\//i.test(trimmed)) return DEFAULT_API_BASE_URL;
		return trimmed;
	} catch (_e) {
		return DEFAULT_API_BASE_URL;
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
    const token = localStorage.getItem('token');

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
      console.error('获取当前用户信息失败:', error);
      return undefined;
    }
  };
  
  // 如果不是登录页面，执行
  const { location } = history;
  if (![loginPath, '/user/register', '/user/register-result'].includes(location.pathname)) {
    const currentUser = await fetchUserInfo();

	// 仅允许管理员登录：如果本地有用户信息但角色不是管理员，直接清理并跳转登录
	if (currentUser && (currentUser as any).role !== 'admin') {
		localStorage.removeItem('token');
		localStorage.removeItem('userInfo');
		sessionStorage.removeItem('token');
		history.push(loginPath);
		return {
			fetchUserInfo,
			currentUser: undefined,
			settings: defaultSettings as Partial<LayoutSettings>,
		};
	}
    
    // 如果没有用户信息且有 token，说明 token 可能过期，跳转登录页
    if (!currentUser && localStorage.getItem('token')) {
      localStorage.removeItem('token');
      localStorage.removeItem('userInfo');
      history.push(loginPath);
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
export const layout: RunTimeLayoutConfig = ({ initialState, setInitialState }) => {
  return {
    actionsRender: () => [<Question key="doc" />],
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
      // 如果没有登录，重定向到 login
      if (!initialState?.currentUser && location.pathname !== loginPath) {
        history.push(loginPath);
      }

		// 仅允许管理员登录：如果已登录但不是管理员，同样强制回到登录页
		if (initialState?.currentUser && (initialState as any).currentUser?.role !== 'admin' && location.pathname !== loginPath) {
			localStorage.removeItem('token');
			localStorage.removeItem('userInfo');
			sessionStorage.removeItem('token');
			history.push(loginPath);
		}
    },
    bgLayoutImgList: [
      {
        src: 'https://mdn.alipayobjects.com/yuyan_qk0oxh/afts/img/D2LWSqNny4sAAAAAAAAAAAAAFl94AQBr',
        left: 85,
        bottom: 100,
        height: '303px',
      },
      {
        src: 'https://mdn.alipayobjects.com/yuyan_qk0oxh/afts/img/C2TWRpJpiC0AAAAAAAAAAAAAFl94AQBr',
        bottom: -68,
        right: -45,
        height: '303px',
      },
      {
        src: 'https://mdn.alipayobjects.com/yuyan_qk0oxh/afts/img/F6vSTbj8KpYAAAAAAAAAAAAAFl94AQBr',
        bottom: 0,
        left: 0,
        width: '331px',
      },
    ],
    links: isDev
      ? [
          <Link key="openapi" to="/umi/plugin/openapi" target="_blank">
            <LinkOutlined />
            <span>OpenAPI 文档</span>
          </Link>,
        ]
      : [],
    menuHeaderRender: undefined,
    // 自定义 403 页面
    // unAccessible: <div>unAccessible</div>,
    // 增加一个 loading 的状态
    childrenRender: (children) => {
      // if (initialState?.loading) return <PageLoading />;
      return (
        <>
          {children}
          {isDev && (
            <SettingDrawer
              disableUrlParams
              enableDarkTheme
              settings={initialState?.settings}
              onSettingChange={(settings) => {
                setInitialState((preInitialState) => ({
                  ...preInitialState,
                  settings,
                }));
              }}
            />
          )}
        </>
      );
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
