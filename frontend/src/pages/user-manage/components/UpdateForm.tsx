import { postAdminUsersId } from '@/services/api/guanliyuan';
import {
  ProFormDigit,
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
  const run = async (data: API.UpdateUserRequest) => {
    if (!values.id) {
      throw new Error('缺少用户ID');
    }

    await postAdminUsersId({ id: Number(values.id) }, data);
    messageApi.success('更新用户成功');
  };
  const onCancel = useCallback(() => {
    setOpen(false);
  }, []);
  const onOpen = useCallback(() => {
    setOpen(true);
  }, []);
  const onFinish = useCallback(
    async (values?: any) => {
      await run(values as API.UpdateUserRequest);
      onOk?.();
      onCancel();
    },
    [onCancel, onOk],
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
          <ProFormText name="avatar" label={'头像 URL'} width="md" />
        </StepsForm.StepForm>
        <StepsForm.StepForm initialValues={values} title={'角色与会员信息'}>
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
          <ProFormSelect
            name="level"
            width="md"
            label={'会员等级'}
            valueEnum={{
              normal: '普通会员',
              silver: '白银会员',
              gold: '黄金会员',
              platinum: '铂金会员',
            }}
          />
          <ProFormDigit name="points" label={'积分'} width="md" min={0} />
          <ProFormDigit name="balance" label={'余额'} width="md" min={0} />
          <ProFormDigit name="total_spend" label={'累计消费'} width="md" min={0} />
        </StepsForm.StepForm>
      </StepsForm>
    </>
  );
};
export default UpdateForm;
