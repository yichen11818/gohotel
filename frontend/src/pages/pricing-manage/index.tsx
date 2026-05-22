import { PageContainer, ProTable } from '@ant-design/pro-components';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { Button, Modal, Form, Input, Select, DatePicker, InputNumber, message, Switch } from 'antd';
import React, { useRef, useState } from 'react';
import { getAdminPricingRules as listPricingRules, postAdminPricingRules as createPricingRule } from '@/services/api/dingjiaguanli';
import { type Dayjs } from 'dayjs';
import { getBackendErrorMessage } from '@/utils/backendError';
import { useRoomCategoryOptions } from '@/utils/roomCategory';

type CreatePricingRulePayload = {
  adjustment: number;
  end_date: string;
  is_percent?: boolean;
  name: string;
  priority?: number;
  room_type?: string;
  start_date: string;
  type: string;
};

type PricingFormValues = Omit<CreatePricingRulePayload, 'start_date' | 'end_date'> & {
  dateRange: [Dayjs, Dayjs];
};

const PricingManage: React.FC = () => {
  const actionRef = useRef<ActionType | null>(null);
  const [createModalVisible, setCreateModalVisible] = useState<boolean>(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [form] = Form.useForm<PricingFormValues>();
  const { options: roomCategoryOptions } = useRoomCategoryOptions();

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

  const resetForm = () => {
    form.resetFields();
  };

  const openCreateModal = () => {
    resetForm();
    setCreateModalVisible(true);
  };

  const closeCreateModal = () => {
    setCreateModalVisible(false);
    resetForm();
  };

  const handleCreate = async (values: PricingFormValues) => {
    const { dateRange, ...rest } = values;
    if (!dateRange || dateRange.length !== 2) {
      message.error('请选择有效期');
      return;
    }

    const [start, end] = dateRange;
    if (!start || !end) {
      message.error('请选择完整的有效期区间');
      return;
    }

    if (end.isBefore(start, 'day')) {
      message.error('结束日期不能早于开始日期');
      return;
    }

    setIsSubmitting(true);
    try {
      const response = await createPricingRule({
        ...rest,
        start_date: start.format('YYYY-MM-DD'),
        end_date: end.format('YYYY-MM-DD'),
        adjustment: Number(rest.adjustment),
        priority: rest.priority ?? 0,
        is_percent: Boolean(rest.is_percent),
      });
      if (response.success) {
        message.success('规则创建成功');
        closeCreateModal();
        actionRef.current?.reload();
      } else {
        message.error(response.message || '创建失败');
      }
    } catch (error) {
      message.error(getBackendErrorMessage(error, '创建失败'));
    } finally {
      setIsSubmitting(false);
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
          <Button key="create" type="primary" onClick={openCreateModal}>
            新建规则
          </Button>,
        ]}
        request={async () => {
          try {
            const res = await listPricingRules();
            return {
              data: res.data || [],
              success: res.success ?? true,
            };
          } catch (error) {
            message.error(getBackendErrorMessage(error, '获取失败'));
            return {
              data: [],
              success: false,
            };
          }
        }}
        columns={columns}
      />

      <Modal
        title="新建定价规则"
        open={createModalVisible}
        onCancel={closeCreateModal}
        onOk={() => form.submit()}
        destroyOnHidden
        okButtonProps={{ loading: isSubmitting }}
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
            <Select allowClear placeholder="全部房型" options={roomCategoryOptions} />
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
