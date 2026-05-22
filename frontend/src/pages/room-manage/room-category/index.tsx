import { DeleteOutlined, PlusOutlined } from '@ant-design/icons';
import { PageContainer, ProTable } from '@ant-design/pro-components';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { Button, Form, Image, Input, message, Modal, Popconfirm, Select, Space, Tag } from 'antd';
import React, { useMemo, useRef, useState } from 'react';

import {
  createAdminRoomCategory,
  deleteAdminRoomCategory,
  getAdminRoomCategories,
  updateAdminRoomCategory,
} from '@/services/roomCategories';
import { getBackendErrorMessage } from '@/utils/backendError';
import { parseJsonArray, stringifyJsonArray } from '@/utils/roomCategory';

type CategoryFormValues = {
  description?: string;
  facilities?: string[];
  images?: string[];
  name: string;
};

const RoomCategoryPage: React.FC = () => {
  const actionRef = useRef<ActionType | null>(null);
  const [messageApi, contextHolder] = message.useMessage();
  const [form] = Form.useForm<CategoryFormValues>();
  const [modalOpen, setModalOpen] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [editingCategory, setEditingCategory] = useState<API.RoomCategory | undefined>();

  const openCreateModal = () => {
    setEditingCategory(undefined);
    form.resetFields();
    setModalOpen(true);
  };

  const openEditModal = (category: API.RoomCategory) => {
    setEditingCategory(category);
    form.setFieldsValue({
      name: category.name || '',
      description: category.description || '',
      facilities: parseJsonArray(category.facilities),
      images: parseJsonArray(category.images),
    });
    setModalOpen(true);
  };

  const closeModal = () => {
    setModalOpen(false);
    setEditingCategory(undefined);
    form.resetFields();
  };

  const columns: ProColumns<API.RoomCategory>[] = useMemo(
    () => [
      {
        title: '房型名称',
        dataIndex: 'name',
        render: (_, record) => <a onClick={() => openEditModal(record)}>{record.name || '-'}</a>,
      },
      {
        title: '预览图',
        dataIndex: 'images',
        hideInSearch: true,
        render: (_, record) => {
          const images = parseJsonArray(record.images);
          if (!images.length) {
            return '-';
          }

          return (
            <Space>
              <Image
                src={images[0]}
                width={72}
                height={48}
                style={{ objectFit: 'cover', borderRadius: 8 }}
              />
              <span>{images.length} 张</span>
            </Space>
          );
        },
      },
      {
        title: '房型描述',
        dataIndex: 'description',
        ellipsis: true,
        hideInSearch: true,
      },
      {
        title: '设施',
        dataIndex: 'facilities',
        hideInSearch: true,
        render: (_, record) => {
          const facilities = parseJsonArray(record.facilities);
          if (!facilities.length) {
            return '-';
          }

          return facilities.slice(0, 4).map((item) => <Tag key={item}>{item}</Tag>);
        },
      },
      {
        title: '关联房间数',
        dataIndex: 'usage_count',
        hideInSearch: true,
      },
      {
        title: '更新时间',
        dataIndex: 'updated_at',
        valueType: 'dateTime',
        hideInSearch: true,
      },
      {
        title: '操作',
        dataIndex: 'option',
        valueType: 'option',
        render: (_, record) => [
          <a key="edit" onClick={() => openEditModal(record)}>
            编辑
          </a>,
          <Popconfirm
            key="delete"
            title="确定删除这个房型分类吗？"
            description="如果还有房间或库存使用它，系统会阻止删除。"
            onConfirm={async () => {
              try {
                if (!record.id) {
                  return;
                }
                await deleteAdminRoomCategory(record.id);
                messageApi.success('房型分类删除成功');
                actionRef.current?.reload();
              } catch (error) {
                messageApi.error(getBackendErrorMessage(error, '删除房型分类失败'));
              }
            }}
          >
            <Button type="link" danger icon={<DeleteOutlined />} style={{ paddingInline: 0 }}>
              删除
            </Button>
          </Popconfirm>,
        ],
      },
    ],
    [messageApi],
  );

  const handleSubmit = async (values: CategoryFormValues) => {
    setSubmitting(true);
    try {
      const payload = {
        name: values.name.trim(),
        description: values.description?.trim(),
        facilities: stringifyJsonArray(values.facilities),
        images: stringifyJsonArray(values.images),
      };

      if (editingCategory?.id) {
        await updateAdminRoomCategory(editingCategory.id, payload);
        messageApi.success('房型分类更新成功');
      } else {
        await createAdminRoomCategory(payload);
        messageApi.success('房型分类创建成功');
      }

      closeModal();
      actionRef.current?.reload();
      return true;
    } catch (error) {
      messageApi.error(getBackendErrorMessage(error, '保存房型分类失败'));
      return false;
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <PageContainer>
      {contextHolder}
      <ProTable<API.RoomCategory>
        rowKey="id"
        actionRef={actionRef}
        headerTitle="房型分类"
        toolBarRender={() => [
          <Button key="create" type="primary" icon={<PlusOutlined />} onClick={openCreateModal}>
            新建房型分类
          </Button>,
        ]}
        request={async (params) => {
          try {
            const response = await getAdminRoomCategories({
              page: params.current || 1,
              page_size: params.pageSize || 10,
            });

            return {
              data: response.data || [],
              success: response.success ?? true,
              total: response.page?.total || 0,
            };
          } catch (error) {
            messageApi.error(getBackendErrorMessage(error, '获取房型分类失败'));
            return {
              data: [],
              success: false,
              total: 0,
            };
          }
        }}
        columns={columns}
      />

      <Modal
        title={editingCategory ? '编辑房型分类' : '新建房型分类'}
        open={modalOpen}
        onCancel={closeModal}
        onOk={() => form.submit()}
        okButtonProps={{ loading: submitting }}
        destroyOnHidden
      >
        <Form<CategoryFormValues> form={form} layout="vertical" onFinish={handleSubmit}>
          <Form.Item name="name" label="房型名称" rules={[{ required: true, message: '请输入房型名称' }]}>
            <Input placeholder="例如：豪华套房" />
          </Form.Item>
          <Form.Item name="description" label="房型描述">
            <Input.TextArea rows={4} placeholder="输入统一展示给该房型的描述文案" />
          </Form.Item>
          <Form.Item name="facilities" label="设施">
            <Select
              mode="tags"
              placeholder="输入设施后回车，例如 WiFi、浴缸、迷你吧"
              tokenSeparators={[',', '，']}
            />
          </Form.Item>
          <Form.Item name="images" label="预览图链接">
            <Select
              mode="tags"
              placeholder="输入图片链接后回车，可维护多张"
              tokenSeparators={[',']}
            />
          </Form.Item>
        </Form>
      </Modal>
    </PageContainer>
  );
};

export default RoomCategoryPage;
