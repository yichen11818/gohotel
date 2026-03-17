import { PageContainer, ProTable } from '@ant-design/pro-components';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { Card, message, Tag, Button, Modal, Form, Input, Select, Popconfirm } from 'antd';
import React, { useRef, useState } from 'react';
import { getAdminWorkOrdersRepairs as listMaintenance, postAdminWorkOrdersRepairIdComplete as completeRepair, postAdminWorkOrdersRepair as createRepairRequest } from '@/services/api/gongdanguanli';

const RepairManage: React.FC = () => {
  const actionRef = useRef<ActionType>(undefined);
  const [createModalVisible, setCreateModalVisible] = useState<boolean>(false);
  const [completeModalVisible, setCompleteModalVisible] = useState<boolean>(false);
  const [currentRow, setCurrentRow] = useState<API.Maintenance>();
  const [form] = Form.useForm();
  const [completeForm] = Form.useForm();

  const columns: ProColumns<API.Maintenance>[] = [
    {
      title: 'ID',
      dataIndex: 'id',
      hideInSearch: true,
      width: 60,
    },
    {
      title: '房间号',
      dataIndex: ['room', 'room_number'],
      render: (dom, entity) => dom || entity.room_id,
    },
    {
      title: '报修类型',
      dataIndex: 'type',
      valueEnum: {
        plumbing: { text: '水路', status: 'Default' },
        electrical: { text: '电路', status: 'Default' },
        furniture: { text: '家具', status: 'Default' },
        appliance: { text: '家电', status: 'Default' },
        other: { text: '其他', status: 'Default' },
      },
    },
    {
      title: '描述',
      dataIndex: 'description',
      ellipsis: true,
    },
    {
      title: '状态',
      dataIndex: 'status',
      valueEnum: {
        pending: { text: '待处理', status: 'Warning' },
        in_progress: { text: '维修中', status: 'Processing' },
        completed: { text: '已完成', status: 'Success' },
        cancelled: { text: '已取消', status: 'Error' },
      },
    },
    {
      title: '维修人备注',
      dataIndex: 'remark',
      hideInSearch: true,
    },
    {
      title: '提交时间',
      dataIndex: 'created_at',
      valueType: 'dateTime',
      sorter: true,
    },
    {
      title: '操作',
      dataIndex: 'option',
      valueType: 'option',
      render: (_, record) => [
        record.status === 'pending' && (
          <Button
            key="complete"
            type="link"
            onClick={() => {
              setCurrentRow(record);
              setCompleteModalVisible(true);
            }}
          >
            完成维修
          </Button>
        ),
      ],
    },
  ];

  const handleCreate = async (values: any) => {
    try {
      const res = await createRepairRequest(values);
      if (res.success) {
        message.success('报修申请已提交');
        setCreateModalVisible(false);
        actionRef.current?.reload();
      }
    } catch (error) {
      message.error('提交失败');
    }
  };

  const handleComplete = async (values: any) => {
    if (!currentRow?.id) return;
    try {
      const res = await completeRepair({ id: currentRow.id }, values);
      if (res.success) {
        message.success('维修任务已完成');
        setCompleteModalVisible(false);
        actionRef.current?.reload();
      }
    } catch (error) {
      message.error('操作失败');
    }
  };

  return (
    <PageContainer>
      <ProTable<API.Maintenance>
        headerTitle="维修工单管理"
        actionRef={actionRef}
        rowKey="id"
        search={{ labelWidth: 120 }}
        toolBarRender={() => [
          <Button key="create" type="primary" onClick={() => setCreateModalVisible(true)}>
            新建报修
          </Button>,
        ]}
        request={async (params) => {
          const res = await listMaintenance({ status: params.status });
          return {
            data: res.data || [],
            success: true,
          };
        }}
        columns={columns}
      />

      <Modal
        title="新建报修"
        open={createModalVisible}
        onCancel={() => setCreateModalVisible(false)}
        onOk={() => form.submit()}
        destroyOnClose
      >
        <Form form={form} onFinish={handleCreate} layout="vertical">
          <Form.Item name="room_id" label="房间ID" rules={[{ required: true }]}>
            <Input placeholder="输入房间ID" />
          </Form.Item>
          <Form.Item name="type" label="报修类型" rules={[{ required: true }]}>
            <Select>
              <Select.Option value="plumbing">水路</Select.Option>
              <Select.Option value="electrical">电路</Select.Option>
              <Select.Option value="furniture">家具</Select.Option>
              <Select.Option value="appliance">家电</Select.Option>
              <Select.Option value="other">其他</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={4} />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="确认完成维修"
        open={completeModalVisible}
        onCancel={() => setCompleteModalVisible(false)}
        onOk={() => completeForm.submit()}
        destroyOnClose
      >
        <Form form={completeForm} onFinish={handleComplete} layout="vertical">
          <Form.Item name="remark" label="维修备注">
            <Input.TextArea rows={4} placeholder="描述维修结果..." />
          </Form.Item>
        </Form>
      </Modal>
    </PageContainer>
  );
};

export default RepairManage;
