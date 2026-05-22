import { PlusOutlined } from '@ant-design/icons';
import {
  ModalForm,
  ProFormText,
  ProFormDigit,
  ProFormDateTimePicker,
  ProFormInstance,
} from '@ant-design/pro-components';
import { useRequest } from '@umijs/max';
import { Button, message } from 'antd';
import dayjs from 'dayjs';
import type { FC } from 'react';
import { useRef } from 'react';
import { postAdminNotices } from '@/services/api/gonggaoguanli';
import { getBackendErrorMessage } from '@/utils/backendError';


interface CreateNoticeFormProps {
  reload?: () => void;
}

const CreateNoticeForm: FC<CreateNoticeFormProps> = (props) => {
  const { reload } = props;
  const [messageApi, contextHolder] = message.useMessage();
  const formRef = useRef<ProFormInstance>(null);
  const normalizeDateTime = (value: any) => {
    if (value === undefined || value === null || value === '') {
      return undefined;
    }
    if (typeof value === 'string') {
      return value;
    }
    return dayjs(value).format('YYYY-MM-DD HH:mm:ss');
  };

  // 表单提交
  const { run, loading } = useRequest(
    async (data: any) => {
      try {
        const payload: Parameters<typeof postAdminNotices>[0] = {
          title: data.title,
          link_url: data.link_url,
          sort: data.sort,
        };
        const startTime = normalizeDateTime((data as any)?.start_time);
        const endTime = normalizeDateTime((data as any)?.end_time);
        if (startTime) {
          payload.start_time = startTime;
        } else {
          delete payload.start_time;
        }
        if (endTime) {
          payload.end_time = endTime;
        } else {
          delete payload.end_time;
        }
        await postAdminNotices(payload);
        messageApi.success('公告创建成功');
        return true;
      } catch (error) {
        const errorMessage = getBackendErrorMessage(error, '公告创建失败，请重试');
        messageApi.error(errorMessage);
        throw error;
      }
    },
    {
      manual: true,
    }
  );

  // 清空表单内容
  const resetForm = () => {
    // 清空表单字段
    formRef.current?.resetFields();
  };

  return (
    <>
      {contextHolder}
      <ModalForm
        title={'新建公告'}
        trigger={
          <Button type="primary" icon={<PlusOutlined />}>
            新增公告
          </Button>
        }
        width="600px"
        formRef={formRef}
        onOpenChange={(open) => {
          if (open) {
            // 打开模态框时清空表单
            resetForm();
          }
        }}
        modalProps={{
          okButtonProps: {
            loading,
          },
        }}
        onFinish={async (value) => {
          try {
            await run({
              ...value,
            });
            
            if (reload) {
              reload();
            }
            return true;
          } catch (_error) {
            return false;
          }
        }}
      >
        {/* 标题 */}
        <ProFormText
          rules={[
            {
              required: true,
              message: '标题为必填项',
            },
          ]}
          width="md"
          name="title"
          label="标题"
          placeholder="请输入公告标题"
        />

        {/* 跳转链接 */}
        <ProFormText
          width="md"
          name="link_url"
          label="跳转链接"
          placeholder="请输入点击公告后的跳转链接"
        />



        {/* 排序 */}
        <ProFormDigit
          width="md"
          name="sort"
          label="排序"
          min={0}
          fieldProps={{ precision: 0 }}
          initialValue={0}
        />

        {/* 开始时间 */}
        <ProFormDateTimePicker
          width="md"
          name="start_time"
          label="开始时间"
          placeholder="请选择公告开始时间"
        />

        {/* 结束时间 */}
        <ProFormDateTimePicker
          width="md"
          name="end_time"
          label="结束时间"
          placeholder="请选择公告结束时间"
        />
      </ModalForm>
    </>
  );
};

export default CreateNoticeForm;
