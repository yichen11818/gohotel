import { postAuthLogin } from '@/services/api/renzheng';
import { App } from 'antd';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import React from 'react';
import Login from './index';

jest.mock('@/services/api/renzheng', () => ({
  postAuthLogin: jest.fn(),
}));

jest.mock('@/services/ant-design-pro/login', () => ({
  getFakeCaptcha: jest.fn().mockResolvedValue(true),
}));

jest.mock('@umijs/max', () => ({
  Helmet: (props: any) => props.children || null,
  history: {
    push: jest.fn(),
  },
  useModel: jest.fn(),
}));

const mockHistory = jest.requireMock('@umijs/max').history as { push: jest.Mock };
const mockUseModel = jest.requireMock('@umijs/max').useModel as jest.Mock;

let mockSetInitialState: jest.Mock;

const renderLogin = () =>
  render(React.createElement(App, null, React.createElement(Login)));

describe('Login Page', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    mockSetInitialState = jest.fn();
    mockUseModel.mockReturnValue({
      initialState: undefined,
      setInitialState: mockSetInitialState,
    });
    window.history.pushState({}, '', '/user/login');
  });

  it('should show current login form', async () => {
    renderLogin();

    expect(await screen.findByText('账户密码登录')).toBeTruthy();
    expect(screen.getByPlaceholderText('用户名/手机号/邮箱')).toBeTruthy();
    expect(screen.getByPlaceholderText('请输入密码')).toBeTruthy();
    expect(screen.getByText('GoHotel 管理后台')).toBeTruthy();
    expect(screen.getByText('酒店预订与运营管理后台')).toBeTruthy();
  });

  it('should login success for admin', async () => {
    (postAuthLogin as any).mockResolvedValue({
      success: true,
      message: '登录成功！',
      data: {
        token: 'mock-token',
        user: {
          id: '1',
          username: 'admin',
          role: 'admin',
        },
      },
    });

    renderLogin();

    fireEvent.change(await screen.findByPlaceholderText('用户名/手机号/邮箱'), {
      target: { value: 'admin' },
    });
    fireEvent.change(screen.getByPlaceholderText('请输入密码'), {
      target: { value: 'Admin@123456' },
    });
    fireEvent.click(screen.getByRole('button', { name: /登\s*录/ }));

    await waitFor(() =>
      expect(postAuthLogin).toHaveBeenCalledWith({
        username: 'admin',
        password: 'Admin@123456',
      }),
    );

    expect(mockHistory.push).toHaveBeenCalledWith('/');
    expect(mockSetInitialState).toHaveBeenCalled();
  });
});
