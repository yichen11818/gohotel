import { PageContainer } from '@ant-design/pro-components';
import { Link } from '@umijs/max';
import {
  Button,
  Card,
  Col,
  List,
  Progress,
  Row,
  Space,
  Statistic,
  Tag,
  Typography,
} from 'antd';
import React from 'react';

const { Paragraph, Text, Title } = Typography;

const stats = [
  {
    title: '预计入住率',
    value: '88%',
    description: '较昨日提升 3%，预计晚间 94% 可售',
    extra: '13:00 入住高峰',
  },
  {
    title: '当前在住',
    value: '128 / 150 间',
    description: '其中包含 8 位会员和 12 位公务团队',
    extra: '可售 22 间（含 5 间升级间）',
  },
  {
    title: '待处理工单',
    value: '9 条',
    description: '维修 3 / 清洁 6，3 条需协同供应商',
    extra: '预计 16:30 前完成',
  },
];

const quickActions = [
  {
    title: '今日订单盘点',
    description: '6 个预订待确认，包含 2 个大床升级与 1 个团体。',
    path: '/order-manage',
    tag: '账务',
  },
  {
    title: '房态库存矩阵',
    description: '查看未来 7 天库存并锁定高价值客房资源。',
    path: '/room-manage/inventory',
    tag: '房务',
  },
  {
    title: '前台接待',
    description: '集中完成 12 个入住与 8 个退房流程，掌控账务同步。',
    path: '/hotel-front/check-in',
    tag: '前台',
  },
  {
    title: '服务与工单',
    description: '跟进维修记录与清洁任务，重点关注冷气与水压反馈。',
    path: '/work-order/repair',
    tag: '服务',
  },
];

const dailyReminders = [
  {
    title: '收益提醒',
    content: '今日下午 16:00 前调整假日价格策略，确保潜在增长 6%。',
  },
  {
    title: '运营敏捷',
    content: '确认所有 VIP 入住偏好录入 OMS，避免重复沟通与投诉。',
  },
  {
    title: '系统同步',
    content: '检查后端日志，已知订单回写延迟在 22:00 前需复核。',
  },
];

const OperationCard: React.FC<{
  title: string;
  description: string;
  path: string;
  tag: string;
}> = ({ title, description, path, tag }) => (
  <Card bordered={false} style={{ borderRadius: 10, minHeight: 170 }}>
    <Space direction="vertical" size={12} style={{ width: '100%' }}>
      <Tag color="processing">{tag}</Tag>
      <Title level={5} style={{ margin: 0 }}>
        {title}
      </Title>
      <Text type="secondary">{description}</Text>
      <Button type="primary" size="small">
        <Link to={path}>打开</Link>
      </Button>
    </Space>
  </Card>
);

const Welcome: React.FC = () => (
  <PageContainer
    title="酒店运营总览"
    subTitle="跨前台、房务、营运与收益的实时指引"
  >
    <Card
      bordered={false}
      style={{
        borderRadius: 10,
        marginBottom: 16,
        background:
          'linear-gradient(135deg, rgba(255,255,255,0.96) 0%, rgba(236,245,255,0.98) 55%, rgba(226,238,252,0.95) 100%)',
        color: '#17324a',
        border: '1px solid rgba(11, 95, 146, 0.12)',
        boxShadow: '0 16px 36px rgba(16, 51, 82, 0.08)',
      }}
    >
      <Row gutter={[24, 16]}>
        <Col xs={24} lg={16}>
          <Title
            level={3}
            style={{ marginBottom: 12, color: '#17324a', fontWeight: 600 }}
          >
            今日运营主线
          </Title>
          <Paragraph
            style={{
              color: '#5d7286',
              maxWidth: 680,
              marginBottom: 16,
            }}
          >
            以“接待 · 房态 · 服务 · 收益”四步闭环贯穿当天工作。先确认预订与入住，再锁定可售房源，随后调度服务工单，最终评估价格与宣传效果。
          </Paragraph>
          <Space wrap>
            <Button type="primary">
              <Link to="/order-manage">查看今日订单</Link>
            </Button>
            <Button type="default">
              <Link to="/room-manage/inventory">立即查看房态矩阵</Link>
            </Button>
          </Space>
        </Col>
        <Col xs={24} lg={8}>
          <Space direction="vertical" size={12} style={{ width: '100%' }}>
            <Statistic
              title="预计入住率"
              value="88%"
              valueStyle={{ color: '#17324a' }}
              suffix={<Text style={{ color: '#6f869c' }}>计划 150 间</Text>}
            />
            <Progress percent={88} status="active" showInfo={false} strokeColor="#0b5f92" />
            <Text style={{ color: '#6f869c', fontSize: 12 }}>
              13:00~15:00 为今日入住高峰，附带加床需求 2 单。
            </Text>
          </Space>
        </Col>
      </Row>
    </Card>

    <Row gutter={[16, 16]} style={{ marginBottom: 16 }}>
      {stats.map((stat) => (
        <Col xs={24} sm={12} md={8} key={stat.title}>
          <Card bordered={false} style={{ borderRadius: 10, minHeight: 160 }}>
            <Statistic title={stat.title} value={stat.value} />
            <Text type="secondary">{stat.description}</Text>
            <div style={{ marginTop: 12 }}>
              <Text type="secondary" style={{ fontSize: 12 }}>
                {stat.extra}
              </Text>
            </div>
          </Card>
        </Col>
      ))}
    </Row>

    <Row gutter={[16, 16]} style={{ marginBottom: 16 }}>
      {quickActions.map((action) => (
        <Col xs={24} sm={12} lg={6} key={action.title}>
          <OperationCard {...action} />
        </Col>
      ))}
    </Row>

    <Row gutter={[16, 16]}>
      <Col xs={24} lg={14}>
        <Card
          bordered={false}
          style={{ borderRadius: 10, minHeight: 220 }}
          title="今日提醒"
        >
          <List
            split={false}
            dataSource={dailyReminders}
            renderItem={(item) => (
              <List.Item>
                <List.Item.Meta
                  title={<Text strong>{item.title}</Text>}
                  description={<Text type="secondary">{item.content}</Text>}
                />
              </List.Item>
            )}
          />
        </Card>
      </Col>
      <Col xs={24} lg={10}>
        <Card
          bordered={false}
          style={{ borderRadius: 10, minHeight: 220 }}
          title="房态与服务看板"
        >
          <Space direction="vertical" size={12}>
            <Text>房型占比：豪华套房 28% · 标准房 52% · 商务间 20%</Text>
            <Text>
              本周主题活动：周末亲子套餐（已上线）、会员点亮计划（待发布）
            </Text>
            <Text>
              设备告警： 默认冷气 2 号机组信号不稳定，已派人复核。
            </Text>
          </Space>
        </Card>
      </Col>
    </Row>
  </PageContainer>
);

export default Welcome;
