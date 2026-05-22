import { PageContainer } from '@ant-design/pro-components';
import { history, useModel } from '@umijs/max';
import {
  Button,
  Card,
  Col,
  Form,
  Input,
  Modal,
  Row,
  Space,
  Tag,
  Typography,
  message,
} from 'antd';
import React, { useEffect, useState } from 'react';
import { getUsersProfile, postUsersPassword, postUsersProfile } from '@/services/api/yonghu';

type ProfileFormValues = {
  real_name?: string;
  phone?: string;
  avatar?: string;
};

type PasswordFormValues = {
  old_password: string;
  new_password: string;
  confirm_password: string;
};

const getErrorMessage = (error: any, fallback: string) => {
  const msg =
    error?.data?.error?.message ||
    error?.data?.message ||
    error?.message ||
    fallback;
  return typeof msg === 'string' ? msg : fallback;
};

const AdminProfilePage: React.FC = () => {
  const [messageApi, contextHolder] = message.useMessage();
  const { initialState, setInitialState } = useModel('@@initialState');
  const [profileForm] = Form.useForm<ProfileFormValues>();
  const [passwordForm] = Form.useForm<PasswordFormValues>();
  const [loadingProfile, setLoadingProfile] = useState(false);
  const [savingProfile, setSavingProfile] = useState(false);
  const [passwordVisible, setPasswordVisible] = useState(false);
  const [changingPassword, setChangingPassword] = useState(false);
  const [forceReloginVisible, setForceReloginVisible] = useState(false);
  const currentUser = initialState?.currentUser as API.CurrentUser | undefined;

  const syncCurrentUser = (profile?: API.User) => {
    const nextName = profile?.username || profile?.real_name || currentUser?.name || '用户';
    setInitialState((prev) => ({
      ...prev,
      currentUser: {
        ...(prev?.currentUser || {}),
        name: nextName,
        avatar: profile?.avatar || prev?.currentUser?.avatar,
        email: profile?.email || prev?.currentUser?.email,
        phone: profile?.phone || prev?.currentUser?.phone,
      },
    }));
  };

  const loadProfile = async () => {
    setLoadingProfile(true);
    try {
      const resp = await getUsersProfile();
      if (!resp?.success || !resp?.data) {
        messageApi.error(resp?.message || '获取个人资料失败');
        return;
      }

      const profile = resp.data as API.User;
      profileForm.setFieldsValue({
        real_name: profile.real_name || '',
        phone: profile.phone || '',
        avatar: profile.avatar || '',
      });
      syncCurrentUser(profile);
    } catch (error) {
      messageApi.error(getErrorMessage(error, '获取个人资料失败'));
    } finally {
      setLoadingProfile(false);
    }
  };

  useEffect(() => {
    loadProfile();
  }, []);

  const handleSaveProfile = async (values: ProfileFormValues) => {
    setSavingProfile(true);
    try {
      const resp = await postUsersProfile({
        real_name: values.real_name?.trim() || '',
        phone: values.phone?.trim() || '',
        avatar: values.avatar?.trim() || '',
      });
      if (!resp?.success) {
        messageApi.error(resp?.message || '保存个人资料失败');
        return;
      }

      const profile = resp.data as API.User | undefined;
      if (profile) {
        syncCurrentUser(profile);
        profileForm.setFieldsValue({
          real_name: profile.real_name || '',
          phone: profile.phone || '',
          avatar: profile.avatar || '',
        });
      }
      messageApi.success(resp?.message || '个人资料保存成功');
    } catch (error) {
      messageApi.error(getErrorMessage(error, '保存个人资料失败'));
    } finally {
      setSavingProfile(false);
    }
  };

  const logoutAfterPasswordChanged = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('userInfo');
    sessionStorage.removeItem('token');
    history.push('/user/login');
  };

  const handleChangePassword = async (values: PasswordFormValues) => {
    setChangingPassword(true);
    try {
      const resp = await postUsersPassword({
        old_password: values.old_password,
        new_password: values.new_password,
      });
      if (!resp?.success) {
        messageApi.error(resp?.message || '修改密码失败');
        return;
      }

      setPasswordVisible(false);
      passwordForm.resetFields();
      messageApi.success(resp?.message || '密码修改成功');
      setForceReloginVisible(true);
    } catch (error) {
      messageApi.error(getErrorMessage(error, '修改密码失败'));
    } finally {
      setChangingPassword(false);
    }
  };

  return (
    <PageContainer title="个人资料">
      {contextHolder}
      <Row gutter={[16, 16]}>
        <Col xs={24} lg={16}>
          <Card title="资料编辑" loading={loadingProfile}>
            <Form
              form={profileForm}
              layout="vertical"
              onFinish={handleSaveProfile}
              autoComplete="off"
            >
              <Form.Item label="用户名">
                <Input value={currentUser?.name || '-'} disabled />
              </Form.Item>
              <Form.Item label="用户 ID">
                <Input value={currentUser?.userid || '-'} disabled />
              </Form.Item>
              <Form.Item
                label="姓名"
                name="real_name"
                rules={[{ max: 32, message: '姓名长度不能超过 32 个字符' }]}
              >
                <Input placeholder="请输入姓名" />
              </Form.Item>
              <Form.Item
                label="手机号"
                name="phone"
                rules={[
                  { pattern: /^$|^1[3-9]\d{9}$/, message: '手机号格式不正确' },
                ]}
              >
                <Input placeholder="请输入手机号" />
              </Form.Item>
              <Form.Item
                label="头像地址"
                name="avatar"
                rules={[{ type: 'url', warningOnly: true, message: '这不是一个标准 URL' }]}
              >
                <Input placeholder="https://example.com/avatar.png" />
              </Form.Item>
              <Form.Item>
                <Space>
                  <Button type="primary" htmlType="submit" loading={savingProfile}>
                    保存资料
                  </Button>
                  <Button onClick={loadProfile} disabled={loadingProfile}>
                    重新加载
                  </Button>
                </Space>
              </Form.Item>
            </Form>
          </Card>
        </Col>

        <Col xs={24} lg={8}>
          <Card title="账号安全">
            <Space direction="vertical" size="middle" style={{ width: '100%' }}>
              <div>
                <Typography.Text type="secondary">当前登录邮箱</Typography.Text>
                <div>{currentUser?.email || '-'}</div>
              </div>
              <div>
                <Typography.Text type="secondary">当前手机号</Typography.Text>
                <div>{currentUser?.phone || '-'}</div>
              </div>
              <div>
                <Typography.Text type="secondary">角色</Typography.Text>
                <div>
                  <Tag color={(currentUser as any)?.role === 'admin' ? 'gold' : 'default'}>
                    {(currentUser as any)?.role || 'unknown'}
                  </Tag>
                </div>
              </div>
              <Button block onClick={() => setPasswordVisible(true)}>
                修改密码
              </Button>
            </Space>
          </Card>
        </Col>
      </Row>

      <Modal
        title="修改密码"
        open={passwordVisible}
        onCancel={() => {
          setPasswordVisible(false);
          passwordForm.resetFields();
        }}
        footer={null}
        destroyOnHidden
      >
        <Form
          form={passwordForm}
          layout="vertical"
          autoComplete="off"
          onFinish={handleChangePassword}
        >
          <Form.Item
            name="old_password"
            label="当前密码"
            rules={[{ required: true, message: '请输入当前密码' }]}
          >
            <Input.Password placeholder="请输入当前密码" />
          </Form.Item>
          <Form.Item
            name="new_password"
            label="新密码"
            rules={[
              { required: true, message: '请输入新密码' },
              { min: 6, message: '新密码至少 6 位' },
            ]}
          >
            <Input.Password placeholder="请输入新密码" />
          </Form.Item>
          <Form.Item
            name="confirm_password"
            label="确认新密码"
            dependencies={['new_password']}
            rules={[
              { required: true, message: '请再次输入新密码' },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue('new_password') === value) {
                    return Promise.resolve();
                  }
                  return Promise.reject(new Error('两次输入的新密码不一致'));
                },
              }),
            ]}
          >
            <Input.Password placeholder="请再次输入新密码" />
          </Form.Item>
          <Form.Item style={{ marginBottom: 0 }}>
            <Space>
              <Button
                type="primary"
                htmlType="submit"
                loading={changingPassword}
              >
                确认修改
              </Button>
              <Button
                onClick={() => {
                  setPasswordVisible(false);
                  passwordForm.resetFields();
                }}
              >
                取消
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="密码已更新"
        open={forceReloginVisible}
        closable={false}
        maskClosable={false}
        okText="重新登录"
        cancelText="稍后手动登录"
        onOk={logoutAfterPasswordChanged}
        onCancel={() => {
          setForceReloginVisible(false);
          logoutAfterPasswordChanged();
        }}
      >
        <Typography.Text>
          为了账号安全，请使用新密码重新登录。
        </Typography.Text>
      </Modal>
    </PageContainer>
  );
};

export default AdminProfilePage;
