import {
  ProFormSelect,
  ProFormText,
  ProFormDigit,
  StepsForm,
} from '@ant-design/pro-components';
import { history } from '@umijs/max';
import { Alert, Button, Modal, message } from 'antd';
import React, { cloneElement, useCallback, useState, useEffect } from 'react';
import { postRoomsId } from '@/services/api/guanliyuan';
import { getBackendErrorMessage } from '@/utils/backendError';

export type FormValueType = Partial<API.Room>;

export type UpdateFormProps = {
  trigger?: React.ReactElement<any>;
  onOk?: () => void;
  onCancel?: () => void;
  roomCategories?: API.RoomCategory[];
  values: Partial<API.Room>;
  visible?: boolean; // 受控模式
};

const UpdateForm: React.FC<UpdateFormProps> = (props) => {
  const { onOk, values, trigger, visible, onCancel: onCancelProp, roomCategories = [] } = props;
  const [open, setOpen] = useState(false);
  const [messageApi, contextHolder] = message.useMessage();
  const hasRoomCategories = roomCategories.some((item) => item.name);

  // 受控模式：当 visible 变化时同步内部状态
  useEffect(() => {
    if (visible !== undefined) {
      setOpen(visible);
    }
  }, [visible]);

  const run = async (data: any) => {
    try {
      if (values.id) {
        await postRoomsId({ id: Number(values.id) }, data as API.UpdateRoomRequest);
        messageApi.success('房间更新成功');
        if (onOk) {
          onOk();
        }
      }
    } catch (error) {
      messageApi.error(getBackendErrorMessage(error, '房间更新失败'));
      throw error;
    }
  };

  const onCancel = useCallback(() => {
    setOpen(false);
    if (onCancelProp) {
      onCancelProp();
    }
  }, [onCancelProp]);

  const onOpen = useCallback(() => {
    if (!hasRoomCategories) {
      messageApi.warning('请先创建房型分类，再编辑房间');
      return;
    }
    setOpen(true);
  }, [hasRoomCategories, messageApi]);

  const onFinish = useCallback(
    async (values?: any) => {
      await run(values);
      onCancel();
    },
    [onCancel],
  );

  return (
    <>
      {contextHolder}
      {trigger
        ? cloneElement(trigger, {
            onClick: onOpen,
          })
        : null}
      <StepsForm
        stepsProps={{
          size: 'small',
        }}
        stepsFormRender={(dom, submitter) => {
          return (
            <Modal
              width={640}
              styles={{
                body: {
                  padding: '32px 40px 48px',
                },
              }}
              destroyOnHidden
              title={'编辑房间'}
              open={open}
              footer={submitter}
              onCancel={onCancel}
            >
              {dom}
            </Modal>
          );
        }}
        onFinish={onFinish}
      >
        <StepsForm.StepForm initialValues={values} title={'基本信息'}>
          <ProFormText
            name="room_number"
            label={'房间号'}
            width="md"
            rules={[
              {
                required: true,
                message: '请输入房间号！',
              },
            ]}
          />
          <ProFormSelect
            name="room_type"
            label={'房型'}
            width="md"
            rules={[
              {
                required: true,
                message: '请选择房型！',
              },
            ]}
            options={roomCategories
              .filter((item) => item.name)
              .map((item) => ({ label: item.name as string, value: item.name as string }))}
          />
          <ProFormDigit
            name="floor"
            label={'楼层'}
            width="md"
            min={1}
            max={100}
            fieldProps={{ precision: 0 }}
            rules={[
              {
                required: true,
                message: '请输入楼层！',
              },
            ]}
          />
          <ProFormDigit
            name="capacity"
            label={'可住人数'}
            width="md"
            min={1}
            max={10}
            fieldProps={{ precision: 0 }}
            rules={[
              {
                required: true,
                message: '请输入可住人数！',
              },
            ]}
          />
          <ProFormDigit
            name="area"
            label={'面积(m²)'}
            width="md"
            min={0}
            fieldProps={{ precision: 1 }}
          />
          <ProFormSelect
            name="bed_type"
            label={'床型'}
            width="md"
            valueEnum={{
              '单人床': '单人床',
              '双人床': '双人床',
              '大床': '大床',
              '两张单人床': '两张单人床',
            }}
          />
        </StepsForm.StepForm>
        <StepsForm.StepForm initialValues={values} title={'价格和状态'}>
          <ProFormDigit
            name="price"
            label={'价格(每晚)'}
            width="md"
            min={0}
            fieldProps={{
              precision: 2,
              addonBefore: '¥',
            }}
            rules={[
              {
                required: true,
                message: '请输入价格！',
              },
            ]}
          />
          <ProFormDigit
            name="original_price"
            label={'原价'}
            width="md"
            min={0}
            fieldProps={{
              precision: 2,
              addonBefore: '¥',
            }}
          />
          <ProFormSelect
            name="status"
            width="md"
            label={'状态'}
            valueEnum={{
              available: '空闲',
              occupied: '在住',
              maintenance: '维修中',
            }}
          />
          <ProFormSelect
            name="clean_status"
            width="md"
            label={'清洁状态'}
            valueEnum={{
              clean: '干净',
              dirty: '脏房',
              inspecting: '待查',
            }}
          />
        </StepsForm.StepForm>
        <StepsForm.StepForm initialValues={values} title={'分类继承'}>
          <Alert
            type="info"
            showIcon
            message="房型描述、预览图、设施来自房型分类"
            description="如果需要统一修改同类房间的展示信息，请前往“房型分类”页面维护。"
            action={
              <Button size="small" type="link" onClick={() => history.push('/room-manage/category')}>
                去维护房型分类
              </Button>
            }
          />
        </StepsForm.StepForm>
      </StepsForm>
    </>
  );
};

export default UpdateForm;










