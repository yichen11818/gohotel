import { PageContainer } from '@ant-design/pro-components';
import { useModel } from '@umijs/max';
import { Card, Descriptions, Tag } from 'antd';
import React from 'react';

const AdminProfilePage: React.FC = () => {
  const { initialState } = useModel('@@initialState');
  const currentUser = initialState?.currentUser as API.CurrentUser | undefined;

  return (
    <PageContainer title="个人资料">
      <Card>
        <Descriptions column={1} labelStyle={{ width: 120 }}>
          <Descriptions.Item label="姓名">
            {currentUser?.name || '-'}
          </Descriptions.Item>
          <Descriptions.Item label="用户 ID">
            {currentUser?.userid || '-'}
          </Descriptions.Item>
          <Descriptions.Item label="邮箱">
            {currentUser?.email || '-'}
          </Descriptions.Item>
          <Descriptions.Item label="手机号">
            {currentUser?.phone || '-'}
          </Descriptions.Item>
          <Descriptions.Item label="角色">
            <Tag color={(currentUser as any)?.role === 'admin' ? 'gold' : 'default'}>
              {(currentUser as any)?.role || 'unknown'}
            </Tag>
          </Descriptions.Item>
        </Descriptions>
      </Card>
    </PageContainer>
  );
};

export default AdminProfilePage;
