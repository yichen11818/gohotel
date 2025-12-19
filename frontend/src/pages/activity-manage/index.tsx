import { getAdminBanners, postAdminBannersIdOpenApiDelete } from '@/services/api/huodongguanli';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { PageContainer, ProTable } from '@ant-design/pro-components';
import { message, Popconfirm } from 'antd';
import React, { useRef, useState } from 'react';
import CreateActivityForm from './components/CreateActivityForm';
import EditActivityForm from './components/EditActivityForm';

const ActivityManage: React.FC = () => {
  const actionRef = useRef<ActionType | null>(null);
  const [messageApi, contextHolder] = message.useMessage();
  
  // 编辑功能状态管理
  const [editVisible, setEditVisible] = useState<boolean>(false);
  const [currentActivity, setCurrentActivity] = useState<API.Banner | null>(null);
  
  // 打开编辑模态框
  const handleEdit = (record: API.Banner) => {
    setCurrentActivity(record);
    setEditVisible(true);
  };
  
  // 关闭编辑模态框
  const handleEditCancel = () => {
    setEditVisible(false);
    setCurrentActivity(null);
  };
  
  // 编辑成功回调
  const handleEditSuccess = () => {
    setEditVisible(false);
    setCurrentActivity(null);
    actionRef.current?.reload();
  };

  // 表格列配置
  const columns: ProColumns<API.Banner>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 60,
      hideInSearch: true,
    },
    {
      title: '图片',
      dataIndex: 'image_url',
      width: 100,
      hideInSearch: true,
      render: (imageUrl, record) => {
        if (!imageUrl) return null;
        // 确保imageUrl是字符串类型
        const src = typeof imageUrl === 'string' ? imageUrl : String(imageUrl);
        return (
          <div style={{ width: 80, height: 60, overflow: 'hidden', borderRadius: 4 }}>
            <img 
              src={src} 
              alt={record.title || ''} 
              style={{ width: '100%', height: '100%', objectFit: 'cover' }}
            />
          </div>
        );
      },
    },

    {
      title: '标题',
      dataIndex: 'title',
      width: 200,
    },
    {
      title: '副标题',
      dataIndex: 'subtitle',
      width: 300,
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
          title="确定要删除这个活动吗？"
          onConfirm={async () => {
            try {
              if (record.id) {
                await postAdminBannersIdOpenApiDelete({ id: String(record.id) });
                messageApi.success('删除成功');
                actionRef.current?.reload();
              }
            } catch (error) {
              messageApi.error('删除失败');
              console.error('删除活动失败:', error);
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
      const response = await getAdminBanners({
        page: current,
        page_size: pageSize,
        ...rest,
      });

      // 处理返回数据
      if (response && typeof response === 'object') {
        const responseObj = response as any;
        if (responseObj.banners && Array.isArray(responseObj.banners)) {
          return {
            data: responseObj.banners,
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
      messageApi.error('获取活动列表失败');
      console.error('获取活动列表失败:', error);
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
      <ProTable<API.Banner>
        headerTitle={'活动列表'}
        actionRef={actionRef}
        rowKey="id"
        search={{
          labelWidth: 120,
        }}
        toolBarRender={() => [<CreateActivityForm key="create" reload={() => actionRef.current?.reload?.()} />]}
        request={fetchActivities}
        columns={columns}
        pagination={{
          pageSize: 10,
        }}
      />
      {/* 编辑活动表单 */}
      {currentActivity && (
        <EditActivityForm
          visible={editVisible}
          onCancel={handleEditCancel}
          onSuccess={handleEditSuccess}
          activityData={currentActivity}
        />
      )}
    </PageContainer>
  );
};

export default ActivityManage;

