import { PageContainer, ProTable } from '@ant-design/pro-components';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { message, Button, Modal, Form, Select } from 'antd';
import React, { useEffect, useRef, useState } from 'react';
import { getAdminWorkOrdersCleanings as listHousekeeping, postAdminWorkOrdersCleaningIdComplete as completeCleaning, postAdminWorkOrdersCleaning as createCleaningTask, postAdminWorkOrdersCleaningIdAssign as assignStaff } from '@/services/api/gongdanguanli';
import { getAdminUsers } from '@/services/api/guanliyuan';
import { getRooms } from '@/services/api/fangjian';
import { getBackendErrorMessage } from '@/utils/backendError';

const CleaningManage: React.FC = () => {
  const actionRef = useRef<ActionType | null>(null);
  const [createModalVisible, setCreateModalVisible] = useState<boolean>(false);
  const [assignModalVisible, setAssignModalVisible] = useState<boolean>(false);
  const [currentRow, setCurrentRow] = useState<API.Housekeeping>();
  const [roomOptions, setRoomOptions] = useState<{ label: string; value: number }[]>([]);
  const [staffOptions, setStaffOptions] = useState<{ label: string; value: string }[]>([]);
  const [optionsLoading, setOptionsLoading] = useState(false);
  const [form] = Form.useForm();
  const [assignForm] = Form.useForm();

  useEffect(() => {
    const loadOptions = async () => {
      setOptionsLoading(true);
      try {
        const [roomResponse, userResponse] = await Promise.all([
          getRooms({ page_size: 200 } as API.getRoomsParams),
          getAdminUsers({ page_size: 200 } as API.getAdminUsersParams),
        ]);

        const roomList = Array.isArray(roomResponse)
          ? roomResponse
          : ((roomResponse as any)?.data ?? []);
        const userList = userResponse?.data ?? [];

        setRoomOptions(
          roomList
            .filter((room: API.Room) => room.id !== undefined && room.id !== null)
            .map((room: API.Room) => ({
              label: `${room.room_number} · ${room.room_type || '未设置房型'} · ${room.floor ?? '-'}楼`,
              value: Number(room.id),
            }))
            .filter((option: { label: string; value: number }) => Number.isFinite(option.value)),
        );

        const activeAdmins = userList.filter(
          (user: API.User) =>
            user.id !== undefined && user.id !== null && user.status === 'active' && user.role === 'admin',
        );

        setStaffOptions(
          activeAdmins.map((user: API.User) => ({
            label: `${user.real_name || user.username} (${user.username})`,
            value: String(user.id),
          })),
        );
      } catch (error) {
        message.error(getBackendErrorMessage(error, '加载房间和人员数据失败'));
      } finally {
        setOptionsLoading(false);
      }
    };

    void loadOptions();
  }, []);

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
      title: '清洁人员',
      dataIndex: 'staff_id',
      hideInSearch: true,
      render: (_, record) => {
        if (!record.staff_id) {
          return '-';
        }
        const matchedStaff = staffOptions.find((staff) => staff.value === String(record.staff_id));
        return matchedStaff?.label ?? record.staff_id;
      },
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
                    await completeCleaning({ id: Number(record.id) });
                    message.success('清洁任务已完成');
                    actionRef.current?.reload();
                  }
                } catch (error) {
                  message.error(getBackendErrorMessage(error, '操作失败'));
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
      const res = await createCleaningTask({
        ...values,
        room_id: Number(values.room_id),
      });
      if (res.success) {
        message.success('清洁任务已创建');
        setCreateModalVisible(false);
        form.resetFields();
        actionRef.current?.reload();
      }
    } catch (error) {
      message.error(getBackendErrorMessage(error, '创建失败'));
    }
  };

  const handleAssign = async (values: any) => {
    if (!currentRow?.id) return;
    try {
      const res = await assignStaff(
        { id: Number(currentRow.id) },
        {
          ...values,
          staff_id: values.staff_id,
        },
      );
      if (res.success) {
        message.success('人员已指派');
        setAssignModalVisible(false);
        assignForm.resetFields();
        actionRef.current?.reload();
      }
    } catch (error) {
      message.error(getBackendErrorMessage(error, '操作失败'));
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
        onCancel={() => {
          setCreateModalVisible(false);
          form.resetFields();
        }}
        onOk={() => form.submit()}
        destroyOnHidden
      >
        <Form form={form} onFinish={handleCreate} layout="vertical">
          <Form.Item name="room_id" label="房间" rules={[{ required: true, message: '请选择房间' }]}>
            <Select
              showSearch
              placeholder="请选择房间"
              options={roomOptions}
              loading={optionsLoading}
              optionFilterProp="label"
            />
          </Form.Item>
          <Form.Item name="type" label="清洁类型" rules={[{ required: true, message: '请选择清洁类型' }]}>
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
        onCancel={() => {
          setAssignModalVisible(false);
          assignForm.resetFields();
        }}
        onOk={() => assignForm.submit()}
        destroyOnHidden
      >
        <Form form={assignForm} onFinish={handleAssign} layout="vertical">
          <Form.Item name="staff_id" label="清洁人员" rules={[{ required: true, message: '请选择清洁人员' }]}>
            <Select
              showSearch
              placeholder="请选择清洁人员"
              options={staffOptions}
              loading={optionsLoading}
              optionFilterProp="label"
            />
          </Form.Item>
        </Form>
      </Modal>
    </PageContainer>
  );
};

export default CleaningManage;
