import { getAdminUsers, postAdminUsersBatch } from '@/services/api/guanliyuan';
import type { ActionType, ProColumns, ProDescriptionsItemProps } from '@ant-design/pro-components';
import {
  FooterToolbar,
  PageContainer,
  ProDescriptions,
  ProTable,
} from '@ant-design/pro-components';
import { Button, Drawer, message } from 'antd';
import React, { useRef, useState } from 'react';
import CreateForm from './components/CreateForm';
import UpdateForm from './components/UpdateForm';
const TableList: React.FC = () => {
  const actionRef = useRef<ActionType | null>(null);
  const [showDetail, setShowDetail] = useState<boolean>(false);
  const [currentRow, setCurrentRow] = useState<API.User>();
  const [selectedRowsState, setSelectedRows] = useState<API.User[]>([]);

  /**
   * @en-US International configuration
   * @zh-CN 国际化配置
   * */

  const [messageApi, contextHolder] = message.useMessage();
  
  const columns: ProColumns<API.User>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
      hideInSearch: true,
    },
    {
      title: '用户名',
      dataIndex: 'username',
      render: (dom, entity) => {
        return (
          <a
            onClick={() => {
              setCurrentRow(entity);
              setShowDetail(true);
            }}
          >
            {dom}
          </a>
        );
      },
    },
    {
      title: '真实姓名',
      dataIndex: 'real_name',
    },
    {
      title: '邮箱',
      dataIndex: 'email',
    },
    {
      title: '手机号',
      dataIndex: 'phone',
    },
    {
      title: '角色',
      dataIndex: 'role',
      valueEnum: {
        user: {
          text: '用户',
          status: 'Default',
        },
        admin: {
          text: '管理员',
          status: 'Success',
        },
      },
    },
    {
      title: '状态',
      dataIndex: 'status',
      valueEnum: {
        active: {
          text: '活跃',
          status: 'Success',
        },
        blocked: {
          text: '已封禁',
          status: 'Error',
        },
      },
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      valueType: 'dateTime',
      hideInSearch: true,
      sorter: true,
    },
    {
      title: '更新时间',
      dataIndex: 'updated_at',
      valueType: 'dateTime',
      hideInSearch: true,
      sorter: true,
    },
    {
      title: '操作',
      dataIndex: 'option',
      valueType: 'option',
      render: (_, record) => [
        <UpdateForm
          trigger={<a>编辑</a>}
          key="edit"
          onOk={actionRef.current?.reload}
          values={record}
        />,
      ],
    },
  ];

  // 请求函数，将 ProTable 的分页参数转换为后端需要的格式
  const fetchUsers = async (params: any) => {
    try {
      const { current = 1, pageSize = 10, ...rest } = params;
      const response = await getAdminUsers({
        page: current,
        page_size: pageSize,
        ...rest,
      });
      
      // 如果后端返回的是数组（根据 Swagger 定义，返回的是 User[]）
      if (Array.isArray(response)) {
        return {
          data: response,
          success: true,
          total: response.length, // 注意：如果后端没有返回总数，这里使用数组长度（不准确）
        };
      }
      
      // 后端实际返回格式：{ success: true, data: User[], page: { total: number, ... } }
      const responseObj = response as any;
      if (responseObj && typeof responseObj === 'object' && 'data' in responseObj) {
        // 从 page.total 中获取总数
        const total = responseObj.page?.total || responseObj.total || 0;
        return {
          data: responseObj.data || [],
          success: true,
          total: total,
        };
      }
      
      return {
        data: [],
        success: true,
        total: 0,
      };
    } catch (error) {
      messageApi.error('获取用户列表失败');
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
      <ProTable<API.User>
        headerTitle={'用户管理'}
        actionRef={actionRef}
        rowKey="id"
        search={{
          labelWidth: 120,
        }}
        toolBarRender={() => [<CreateForm key="create" reload={actionRef.current?.reload} />]}
        request={fetchUsers}
        columns={columns}
        rowSelection={{
          onChange: (_, selectedRows) => {
            setSelectedRows(selectedRows);
          },
        }}
      />
      {selectedRowsState?.length > 0 && (
        <FooterToolbar
          extra={
            <div>
              已选择{' '}
              <a
                style={{
                  fontWeight: 600,
                }}
              >
                {selectedRowsState.length}
              </a>{' '}
              项
            </div>
          }
        >
          <Button
            danger
            onClick={async () => {
              try {
                // 获取选中用户的ID列表，并过滤掉可能的undefined值
                const userIds = selectedRowsState
                  .map(row => row.id)
                  .filter((id): id is number => id !== undefined && id !== null);
                
                // 确保有有效ID再调用API
                if (userIds.length === 0) {
                  messageApi.warning('未找到有效的用户ID');
                  return;
                }
                
                // 调用删除API
                await postAdminUsersBatch({ user_ids: userIds });
                // 显示成功消息
                messageApi.success('删除成功');
                // 刷新表格数据
                actionRef.current?.reload();
                // 清空选中状态
                setSelectedRows([]);
              } catch (error) {
                // 显示错误消息
                messageApi.error('删除失败，请重试');
                console.error('删除用户失败:', error);
              }
            }}
          >
            批量删除
          </Button>
        </FooterToolbar>
      )}

      <Drawer
        width={600}
        open={showDetail}
        onClose={() => {
          setCurrentRow(undefined);
          setShowDetail(false);
        }}
        closable={false}
      >
        {currentRow?.username && (
          <ProDescriptions<API.User>
            column={2}
            title={currentRow?.username}
            request={async () => ({
              data: currentRow || {},
            })}
            params={{
              id: currentRow?.id,
            }}
            columns={columns as ProDescriptionsItemProps<API.User>[]}
          />
        )}
      </Drawer>
    </PageContainer>
  );
};
export default TableList;
