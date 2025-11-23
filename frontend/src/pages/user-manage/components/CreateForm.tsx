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
import { postAdminUsersUser } from '@/services/api/guanliyuan';
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

  const { run, loading } = useRequest(
    async (data: { data: API.AddUserRequest }) => {
      try {
        await postAdminUsersUser(data.data);
        messageApi.success('用户创建成功');
        return true;
      } catch (error) {
        messageApi.error('用户创建失败，请重试');
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
          try {
            await run({
              data: value as API.AddUserRequest,
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
