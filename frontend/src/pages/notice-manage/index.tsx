import { getAdminNotices } from '@/services/api/gonggaoguanli';
import { postAdminNoticesIdOpenApiDelete } from '@/services/api/huodongguanli';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { PageContainer, ProTable } from '@ant-design/pro-components';
import { message, Popconfirm } from 'antd';
import React, { useRef, useState } from 'react';
import CreateNoticeForm from './components/CreateNoticeForm';
import EditNoticeForm from './components/EditNoticeForm';
import { getBackendErrorMessage } from '@/utils/backendError';

const NoticeManage: React.FC = () => {
  const actionRef = useRef<ActionType | null>(null);
  const [messageApi, contextHolder] = message.useMessage();
  
  // 编辑功能状态管理
  const [editVisible, setEditVisible] = useState<boolean>(false);
  const [currentNotice, setCurrentNotice] = useState<API.Notice | null>(null);
  
  // 打开编辑模态框
  const handleEdit = (record: API.Notice) => {
    setCurrentNotice(record);
    setEditVisible(true);
  };
  
  // 关闭编辑模态框
  const handleEditCancel = () => {
    setEditVisible(false);
    setCurrentNotice(null);
  };
  
  // 编辑成功回调
  const handleEditSuccess = () => {
    setEditVisible(false);
    setCurrentNotice(null);
    actionRef.current?.reload();
  };

  // 表格列配置
  const columns: ProColumns<API.Notice>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 60,
      hideInSearch: true,
    },
    
    {
      title: '标题',
      dataIndex: 'title',
      width: 200,
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      valueEnum: {
        active: {
          text: '激活',
          status: 'Success',
        },
        inactive: {
          text: '未激活',
          status: 'Default',
        },
      },
    },
    {
      title: '排序',
      dataIndex: 'sort',
      width: 80,
      hideInSearch: true,
    },
    {
      title: '开始时间',
      dataIndex: 'start_time',
      valueType: 'dateTime',
      width: 180,
    },
    {
      title: '结束时间',
      dataIndex: 'end_time',
      valueType: 'dateTime',
      width: 180,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      valueType: 'dateTime',
      hideInSearch: true,
      width: 180,
    },
    {
      title: '更新时间',
      dataIndex: 'updated_at',
      valueType: 'dateTime',
      hideInSearch: true,
      width: 180,
    },
    {
      title: '操作',
      dataIndex: 'option',
      valueType: 'option',
      width: 150,
      render: (_, record) => [
        <a key="edit" style={{ marginRight: 8 }} onClick={() => handleEdit(record)}>
          编辑
        </a>,
        <Popconfirm
          key="delete"
          title="确定要删除这个公告吗？"  
          onConfirm={async () => {
            try {
              if (record.id) {
                await postAdminNoticesIdOpenApiDelete({ id: String(record.id) });
                messageApi.success('删除成功');
                actionRef.current?.reload();
              }
            } catch (error) {
              const errorMessage = getBackendErrorMessage(error, '删除失败');
              messageApi.error(errorMessage);
              console.error('删除公告失败:', error);
            }
          }}
          okText="确定"
          cancelText="取消"
        >
          <a style={{ color: 'red' }}>
            删除
          </a>
        </Popconfirm>,
      ],
    },
  ];

  // 数据请求函数
  const fetchActivities = async (params: any) => {
    try {
      const { current = 1, pageSize = 10, ...rest } = params;
      const response = await getAdminNotices({
        page: current,
        pageSize,
        ...rest,
      });

      // 处理返回数据
      if (response && typeof response === 'object') {
        const responseObj = response as any;
        if (responseObj.notices && Array.isArray(responseObj.notices)) {
          return {
            data: responseObj.notices,
            success: true,
            total: responseObj.total || 0,
          };
        }
      }

      return {
        data: [],
        success: true,
        total: 0,
      };
    } catch (error) {
      const errorMessage = getBackendErrorMessage(error, '获取公告列表失败');
      messageApi.error(errorMessage);
      console.error('获取公告列表失败:', error);
      return {
        data: [],
        success: false,
        total: 0,
      };
    }
  };

  return (
    <PageContainer>
      {contextHolder}
      <ProTable<API.Notice>
        headerTitle={'公告列表'}
        actionRef={actionRef}
        rowKey="id"
        search={{
          labelWidth: 120,
        }}
        toolBarRender={() => [<CreateNoticeForm key="create" reload={() => actionRef.current?.reload?.()} />]}
        request={fetchActivities}
        columns={columns}
        pagination={{
          pageSize: 10,
        }}
      />
      {/* 编辑活动表单 */}
      {currentNotice && (
        <EditNoticeForm
          visible={editVisible}
          onCancel={handleEditCancel}
          onSuccess={handleEditSuccess}
          noticeData={currentNotice}
        />
      )}
    </PageContainer>
  );
};

export default NoticeManage;

