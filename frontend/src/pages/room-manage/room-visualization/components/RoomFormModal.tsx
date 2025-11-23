import React, { useEffect } from 'react';
import { Modal, Form, Input, InputNumber, Select, message } from 'antd';

const { TextArea } = Input;
const { Option } = Select;

interface RoomFormModalProps {
  visible: boolean;
  room?: API.Room | null;
  onCancel: () => void;
  onSuccess: () => void;
}

const RoomFormModal: React.FC<RoomFormModalProps> = ({ visible, room, onCancel, onSuccess }) => {
  const [form] = Form.useForm();

  useEffect(() => {
    if (visible && room) {
      // 编辑模式,填充表单数据
      form.setFieldsValue({
        room_number: room.room_number,
        room_type: room.room_type,
        floor: room.floor,
        price: room.price,
        original_price: room.original_price,
        capacity: room.capacity,
        area: room.area,
        bed_type: room.bed_type,
        status: room.status,
        description: room.description,
        facilities: room.facilities,
      });
    } else if (visible) {
      // 新建模式,清空表单
      form.resetFields();
    }
  }, [visible, room, form]);

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      console.log('Form values:', values);
      
      // 这里应该调用创建或更新的 API
      if (room) {
        // 更新房间
        // await updateRoom({ id: room.id, ...values });
        message.success('房间更新成功(模拟)');
      } else {
        // 创建房间
        // await createRoom(values);
        message.success('房间创建成功(模拟)');
      }
      
      onSuccess();
      form.resetFields();
    } catch (error) {
      console.error('Validation failed:', error);
    }
  };

  return (
    <Modal
      title={room ? '编辑房间' : '新建房间'}
      open={visible}
      onCancel={onCancel}
      onOk={handleSubmit}
      width={600}
      destroyOnClose
    >
      <Form
        form={form}
        layout="vertical"
        initialValues={{
          status: 'available',
          capacity: 1,
          floor: 1,
        }}
      >
        <Form.Item
          name="room_number"
          label="房间号"
          rules={[{ required: true, message: '请输入房间号' }]}
        >
          <Input placeholder="例如: 101" />
        </Form.Item>

        <Form.Item
          name="room_type"
          label="房型"
          rules={[{ required: true, message: '请输入房型' }]}
        >
          <Select placeholder="请选择房型">
            <Option value="单人间">单人间</Option>
            <Option value="双人间">双人间</Option>
            <Option value="豪华套房">豪华套房</Option>
            <Option value="总统套房">总统套房</Option>
            <Option value="商务套房">商务套房</Option>
          </Select>
        </Form.Item>

        <Form.Item
          name="floor"
          label="楼层"
          rules={[{ required: true, message: '请输入楼层' }]}
        >
          <InputNumber min={1} max={100} style={{ width: '100%' }} />
        </Form.Item>

        <Form.Item
          name="price"
          label="价格(每晚)"
          rules={[{ required: true, message: '请输入价格' }]}
        >
          <InputNumber
            min={0}
            precision={2}
            style={{ width: '100%' }}
            prefix="¥"
            placeholder="0.00"
          />
        </Form.Item>

        <Form.Item name="original_price" label="原价">
          <InputNumber
            min={0}
            precision={2}
            style={{ width: '100%' }}
            prefix="¥"
            placeholder="0.00"
          />
        </Form.Item>

        <Form.Item
          name="capacity"
          label="可住人数"
          rules={[{ required: true, message: '请输入可住人数' }]}
        >
          <InputNumber min={1} max={10} style={{ width: '100%' }} />
        </Form.Item>

        <Form.Item name="area" label="面积(m²)">
          <InputNumber min={0} precision={1} style={{ width: '100%' }} />
        </Form.Item>

        <Form.Item name="bed_type" label="床型">
          <Select placeholder="请选择床型">
            <Option value="单人床">单人床</Option>
            <Option value="双人床">双人床</Option>
            <Option value="大床">大床</Option>
            <Option value="两张单人床">两张单人床</Option>
          </Select>
        </Form.Item>

        <Form.Item name="status" label="状态">
          <Select>
            <Option value="available">可用</Option>
            <Option value="occupied">占用</Option>
            <Option value="maintenance">维护中</Option>
          </Select>
        </Form.Item>

        <Form.Item name="description" label="房间描述">
          <TextArea rows={3} placeholder="请输入房间描述" />
        </Form.Item>

        <Form.Item name="facilities" label="设施(JSON格式)">
          <TextArea
            rows={2}
            placeholder='例如: ["WiFi", "空调", "电视"]'
          />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default RoomFormModal;

