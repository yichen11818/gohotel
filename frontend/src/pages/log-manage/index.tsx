import { getAdminLogs } from '@/services/api/rizhi';
import { PageContainer, ProTable } from '@ant-design/pro-components';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { message, Tag } from 'antd';
import React, { useRef } from 'react';

const levelColorMap: Record<string, string> = {
  debug: 'default',
  info: 'processing',
  warn: 'warning',
  error: 'error',
};

const LogManage: React.FC = () => {
  const actionRef = useRef<ActionType | null>(null);
  const [messageApi, contextHolder] = message.useMessage();

  const columns: ProColumns<API.Log>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
      hideInSearch: true,
    },
    {
      title: '级别',
      dataIndex: 'level',
      width: 100,
      valueEnum: {
        debug: { text: 'Debug' },
        info: { text: 'Info' },
        warn: { text: 'Warn' },
        error: { text: 'Error' },
      },
      render: (_, record) => {
        const level = record.level || 'debug';
        return <Tag color={levelColorMap[level] || 'default'}>{level.toUpperCase()}</Tag>;
      },
    },
    {
      title: '日志内容',
      dataIndex: 'message',
      ellipsis: true,
      copyable: true,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      valueType: 'dateTime',
      width: 180,
      hideInSearch: true,
    },
    {
      title: '更新时间',
      dataIndex: 'updated_at',
      valueType: 'dateTime',
      width: 180,
      hideInSearch: true,
    },
  ];

  return (
    <PageContainer>
      {contextHolder}
      <ProTable<API.Log>
        headerTitle={'日志管理'}
        actionRef={actionRef}
        rowKey="id"
        request={async (params) => {
          try {
            const response = await getAdminLogs({
              page: params.current,
              page_size: params.pageSize,
            });

            return {
              data: response.data || [],
              success: response.success,
              total: response.page?.total || 0,
            };
          } catch (error) {
            messageApi.error('获取日志列表失败');
            return {
              data: [],
              success: false,
              total: 0,
            };
          }
        }}
        columns={columns}
      />
    </PageContainer>
  );
};

export default LogManage;
