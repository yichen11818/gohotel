import { PageContainer } from '@ant-design/pro-components';
import { Button, Card, Form, Input, InputNumber, Modal, Select, Space, TimePicker, message } from 'antd';
import React from 'react';
import dayjs from 'dayjs';
import { getAdminHotels, getAdminSettings, postAdminHotels, postAdminSettingsSave } from '@/services/api/jiudianguanli';

const SystemSetting: React.FC = () => {
  const [hotels, setHotels] = React.useState<API.Hotel[]>([]);
  const [hotelsLoading, setHotelsLoading] = React.useState(false);
  const [settingsLoading, setSettingsLoading] = React.useState(false);
  const [selectedHotelId, setSelectedHotelId] = React.useState<number | undefined>(undefined);
  const [form] = Form.useForm();
  const [createHotelForm] = Form.useForm();
  const [createHotelVisible, setCreateHotelVisible] = React.useState(false);

  const selectedHotelStorageKey = 'gohotel_admin_selected_hotel_id';

  const loadHotels = async () => {
    setHotelsLoading(true);
    try {
      const resp: any = await getAdminHotels({ page: 1, page_size: 200 });
      const list = Array.isArray(resp?.hotels) ? resp.hotels : [];
      setHotels(list);

      let nextId: number | undefined = undefined;
      const rawSelected = localStorage.getItem(selectedHotelStorageKey);
      const parsed = rawSelected ? Number(rawSelected) : NaN;
      if (!Number.isNaN(parsed) && list.some((h: any) => Number(h?.id) === parsed)) {
        nextId = parsed;
      } else if (list.length > 0) {
        nextId = Number((list[0] as any).id);
      }

      if (nextId !== undefined) {
        setSelectedHotelId(nextId);
      }
    } catch (_e) {
      message.error('加载酒店列表失败');
    } finally {
      setHotelsLoading(false);
    }
  };

  const loadSettings = async (hotelId: number) => {
    setSettingsLoading(true);
    try {
      const resp: any = await getAdminSettings({ hotel_id: hotelId });
      const settings = resp?.data && typeof resp.data === 'object' ? resp.data : {};

      const ci = settings?.booking_rules?.check_in_time;
      const co = settings?.booking_rules?.check_out_time;
      const booking_rules = {
        ...(settings?.booking_rules || {}),
        check_in_time: typeof ci === 'string' && ci ? dayjs(ci, 'HH:mm') : undefined,
        check_out_time: typeof co === 'string' && co ? dayjs(co, 'HH:mm') : undefined,
      };

      form.setFieldsValue({
        ...settings,
        booking_rules,
      });
    } catch (_e) {
      message.error('加载设置失败');
    } finally {
      setSettingsLoading(false);
    }
  };

  React.useEffect(() => {
    loadHotels();
  }, []);

  React.useEffect(() => {
    if (!selectedHotelId) return;
    localStorage.setItem(selectedHotelStorageKey, String(selectedHotelId));
    loadSettings(selectedHotelId);
  }, [selectedHotelId]);

  const onSave = async () => {
    if (!selectedHotelId) {
      message.error('请先选择酒店');
      return;
    }
    const values = await form.validateFields();
    const ci = values?.booking_rules?.check_in_time;
    const co = values?.booking_rules?.check_out_time;

    const payload = {
      ...values,
      booking_rules: {
        ...(values?.booking_rules || {}),
        check_in_time: ci ? dayjs(ci).format('HH:mm') : '',
        check_out_time: co ? dayjs(co).format('HH:mm') : '',
      },
    };

    try {
      await postAdminSettingsSave({
        hotel_id: selectedHotelId,
        settings: payload,
      });
      message.success('保存成功');
      await loadSettings(selectedHotelId);
    } catch (_e) {
      message.error('保存失败');
    }
  };

  const openCreateHotel = () => {
    createHotelForm.resetFields();
    setCreateHotelVisible(true);
  };

  const onCreateHotel = async () => {
    const values = await createHotelForm.validateFields();
    try {
      const created: any = await postAdminHotels({
        name: values.name,
        status: values.status,
      });
      message.success('创建成功');
      setCreateHotelVisible(false);
      await loadHotels();
      const newId = Number(created?.id);
      if (!Number.isNaN(newId) && newId > 0) {
        setSelectedHotelId(newId);
      }
    } catch (_e) {
      message.error('创建失败');
    }
  };

  return (
    <PageContainer>
      <Card>
        <Space style={{ width: '100%', justifyContent: 'space-between' }}>
          <Space>
            <Select
              style={{ width: 260 }}
              loading={hotelsLoading}
              value={selectedHotelId}
              placeholder="选择酒店"
              options={hotels.map((h) => ({ label: `${h.name} (#${h.id})`, value: Number(h.id) }))}
              onChange={(v) => setSelectedHotelId(v)}
            />
            <Button onClick={openCreateHotel}>新增酒店</Button>
          </Space>
          <Button type="primary" loading={settingsLoading} onClick={onSave}>
            保存
          </Button>
        </Space>

        <Form form={form} layout="vertical" style={{ marginTop: 24 }} disabled={!selectedHotelId}>
          <Card size="small" title="酒店基础信息" style={{ marginBottom: 16 }} loading={settingsLoading}>
            <Form.Item label="酒店名称" name={['hotel_profile', 'name']} rules={[{ required: true, message: '请输入酒店名称' }]}>
              <Input />
            </Form.Item>
            <Form.Item label="酒店简称/品牌名" name={['hotel_profile', 'brand_name']}>
              <Input />
            </Form.Item>
            <Form.Item label="酒店LOGO" name={['hotel_profile', 'logo_url']}>
              <Input />
            </Form.Item>
            <Form.Item label="酒店封面图（JSON数组）" name={['hotel_profile', 'cover_images']}>
              <Input.TextArea rows={3} placeholder='例如: ["https://...","https://..."]' />
            </Form.Item>
            <Form.Item label="酒店地址" name={['hotel_profile', 'address']}>
              <Input />
            </Form.Item>
            <Space style={{ width: '100%' }}>
              <Form.Item label="纬度" name={['hotel_profile', 'geo', 'lat']} style={{ flex: 1 }}>
                <InputNumber style={{ width: '100%' }} />
              </Form.Item>
              <Form.Item label="经度" name={['hotel_profile', 'geo', 'lng']} style={{ flex: 1 }}>
                <InputNumber style={{ width: '100%' }} />
              </Form.Item>
            </Space>
            <Form.Item label="联系电话（前台）" name={['hotel_profile', 'front_desk_phone']}>
              <Input />
            </Form.Item>
            <Form.Item label="营业时间/前台服务时间" name={['hotel_profile', 'service_time_text']}>
              <Input />
            </Form.Item>
            <Form.Item label="酒店介绍/公告（富文本HTML）" name={['hotel_profile', 'intro_html']}>
              <Input.TextArea rows={4} />
            </Form.Item>
          </Card>

          <Card size="small" title="预订规则" loading={settingsLoading}>
            <Space style={{ width: '100%' }}>
              <Form.Item label="入住时间" name={['booking_rules', 'check_in_time']} style={{ flex: 1 }}>
                <TimePicker format="HH:mm" style={{ width: '100%' }} />
              </Form.Item>
              <Form.Item label="退房时间" name={['booking_rules', 'check_out_time']} style={{ flex: 1 }}>
                <TimePicker format="HH:mm" style={{ width: '100%' }} />
              </Form.Item>
            </Space>

            <Space style={{ width: '100%' }}>
              <Form.Item label="最早可预订天数" name={['booking_rules', 'min_advance_days']} style={{ flex: 1 }}>
                <InputNumber style={{ width: '100%' }} min={0} />
              </Form.Item>
              <Form.Item label="最远可预订天数" name={['booking_rules', 'max_advance_days']} style={{ flex: 1 }}>
                <InputNumber style={{ width: '100%' }} min={0} />
              </Form.Item>
            </Space>

            <Space style={{ width: '100%' }}>
              <Form.Item label="最少入住晚数" name={['booking_rules', 'min_nights']} style={{ flex: 1 }}>
                <InputNumber style={{ width: '100%' }} min={1} />
              </Form.Item>
              <Form.Item label="最多入住晚数" name={['booking_rules', 'max_nights']} style={{ flex: 1 }}>
                <InputNumber style={{ width: '100%' }} min={1} />
              </Form.Item>
            </Space>

            <Form.Item label="可否当天预订" name={['booking_rules', 'allow_same_day_booking']}>
              <Select
                options={[
                  { label: '允许', value: true },
                  { label: '不允许', value: false },
                ]}
              />
            </Form.Item>
          </Card>
        </Form>

        <Modal
          title="新增酒店"
          open={createHotelVisible}
          onCancel={() => setCreateHotelVisible(false)}
          onOk={onCreateHotel}
          okText="创建"
          cancelText="取消"
        >
          <Form form={createHotelForm} layout="vertical">
            <Form.Item label="酒店名称" name="name" rules={[{ required: true, message: '请输入酒店名称' }]}>
              <Input />
            </Form.Item>
            <Form.Item label="状态" name="status" initialValue="active">
              <Select
                options={[
                  { label: 'active', value: 'active' },
                  { label: 'inactive', value: 'inactive' },
                ]}
              />
            </Form.Item>
          </Form>
        </Modal>
      </Card>
    </PageContainer>
  );
};

export default SystemSetting;

