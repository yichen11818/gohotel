import { PageContainer } from '@ant-design/pro-components';
import { Typography, Button, Row, Col } from 'antd';
import React from 'react';
import { Link } from 'umi';

const { Title, Paragraph } = Typography;

const HotelFrontPage: React.FC = () => {
  return (
    <PageContainer title="酒店前台管理">
      <div style={{ textAlign: 'center', padding: '40px 0' }}>
        <Title level={2}>欢迎使用酒店前台管理系统</Title>
        <Paragraph style={{ margin: '20px 0 40px' }}>
          请选择您要进行的操作
        </Paragraph>
        <Row gutter={[32, 32]} justify="center">
          <Col>
            <Link to="/hotel-front/check-in">
              <Button type="primary" size="large" style={{ width: 200, height: 60, fontSize: 18 }}>
                办理入住
              </Button>
            </Link>
          </Col>
          <Col>
            <Link to="/hotel-front/check-out">
              <Button type="primary" size="large" style={{ width: 200, height: 60, fontSize: 18 }}>
                办理退房
              </Button>
            </Link>
          </Col>
        </Row>
      </div>
    </PageContainer>
  );
};

export default HotelFrontPage;

