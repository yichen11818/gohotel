import { getAdminBookings } from '@/services/api/guanliyuan';
import type { ActionType, ProColumns, ProDescriptionsItemProps } from '@ant-design/pro-components';
import {
  FooterToolbar,
  PageContainer,
  ProDescriptions,
  ProTable,
} from '@ant-design/pro-components';
import { Button, Drawer, message, Tag } from 'antd';
import React, { useRef, useState } from 'react';
import CreateForm from './components/CreateForm';
import UpdateForm from './components/UpdateForm';

// 定义订单接口类型
interface BookingType {
  id: string;
  booking_number: string;
  user_id: string;
  room_id: number;
  check_in: string;
  check_out: string;
  total_days: number;
  total_price: number;
  guest_name: string;
  guest_phone: string;
  guest_id_card?: string;
  special_request?: string;
  status: string;
  payment_status: string;
  payment_method?: string;
  cancel_reason?: string;
  created_at: string;
  updated_at: string;
  user?: {
    id: string;
    username: string;
    email?: string;
    phone?: string;
  };
  room?: {
    id: number;
    room_number: string;
    room_type?: string;
    price?: number;
  };
}

const TableList: React.FC = () => {
  const actionRef = useRef<ActionType | null>(null);
  const [showDetail, setShowDetail] = useState<boolean>(false);
  const [currentRow, setCurrentRow] = useState<BookingType>();
  const [selectedRowsState, setSelectedRows] = useState<BookingType[]>([]);

  const [messageApi, contextHolder] = message.useMessage();

  // 状态标签渲染
  const renderStatusTag = (status: string) => {
    const statusMap: Record<string, { color: string; text: string }> = {
      pending: { color: 'default', text: '待确认' },
      confirmed: { color: 'blue', text: '已确认' },
      checkin: { color: 'processing', text: '已入住' },
      checkout: { color: 'success', text: '已退房' },
      cancelled: { color: 'error', text: '已取消' },
    };
    const config = statusMap[status] || { color: 'default', text: status };
    return <Tag color={config.color}>{config.text}</Tag>;
  };

  // 支付状态标签渲染
  const renderPaymentStatusTag = (paymentStatus: string) => {
    const statusMap: Record<string, { color: string; text: string }> = {
      unpaid: { color: 'warning', text: '未支付' },
      paid: { color: 'success', text: '已支付' },
      refunded: { color: 'default', text: '已退款' },
    };
    const config = statusMap[paymentStatus] || { color: 'default', text: paymentStatus };
    return <Tag color={config.color}>{config.text}</Tag>;
  };

  // 表格列配置 
  const columns: ProColumns<BookingType>[] = [ 
    {
      title: '预订单号',
      dataIndex: 'booking_number',
      width: 180,
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
      title: '入住人姓名',
      dataIndex: 'guest_name',
      width: 120,
    },
    {
      title: '手机号',
      dataIndex: 'guest_phone',
      width: 130,
    },
    {
      title: '用户名',
      dataIndex: ['user', 'username'],
      width: 120,
      hideInSearch: true,
    },
    {
      title: '房间号',
      dataIndex: ['room', 'room_number'],
      width: 100,
    },
    {
      title: '入住日期',
      dataIndex: 'check_in',
      valueType: 'date',
      width: 120,
      sorter: true,
    },
    {
      title: '退房日期',
      dataIndex: 'check_out',
      valueType: 'date',
      width: 120,
      sorter: true,
    },
    {
      title: '总天数',
      dataIndex: 'total_days',
      width: 90,
      hideInSearch: true,
    },
    {
      title: '总价',
      dataIndex: 'total_price',
      width: 100,
      hideInSearch: true,
      render: (text) => `¥${text}`,
    },
    {
      title: '订单状态',
      dataIndex: 'status',
      width: 110,
      valueType: 'select',
      valueEnum: {
        pending: { text: '待确认', status: 'Default' },
        confirmed: { text: '已确认', status: 'Processing' },
        checkin: { text: '已入住', status: 'Processing' },
        checkout: { text: '已退房', status: 'Success' },
        cancelled: { text: '已取消', status: 'Error' },
      },
      render: (_, record) => renderStatusTag(record.status),
    },
    {
      title: '支付状态',
      dataIndex: 'payment_status',
      width: 110,
      valueType: 'select',
      valueEnum: {
        unpaid: { text: '未支付', status: 'Warning' },
        paid: { text: '已支付', status: 'Success' },
        refunded: { text: '已退款', status: 'Default' },
      },
      render: (_, record) => renderPaymentStatusTag(record.payment_status),
    },
    {
      title: '支付方式',
      dataIndex: 'payment_method',
      width: 110,
      hideInSearch: true,
      valueEnum: {
        wechat: { text: '微信支付' },
        alipay: { text: '支付宝' },
        card: { text: '银行卡' },
      },
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      valueType: 'dateTime',
      width: 160,
      hideInSearch: true,
      sorter: true,
    },
    {
      title: '操作',
      dataIndex: 'option',
      valueType: 'option',
      width: 100,
      fixed: 'right',
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
  const fetchBookings = async (params: any) => {
    try {
      const { current = 1, pageSize = 10, ...rest } = params;
      const response = await getAdminBookings({
        page: current,
        page_size: pageSize,
        ...rest,
      });
      
      // 如果后端返回的是数组（根据 Swagger 定义，返回的是 Booking[]）
      if (Array.isArray(response)) {
        return {
          data: response,
          success: true,
          total: response.length, // 注意：如果后端没有返回总数，这里使用数组长度（不准确）
        };
      }
      
      // 后端实际返回格式：{ success: true, data: Booking[], page: { total: number, ... } }
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
      messageApi.error('获取订单列表失败');
      return {
        data: [],
        success: false,
        total: 0,
      };
    }
  };
  return (
    <PageContainer>
      {/* 消息提示 */}
      {contextHolder}
      {/* 表格 */}
      <ProTable<BookingType>
        headerTitle={'订单管理'}
        actionRef={actionRef}
        rowKey="id"
        search={{
          labelWidth: 120,
        }}
        scroll={{ x: 1800 }}
        toolBarRender={() => [<CreateForm key="create" reload={actionRef.current?.reload} />]}
        request={fetchBookings}
        columns={columns}
        rowSelection={{
          onChange: (_, selectedRows) => {
            setSelectedRows(selectedRows);
          },
        }}
      />
      {/* 批量删除 */}
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
            onClick={() => {
              messageApi.warning('批量删除功能待实现');
            }}
          >
            批量删除
          </Button>
        </FooterToolbar>
      )}

      {/* 详情 */}
      <Drawer
        width={700}
        open={showDetail}
        onClose={() => {
          setCurrentRow(undefined);
          setShowDetail(false);
        }}
        closable={false}
      >
        {currentRow?.booking_number && (
          <ProDescriptions<BookingType>
            column={2}
            title={`订单详情 - ${currentRow?.booking_number}`}
            request={async () => ({
              data: currentRow || {},
            })}
            params={{
              id: currentRow?.id,
            }}
            columns={columns as ProDescriptionsItemProps<BookingType>[]}
          />
        )}
      </Drawer>
    </PageContainer>
  );
};
export default TableList;
