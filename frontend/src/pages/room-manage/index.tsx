import { getRooms } from '@/services/api/fangjian';
import { postRoomsIdOpenApiDelete } from '@/services/api/guanliyuan';
import type { ActionType, ProColumns, ProDescriptionsItemProps } from '@ant-design/pro-components';
import {
  FooterToolbar,
  PageContainer,
  ProDescriptions,
  ProTable,
} from '@ant-design/pro-components';
import { Button, Drawer, message, Popconfirm } from 'antd';
import React, { useRef, useState } from 'react';
import CreateForm from './components/CreateForm';
import UpdateForm from './components/UpdateForm';

const RoomList: React.FC = () => {
  const actionRef = useRef<ActionType | null>(null);
  const [showDetail, setShowDetail] = useState<boolean>(false);
  const [currentRow, setCurrentRow] = useState<API.Room>();
  const [selectedRowsState, setSelectedRows] = useState<API.Room[]>([]);

  const [messageApi, contextHolder] = message.useMessage();

  const columns: ProColumns<API.Room>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 60,
      hideInSearch: true,
    },
    {
      title: '房间号',
      dataIndex: 'room_number',
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
      title: '房型',
      dataIndex: 'room_type',
      valueEnum: {
        '单人间': {
          text: '单人间',
        },
        '双人间': {
          text: '双人间',
        },
        '豪华套房': {
          text: '豪华套房',
        },
        '总统套房': {
          text: '总统套房',
        },
        '商务套房': {
          text: '商务套房',
        },
      },
    },
    {
      title: '楼层',
      dataIndex: 'floor',
      hideInSearch: true,
      sorter: true,
    },
    {
      title: '价格',
      dataIndex: 'price',
      hideInSearch: true,
      render: (dom, entity) => `¥${entity.price}`,
    },
    {
      title: '原价',
      dataIndex: 'original_price',
      hideInSearch: true,
      render: (dom, entity) => entity.original_price ? `¥${entity.original_price}` : '-',
    },
    {
      title: '可住人数',
      dataIndex: 'capacity',
      hideInSearch: true,
    },
    {
      title: '面积(m²)',
      dataIndex: 'area',
      hideInSearch: true,
    },
    {
      title: '床型',
      dataIndex: 'bed_type',
      hideInSearch: true,
    },
    {
      title: '状态',
      dataIndex: 'status',
      valueEnum: {
        available: {
          text: '可用',
          status: 'Success',
        },
        occupied: {
          text: '占用',
          status: 'Error',
        },
        maintenance: {
          text: '维护中',
          status: 'Warning',
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
        <Popconfirm
          key="delete"
          title="确定要删除这个房间吗？"
          onConfirm={async () => {
            try {
              if (record.id) {
                await postRoomsIdOpenApiDelete({ id: record.id });
                messageApi.success('删除成功');
                actionRef.current?.reload();
              }
            } catch (error) {
              messageApi.error('删除失败');
            }
          }}
          okText="确定"
          cancelText="取消"
        >
          <a style={{ color: 'red' }}>删除</a>
        </Popconfirm>,
      ],
    },
  ];

  // 请求函数，将 ProTable 的分页参数转换为后端需要的格式
  const fetchRooms = async (params: any) => {
    try {
      const { current = 1, pageSize = 10, ...rest } = params;
      const response = await getRooms({
        page: current,
        page_size: pageSize,
        ...rest,
      });

      // 如果后端返回的是数组
      if (Array.isArray(response)) {
        return {
          data: response,
          success: true,
          total: response.length,
        };
      }

      // 后端实际返回格式：{ success: true, data: Room[], page: { total: number, ... } }
      const responseObj = response as any;
      if (responseObj && typeof responseObj === 'object' && 'data' in responseObj) {
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
      messageApi.error('获取房间列表失败');
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
      <ProTable<API.Room>
        headerTitle={'房间列表'}
        actionRef={actionRef}
        rowKey="id"
        search={{
          labelWidth: 120,
        }}
        toolBarRender={() => [<CreateForm key="create" reload={actionRef.current?.reload} />]}
        request={fetchRooms}
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
                // 获取选中房间的ID列表
                const roomIds = selectedRowsState
                  .map((row) => row.id)
                  .filter((id): id is number => id !== undefined && id !== null);

                if (roomIds.length === 0) {
                  messageApi.warning('未找到有效的房间ID');
                  return;
                }

                // 批量删除
                await Promise.all(
                  roomIds.map((id) => postRoomsIdOpenApiDelete({ id }))
                );

                messageApi.success('删除成功');
                actionRef.current?.reload();
                setSelectedRows([]);
              } catch (error) {
                messageApi.error('删除失败，请重试');
                console.error('删除房间失败:', error);
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
        {currentRow?.room_number && (
          <ProDescriptions<API.Room>
            column={2}
            title={`房间 ${currentRow?.room_number}`}
            request={async () => ({
              data: currentRow || {},
            })}
            params={{
              id: currentRow?.id,
            }}
            columns={columns as ProDescriptionsItemProps<API.Room>[]}
          />
        )}
      </Drawer>
    </PageContainer>
  );
};

export default RoomList;

