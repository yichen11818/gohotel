import {
  ModalForm,
  ProFormText,
  ProFormDigit,
  ProFormDateTimePicker,
  ProFormInstance,
} from '@ant-design/pro-components';
import { useRequest } from '@umijs/max';
import { message } from 'antd';
import dayjs from 'dayjs';
import type { FC } from 'react';
import { useRef, useEffect } from 'react';
import { postAdminNoticesId } from '@/services/api/gonggaoguanli';
import { getBackendErrorMessage } from '@/utils/backendError';


interface EditNoticeFormProps {
  visible: boolean;
  onCancel: () => void;
  onSuccess: () => void;
  noticeData: API.Notice;
}

const EditNoticeForm: FC<EditNoticeFormProps> = (props) => {
  const { visible, onCancel, onSuccess, noticeData } = props;
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

  // 监听活动数据变化，初始化表单
  useEffect(() => {
    if (visible && noticeData) {
      // 初始化表单数据
      formRef.current?.setFieldsValue({
        title: noticeData.title || '',
        link_url: noticeData.link_url || '',
        sort: noticeData.sort || 0,
        start_time: noticeData.start_time ? dayjs(noticeData.start_time) : undefined,
        end_time: noticeData.end_time ? dayjs(noticeData.end_time) : undefined,
      });

    }
  }, [visible, noticeData]);


  // 表单提交
  const { run, loading } = useRequest(
    async (data: any) => {
      try {
        const payload: Record<string, any> = {
          ...data,
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
        await postAdminNoticesId(
          {
            id: String(noticeData.id),
          },
          payload,
        );
        messageApi.success('公告更新成功');
        return true;
      } catch (error) {
        const errorMessage = getBackendErrorMessage(error, '公告更新失败，请重试');
        messageApi.error(errorMessage);
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
        title={'编辑公告'}    
        open={visible}
        width="600px"
        formRef={formRef}
        onOpenChange={(open) => {
          if (!open) {
            onCancel();
          }
        }}
        modalProps={{
          okButtonProps: {
            loading,
          },
          onCancel,
        }}
        onFinish={async (value) => {
          try {
            await run({
              ...value,
            });
            
            // 调用成功回调
            onSuccess();
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

export default EditNoticeForm;
