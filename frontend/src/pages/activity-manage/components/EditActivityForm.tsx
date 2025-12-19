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
import { useState, useRef, useEffect } from 'react';
import { postAdminBannersId } from '@/services/api/huodongguanli';
import { postUploadImage } from '@/services/api/wenjianshangchuan';

interface EditActivityFormProps {
  visible: boolean;
  onCancel: () => void;
  onSuccess: () => void;
  activityData: API.Banner;
}

const EditActivityForm: FC<EditActivityFormProps> = (props) => {
  const { visible, onCancel, onSuccess, activityData } = props;
  const [messageApi, contextHolder] = message.useMessage();
  const [tempImageUrl, setTempImageUrl] = useState<string>('');
  const [_uploading, setUploading] = useState<boolean>(false);
  const [uploadedFile, setUploadedFile] = useState<any>(null);
  const formRef = useRef<ProFormInstance>(null);

  // 监听活动数据变化，初始化表单
  useEffect(() => {
    if (visible && activityData) {
      // 初始化表单数据
      formRef.current?.setFieldsValue({
        title: activityData.title || '',
        subtitle: activityData.subtitle || '',
        link_url: activityData.link_url || '',
        sort: activityData.sort || 0,
        start_time: activityData.start_time || '',
        end_time: activityData.end_time || '',
      });

      // 初始化图片
      if (activityData.image_url) {
        // 对于已有的图片，我们不设置tempImageUrl，只有新上传的图片才会有tempImageUrl
        // 这样可以区分是使用原有图片还是新上传的图片
        setTempImageUrl('');
        
        // 创建已上传文件对象
        const existingFile = {
          uid: String(Date.now()),
          name: `existing_${activityData.id}`,
          status: 'done' as const,
          url: activityData.image_url,
        };
        
        setUploadedFile(existingFile);
      } else {
        setUploadedFile(null);
        setTempImageUrl('');
      }
    }
  }, [visible, activityData]);

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
        // 调用API更新活动
        await postAdminBannersId({
          id: String(activityData.id),
        }, {
          ...data,
          temp_url: tempImageUrl,
        });
        messageApi.success('活动更新成功');
        return true;
      } catch (error) {
        messageApi.error('活动更新失败，请重试');
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
        title={'编辑活动横幅'}
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
            
            // 重置图片URL和上传文件
            setTempImageUrl('');
            setUploadedFile(null);
            
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

export default EditActivityForm;