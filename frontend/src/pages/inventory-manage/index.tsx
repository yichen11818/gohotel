import { PageContainer } from '@ant-design/pro-components';
import { Card, DatePicker, Table, Tag, message, Space, Button, Modal, Form, InputNumber, Select } from 'antd';
import React, { useEffect, useState } from 'react';
import dayjs from 'dayjs';
import { getAdminInventoryGrid as getInventoryGrid, postAdminInventoryInit as initInventory, postAdminInventoryUpdate as updateInventory } from '@/services/api/kucunguanli';
import { getBackendErrorMessage } from '@/utils/backendError';
import { useRoomCategoryOptions } from '@/utils/roomCategory';

const { RangePicker } = DatePicker;

const InventoryManage: React.FC = () => {
  const [loading, setLoading] = useState<boolean>(false);
  const [dateRange, setDateRange] = useState<[dayjs.Dayjs, dayjs.Dayjs]>([
    dayjs(),
    dayjs().add(14, 'day'),
  ]);
  const [data, setData] = useState<Record<string, API.RoomInventory[]>>({});
  const [initModalVisible, setInitModalVisible] = useState<boolean>(false);
  const [form] = Form.useForm();
  const { options: roomCategoryOptions } = useRoomCategoryOptions();

  const fetchData = async () => {
    setLoading(true);
    try {
      const res = await getInventoryGrid({
        start_date: dateRange[0].format('YYYY-MM-DD'),
        end_date: dateRange[1].format('YYYY-MM-DD'),
      });
      if (res.success) {
        setData(res.data || {});
      }
    } catch (error) {
      message.error(getBackendErrorMessage(error, '获取房态库存失败'));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, [dateRange]);

  const columns = [
    {
      title: '房型',
      dataIndex: 'roomType',
      key: 'roomType',
      fixed: 'left' as const,
      width: 120,
    },
    ...Array.from({ length: dateRange[1].diff(dateRange[0], 'day') + 1 }).map((_, index) => {
      const date = dateRange[0].add(index, 'day');
      const dateStr = date.format('MM-DD');
      const isWeekend = date.day() === 0 || date.day() === 6;
      return {
        title: (
          <div style={{ color: isWeekend ? '#ff4d4f' : 'inherit', textAlign: 'center' }}>
            {dateStr}
            <br />
            {['日', '一', '二', '三', '四', '五', '六'][date.day()]}
          </div>
        ),
        dataIndex: date.format('YYYY-MM-DD'),
        key: date.format('YYYY-MM-DD'),
        width: 100,
        render: (_: any, record: any) => {
          const inv = record.dates[date.format('YYYY-MM-DD')];
          if (!inv) return '-';
          const remaining = inv.total_count - inv.booked_count;
          const statusColor = remaining <= 0 ? 'error' : remaining < 3 ? 'warning' : 'success';
          return (
            <div style={{ textAlign: 'center' }}>
              <Tag color={statusColor} style={{ margin: 0, width: '100%' }}>
                余 {remaining}
              </Tag>
              <div style={{ fontSize: '12px', color: '#999', marginTop: '4px' }}>
                ¥{inv.price}
              </div>
            </div>
          );
        },
      };
    }),
  ];

  const tableData = Object.keys(data).map((roomType) => {
    const dates: Record<string, API.RoomInventory> = {};
    data[roomType].forEach((inv) => {
      if (inv.date) {
        dates[dayjs(inv.date).format('YYYY-MM-DD')] = inv;
      }
    });
    return {
      key: roomType,
      roomType,
      dates,
    };
  });

  const handleInit = async (values: any) => {
    try {
      const res = await initInventory(values);
      if (res.success) {
        message.success('初始化成功');
        setInitModalVisible(false);
        fetchData();
      }
    } catch (error) {
      message.error(getBackendErrorMessage(error, '初始化失败'));
    }
  };

  return (
    <PageContainer>
      <Card
        title="房态库存矩阵"
        extra={
          <Space>
            <RangePicker
              value={dateRange}
              onChange={(dates) => dates && setDateRange([dates[0]!, dates[1]!])}
              allowClear={false}
            />
            <Button type="primary" onClick={() => setInitModalVisible(true)}>
              初始化库存
            </Button>
            <Button onClick={fetchData}>刷新</Button>
          </Space>
        }
      >
        <Table
          columns={columns}
          dataSource={tableData}
          loading={loading}
          scroll={{ x: 'max-content' }}
          pagination={false}
          bordered
        />
      </Card>

      <Modal
        title="初始化库存"
        open={initModalVisible}
        onCancel={() => setInitModalVisible(false)}
        onOk={() => form.submit()}
        destroyOnHidden
      >
        <Form form={form} onFinish={handleInit} layout="vertical">
          <Form.Item name="room_type" label="房型" rules={[{ required: true }]}>
            <Select options={roomCategoryOptions} />
          </Form.Item>
          <Form.Item name="total_count" label="总房间数" rules={[{ required: true }]}>
            <InputNumber min={1} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="price" label="基础价格" rules={[{ required: true }]}>
            <InputNumber min={0} precision={2} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="days" label="初始化天数" initialValue={30} rules={[{ required: true }]}>
            <InputNumber min={1} max={365} style={{ width: '100%' }} />
          </Form.Item>
        </Form>
      </Modal>
    </PageContainer>
  );
};

export default InventoryManage;
