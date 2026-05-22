import { getRooms } from '@/services/api/fangjian';
import { postRoomsIdOpenApiDelete } from '@/services/api/guanliyuan';
import type { ActionType, ProColumns, ProDescriptionsItemProps } from '@ant-design/pro-components';
import {
  FooterToolbar,
  PageContainer,
  ProDescriptions,
  ProTable,
} from '@ant-design/pro-components';
import { history } from '@umijs/max';
import { Alert, Button, Drawer, Image, Tag, Typography, message, Popconfirm } from 'antd';
import React, { useRef, useState } from 'react';
import { getBackendErrorMessage } from '@/utils/backendError';
import { parseJsonArray, useRoomCategoryOptions } from '@/utils/roomCategory';
import BatchCreateForm from './components/BatchCreateForm';
import CreateForm from './components/CreateForm';
import UpdateForm from './components/UpdateForm';

const { Paragraph } = Typography;

const RoomList: React.FC = () => {
  const actionRef = useRef<ActionType | null>(null);
  const [showDetail, setShowDetail] = useState<boolean>(false);
  const [currentRow, setCurrentRow] = useState<API.Room>();
  const [selectedRowsState, setSelectedRows] = useState<API.Room[]>([]);

  const [messageApi, contextHolder] = message.useMessage();
  const { categories, error: roomCategoryError, loading: roomCategoryLoading, valueEnum } =
    useRoomCategoryOptions();
  const triggerReload: ActionType['reload'] = async (resetPageIndex?: boolean) => {
    await actionRef.current?.reload?.(resetPageIndex);
  };

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
      valueEnum,
    },
    {
      title: '预览图',
      dataIndex: 'images',
      hideInSearch: true,
      render: (_, entity) => {
        const images = parseJsonArray(entity.images);
        if (!images.length) {
          return '-';
        }

        return (
          <Image
            src={images[0]}
            width={72}
            height={48}
            style={{ objectFit: 'cover', borderRadius: 8 }}
          />
        );
      },
    },
    {
      title: '房型描述',
      dataIndex: 'description',
      hideInSearch: true,
      render: (_, entity) => (
        <Paragraph ellipsis={{ rows: 2, tooltip: entity.description }}>{entity.description || '-'}</Paragraph>
      ),
    },
    {
      title: '设施',
      dataIndex: 'facilities',
      hideInSearch: true,
      render: (_, entity) => {
        const facilities = parseJsonArray(entity.facilities);
        if (!facilities.length) {
          return '-';
        }

        return facilities.slice(0, 4).map((item) => <Tag key={item}>{item}</Tag>);
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
          text: '空闲',
          status: 'Success',
        },
        occupied: {
          text: '在住',
          status: 'Error',
        },
        maintenance: {
          text: '维修中',
          status: 'Warning',
        },
      },
    },
    {
      title: '清洁状态',
      dataIndex: 'clean_status',
      valueEnum: {
        clean: {
          text: '干净',
          status: 'Success',
        },
        dirty: {
          text: '脏房',
          status: 'Error',
        },
        inspecting: {
          text: '待查',
          status: 'Processing',
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
          onOk={triggerReload}
          roomCategories={categories}
          values={record}
        />,
        <Popconfirm
          key="delete"
          title="确定要删除这个房间吗？"
            onConfirm={async () => {
              try {
                if (record.id) {
                  await postRoomsIdOpenApiDelete({ id: Number(record.id) });
                  messageApi.success('删除成功');
                  actionRef.current?.reload();
                }
              } catch (error) {
                messageApi.error(getBackendErrorMessage(error, '删除失败'));
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
      messageApi.error(getBackendErrorMessage(error, '获取房间列表失败'));
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
      {!roomCategoryLoading && !categories.length && (
        <Alert
          type={roomCategoryError ? 'error' : 'warning'}
          showIcon
          style={{ marginBottom: 16 }}
          message={roomCategoryError ? '房型分类加载失败' : '还没有可用的房型分类'}
          description={
            roomCategoryError
              ? '请检查房型分类接口或登录状态，然后重试。'
              : '请先去“房型分类”创建房型，再新增房间。'
          }
          action={
            <Button size="small" type="link" onClick={() => history.push('/room-manage/category')}>
              去房型分类
            </Button>
          }
        />
      )}
      <ProTable<API.Room>
        headerTitle={'房间列表'}
        actionRef={actionRef}
        rowKey="id"
        search={{
          labelWidth: 120,
        }}
        toolBarRender={() => [
          <Button key="category" onClick={() => history.push('/room-manage/category')}>
            房型分类
          </Button>,
          <BatchCreateForm key="batch-create" reload={triggerReload} roomCategories={categories} />,
          <CreateForm key="create" reload={triggerReload} roomCategories={categories} />,
        ]}
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
                  roomIds.map((id) => postRoomsIdOpenApiDelete({ id: Number(id) }))
                );

                messageApi.success('删除成功');
                actionRef.current?.reload();
                setSelectedRows([]);
              } catch (error) {
                messageApi.error(getBackendErrorMessage(error, '删除失败，请重试'));
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

