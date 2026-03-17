import { PageContainer, ProTable } from '@ant-design/pro-components';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { Button, Modal, Form, Input, Select, DatePicker, InputNumber, message, Switch } from 'antd';
import React, { useRef, useState } from 'react';
import { getAdminPricingRules as listPricingRules, postAdminPricingRules as createPricingRule } from '@/services/api/dingjiaguanli';
import dayjs from 'dayjs';

const PricingManage: React.FC = () => {
  const actionRef = useRef<ActionType | null>(null);
  const [createModalVisible, setCreateModalVisible] = useState<boolean>(false);
  const [form] = Form.useForm();

  const columns: ProColumns<API.PricingRule>[] = [
    {
      title: '规则名称',
      dataIndex: 'name',
    },
    {
      title: '类型',
      dataIndex: 'type',
      valueEnum: {
        holiday: { text: '节假日', status: 'Error' },
        weekend: { text: '周末', status: 'Processing' },
        special: { text: '特殊促销', status: 'Success' },
      },
    },
    {
      title: '适用房型',
      dataIndex: 'room_type',
      render: (text) => text || '全部房型',
    },
    {
      title: '开始日期',
      dataIndex: 'start_date',
      valueType: 'date',
    },
    {
      title: '结束日期',
      dataIndex: 'end_date',
      valueType: 'date',
    },
    {
      title: '调整幅度',
      dataIndex: 'adjustment',
      render: (text, record) => {
        const prefix = (record.adjustment || 0) > 0 ? '+' : '';
        return `${prefix}${text}${record.is_percent ? '%' : '元'}`;
      },
    },
    {
      title: '优先级',
      dataIndex: 'priority',
      sorter: true,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      valueType: 'dateTime',
      hideInSearch: true,
    },
  ];

  const handleCreate = async (values: any) => {
    try {
      const { dateRange, ...rest } = values;
      const res = await createPricingRule({
        ...rest,
        start_date: dateRange[0].format('YYYY-MM-DD'),
        end_date: dateRange[1].format('YYYY-MM-DD'),
      });
      if (res.success) {
        message.success('规则创建成功');
        setCreateModalVisible(false);
        actionRef.current?.reload();
      }
    } catch (error) {
      message.error('创建失败');
    }
  };

  return (
    <PageContainer>
      <ProTable<API.PricingRule>
        headerTitle="动态定价规则"
        actionRef={actionRef}
        rowKey="id"
        search={{ labelWidth: 120 }}
        toolBarRender={() => [
          <Button key="create" type="primary" onClick={() => setCreateModalVisible(true)}>
            新建规则
          </Button>,
        ]}
        request={async () => {
          const res = await listPricingRules();
          return {
            data: res.data || [],
            success: true,
          };
        }}
        columns={columns}
      />

      <Modal
        title="新建定价规则"
        open={createModalVisible}
        onCancel={() => setCreateModalVisible(false)}
        onOk={() => form.submit()}
        destroyOnClose
      >
        <Form form={form} onFinish={handleCreate} layout="vertical">
          <Form.Item name="name" label="规则名称" rules={[{ required: true }]}>
            <Input placeholder="如：国庆涨价" />
          </Form.Item>
          <Form.Item name="type" label="类型" rules={[{ required: true }]}>
            <Select>
              <Select.Option value="holiday">节假日</Select.Option>
              <Select.Option value="weekend">周末</Select.Option>
              <Select.Option value="special">特殊促销</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="room_type" label="适用房型">
            <Select allowClear placeholder="全部房型">
              <Select.Option value="标准间">标准间</Select.Option>
              <Select.Option value="单人间">单人间</Select.Option>
              <Select.Option value="双人间">双人间</Select.Option>
              <Select.Option value="豪华套房">豪华套房</Select.Option>
              <Select.Option value="总统套房">总统套房</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="dateRange" label="有效期" rules={[{ required: true }]}>
            <DatePicker.RangePicker style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="adjustment" label="调整幅度" rules={[{ required: true }]}>
            <InputNumber style={{ width: '100%' }} placeholder="正数为涨价，负数为降价" />
          </Form.Item>
          <Form.Item name="is_percent" label="按百分比调整" valuePropName="checked">
            <Switch />
          </Form.Item>
          <Form.Item name="priority" label="优先级" initialValue={0}>
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>
        </Form>
      </Modal>
    </PageContainer>
  );
};

export default PricingManage;
