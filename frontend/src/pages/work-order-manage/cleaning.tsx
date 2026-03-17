import { PageContainer, ProTable } from '@ant-design/pro-components';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { message, Button, Modal, Form, Input, Select } from 'antd';
import React, { useRef, useState } from 'react';
import { getAdminWorkOrdersCleanings as listHousekeeping, postAdminWorkOrdersCleaningIdComplete as completeCleaning, postAdminWorkOrdersCleaning as createCleaningTask, postAdminWorkOrdersCleaningIdAssign as assignStaff } from '@/services/api/gongdanguanli';

const CleaningManage: React.FC = () => {
  const actionRef = useRef<ActionType>(undefined);
  const [createModalVisible, setCreateModalVisible] = useState<boolean>(false);
  const [assignModalVisible, setAssignModalVisible] = useState<boolean>(false);
  const [currentRow, setCurrentRow] = useState<API.Housekeeping>();
  const [form] = Form.useForm();
  const [assignForm] = Form.useForm();

  const columns: ProColumns<API.Housekeeping>[] = [
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
      title: '清洁类型',
      dataIndex: 'type',
      valueEnum: {
        daily: { text: '续住清' },
        checkout: { text: '退房清' },
        deep: { text: '深清' },
      },
    },
    {
      title: '状态',
      dataIndex: 'status',
      valueEnum: {
        pending: { text: '待指派', status: 'Default' },
        in_progress: { text: '清洁中', status: 'Processing' },
        completed: { text: '已完成', status: 'Success' },
      },
    },
    {
      title: '清洁员ID',
      dataIndex: 'staff_id',
      hideInSearch: true,
    },
    {
      title: '开始时间',
      dataIndex: 'start_time',
      valueType: 'dateTime',
      hideInSearch: true,
    },
    {
      title: '完成时间',
      dataIndex: 'end_time',
      valueType: 'dateTime',
      hideInSearch: true,
    },
    {
      title: '创建时间',
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
            key="assign"
            type="link"
            onClick={() => {
              setCurrentRow(record);
              setAssignModalVisible(true);
            }}
          >
            指派人员
          </Button>
        ),
        record.status === 'in_progress' && (
          <Button
            key="complete"
            type="link"
            onClick={async () => {
              try {
                if (record.id) {
                  await completeCleaning({ id: record.id });
                  message.success('清洁任务已完成');
                  actionRef.current?.reload();
                }
              } catch (error) {
                message.error('操作失败');
              }
            }}
          >
            完成清洁
          </Button>
        ),
      ],
    },
  ];

  const handleCreate = async (values: any) => {
    try {
      const res = await createCleaningTask(values);
      if (res.success) {
        message.success('清洁任务已创建');
        setCreateModalVisible(false);
        actionRef.current?.reload();
      }
    } catch (error) {
      message.error('创建失败');
    }
  };

  const handleAssign = async (values: any) => {
    if (!currentRow?.id) return;
    try {
      const res = await assignStaff({ id: currentRow.id }, values);
      if (res.success) {
        message.success('人员已指派');
        setAssignModalVisible(false);
        actionRef.current?.reload();
      }
    } catch (error) {
      message.error('操作失败');
    }
  };

  return (
    <PageContainer>
      <ProTable<API.Housekeeping>
        headerTitle="清洁任务管理"
        actionRef={actionRef}
        rowKey="id"
        search={{ labelWidth: 120 }}
        toolBarRender={() => [
          <Button key="create" type="primary" onClick={() => setCreateModalVisible(true)}>
            手动新建任务
          </Button>,
        ]}
        request={async (params) => {
          const res = await listHousekeeping({ status: params.status });
          return {
            data: res.data || [],
            success: true,
          };
        }}
        columns={columns}
      />

      <Modal
        title="创建清洁任务"
        open={createModalVisible}
        onCancel={() => setCreateModalVisible(false)}
        onOk={() => form.submit()}
        destroyOnClose
      >
        <Form form={form} onFinish={handleCreate} layout="vertical">
          <Form.Item name="room_id" label="房间ID" rules={[{ required: true }]}>
            <Input placeholder="输入房间ID" />
          </Form.Item>
          <Form.Item name="type" label="清洁类型" rules={[{ required: true }]}>
            <Select>
              <Select.Option value="daily">续住清</Select.Option>
              <Select.Option value="checkout">退房清</Select.Option>
              <Select.Option value="deep">深清</Select.Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="指派清洁人员"
        open={assignModalVisible}
        onCancel={() => setAssignModalVisible(false)}
        onOk={() => assignForm.submit()}
        destroyOnClose
      >
        <Form form={assignForm} onFinish={handleAssign} layout="vertical">
          <Form.Item name="staff_id" label="清洁员ID" rules={[{ required: true }]}>
            <Input placeholder="输入员工ID" />
          </Form.Item>
        </Form>
      </Modal>
    </PageContainer>
  );
};

export default CleaningManage;
