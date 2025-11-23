// import { addRule } from '@/services/ant-design-pro/api';
// TODO: 等待后端提供创建用户的 API
import { PlusOutlined } from '@ant-design/icons';
import {
  type ActionType,
  ModalForm,
  ProFormText,
  ProFormTextArea,
} from '@ant-design/pro-components';
import { useRequest } from '@umijs/max';
import { Button, message } from 'antd';
import type { FC } from 'react';
interface CreateFormProps {
  reload?: ActionType['reload'];
}
const CreateForm: FC<CreateFormProps> = (props) => {
  const { reload } = props;
  const [messageApi, contextHolder] = message.useMessage();
  /**
   * @en-US International configuration
   * @zh-CN 国际化配置
   * */

  // TODO: 等待后端提供创建用户的 API
  const loading = false;
  const run = async (data: any) => {
    messageApi.warning('创建用户功能待实现');
    return Promise.resolve();
  };
  return (
    <>
      {contextHolder}
      <ModalForm
        title={'新建用户'}
        trigger={
          <Button type="primary" icon={<PlusOutlined />}>
            新建
          </Button>
        }
        width="400px"
        modalProps={{
          okButtonProps: {
            loading,
          },
        }}
        onFinish={async (value) => {
          await run({
            data: value as API.User,
          });
          reload?.();
          return true;
        }}
      >
        <ProFormText
          rules={[
            {
              required: true,
              message: '用户名为必填项',
            },
          ]}
          width="md"
          name="username"
          label="用户名"
        />
        <ProFormText
          rules={[
            {
              required: true,
              message: '邮箱为必填项',
            },
            {
              type: 'email',
              message: '请输入有效的邮箱地址',
            },
          ]}
          width="md"
          name="email"
          label="邮箱"
        />
        <ProFormText width="md" name="real_name" label="真实姓名" />
        <ProFormText width="md" name="phone" label="手机号" />
      </ModalForm>
    </>
  );
};
export default CreateForm;
