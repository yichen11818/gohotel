import { PlusOutlined } from '@ant-design/icons';
import {
  type ActionType,
  ModalForm,
  ProFormText,
  ProFormDigit,
  ProFormSelect,
  ProFormTextArea,
} from '@ant-design/pro-components';
import { useRequest } from '@umijs/max';
import { Button, message } from 'antd';
import type { FC } from 'react';
import { postRooms } from '@/services/api/guanliyuan';

interface CreateFormProps {
  reload?: ActionType['reload'];
}

const CreateForm: FC<CreateFormProps> = (props) => {
  const { reload } = props;
  const [messageApi, contextHolder] = message.useMessage();

  const { run, loading } = useRequest(
    async (data: { data: API.CreateRoomRequest }) => {
      try {
        await postRooms(data.data);
        messageApi.success('房间创建成功');
        return true;
      } catch (error) {
        messageApi.error('房间创建失败，请重试');
        throw error;
      }
    },
    {
      manual: true,
    }
  );

  return (
    <>
      {contextHolder}
      <ModalForm
        title={'新建房间'}
        trigger={
          <Button type="primary" icon={<PlusOutlined />}>
            新建
          </Button>
        }
        width="600px"
        modalProps={{
          okButtonProps: {
            loading,
          },
        }}
        onFinish={async (value) => {
          try {
            await run({
              data: value as API.CreateRoomRequest,
            });
            if (reload) {
              reload();
            }
            return true;
          } catch (error) {
            return false;
          }
        }}
      >
        <ProFormText
          rules={[
            {
              required: true,
              message: '房间号为必填项',
            },
          ]}
          width="md"
          name="room_number"
          label="房间号"
          placeholder="例如: 101"
        />

        <ProFormSelect
          rules={[
            {
              required: true,
              message: '房型为必填项',
            },
          ]}
          width="md"
          name="room_type"
          label="房型"
          valueEnum={{
            '单人间': '单人间',
            '双人间': '双人间',
            '豪华套房': '豪华套房',
            '总统套房': '总统套房',
            '商务套房': '商务套房',
          }}
          placeholder="请选择房型"
        />

        <ProFormDigit
          rules={[
            {
              required: true,
              message: '楼层为必填项',
            },
          ]}
          width="md"
          name="floor"
          label="楼层"
          min={1}
          max={100}
          fieldProps={{ precision: 0 }}
        />

        <ProFormDigit
          rules={[
            {
              required: true,
              message: '价格为必填项',
            },
          ]}
          width="md"
          name="price"
          label="价格(每晚)"
          min={0}
          fieldProps={{
            precision: 2,
            addonBefore: '¥',
          }}
          placeholder="0.00"
        />

        <ProFormDigit
          width="md"
          name="original_price"
          label="原价"
          min={0}
          fieldProps={{
            precision: 2,
            addonBefore: '¥',
          }}
          placeholder="0.00"
        />

        <ProFormDigit
          rules={[
            {
              required: true,
              message: '可住人数为必填项',
            },
          ]}
          width="md"
          name="capacity"
          label="可住人数"
          min={1}
          max={10}
          fieldProps={{ precision: 0 }}
        />

        <ProFormDigit
          width="md"
          name="area"
          label="面积(m²)"
          min={0}
          fieldProps={{ precision: 1 }}
        />

        <ProFormSelect
          width="md"
          name="bed_type"
          label="床型"
          valueEnum={{
            '单人床': '单人床',
            '双人床': '双人床',
            '大床': '大床',
            '两张单人床': '两张单人床',
          }}
          placeholder="请选择床型"
        />

        <ProFormSelect
          width="md"
          name="status"
          label="状态"
          valueEnum={{
            available: '可用',
            occupied: '占用',
            maintenance: '维护中',
          }}
        />

        <ProFormTextArea
          width="md"
          name="description"
          label="房间描述"
          placeholder="请输入房间描述"
          fieldProps={{
            rows: 3,
          }}
        />

        <ProFormTextArea
          width="md"
          name="facilities"
          label="设施(JSON格式)"
          placeholder='例如: ["WiFi", "空调", "电视"]'
          fieldProps={{
            rows: 2,
          }}
        />
      </ModalForm>
    </>
  );
};

export default CreateForm;

