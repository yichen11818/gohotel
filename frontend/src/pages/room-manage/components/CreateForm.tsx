import { PlusOutlined } from '@ant-design/icons';
import {
  type ActionType,
  ModalForm,
  ProFormText,
  ProFormDigit,
  ProFormSelect,
} from '@ant-design/pro-components';
import { history, useRequest } from '@umijs/max';
import { Alert, Button, message } from 'antd';
import type { FC } from 'react';
import { postRooms } from '@/services/api/guanliyuan';
import { getBackendErrorMessage } from '@/utils/backendError';

interface CreateFormProps {
  reload?: ActionType['reload'];
  roomCategories?: API.RoomCategory[];
}

const CreateForm: FC<CreateFormProps> = (props) => {
  const { reload, roomCategories = [] } = props;
  const [messageApi, contextHolder] = message.useMessage();
  const hasRoomCategories = roomCategories.some((item) => item.name);

  const { run, loading } = useRequest(
    async (data: { data: API.CreateRoomRequest }) => {
      try {
        await postRooms(data.data);
        messageApi.success('房间创建成功');
        return true;
      } catch (error) {
        messageApi.error(getBackendErrorMessage(error, '房间创建失败，请重试'));
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
          if (!hasRoomCategories) {
            messageApi.warning('请先创建房型分类后再新增房间');
            return false;
          }
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
        {!hasRoomCategories && (
          <Alert
            type="warning"
            showIcon
            message="请先创建房型分类"
            description="房间必须绑定已有房型分类后才能创建。"
            action={
              <Button size="small" type="link" onClick={() => history.push('/room-manage/category')}>
                去维护房型分类
              </Button>
            }
          />
        )}
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
          options={roomCategories
            .filter((item) => item.name)
            .map((item) => ({ label: item.name as string, value: item.name as string }))}
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

        <Alert
          type="info"
          showIcon
          style={{ marginTop: 12 }}
          message="房型描述、预览图、设施由“房型分类”统一维护"
          description="当前房间只维护房号、楼层、价格、可住人数等差异化信息。"
        />
      </ModalForm>
    </>
  );
};

export default CreateForm;

