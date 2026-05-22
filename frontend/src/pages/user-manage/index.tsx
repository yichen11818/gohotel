import { getAdminUsers, postAdminUsersBatch } from '@/services/api/guanliyuan';
import { useModel } from '@umijs/max';
import type { ActionType, ProColumns, ProDescriptionsItemProps } from '@ant-design/pro-components';
import {
  FooterToolbar,
  PageContainer,
  ProDescriptions,
  ProTable,
} from '@ant-design/pro-components';
import { Button, Drawer, Modal, Space, Tag, message } from 'antd';
import React, { useRef, useState } from 'react';
import { getBackendErrorMessage } from '@/utils/backendError';
import CreateForm from './components/CreateForm';
import UpdateForm from './components/UpdateForm';
const TableList: React.FC = () => {
  const { initialState } = useModel('@@initialState');
  const actionRef = useRef<ActionType | null>(null);
  const [showDetail, setShowDetail] = useState<boolean>(false);
  const [currentRow, setCurrentRow] = useState<API.User>();
  const [selectedRowsState, setSelectedRows] = useState<API.User[]>([]);

  /**
   * @en-US International configuration
   * @zh-CN 国际化配置
   * */

  const [messageApi, contextHolder] = message.useMessage();
  const currentUserId = String(initialState?.currentUser?.userid || '');

  const isCurrentUser = (record?: Partial<API.User>) =>
    currentUserId !== '' && String(record?.id || '') === currentUserId;

  const handleDeleteUsers = async (rows: API.User[]) => {
    const rowsToDelete = rows.filter((row) => !isCurrentUser(row));
    const userIds = rowsToDelete
      .map((row) => row.id?.toString())
      .filter((id): id is string => Boolean(id));

    if (!userIds.length) {
      messageApi.warning('请选择可删除的用户');
      return false;
    }

    try {
      await postAdminUsersBatch({ user_ids: userIds });
      messageApi.success(rowsToDelete.length > 1 ? '批量删除成功' : '删除成功');
      setSelectedRows((prev) =>
        prev.filter((row) => !userIds.includes(String(row.id || ''))),
      );
      actionRef.current?.reload?.();
      return true;
    } catch (error) {
      messageApi.error(getBackendErrorMessage(error, '删除失败，请重试'));
      return false;
    }
  };

  const confirmDeleteUsers = (rows: API.User[]) => {
    const containsCurrentUser = rows.some((row) => isCurrentUser(row));
    const deletableRows = rows.filter((row) => !isCurrentUser(row));

    Modal.confirm({
      title: rows.length > 1 ? '确认批量删除这些用户吗？' : '确认删除该用户吗？',
      content: containsCurrentUser
        ? '已自动排除当前登录账号，其他选中用户会被删除。'
        : '删除后将无法恢复，请谨慎操作。',
      okText: '确认删除',
      cancelText: '取消',
      okButtonProps: { danger: true },
      onOk: async () => {
        if (!deletableRows.length) {
          messageApi.warning('当前登录账号不能删除');
          return;
        }
        await handleDeleteUsers(deletableRows);
      },
    });
  };
  
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
          <Space size={8}>
            <a
              onClick={() => {
                setCurrentRow(entity);
                setShowDetail(true);
              }}
            >
              {dom}
            </a>
            {isCurrentUser(entity) ? <Tag color="blue">当前账号</Tag> : null}
          </Space>
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
      title: '等级',
      dataIndex: 'level',
      valueEnum: {
        normal: { text: '普通会员', status: 'Default' },
        silver: { text: '白银会员', status: 'Processing' },
        gold: { text: '黄金会员', status: 'Warning' },
        platinum: { text: '铂金会员', status: 'Error' },
      },
    },
    {
      title: '积分',
      dataIndex: 'points',
      hideInSearch: true,
      sorter: true,
    },
    {
      title: '累计消费',
      dataIndex: 'total_spend',
      hideInSearch: true,
      sorter: true,
      render: (text) => `¥${text || 0}`,
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
          onOk={() => actionRef.current?.reload?.()}
          values={record}
        />,
        isCurrentUser(record) ? (
          <a
            key="current"
            style={{ color: '#bfbfbf', cursor: 'not-allowed' }}
            onClick={(event) => event.preventDefault()}
          >
            当前账号
          </a>
        ) : (
          <a
            key="delete"
            style={{ color: '#cf1322' }}
            onClick={() => confirmDeleteUsers([record])}
          >
            删除
          </a>
        ),
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
      messageApi.error(getBackendErrorMessage(error, '获取用户列表失败'));
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
        toolBarRender={() => [
          <CreateForm key="create" reload={() => actionRef.current?.reload?.()} />,
          <Button
            key="batchDelete"
            danger
            disabled={selectedRowsState.length === 0}
            onClick={() => confirmDeleteUsers(selectedRowsState)}
          >
            {selectedRowsState.length > 0
              ? `批量删除 (${selectedRowsState.length})`
              : '批量删除'}
          </Button>,
        ]}
        request={fetchUsers}
        columns={columns}
        rowSelection={{
          getCheckboxProps: (record) => ({
            disabled: isCurrentUser(record),
          }),
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
            onClick={() => confirmDeleteUsers(selectedRowsState)}
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
