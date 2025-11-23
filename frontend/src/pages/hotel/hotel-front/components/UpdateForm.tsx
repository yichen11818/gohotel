// import { updateRule } from '@/services/ant-design-pro/api';
// TODO: 等待后端提供更新用户的 API
import {
  ProFormSelect,
  ProFormText,
  StepsForm,
} from '@ant-design/pro-components';
import { Modal, message } from 'antd';
import React, { cloneElement, useCallback, useState } from 'react';
export type FormValueType = Partial<API.User>;
export type UpdateFormProps = {
  trigger?: React.ReactElement<any>;
  onOk?: () => void;
  values: Partial<API.User>;
};
const UpdateForm: React.FC<UpdateFormProps> = (props) => {
  const { onOk, values, trigger } = props;
  const [open, setOpen] = useState(false);
  const [messageApi, contextHolder] = message.useMessage();
  // TODO: 等待后端提供更新用户的 API
  const run = async (data: any) => {
    messageApi.warning('更新用户功能待实现');
    return Promise.resolve();
  };
  const onCancel = useCallback(() => {
    setOpen(false);
  }, []);
  const onOpen = useCallback(() => {
    setOpen(true);
  }, []);
  const onFinish = useCallback(
    async (values?: any) => {
      await run({
        data: values,
      });
      onCancel();
    },
    [onCancel, run],
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
              title={'编辑用户'}
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
            name="username"
            label={'用户名'}
            width="md"
            rules={[
              {
                required: true,
                message: '请输入用户名！',
              },
            ]}
          />
          <ProFormText
            name="email"
            label={'邮箱'}
            width="md"
            rules={[
              {
                required: true,
                message: '请输入邮箱！',
              },
              {
                type: 'email',
                message: '请输入有效的邮箱地址',
              },
            ]}
          />
          <ProFormText name="real_name" label={'真实姓名'} width="md" />
          <ProFormText name="phone" label={'手机号'} width="md" />
        </StepsForm.StepForm>
        <StepsForm.StepForm initialValues={values} title={'角色和状态'}>
          <ProFormSelect
            name="role"
            width="md"
            label={'角色'}
            valueEnum={{
              user: '用户',
              admin: '管理员',
            }}
          />
          <ProFormSelect
            name="status"
            width="md"
            label={'状态'}
            valueEnum={{
              active: '活跃',
              blocked: '已封禁',
            }}
          />
        </StepsForm.StepForm>
      </StepsForm>
    </>
  );
};
export default UpdateForm;
