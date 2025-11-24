import React, { useState, useEffect } from 'react';
import {
  Steps,
  Form,
  Input,
  Button,
  message,
  Card,
  Descriptions,
  Tag,
  Spin,
} from 'antd';
import { SearchOutlined, CheckCircleOutlined, UserOutlined } from '@ant-design/icons';
import type { StepProps } from 'antd';
import { getAdminBookingsSearch } from '@/services/api/guanliyuan';

const { Step } = Steps;

interface BookingInfo {
  id: string;
  bookingCode: string;
  guestName: string;
  guestPhone: string;
  roomNumber: string;
  roomType: string;
  checkInDate: string;
  checkOutDate: string;
  status: string;
  totalAmount: number;
}

const CheckInForm: React.FC = () => {
  const [currentStep, setCurrentStep] = useState<number>(0);
  const [form] = Form.useForm();
  const [bookingInfo, setBookingInfo] = useState<BookingInfo | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [submitting, setSubmitting] = useState<boolean>(false);

  const steps: StepProps[] = [
    {
      title: '查询预订',
      description: '输入预订号或客人信息',
    },
    {
      title: '确认信息',
      description: '核对预订详情',
    },
    {
      title: '办理入住',
      description: '完成入住手续',
    },
  ];

  // 查询预订信息
  const handleSearchBooking = async () => {
    try {
      // 表单验证
      await form.validateFields();
      
      setLoading(true);
      const values = form.getFieldsValue();
      
      // 验证至少提供客人姓名或手机号中的一项
      if (!values.guestName && !values.guestPhone) {
        message.error('请至少输入客人姓名或手机号码中的一项');
        setLoading(false);
        return;
      }
      
      // 调用真实的API接口
      const response = await getAdminBookingsSearch({
        guest_name: values.guestName,
        guest_phone: values.guestPhone,
      });
      
      // 处理响应数据
      if (response.data && response.data.length > 0) {
        // 取第一个匹配的预订信息
        const booking = response.data[0];
        
        // 转换日期格式
        const formatDate = (dateString: string) => {
          if (!dateString) return '';
          const date = new Date(dateString);
          return date.toLocaleDateString('zh-CN', {
            year: 'numeric',
            month: '2-digit',
            day: '2-digit'
          });
        };
        
        // 转换为组件需要的格式
        const formattedBooking: BookingInfo = {
          id: String(booking.id),
          bookingCode: booking.booking_number || booking.booking_code || booking.bookingCode,
          guestName: booking.guest_name || booking.guestName,
          guestPhone: booking.guest_phone || booking.guestPhone,
          roomNumber: booking.room?.room_number || booking.roomNumber || '-',
          roomType: booking.room?.room_type || booking.roomType || '-',
          checkInDate: formatDate(booking.check_in) || booking.check_in_date || booking.checkInDate || '-',
          checkOutDate: formatDate(booking.check_out) || booking.check_out_date || booking.checkOutDate || '-',
          status: booking.status,
          totalAmount: booking.total_price || booking.total_amount || booking.totalAmount || 0,
        };
        
        setBookingInfo(formattedBooking);
        setCurrentStep(1);
        
        // 如果预订状态不是已确认，给出提示
        if (formattedBooking.status !== 'confirmed') {
          message.warning(`该预订当前状态为${
            formattedBooking.status === 'checkin' ? '已入住' : 
            formattedBooking.status === 'checkout' ? '已退房' : '未确认'
          }，可能无法办理入住`);
        }
      } else {
        message.error('未找到匹配的预订信息');
      }
    } catch (error: any) {
      console.error('查询预订失败:', error);
      // 处理不同类型的错误
      if (error.response?.status === 400) {
        message.error(error.response?.data?.error || '查询参数无效');
      } else if (error.response?.status === 401) {
        message.error('请先登录后再操作');
      } else if (error.response?.status === 403) {
        message.error('您没有权限执行此操作');
      } else {
        message.error('查询失败，请稍后重试');
      }
    } finally {
      setLoading(false);
    }
  };

  // 办理入住
  const handleCheckIn = async () => {
    if (!bookingInfo) return;
    
    try {
      setSubmitting(true);
      
      // 调用办理入住接口
      // const response = await request(`/api/bookings/${bookingInfo.id}/checkin`, {
      //   method: 'POST',
      // });
      
      // Mock成功响应
      setTimeout(() => {
        message.success('入住办理成功！');
        setCurrentStep(2);
        setSubmitting(false);
      }, 800);
    } catch (error) {
      message.error('入住办理失败，请稍后重试');
      setSubmitting(false);
    }
  };

  // 重置表单并返回第一步
  const handleReset = () => {
    form.resetFields();
    setBookingInfo(null);
    setCurrentStep(0);
  };

  // 渲染第一步：查询表单
  const renderSearchStep = () => (
    <Card>
      <Form form={form} layout="vertical" className="search-form">
        <Form.Item
          name="bookingCode"
          label="预订号"
          rules={[{ required: false }]}
        >
          <Input prefix={<SearchOutlined />} placeholder="暂不支持通过预订号查询" disabled />
        </Form.Item>
        
        <Form.Item
          name="guestName"
          label="客人姓名"
          rules={[{ required: false, message: '请输入客人姓名' }]}
        >
          <Input prefix={<UserOutlined />} placeholder="请输入客人姓名（与手机号至少填一项）" />
        </Form.Item>
        
        <Form.Item
          name="guestPhone"
          label="手机号码"
          rules={[{ required: false, message: '请输入手机号码' }]}
        >
          <Input prefix={<UserOutlined />} placeholder="请输入手机号码（与姓名至少填一项）" />
        </Form.Item>
        
        <Form.Item>
          <Button
            type="primary"
            onClick={handleSearchBooking}
            loading={loading}
            style={{ width: '100%' }}
          >
            查询预订
          </Button>
        </Form.Item>
      </Form>
    </Card>
  );

  // 渲染第二步：确认信息
  const renderConfirmStep = () => (
    <Card>
      {loading ? (
        <Spin size="large" style={{ display: 'block', margin: '40px auto' }} />
      ) : bookingInfo ? (
        <>
          <Descriptions column={1} bordered>
            <Descriptions.Item label="预订号">{bookingInfo.bookingCode || '-'}</Descriptions.Item>
            <Descriptions.Item label="客人姓名">{bookingInfo.guestName || '-'}</Descriptions.Item>
            <Descriptions.Item label="手机号码">{bookingInfo.guestPhone || '-'}</Descriptions.Item>
            <Descriptions.Item label="房间号">{bookingInfo.roomNumber || '-'}</Descriptions.Item>
            <Descriptions.Item label="房间类型">{bookingInfo.roomType || '-'}</Descriptions.Item>
            <Descriptions.Item label="入住日期">{bookingInfo.checkInDate || '-'}</Descriptions.Item>
            <Descriptions.Item label="退房日期">{bookingInfo.checkOutDate || '-'}</Descriptions.Item>
            <Descriptions.Item label="预订状态">
              <Tag color={bookingInfo.status === 'confirmed' ? 'green' : 'blue'}>
                {bookingInfo.status === 'confirmed' ? '已确认' : 
                 bookingInfo.status === 'checkin' ? '已入住' : 
                 bookingInfo.status === 'checkout' ? '已退房' : '未确认'}
              </Tag>
            </Descriptions.Item>
            <Descriptions.Item label="总金额">¥{bookingInfo.totalAmount ? bookingInfo.totalAmount.toFixed(2) : '0.00'}</Descriptions.Item>
          </Descriptions>
          
          <div style={{ marginTop: 24, textAlign: 'center' }}>
            <Button onClick={() => setCurrentStep(0)} style={{ marginRight: 8 }}>
              返回修改
            </Button>
            <Button 
              type="primary" 
              onClick={handleCheckIn}
              loading={submitting}
              disabled={bookingInfo.status !== 'confirmed'}
            >
              确认办理入住
            </Button>
          </div>
        </>
      ) : null}
    </Card>
  );

  // 渲染第三步：成功结果
  const renderSuccessStep = () => (
    <Card>
      <div style={{ textAlign: 'center', padding: 40 }}>
        <CheckCircleOutlined style={{ fontSize: 64, color: '#52c41a' }} />
        <h2 style={{ marginTop: 16, color: '#52c41a' }}>入住办理成功！</h2>
        <p style={{ marginTop: 8, color: '#666' }}>客人 {bookingInfo?.guestName || '未知'} 已成功办理入住</p>
        <p style={{ marginTop: 4, color: '#666' }}>房间号：{bookingInfo?.roomNumber || '未知'}</p>
        <Button type="primary" onClick={handleReset} style={{ marginTop: 24 }}>
          开始新的入住办理
        </Button>
      </div>
    </Card>
  );

  const renderStepContent = () => {
    switch (currentStep) {
      case 0:
        return renderSearchStep();
      case 1:
        return renderConfirmStep();
      case 2:
        return renderSuccessStep();
      default:
        return null;
    }
  };

  return (
    <div className="checkin-form">
      <Steps current={currentStep} items={steps} style={{ marginBottom: 24 }} />
      {renderStepContent()}
    </div>
  );
};

export default CheckInForm;