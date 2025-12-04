import { PlusOutlined } from '@ant-design/icons';
import {
  type ActionType,
  ModalForm,
  ProFormText,
  ProFormSelect,
  ProFormDigit,
  ProFormDateTimePicker,
  ProFormItem,
  ProFormInstance,
} from '@ant-design/pro-components';
import { useRequest } from '@umijs/max';
import { Button, message, Upload, type UploadProps } from 'antd';
import type { FC } from 'react';
import { useState, useRef } from 'react';
import { postAdminBanners } from '@/services/api/huodongguanli';
import { postUploadImage } from '@/services/api/wenjianshangchuan';

interface CreateActivityFormProps {
  reload?: () => void;
}

const CreateActivityForm: FC<CreateActivityFormProps> = (props) => {
  const { reload } = props;
  const [messageApi, contextHolder] = message.useMessage();
  const [tempImageUrl, setTempImageUrl] = useState<string>('');
  const [_uploading, setUploading] = useState<boolean>(false);
  const [uploadedFile, setUploadedFile] = useState<any>(null);
  const formRef = useRef<ProFormInstance>(null);

  // 图片上传处理
  const handleImageUpload: UploadProps['customRequest'] = async ({ file, onSuccess, onError }) => {
    setUploading(true);
    try {
      // 上传图片到服务器
      const response = await postUploadImage({ type: 'banner' }, file as File);
      
      // 检查上传结果
      if (response?.temp_url) {
        // 更新状态
        setTempImageUrl(response.temp_url);
        messageApi.success('图片上传成功');
        
        // 创建上传文件对象
        const uploadedFileObj = {
          uid: String(Date.now()),
          name: (file as File).name,
          status: 'done' as const,
          url: response.temp_url,
          response,
        };
        
        // 更新上传文件状态
        setUploadedFile(uploadedFileObj);
        // 调用onSuccess时只传递response，第二个参数应为XMLHttpRequest或undefined
        onSuccess?.(response);
      } else {
        // 上传失败处理
        const errorMsg = '图片上传失败';
        messageApi.error(errorMsg);
        onError?.(new Error(errorMsg));
      }
    } catch (error) {
      // 异常处理
      console.error('图片上传异常:', error);
      messageApi.error('图片上传失败，请重试');
      onError?.(error instanceof Error ? error : new Error('图片上传异常'));
    } finally {
      // 无论成功失败，都结束上传状态
      setUploading(false);
    }
  };

  // 文件变化处理
  const handleFileChange: UploadProps['onChange'] = ({ file, fileList }) => {
    // 处理文件上传状态变化
    if (file.status === 'done') {
      messageApi.success('图片上传成功');
    } else if (file.status === 'error') {
      messageApi.error('图片上传失败');
    }
  };

  // 处理图片删除
  const handleRemove: UploadProps['onRemove'] = () => {
    setUploadedFile(null);
    setTempImageUrl('');
    return true;
  };

  // 表单提交
  const { run, loading } = useRequest(
    async (data: any) => {
      try {
        await postAdminBanners(data);
        messageApi.success('活动创建成功');
        return true;
      } catch (error) {
        messageApi.error('活动创建失败，请重试');
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
    // 清空图片相关状态
    setTempImageUrl('');
    setUploadedFile(null);
  };

  return (
    <>
      {contextHolder}
      <ModalForm
        title={'新建活动横幅'}
        trigger={
          <Button type="primary" icon={<PlusOutlined />}>
            新增活动
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
            if (!tempImageUrl) {
              messageApi.error('请先上传图片');
              return false;
            }

            await run({
              ...value,
              temp_url: tempImageUrl,
            });
            
            // 重置图片URL和上传文件
            setTempImageUrl('');
            setUploadedFile(null);
            
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
          placeholder="请输入活动标题"
        />

        {/* 副标题 */}
        <ProFormText
          width="md"
          name="subtitle"
          label="副标题"
          placeholder="请输入活动副标题"
        />

        {/* 图片上传 */}
        <ProFormItem
          label="图片"
        >
          <Upload
            fileList={uploadedFile ? [uploadedFile] : []}
            customRequest={handleImageUpload}
            listType="picture-card"
            maxCount={1}
            showUploadList={true}
            onRemove={handleRemove}
            onChange={handleFileChange}
          >
            <div>
              <PlusOutlined />
              <div style={{ marginTop: 8 }}>上传</div>
            </div>
          </Upload>
        </ProFormItem>

        {/* 跳转链接 */}
        <ProFormText
          width="md"
          name="link_url"
          label="跳转链接"
          placeholder="请输入点击横幅后的跳转链接"
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
          placeholder="请选择活动开始时间"
        />

        {/* 结束时间 */}
        <ProFormDateTimePicker
          width="md"
          name="end_time"
          label="结束时间"
          placeholder="请选择活动结束时间"
        />
      </ModalForm>
    </>
  );
};

export default CreateActivityForm;