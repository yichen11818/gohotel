import {
  ProFormSelect,
  ProFormText,
  ProFormDigit,
  ProFormTextArea,
  StepsForm,
} from '@ant-design/pro-components';
import { Modal, message } from 'antd';
import React, { cloneElement, useCallback, useState } from 'react';
import { postRoomsId } from '@/services/api/guanliyuan';

export type FormValueType = Partial<API.Room>;

export type UpdateFormProps = {
  trigger?: React.ReactElement<any>;
  onOk?: () => void;
  values: Partial<API.Room>;
};

const UpdateForm: React.FC<UpdateFormProps> = (props) => {
  const { onOk, values, trigger } = props;
  const [open, setOpen] = useState(false);
  const [messageApi, contextHolder] = message.useMessage();

  const run = async (data: any) => {
    try {
      if (values.id) {
        await postRoomsId({ id: values.id }, data as API.UpdateRoomRequest);
        messageApi.success('房间更新成功');
        if (onOk) {
          onOk();
        }
      }
    } catch (error) {
      messageApi.error('房间更新失败');
      throw error;
    }
  };

  const onCancel = useCallback(() => {
    setOpen(false);
  }, []);

  const onOpen = useCallback(() => {
    setOpen(true);
  }, []);

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
              bodyStyle={{
                padding: '32px 40px 48px',
              }}
              destroyOnClose
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
            valueEnum={{
              '单人间': '单人间',
              '双人间': '双人间',
              '豪华套房': '豪华套房',
              '总统套房': '总统套房',
              '商务套房': '商务套房',
            }}
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
              available: '可用',
              occupied: '占用',
              maintenance: '维护中',
            }}
          />
        </StepsForm.StepForm>
        <StepsForm.StepForm initialValues={values} title={'详细信息'}>
          <ProFormTextArea
            name="description"
            label={'房间描述'}
            width="md"
            fieldProps={{
              rows: 3,
            }}
          />
          <ProFormTextArea
            name="facilities"
            label={'设施(JSON格式)'}
            width="md"
            placeholder='例如: ["WiFi", "空调", "电视"]'
            fieldProps={{
              rows: 2,
            }}
          />
        </StepsForm.StepForm>
      </StepsForm>
    </>
  );
};

export default UpdateForm;

