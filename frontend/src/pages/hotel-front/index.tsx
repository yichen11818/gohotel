import { PageContainer } from '@ant-design/pro-components';
import { Card } from 'antd';
import React from 'react';
import CheckInForm from './components/CheckInForm';

const HotelFrontPage: React.FC = () => {
  return (
    <PageContainer title="酒店前台管理">
      <Card title="客人入住办理">
        <CheckInForm />
      </Card>
    </PageContainer>
  );
};

export default HotelFrontPage;

